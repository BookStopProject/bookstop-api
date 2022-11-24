CREATE TABLE public."book" (
    id serial PRIMARY KEY,
    title varchar(256) NOT NULL,
    subtitle varchar(256),
    description text,
    published_year integer NOT NULL,
    genre_id integer,
    tradein_credit integer NOT NULL,
    exchange_credit integer NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES public."genre" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."book_author" (
    book_id integer,
    author_id integer,
    PRIMARY KEY (book_id, author_id),
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (author_id) REFERENCES public."author" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."author" (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    description varchar(160),
    date_of_birth date varying(200) NOT NULL,
    date_of_death date varying(200) NOT NULL,
    nationality varchar(200) NOT NULL
);

CREATE TABLE public."genre" (
    id serial PRIMARY KEY,
    name varchar(64) NOT NULL,
    description text NOT NULL
);

CREATE TYPE book_condition AS ENUM (
    'new',
    'like_new',
    'very_good',
    'acceptable'
);

CREATE TABLE pubic."book_copy" (
    id serial PRIMARY KEY,
    book_id integer,
    condition book_condition,
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."user" (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    bio varchar(160),
    profile_picture varchar,
    creation_time timestamp WITHOUT time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    credit integer NOT NULL DEFAULT 0
);

CREATE TABLE public."user_book" (
    id serial PRIMARY KEY,
    book_id integer NOT NULL,
    user_id integer NOT NULL,
    start_date date,
    end_date date,
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES public."user" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);
CREATE TABLE public."post" (
    id serial PRIMARY KEY,
    text varchar(250),
    creation_time timestamp NOT NULL,
    book_id integer NOT NULL,
    user_id integer NOT NULL,
    is_recommending boolean NOT NULL,
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES public."user" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."location" (
    id serial PRIMARY KEY name varchar(128) NOT NULL,
    address varchar(512) NOT NULL,
);

CREATE TABLE public."exchange" (
    id serial PRIMARY KEY,
    book_copy_id integer,
    previous_user_id integer,
    next_user_id integer,
    exchange_time timestamp WITHOUT time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    FOREIGN KEY (previous_user_id) REFERENCES public."user" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (next_user_id) REFERENCES public."user" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."inventory_entry" ();

CREATE TABLE public."event_inventory_entry" ();

CREATE TABLE public."invoice" ();

CREATE TABLE public."invoice_entry" ();

