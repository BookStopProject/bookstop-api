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
    PRIMARY KEY (id),
    FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public.user_book
    ADD COLUMN id_original integer;

ALTER TABLE user_book 
	ADD CONSTRAINT user_book_book_id_original_fkey 
	FOREIGN KEY (id_original) 
	REFERENCES user_book (id);

CREATE TABLE public.exchange
(
    id serial NOT NULL,
    user_id_from integer NOT NULL,
    user_id_to integer NOT NULL,
    user_book_id_from integer NOT NULL,
    user_book_id_original integer NOT NULL,
    exchange_time timestamp without time zone NOT NULL,
    exchange_loc integer NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id_from)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    FOREIGN KEY (user_id_to)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID,
    FOREIGN KEY (user_book_id_from)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (user_book_id_original)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (exchange_loc)
        REFERENCES public.location (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);