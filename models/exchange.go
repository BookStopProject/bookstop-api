package models

import "context"

func DoExchange(ctx context.Context, userID int, bookCopyIDs []int) (*Invoice, error) {
	// TODO: implement procedure
	// This procedure should:
	// 1) Verify that book copies are available at locations (have location_id)
	// 2) Create an invoice
	// 3) For each book copy, create an invoice entry. The invoice entry credit
	// 		will be equal to the book exchange price * book copy condition multiplier (see book_copy.go)).
	// 4) For each book copy, update the book copy location_id to nil.
	// 5) For each book copy, create a user book with the book copy id and user id.
	// 6) Calculate the total credit of the invoice. Deduct the total credit from the user's balance
	//		(must verify that user has enough balance)
	// 7) Return the invoice.
	return nil, nil
}
