CREATE TABLE review (
                        id text PRIMARY KEY,
                        reviewed text,
                        reviewer text,
                        created_at timestamp,
                        FOREIGN KEY(reviewed)
                            REFERENCES student(id),
                        FOREIGN KEY(reviewer)
                            REFERENCES student(id)
);

CREATE TABLE review_tag (
                            review_id text,
                            tag_name text,
                            PRIMARY KEY (review_id, tag_name),
                            FOREIGN KEY (review_id)
                                REFERENCES review(id),
                            FOREIGN KEY (tag_name)
                                REFERENCES tag(name)
);