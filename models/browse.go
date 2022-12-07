package models

import (
	"bookstop/db"
	"context"
)

type Browse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FindAllBrowses(ctx context.Context) ([]*Browse, error) {
	rows, err := db.Conn.Query(ctx, "SELECT id, name, description FROM public.browse")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	browses := make([]*Browse, 0)
	for rows.Next() {
		browse := new(Browse)
		err := rows.Scan(&browse.ID, &browse.Name, &browse.Description)
		if err != nil {
			return nil, err
		}
		browses = append(browses, browse)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return browses, nil
}

func FindBrowseBooks(ctx context.Context, browseID int) ([]*Book, error) {
	// select from browse_book table by id
	// join with book table, then join with author table
	rows, err := db.Conn.Query(ctx, `SELECT 
		b.id,
		b.title,
		b.subtitle,
		b.image_url,
		a.id,
		a.name
	FROM browse_book bb
	JOIN book b ON bb.book_id = b.id
	JOIN author a ON b.author_id = a.id
	WHERE bb.browse_id = $1`, browseID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []*Book{}
	for rows.Next() {
		book := &Book{}
		book.Author = &Author{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Subtitle,
			&book.ImageURL,
			&book.Author.ID,
			&book.Author.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func CreateBrowse(ctx context.Context, browse *Browse) (*Browse, error) {
	err := db.Conn.QueryRow(ctx, "INSERT INTO public.browse (name, description) VALUES ($1, $2) RETURNING id",
		browse.Name, browse.Description).Scan(&browse.ID)
	if err != nil {
		return nil, err
	}
	return browse, nil
}

func UpdateBrowse(ctx context.Context, browse *Browse) (*Browse, error) {
	_, err := db.Conn.Exec(ctx, "UPDATE public.browse SET name = $1, description = $2 WHERE id = $3",
		browse.Name, browse.Description, browse.ID)
	if err != nil {
		return nil, err
	}
	return browse, nil
}

func DeleteBrowse(ctx context.Context, browseID int) error {
	_, err := db.Conn.Exec(ctx, "DELETE FROM public.browse WHERE id = $1", browseID)
	if err != nil {
		return err
	}
	return nil
}

func AddBookToBrowse(ctx context.Context, browseID int, bookID int) error {
	_, err := db.Conn.Exec(ctx, "INSERT INTO public.browse_book (browse_id, book_id) VALUES ($1, $2)",
		browseID, bookID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveBookFromBrowse(ctx context.Context, browseID int, bookID int) error {
	_, err := db.Conn.Exec(ctx, "DELETE FROM public.browse_book WHERE browse_id = $1 AND book_id = $2",
		browseID, bookID)
	if err != nil {
		return err
	}
	return nil
}
