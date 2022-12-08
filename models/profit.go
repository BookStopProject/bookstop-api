package models

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Profit struct {
	BookId          int  `json:"book_id"`
	InvoiceEntrySum *int `json:"invoice_entry_sum"`
	TradeInSum      *int `json:"trade_in_sum"`
	Profit          *int `json:"profit"`
	InvoiceCount    *int `json:"invoice_count"`
	TradeInCount    *int `json:"trade_in_count"`
}

func CalculateProfit(ctx context.Context, conn *pgx.Conn) ([]Profit, error) {
	rows, err := conn.Query(ctx, `SELECT
	bc.book_id,
	SUM(ie.credit),
	SUM(ti.credit),
	COALESCE(SUM(ie.credit),0) - COALESCE(SUM(ti.credit),0) AS profit,
	COUNT(DISTINCT ie.invoice_id) AS invoice_count,
	COUNT(DISTINCT ti.id) AS trade_in_count
FROM
	public."book_copy" bc
	LEFT JOIN public."invoice_entry" ie ON ie.book_copy_id = bc.id
	LEFT JOIN public."trade_in" ti ON ti.book_copy_id = bc.id
GROUP BY
	bc.book_id
ORDER BY
	profit DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profits := make([]Profit, 0)
	for rows.Next() {
		var profit Profit
		err := rows.Scan(&profit.BookId, &profit.InvoiceEntrySum, &profit.TradeInSum, &profit.Profit, &profit.InvoiceCount, &profit.TradeInCount)
		if err != nil {
			return nil, err
		}
		profits = append(profits, profit)
	}

	return profits, nil
}
