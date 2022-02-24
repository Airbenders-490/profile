-- Table: public.student

-- DROP TABLE public.student;

CREATE TABLE IF NOT EXISTS public.student
(
    id character varying(64) COLLATE pg_catalog."default" NOT NULL,
    first_name character varying(64) COLLATE pg_catalog."default",
    last_name character varying(64) COLLATE pg_catalog."default",
    email character varying(64) COLLATE pg_catalog."default",
    current_classes text[] COLLATE pg_catalog."default",
    classes_taken text[] COLLATE pg_catalog."default",
    general_info character varying(1024) COLLATE pg_catalog."default",
    school character varying(64) COLLATE pg_catalog."default",
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT student_pkey PRIMARY KEY (id),
    CONSTRAINT school_pkey FOREIGN KEY (school)
        REFERENCES public.school (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.student
    OWNER to postgres;

INSERT INTO public.student (id, first_name, last_name, email, general_info, created_at, updated_at)
VALUES ('234', 'Also Zubair', 'Nurie', 'mzznurie@msn.com', 'ballerr', now(), now());
INSERT INTO public.student (id, first_name, last_name, email, general_info, created_at, updated_at)
VALUES ('123', 'Zubair', 'Nurie', 'mznurie@msn.com', 'baller', now(), now());
