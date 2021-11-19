package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/auth"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/inventory"
	"bookstop/loader"
	"bookstop/location"
	"bookstop/userbook"
	"context"
	"errors"
	"strconv"
)

func (r *inventoryResolver) UserBook(ctx context.Context, obj *model.Inventory) (*model.UserBook, error) {
	intId, _ := strconv.Atoi(obj.UserBookID)
	ub, err := loader.For(ctx).UserBookById.Load(intId)
	if err != nil {
		return nil, err
	}
	return userbook.ToGraph(ub), nil
}

func (r *inventoryResolver) Location(ctx context.Context, obj *model.Inventory) (*model.Location, error) {
	intId, _ := strconv.Atoi(obj.LocationID)
	loc, err := loader.For(ctx).LocationById.Load(intId)
	if err != nil {
		return nil, err
	}
	return location.ToGraph(loc), nil
}

func (r *inventoryClaimResolver) Inventory(ctx context.Context, obj *model.InventoryClaim) (*model.Inventory, error) {
	intId, _ := strconv.Atoi(obj.InventoryID)
	inv, err := loader.For(ctx).InventoryById.Load(intId)
	if err != nil {
		return nil, err
	}
	return inventory.ToGraph(inv), nil
}

func (r *mutationResolver) InventoryClaimDo(ctx context.Context, id string) (*model.InventoryClaim, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	claim, err := inventory.DoInventoryClaim(ctx, int(usr.ID.Int), intId)
	if err != nil {
		return nil, err
	}
	return inventory.ToGraphClaim(claim), nil
}

func (r *queryResolver) Inventories(ctx context.Context, bookID *string, locationID *string) ([]*model.Inventory, error) {
	var invs []*inventory.Inventory
	var err error
	if bookID != nil {
		invs, err = inventory.FindManyByBookId(ctx, *bookID, true)
	} else if locationID != nil {
		intId, intErr := strconv.Atoi(*locationID)
		if intErr != nil {
			return nil, err
		}
		invs, err = inventory.FindManyByLocationId(ctx, intId, true)
	} else {
		return nil, errors.New("must provide either bookID or locationID")
	}
	if err != nil {
		return nil, err
	}
	results := make([]*model.Inventory, len(invs))
	for i, inv := range invs {
		results[i] = inventory.ToGraph(inv)
	}
	return results, nil
}

func (r *queryResolver) InventoryClaimsMine(ctx context.Context) ([]*model.InventoryClaim, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	invs, err := inventory.FindManyClaimsByUserId(ctx, int(usr.ID.Int))
	if err != nil {
		return nil, err
	}
	results := make([]*model.InventoryClaim, len(invs))
	for i, inv := range invs {
		results[i] = inventory.ToGraphClaim(inv)
	}
	return results, nil
}

func (r *queryResolver) InventoryClaimToken(ctx context.Context, id string) (string, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return "", err
	}
	if usr == nil {
		return "", auth.ErrUnauthorized
	}
	return inventory.GenerateClaimToken(ctx, int(usr.ID.Int), intId)
}

// Inventory returns generated.InventoryResolver implementation.
func (r *Resolver) Inventory() generated.InventoryResolver { return &inventoryResolver{r} }

// InventoryClaim returns generated.InventoryClaimResolver implementation.
func (r *Resolver) InventoryClaim() generated.InventoryClaimResolver {
	return &inventoryClaimResolver{r}
}

type inventoryResolver struct{ *Resolver }
type inventoryClaimResolver struct{ *Resolver }
