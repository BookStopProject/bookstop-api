package inventory

import (
	"bookstop/auth"
	"bookstop/db"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Inventory struct {
	Id         pgtype.Int4
	CreatedAt  pgtype.Timestamp
	RemovedAt  pgtype.Timestamp
	UserBookId pgtype.Int4
	LocationId pgtype.Int4
}

type InventoryClaim struct {
	Id          pgtype.Int4
	UserId      pgtype.Int4
	InventoryId pgtype.Int4
	ClaimedAt   pgtype.Timestamp
}

const allSelects = `public.inventory.id, user_book_id, location_id, created_at, removed_at`
const allSelectsClaim = `id, user_id, inventory_id, claimed_at`

func scanRows(rows *pgx.Rows) (inventories []*Inventory, errs []error) {
	for (*rows).Next() {
		inv := &Inventory{}
		err := (*rows).Scan(
			&inv.Id,
			&inv.UserBookId,
			&inv.LocationId,
			&inv.CreatedAt,
			&inv.RemovedAt,
		)
		if err != nil {
			errs = append(errs, err)
			inventories = append(inventories, nil)
		} else {
			errs = append(errs, nil)
			inventories = append(inventories, inv)
		}
	}
	return
}

const availableOnlyAnd = ` AND NOT EXISTS (SELECT FROM public.inventory_claim WHERE inventory_id = public.inventory.id) 
AND removed_at IS NULL`

func FindManyByIds(ctx context.Context, ids []int) ([]*Inventory, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.inventory WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	return scanRows(&rows)
}

func FindManyByBookId(ctx context.Context, bookId string, availableOnly bool) ([]*Inventory, error) {
	query := `SELECT ` + allSelects + ` 
FROM public.inventory
INNER JOIN public.user_book ON public.inventory.user_book_id = public.user_book.id
WHERE public.user_book.book_id = $1`

	if availableOnly {
		query += availableOnlyAnd
	}

	rows, err := db.Conn.Query(ctx, query, bookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	invs, errs := scanRows(&rows)
	var results []*Inventory
	for idx, inv := range invs {
		if errs[idx] == nil {
			results = append(results, inv)
		}
	}
	return results, nil
}

func FindManyByLocationId(ctx context.Context, locationId int, availableOnly bool) ([]*Inventory, error) {
	query := "SELECT " + allSelects + " FROM public.inventory WHERE location_id = $1"

	if availableOnly {
		query += availableOnlyAnd
	}

	rows, err := db.Conn.Query(ctx, query, locationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	invs, errs := scanRows(&rows)
	var results []*Inventory
	for idx, inv := range invs {
		if errs[idx] == nil {
			results = append(results, inv)
		}
	}
	return results, nil
}

func FindManyClaimsByUserId(ctx context.Context, userId int) ([]*InventoryClaim, error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelectsClaim+" FROM public.inventory_claim WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var invClaims []*InventoryClaim

	for (rows).Next() {
		cl := &InventoryClaim{}
		err := (rows).Scan(
			&cl.Id,
			&cl.UserId,
			&cl.InventoryId,
			&cl.ClaimedAt,
		)
		if err != nil {
			return nil, err
		} else {
			invClaims = append(invClaims, cl)
		}
	}

	return invClaims, nil
}

func DoInventoryClaim(ctx context.Context, userId int, inventoryId int) (*InventoryClaim, error) {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	balance := 0
	err = tx.QueryRow(ctx, "SELECT credit FROM public.user WHERE id=$1 FOR UPDATE", userId).Scan(&balance)
	if err != nil {
		return nil, err
	}

	if balance <= 0 {
		return nil, errors.New("not enough balance")
	}

	err = tx.QueryRow(ctx, "SELECT FROM public.inventory WHERE id = $1"+availableOnlyAnd, inventoryId).Scan()

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("inventory is not available")
		}
		return nil, err
	}

	rows, err := tx.Query(ctx, "UPDATE public.user SET credit = credit - 1 WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}

	rows.Close()

	invClaim := InventoryClaim{}

	err = tx.QueryRow(ctx, "INSERT INTO public.inventory_claim(user_id, inventory_id) VALUES ($1, $2) RETURNING "+allSelectsClaim, userId, inventoryId).Scan(
		&invClaim.Id,
		&invClaim.UserId,
		&invClaim.InventoryId,
		&invClaim.ClaimedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &invClaim, nil
}

func GenerateClaimToken(ctx context.Context, userId int, claimId int) (string, error) {
	invClaim := InventoryClaim{}

	err := db.Conn.QueryRow(ctx, "SELECT "+allSelectsClaim+" FROM public.inventory_claim WHERE id = $1 AND user_id = $2", claimId, userId).Scan(
		&invClaim.Id,
		&invClaim.UserId,
		&invClaim.InventoryId,
		&invClaim.ClaimedAt,
	)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(claimId),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	})

	tokenString, err := token.SignedString(auth.HmacSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyClaimToken(ctx context.Context, tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return auth.HmacSecret, nil
	})

	if err != nil {
		return nil
	}

	return token
}
