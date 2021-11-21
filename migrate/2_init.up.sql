CREATE TABLE public.thought
(
    id serial NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    text character varying(300) NOT NULL,
    book_id character varying(21)
);