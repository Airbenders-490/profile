create table school
(
    id      text not null
        constraint school_pkey
            primary key,
    name    text,
    country text,
    domains text[]
);

alter table school
    owner to postgres;

CREATE TABLE public.confirmation (
    token text PRIMARY KEY NOT NULL,
    st_id text not null REFERENCES student(id),
    sc_id text not null REFERENCES school(id),
    created_at timestamp
);