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