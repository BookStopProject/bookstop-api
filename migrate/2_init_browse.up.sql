CREATE TABLE public.browse
(
    id serial NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(160),
    "from" timestamp without time zone NOT NULL,
    "to" timestamp without time zone NOT NULL,
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