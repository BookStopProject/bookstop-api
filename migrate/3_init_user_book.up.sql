CREATE TABLE public.location
(
    id serial NOT NULL,
    name character varying(100) NOT NULL,
    parent_name character varying(200),
    address_line character varying NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE public.user_book
(
    id serial NOT NULL,
    user_id integer NOT NULL,
    book_id character varying NOT NULL,
    started_at date,
    ended_at date,
    id_original integer,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (id_original)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
