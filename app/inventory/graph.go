package inventory

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func ToGraph(inventory *Inventory) *model.Inventory {
	if inventory == nil {
		return nil
	}

	val := model.Inventory{
		ID:         strconv.Itoa(int(inventory.ID.Int)),
		UserBookID: strconv.Itoa(int(inventory.UserBookID.Int)),
		LocationID: strconv.Itoa(int(inventory.LocationID.Int)),
		Removed:    inventory.RemovedAt.Status == pgtype.Present,
	}

	return &val
}

func ToGraphClaim(claim *InventoryClaim) *model.InventoryClaim {
	if claim == nil {
		return nil
	}

	val := model.InventoryClaim{
		ID:          strconv.Itoa(int(claim.ID.Int)),
		InventoryID: strconv.Itoa(int(claim.InventoryID.Int)),
		ClaimedAt:   claim.ClaimedAt.Time,
	}

	return &val
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*model.Inventory, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, `SELECT `+queryFieldsAll+`
	FROM public.inventory
	WHERE id IN (`+db.ParamRefsStr(len(ids))+`)`, args...)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	resultMap := make(map[int]*model.Inventory)

	for rows.Next() {
		ub, err := scanRow(rows.(pgx.Row))
		if err != nil {
			panic(err)
		}
		resultMap[int(ub.ID.Int)] = ToGraph(ub)
	}

	result := make([]*model.Inventory, len(ids))
	errors := make([]error, len(ids))

	for i, id := range ids {
		ub, ok := resultMap[id]
		if ok {
			result[i] = ub
		}
	}

	return result, errors
}
