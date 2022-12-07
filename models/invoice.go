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
	LocationID   int       `json:"locationId"`
	Location     *Location `json:"location"`
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
	invoice.Location = &Location{}

	err := db.Conn.QueryRow(ctx, `
		SELECT
			i.id,
			i.creation_time,
			i.user_id,
			i.location_id,
			l.id,
			l.name,
			l.address,
		FROM invoice i
		JOIN location l ON i.location_id = l.id
		WHERE i.id = $1
	`, id).Scan(
		&invoice.ID,
		&invoice.CreationTime,
		&invoice.UserID,
		&invoice.LocationID,
		&invoice.Location.ID,
		&invoice.Location.Name,
		&invoice.Location.Address,
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
			i.user_id,
			i.location_id,
			l.id,
			l.name,
			l.address,
		FROM invoice i
		JOIN location l ON i.location_id = l.id
		WHERE i.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invoice Invoice
		invoice.User = &User{}
		invoice.Location = &Location{}

		err := rows.Scan(
			&invoice.ID,
			&invoice.CreationTime,
			&invoice.UserID,
			&invoice.LocationID,
			&invoice.Location.ID,
			&invoice.Location.Name,
			&invoice.Location.Address,
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
			ie.book_copy_id
			bc.id,
			bc.book_id,
			b.id,
			b.title,
			b.subtitle,
			b.image_url,
		FROM invoice_entry ie
		JOIN book_copy bc ON ie.book_copy_id = bc.id
		JOIN book b ON bc.book_id = b.id
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
		)
		if err != nil {
			return nil, err
		}

		invoiceEntries = append(invoiceEntries, &invoiceEntry)
	}

	return invoiceEntries, nil
}
