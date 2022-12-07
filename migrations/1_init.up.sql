CREATE TABLE public."author" (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    description varchar(160),
    date_of_birth date NULL,
    date_of_death date NULL
);

CREATE TABLE public."genre" (
    id serial PRIMARY KEY,
    name varchar(64) NOT NULL,
    description text NOT NULL
);

CREATE TABLE public."book" (
    id serial PRIMARY KEY,
    author_id integer NOT NULL,
    title varchar(256) NOT NULL,
    subtitle varchar(256),
    image_url varchar(256),
    description text,
    published_year integer NOT NULL,
    genre_id integer NOT NULL,
    tradein_credit integer NOT NULL,
    exchange_credit integer NOT NULL,
    FOREIGN KEY (author_id) REFERENCES public."author" (id),
    FOREIGN KEY (genre_id) REFERENCES public."genre" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TYPE book_condition AS ENUM (
    'new',
    'like_new',
    'good',
    'acceptable'
);

CREATE TABLE public."user" (
    id serial PRIMARY KEY,
    oauth_id varchar(50) NOT NULL,
    name varchar(100) NOT NULL,
    bio varchar(160),
    profile_picture varchar,
    creation_time timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    credit integer NOT NULL DEFAULT 0
);

CREATE TABLE public."post" (
    id serial PRIMARY KEY,
    text varchar(250),
    creation_time timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    book_id integer NOT NULL,
    user_id integer NOT NULL,
    is_recommending boolean NOT NULL,
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES public."user" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."location" (
    id serial PRIMARY KEY,
    name varchar(128) NOT NULL,
    address varchar(512) NOT NULL
);

CREATE TABLE public."book_copy" (
    id serial PRIMARY KEY,
    book_id integer,
    condition book_condition,
    location_id integer,
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (location_id) REFERENCES public."location" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
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

CREATE TABLE public."event" (
    id serial PRIMARY KEY,
    name varchar(128) NOT NULL,
    description varchar(512) NOT NULL,
    start_time timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    end_time timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    location_id integer,
    FOREIGN KEY (location_id) REFERENCES public."location" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."event_book_copy" (
    id serial PRIMARY KEY,
    event_id integer,
    book_copy_id integer,
    FOREIGN KEY (event_id) REFERENCES public."event" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (book_copy_id) REFERENCES public."book_copy" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."invoice" (
    id serial PRIMARY KEY,
    user_id integer,
    location_id integer,
    creation_time timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc'),
    FOREIGN KEY (user_id) REFERENCES public."user" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (location_id) REFERENCES public."location" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."invoice_entry" (
    invoice_id integer,
    credit integer,
    book_copy_id integer,
    PRIMARY KEY (invoice_id, book_copy_id),
    FOREIGN KEY (invoice_id) REFERENCES public."invoice" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (book_copy_id) REFERENCES public."book_copy" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

CREATE TABLE public."browse" (
    id serial PRIMARY KEY,
    name varchar(128) NOT NULL,
    description varchar(512) NULL
);

CREATE TABLE public."browse_book" (
    browse_id integer,
    book_id integer,
    PRIMARY KEY (browse_id, book_id),
    FOREIGN KEY (browse_id) REFERENCES public."browse" (id) ON UPDATE NO ACTION ON DELETE RESTRICT,
    FOREIGN KEY (book_id) REFERENCES public."book" (id) ON UPDATE NO ACTION ON DELETE RESTRICT
);

