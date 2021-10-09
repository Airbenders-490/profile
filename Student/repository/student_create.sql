CREATE TABLE IF NOT EXISTS public.student
(
    id character varying(64) COLLATE pg_catalog."default" NOT NULL,
    first_name character varying(64) COLLATE pg_catalog."default",
    last_name character varying(64) COLLATE pg_catalog."default",
    email character varying(64) COLLATE pg_catalog."default",
    general_info character varying(1024) COLLATE pg_catalog."default",
    school character varying(64) COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT student_pkey PRIMARY KEY (id)
)