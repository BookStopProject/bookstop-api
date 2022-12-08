-- Sale people can work with data related to sale
-- and perform trade ins and exchanges of books.
-- Therefore they will need to have access to the following tables:
--     user
--     book
--     book_copy
--     user_book
--     invoice
--     invoice_entry
--     trade_in
CREATE USER sale WITH PASSWORD 'sale';

REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM sale;

GRANT ALL PRIVILEGES ON TABLE public."user" TO sale;

GRANT ALL PRIVILEGES ON TABLE public."book" TO sale;

GRANT ALL PRIVILEGES ON TABLE public."book_copy" TO sale;

GRANT ALL PRIVILEGES ON TABLE public."user_book" TO sale;

GRANT ALL PRIVILEGES ON TABLE public."invoice" TO sale;

GRANT ALL PRIVILEGES ON TABLE public."invoice_entry" TO sale;

GRANT ALL PRIVILEGES ON TABLE public."trade_in" TO sale;

-- Editors can work with data related to website content (browse)
-- and organize events
-- Therefore they will need to have access to the following tables:
--     book
--     genre
--     author
--     location
--     event
--     event_book_copy
--     browse
--     browse_book
-- They should not have access to user or sale related tables
CREATE USER editor WITH PASSWORD 'editor';

REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM editor;

GRANT ALL PRIVILEGES ON TABLE public."book" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."genre" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."author" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."location" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."event" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."event_book_copy" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."browse" TO editor;

GRANT ALL PRIVILEGES ON TABLE public."browse_book" TO editor;

