-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    uid uuid NOT NULL DEFAULT uuid_generate_v4(),
    username character varying(30) COLLATE pg_catalog."default" NOT NULL,
    email character varying(254) COLLATE pg_catalog."default" NOT NULL,
    password_hash text COLLATE pg_catalog."default" NOT NULL,
    role integer NOT NULL DEFAULT 0,
    created_at date NOT NULL DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (uid),
    CONSTRAINT users_email_key UNIQUE (email)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;