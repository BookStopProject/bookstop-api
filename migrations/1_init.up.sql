CREATE TABLE public."user"
(
    id serial NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    oauth_id character varying(50) NOT NULL,
    email character varying(256),
    name character varying(100) NOT NULL,
    description character varying(160),
    "profile_image_url" character varying,
    credit integer NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);

CREATE TABLE public.browse
(
    id serial NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(160),
    started_at timestamp without time zone NOT NULL,
    ended_at timestamp without time zone NOT NULL,
    image_url character varying,
    PRIMARY KEY (id)
);

CREATE TABLE public.browse_book
(
    book_id character varying(21) NOT NULL,
    browse_id integer NOT NULL,
    FOREIGN KEY (browse_id)
        REFERENCES public.browse (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

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
        NOT VALID,
    UNIQUE(user_book_id)
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