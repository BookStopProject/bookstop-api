package inventory

import (
	"bookstop/graph/model"
	"strconv"

	"github.com/jackc/pgtype"
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
