package httputils

type ValidResponse struct {
	Message string `json:"message"`
}

func NewResponse(message string) ValidResponse {
	return ValidResponse{Message: message}
}
