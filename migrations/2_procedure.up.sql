-- A procedure that accepts a user book id and condition
-- and returns user book information and a credit value for the book based on the condition
CREATE OR REPLACE FUNCTION get_user_book_and_credit (user_book_id integer, condition text)
  RETURNS TABLE (
    book_id integer,
    book_title varchar,
    book_subtitle varchar,
    book_credit integer
  )
  AS $$
DECLARE
  user_book RECORD;
  book RECORD;
  credit integer;
BEGIN
  SELECT
    * INTO user_book
  FROM
    public."user_book"
  WHERE
    id = user_book_id;
  SELECT
    * INTO book
  FROM
    public."book"
  WHERE
    id = user_book.book_id;
  IF condition = 'new' THEN
    credit = 1.0 * book.tradein_credit;
  ELSIF condition = 'like_new' THEN
    credit = 0.9 * book.tradein_credit;
  ELSIF condition = 'good' THEN
    credit = 0.8 * book.tradein_credit;
  ELSIF condition = 'acceptable' THEN
    credit = 0.5 * book.tradein_credit;
  END IF;
  book_id = book.id;
  book_title = book.title;
  book_subtitle = book.subtitle;
  book_credit = credit;
  RETURN NEXT;
END;
$$
LANGUAGE 'plpgsql';

-- // do_trade_in
-- // This procedure should:
-- // 1) Create a book copy if user book does not have one and link it to the user book
-- // 2) Update the book copy condition and location
-- // 3) Create a trade in for that book copy. The credit will be equal to
-- //		the book trade in value * condition multiplier (see book_copy.go).
-- // 4) Add the credit to the user's credit balance
-- // 5) Return the trade in id
CREATE OR REPLACE PROCEDURE do_trade_in (IN user_book_id integer, IN curr_condition book_condition, IN curr_location_id integer)
  AS $$
DECLARE
  curr_book_copy_id integer;
  book_copy RECORD;
  user_book RECORD;
  book RECORD;
  trade_in_id integer;
  trade_in_credit integer;
BEGIN
  SELECT
    * INTO user_book
  FROM
    public."user_book"
  WHERE
    id = user_book_id;
  SELECT
    * INTO book
  FROM
    public."book"
  WHERE
    id = user_book.book_id;
  -- Create a book copy if user book does not have one and link it to the user book
  IF user_book.book_copy_id IS NULL THEN
    INSERT INTO public."book_copy" (book_id, condition, location_id)
      VALUES (user_book.book_id, curr_condition, curr_location_id)
    RETURNING
      id INTO curr_book_copy_id;
    UPDATE
      public."user_book"
    SET
      book_copy_id = curr_book_copy_id
    WHERE
      id = user_book_id;
    -- Update the book copy condition and location
  ELSE
    curr_book_copy_id = user_book.book_copy_id;
    UPDATE
      public."book_copy"
    SET
      condition = curr_condition,
      location_id = curr_location_id
    WHERE
      id = curr_book_copy_id;
  END IF;
  -- Calculate credit based on condition
  IF curr_condition = 'new' THEN
    trade_in_credit = 1.0 * book.tradein_credit;
  ELSIF curr_condition = 'like_new' THEN
    trade_in_credit = 0.9 * book.tradein_credit;
  ELSIF curr_condition = 'good' THEN
    trade_in_credit = 0.8 * book.tradein_credit;
  ELSIF curr_condition = 'acceptable' THEN
    trade_in_credit = 0.5 * book.tradein_credit;
  END IF;
  -- Create a trade in for that book copy
  INSERT INTO public."trade_in" (user_id, book_copy_id, credit)
    VALUES (user_book.user_id, curr_book_copy_id, trade_in_credit)
  RETURNING
    id INTO trade_in_id;
  -- Add the credit to the user's credit balance
  UPDATE
    public."user"
  SET
    credit = credit + trade_in_credit
  WHERE
    id = user_book.user_id;
END;
$$
LANGUAGE 'plpgsql';

-- // do_exchange
-- // 1) Verify that book copies are available at locations (have location_id)
-- // 2) Create an invoice
-- // 3) For each book copy, create an invoice entry. The invoice entry credit
-- // 		will be equal to the book exchange price * book copy condition multiplier (see book_copy.go)).
-- // 4) For each book copy, update the book copy location_id to nil.
-- // 5) For each book copy, create a user book with the book copy id and user id.
-- // 6) Calculate the total credit of the invoice. Deduct the total credit from the user's balance
-- //		(must verify that user has enough balance)
-- // 7) Return the invoice.
CREATE OR REPLACE PROCEDURE do_exchange (IN user_id integer, IN book_copy_ids integer[])
  AS $$
DECLARE
  invoice_total_credit integer;
  curr_credit integer;
  book_copy_id integer;
  book_copy RECORD;
  book RECORD;
  usr RECORD;
  invoice_id integer;
BEGIN
  invoice_total_credit = 0;
  -- Create an invoice
  INSERT INTO public."invoice" (user_id)
    VALUES (user_id)
  RETURNING
    id INTO invoice_id;
  FOREACH book_copy_id IN ARRAY book_copy_ids LOOP
    SELECT
      * INTO book_copy
    FROM
      public."book_copy"
    WHERE
      id = book_copy_id;
    SELECT
      * INTO book
    FROM
      public."book"
    WHERE
      id = book_copy.book_id;
    -- Verify that book copies are available at locations (have location_id)
    IF book_copy.location_id IS NULL THEN
      RAISE EXCEPTION 'Book copy % is not available at a location', book_copy_id;
    END IF;
    IF book_copy.condition = 'new' THEN
      curr_credit = 1.0 * book.exchange_credit;
    ELSIF book_copy.condition = 'like_new' THEN
      curr_credit = 0.9 * book.exchange_credit;
    ELSIF book_copy.condition = 'good' THEN
      curr_credit = 0.8 * book.exchange_credit;
    ELSIF book_copy.condition = 'acceptable' THEN
      curr_credit = 0.5 * book.exchange_credit;
    END IF;
    invoice_total_credit + = curr_credit;
    -- For each book copy, create an invoice entry.
    INSERT INTO public."invoice_entry" (invoice_id, book_copy_id, credit)
      VALUES (invoice_id, book_copy_id, curr_credit);
    -- For each book copy, update the book copy location_id to nil.
    UPDATE
      public."book_copy"
    SET
      location_id = NULL
    WHERE
      id = book_copy_id;
    -- For each book copy, create a user book with the book copy id and user id.
    INSERT INTO public."user_book" (user_id, book_id, book_copy_id)
      VALUES (user_id, book_copy.book_id, book_copy_id);
  END LOOP;
  -- Get the current credit of user
  SELECT
    * INTO usr
  FROM
    public."user"
  WHERE
    id = user_id;
  -- Throw an exception if user credit is not enough
  IF usr.credit < invoice_total_credit THEN
    RAISE EXCEPTION 'User % does not have enough credit', user_id;
  END IF;
  -- Deduct the total credit from the user's balance
  UPDATE
    public."user"
  SET
    credit = credit - invoice_total_credit
  WHERE
    id = user_id;
END;
$$
LANGUAGE 'plpgsql';

