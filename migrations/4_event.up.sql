CREATE TABLE public.event
(
    id serial NOT NULL,
    title character varying(200) NOT NULL,
    description character varying(2000) NOT NULL,
    href character varying NOT NULL,
    user_id integer NOT NULL,
    started_at timestamp without time zone NOT NULL,
    ended_at timestamp without time zone NOT NULL,
    PRIMARY KEY (id)
);