package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/airbenders/profile/domain"
	"github.com/airbenders/profile/utils/channelmocks"
	"github.com/airbenders/profile/utils/errors"
	"github.com/streadway/amqp"
	"log"
	"reflect"
	"time"
)


const (
	errorMessage = "No such student with ID %s exists"
	existingStudentError = "Student with ID %s already exists"
	contentType = "text/plain"
	publishErrorMessage = "failed to publish "
	ampqMessageSent     = "student sent to queue"
)


type studentUseCase struct {
	studentRepository domain.StudentRepository
	// add review and possibly tag repositories to fetch their data too
	reviewRepository domain.ReviewRepository
	contextTimeout   time.Duration
	messagingManager *MessagingManager
}

type MessagingManager struct {
	Ch      mocks.Channel
	Created chan domain.Student
	Edited  chan domain.Student
	Deleted chan string
}

func NewMessagingManager(ch mocks.Channel) *MessagingManager {
	return &MessagingManager{
		Ch:      ch,
		Created: make(chan domain.Student),
		Edited:  make(chan domain.Student),
		Deleted: make(chan string),
	}

}

// NewStudentUseCase returns a configured StudentUseCase
func NewStudentUseCase(mm *MessagingManager, sr domain.StudentRepository, rr domain.ReviewRepository, timeout time.Duration) domain.StudentUseCase {
	return &studentUseCase{
		messagingManager: mm,
		studentRepository: sr,
		reviewRepository:  rr,
		contextTimeout:    timeout,
	}
}

// Create stores the student in the db. ID must be provided or return an error
func (s *studentUseCase) Create(c context.Context, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	if st.ID != "" {
		existingStudent, err := s.studentRepository.GetByID(ctx, st.ID)
		if err == nil {
			if !reflect.DeepEqual(*existingStudent, domain.Student{}) {
				return errors.NewConflictError(fmt.Sprintf(existingStudentError, st.ID))
			}
		}
	} else {
		return errors.NewBadRequestError("The student should have an ID from auth service")
	}

	st.CreatedAt = time.Now()
	st.UpdatedAt = time.Now()
	err := s.studentRepository.Create(ctx, st.ID, st)
	if err != nil {
		return err
	}
	s.messagingManager.Created <- *st
	return nil
}

// CreateStudentTopic sends messages for update of student
func (s *studentUseCase) CreateStudentTopic() {
	for student := range s.messagingManager.Created {
		st, err := json.Marshal(student)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = s.messagingManager.Ch.Publish(
			"profile",
			"profile.created",
			false,
			false,
			amqp.Publishing{
				ContentType: contentType,
				Body: st,
			})
		if err != nil {
			log.Println(publishErrorMessage, err)
		}
		log.Println(ampqMessageSent)
	}
}

// GetByID seeks student from repo layer and returns if it exists, else return error
func (s *studentUseCase) GetByID(c context.Context, id string) (*domain.Student, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	student, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(student, &domain.Student{}) {

		return nil, errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}

	reviews, err := s.reviewRepository.GetReviewsFor(ctx, student.ID)
	if err != nil {
		log.Println("Can't get the reviews right now.")
	}
	student.Reviews = reviews

	return student, nil
}

// Update checks if the student exists and updates if so. Otherwise, returns error
func (s *studentUseCase) Update(c context.Context, id string, st *domain.Student) (*domain.Student, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	st.ID = id
	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return nil, errors.NewNotFoundError(fmt.Sprintf(errorMessage, st.ID))
	}
	updateStudent(existingStudent, st)
	err = s.studentRepository.Update(ctx, existingStudent)
	if err != nil {
		return nil, err
	}
	s.messagingManager.Edited <- *existingStudent

	return existingStudent, nil
}

// UpdateStudentTopic sends messages for creation of student
func (s *studentUseCase) UpdateStudentTopic() {
	for student := range s.messagingManager.Edited {
		st, err := json.Marshal(student)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = s.messagingManager.Ch.Publish(
			"profile",
			"profile.updated",
			false,
			false,
			amqp.Publishing{
				ContentType: contentType,
				Body: st,
			})
		if err != nil {
			log.Println(publishErrorMessage, err)
		}
		log.Println(ampqMessageSent)
	}
}

func updateStudent(existing *domain.Student, toUpdate *domain.Student) {
	if toUpdate.FirstName != "" {
		existing.FirstName = toUpdate.FirstName
	}
	if toUpdate.LastName != "" {
		existing.LastName = toUpdate.LastName
	}
	if toUpdate.Email != "" {
		existing.Email = toUpdate.Email
	}
	if toUpdate.GeneralInfo != "" {
		existing.GeneralInfo = toUpdate.GeneralInfo
	}
	existing.UpdatedAt = time.Now()
}

// Delete removes the student if it exists. Otherwise, returns error
func (s *studentUseCase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}

	err = s.studentRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.messagingManager.Deleted <- id
	return nil
}

// DeleteStudentTopic sends messages for update of student
func (s *studentUseCase) DeleteStudentTopic() {
	for studentID := range s.messagingManager.Deleted {
		err := s.messagingManager.Ch.Publish(
			"profile",
			"profile.Deleted",
			false,
			false,
			amqp.Publishing{
				ContentType: contentType,
				Body: []byte(studentID),
			})
		if err != nil {
			log.Println(publishErrorMessage, err)
		}
		log.Println(ampqMessageSent)
	}
}

func (s *studentUseCase) AddClasses(c context.Context, id string, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	st.ID = id
	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}

	st.CurrentClasses = removeDuplicates(append(existingStudent.CurrentClasses, st.CurrentClasses...))
	st.ClassesTaken = removeDuplicates(append(existingStudent.ClassesTaken, st.ClassesTaken...))
	return s.studentRepository.UpdateClasses(ctx, st)
}

func removeDuplicates(slice []string) []string {
	uniques := make(map[string]bool)
	ret := []string{}

	for _, el := range slice {
		if _, value := uniques[el]; !value {
			uniques[el] = true
			ret = append(ret, el)
		}
	}

	return ret
}

func (s *studentUseCase) RemoveClasses(c context.Context, id string, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	st.ID = id
	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}
	st.CurrentClasses = removeClasses(existingStudent.CurrentClasses, st.CurrentClasses)
	st.ClassesTaken = removeClasses(existingStudent.ClassesTaken, st.ClassesTaken)
	return s.studentRepository.UpdateClasses(ctx, st)
}

func removeClasses(existingClasses, classesToRemove []string) []string {
	classRemove := make(map[string]struct{}, len(classesToRemove))
	for _, x := range classesToRemove {
		classRemove[x] = struct{}{}
	}
	var remainingClasses []string
	for _, i := range existingClasses{
		if _, found := classRemove[i]; !found {
			remainingClasses = append(remainingClasses, i)
		}
	}
	return remainingClasses
}

func (s *studentUseCase) CompleteClass(c context.Context, id string, st *domain.Student) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	st.ID = id
	existingStudent, err := s.studentRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(existingStudent, &domain.Student{}) {
		return errors.NewNotFoundError(fmt.Sprintf(errorMessage, id))
	}
	completedClasses := st.CurrentClasses
	st.CurrentClasses = removeClasses(existingStudent.CurrentClasses, st.CurrentClasses)
	st.ClassesTaken = removeDuplicates(append(existingStudent.ClassesTaken, completedClasses...))
	return s.studentRepository.UpdateClasses(ctx, st)
}
