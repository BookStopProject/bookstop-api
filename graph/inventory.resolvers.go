package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/app/inventory"
	"bookstop/auth"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/loader"
	"context"
	"errors"
	"strconv"
)

func (r *inventoryResolver) UserBook(ctx context.Context, obj *model.Inventory) (*model.UserBook, error) {
	intID, _ := strconv.Atoi(obj.UserBookID)
	return loader.For(ctx).UserBookByID.Load(intID)
}

func (r *inventoryResolver) Location(ctx context.Context, obj *model.Inventory) (*model.Location, error) {
	intID, _ := strconv.Atoi(obj.LocationID)
	return loader.For(ctx).LocationByID.Load(intID)
}

func (r *inventoryClaimResolver) Inventory(ctx context.Context, obj *model.InventoryClaim) (*model.Inventory, error) {
	intID, _ := strconv.Atoi(obj.InventoryID)
	return loader.For(ctx).InventoryByID.Load(intID)
}

func (r *mutationResolver) InventoryClaimDo(ctx context.Context, id string) (*model.InventoryClaim, error) {
	intID, _ := strconv.Atoi(id)
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	claim, err := inventory.DoInventoryClaim(ctx, int(usr.ID.Int), intID)
	if err != nil {
		return nil, err
	}
	return inventory.ToGraphClaim(claim), nil
}

func (r *queryResolver) Inventories(ctx context.Context, bookID *string, locationID *string) ([]*model.Inventory, error) {
	var invs []*inventory.Inventory
	var err error
	if bookID != nil {
		invs, err = inventory.FindManyByBookID(ctx, *bookID, true)
	} else if locationID != nil {
		intID, _ := strconv.Atoi(*locationID)
		invs, err = inventory.FindManyByLocationID(ctx, intID, true)
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
	invs, err := inventory.FindManyClaimsByUserID(ctx, int(usr.ID.Int))
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
	intID, _ := strconv.Atoi(id)
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return "", err
	}
	if usr == nil {
		return "", auth.ErrUnauthorized
	}
	return inventory.GenerateClaimToken(ctx, int(usr.ID.Int), intID)
}

// Inventory returns generated.InventoryResolver implementation.
func (r *Resolver) Inventory() generated.InventoryResolver { return &inventoryResolver{r} }

// InventoryClaim returns generated.InventoryClaimResolver implementation.
func (r *Resolver) InventoryClaim() generated.InventoryClaimResolver {
	return &inventoryClaimResolver{r}
}

type inventoryResolver struct{ *Resolver }
type inventoryClaimResolver struct{ *Resolver }
