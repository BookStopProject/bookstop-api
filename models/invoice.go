package models

import (
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Invoice struct {
	ID           int       `json:"id"`
	CreationTime time.Time `json:"creationTime"`
	UserID       int       `json:"userId"`
	User         *User     `json:"user"`
}

type InvoiceEntry struct {
	InvoiceID  int `json:"invoiceId"`
	Credit     int `json:"credit"`
	BookCopyID int `json:"bookCopyId"`
	BookCopy   *BookCopy
}

func FindInvoiceByID(ctx context.Context, id int) (*Invoice, error) {
	var invoice Invoice

	invoice.User = &User{}

	err := db.Conn.QueryRow(ctx, `
		SELECT
			i.id,
			i.creation_time,
			i.user_id
		FROM invoice i
		WHERE i.id = $1
	`, id).Scan(
		&invoice.ID,
		&invoice.CreationTime,
		&invoice.UserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &invoice, nil
}

func FindInvoicesByUserID(ctx context.Context, userID int) ([]*Invoice, error) {
	var invoices []*Invoice

	// Find invoices and join with locations
	rows, err := db.Conn.Query(ctx, `
		SELECT
			i.id,
			i.creation_time,
			i.user_id
		FROM invoice i
		WHERE i.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invoice Invoice
		invoice.User = &User{}

		err := rows.Scan(
			&invoice.ID,
			&invoice.CreationTime,
			&invoice.UserID,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, &invoice)
	}

	return invoices, nil
}

func FindInvoiceEntriesByInvoiceID(ctx context.Context, id int) ([]*InvoiceEntry, error) {
	var invoiceEntries []*InvoiceEntry

	// Find invoice entries and join with inventory entries
	// to get the book copy and book
	rows, err := db.Conn.Query(ctx, `
		SELECT
			ie.invoice_id,
			ie.credit,
			ie.book_copy_id,
			bc.id,
			bc.book_id,
			b.id,
			b.title,
			b.subtitle,
			b.image_url,
			a.id,
			a.name
		FROM invoice_entry ie
		JOIN book_copy bc ON ie.book_copy_id = bc.id
		JOIN book b ON bc.book_id = b.id
		JOIN author a ON b.author_id = a.id
		WHERE ie.invoice_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invoiceEntry InvoiceEntry
		invoiceEntry.BookCopy = &BookCopy{}
		invoiceEntry.BookCopy.Book = &Book{}
		invoiceEntry.BookCopy.Book.Author = &Author{}

		err := rows.Scan(
			&invoiceEntry.InvoiceID,
			&invoiceEntry.Credit,
			&invoiceEntry.BookCopyID,
			&invoiceEntry.BookCopy.ID,
			&invoiceEntry.BookCopy.BookID,
			&invoiceEntry.BookCopy.Book.ID,
			&invoiceEntry.BookCopy.Book.Title,
			&invoiceEntry.BookCopy.Book.Subtitle,
			&invoiceEntry.BookCopy.Book.ImageURL,
			&invoiceEntry.BookCopy.Book.Author.ID,
			&invoiceEntry.BookCopy.Book.Author.Name,
		)
		if err != nil {
			return nil, err
		}

		invoiceEntries = append(invoiceEntries, &invoiceEntry)
	}

	return invoiceEntries, nil
}
