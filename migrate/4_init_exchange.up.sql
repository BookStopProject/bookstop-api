CREATE TABLE public.inventory
(
    id serial NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    removed_at timestamp without time zone,
    user_book_id integer NOT NULL,
    location_id integer NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_book_id)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (location_id)
        REFERENCES public.location (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE public.exchange
(
    id serial NOT NULL,
    user_book_id_old integer NOT NULL,
    user_book_id_new integer NOT NULL,
    user_book_id_original integer NOT NULL,
    exchanged_at timestamp without time zone NOT NULL,
    location_id integer,
    PRIMARY KEY (id),
    FOREIGN KEY (user_book_id_old)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (user_book_id_new)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (user_book_id_original)
        REFERENCES public.user_book (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (location_id)
        REFERENCES public.location (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE public.inventory_claim
(
    id serial NOT NULL,
    user_id integer NOT NULL,
    inventory_id integer NOT NULL,
    claimed_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    FOREIGN KEY (inventory_id)
        REFERENCES public.inventory (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);