package inventory

import (
	"bookstop/auth"
	"bookstop/book"
	"bookstop/db"
	"bookstop/user"
	"bookstop/userbook"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Inventory struct {
	ID         pgtype.Int4
	CreatedAt  pgtype.Timestamp
	RemovedAt  pgtype.Timestamp
	UserBookID pgtype.Int4
	LocationID pgtype.Int4
}

type InventoryClaim struct {
	ID          pgtype.Int4
	UserID      pgtype.Int4
	InventoryID pgtype.Int4
	ClaimedAt   pgtype.Timestamp
}

const allSelects = `inventory.id, user_book_id, location_id, created_at, removed_at`
const allSelectsClaim = `id, user_id, inventory_id, claimed_at`

func scanRows(rows *pgx.Rows) (inventories []*Inventory, errs []error) {
	for (*rows).Next() {
		inv := &Inventory{}
		err := (*rows).Scan(
			&inv.ID,
			&inv.UserBookID,
			&inv.LocationID,
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

func scanRow(row *pgx.Row) (*Inventory, error) {
	inv := Inventory{}
	err := (*row).Scan(&inv.ID, &inv.UserBookID, &inv.LocationID, &inv.CreatedAt, &inv.RemovedAt)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

const availableOnlyAnd = ` AND NOT EXISTS (SELECT FROM public.inventory_claim WHERE inventory_id = public.inventory.id) 
AND removed_at IS NULL`

func LoadManyByIDs(ctx context.Context, ids []int) ([]*Inventory, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.inventory WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result, errs := scanRows(&rows)

	idToPos := make(map[int]int)

	for i, ub := range result {
		idToPos[int(ub.ID.Int)] = i
	}

	sortedResult := make([]*Inventory, len(ids))
	sortedErrs := make([]error, len(ids))

	for i, id := range ids {
		pos := idToPos[id]
		sortedResult[i] = result[pos]
		sortedErrs[i] = errs[pos]
	}

	return sortedResult, sortedErrs
}

func FindManyByBookID(ctx context.Context, bookID string, availableOnly bool) ([]*Inventory, error) {
	query := `SELECT ` + allSelects + ` 
FROM public.inventory
INNER JOIN public.user_book ON public.inventory.user_book_id = public.user_book.id
WHERE public.user_book.book_id = $1`

	if availableOnly {
		query += availableOnlyAnd
	}

	rows, err := db.Conn.Query(ctx, query, bookID)
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

func FindManyByLocationID(ctx context.Context, locationID int, availableOnly bool) ([]*Inventory, error) {
	query := "SELECT " + allSelects + " FROM public.inventory WHERE location_id = $1 ORDER BY id DESC"

	if availableOnly {
		query += availableOnlyAnd
	}

	rows, err := db.Conn.Query(ctx, query, locationID)
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

func FindManyClaimsByUserID(ctx context.Context, userID int) ([]*InventoryClaim, error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelectsClaim+" FROM public.inventory_claim WHERE user_id = $1 ORDER BY id DESC", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var invClaims []*InventoryClaim

	for (rows).Next() {
		cl := &InventoryClaim{}
		err := (rows).Scan(
			&cl.ID,
			&cl.UserID,
			&cl.InventoryID,
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

func InsertInventoryAndReward(ctx context.Context, ubID int, locID int) (*Inventory, error) {
	ub, err := userbook.FindByID(ctx, ubID)
	if err != nil {
		return nil, err
	}

	if ub == nil {
		return nil, errors.New("cannot find user book")
	}

	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, "INSERT INTO public.inventory(user_book_id, location_id) VALUES ($1, $2) RETURNING "+allSelects, ubID, locID)
	inv, err := scanRow(&row)
	if err != nil {
		return nil, err
	}

	var resultedCredit *int
	err = tx.QueryRow(ctx, "UPDATE public.user SET credit = credit + 1 WHERE id = $1 RETURNING credit", int(ub.UserID.Int)).Scan(&resultedCredit)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return inv, nil
}

func DoInventoryClaim(ctx context.Context, userID int, inventoryID int) (*InventoryClaim, error) {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	balance := 0
	err = tx.QueryRow(ctx, "SELECT credit FROM public.user WHERE id=$1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		return nil, err
	}

	if balance <= 0 {
		return nil, errors.New("not enough balance")
	}

	var otherUserID *int
	err = tx.QueryRow(ctx, "SELECT user_book.user_id FROM public.inventory INNER JOIN public.user_book ON public.inventory.user_book_id = public.user_book.id WHERE public.inventory.id = $1"+availableOnlyAnd, inventoryID).Scan(&otherUserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("inventory is not available")
		}
		return nil, err
	}

	if *otherUserID == userID {
		return nil, errors.New("cannot exchange your own book")
	}

	var resultedCredit *int
	err = tx.QueryRow(ctx, "UPDATE public.user SET credit = credit - 1 WHERE id = $1 RETURNING credit", userID).Scan(&resultedCredit)

	if err != nil {
		return nil, err
	}

	if *resultedCredit < 0 {
		return nil, errors.New("invalid credit")
	}

	invClaim := InventoryClaim{}

	err = tx.QueryRow(ctx, "INSERT INTO public.inventory_claim(user_id, inventory_id) VALUES ($1, $2) RETURNING "+allSelectsClaim, userID, inventoryID).Scan(
		&invClaim.ID,
		&invClaim.UserID,
		&invClaim.InventoryID,
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

func GenerateClaimToken(ctx context.Context, userID int, claimID int) (string, error) {
	invClaim := InventoryClaim{}

	err := db.Conn.QueryRow(ctx, "SELECT "+allSelectsClaim+" FROM public.inventory_claim WHERE id = $1 AND user_id = $2", claimID, userID).Scan(
		&invClaim.ID,
		&invClaim.UserID,
		&invClaim.InventoryID,
		&invClaim.ClaimedAt,
	)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(claimID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	})

	tokenString, err := token.SignedString(auth.HmacSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type InventoryClaimDetailed struct {
	InventoryClaim
	BookID             string
	LocationID         int
	Book               *book.Book
	User               *user.User
	UserBookIDOriginal *int
}

func VerifyClaimToken(ctx context.Context, tokenString string) (*InventoryClaimDetailed, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return auth.HmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claimID := claims["sub"].(string)

	detailedInvClaim := InventoryClaimDetailed{}
	err = db.Conn.QueryRow(ctx, `SELECT inventory_claim.id, inventory_claim.user_id, inventory_id, claimed_at, book_id, location_id, user_book.id_original AS user_book_id_original
	FROM public.inventory_claim
	INNER JOIN public.inventory ON public.inventory_claim.inventory_id = public.inventory.id
	INNER JOIN public.user_book ON public.inventory.user_book_id = public.user_book.id
	WHERE inventory_claim.id = $1 AND inventory.removed_at IS NULL
`, claimID).Scan(
		&detailedInvClaim.ID,
		&detailedInvClaim.UserID,
		&detailedInvClaim.InventoryID,
		&detailedInvClaim.ClaimedAt,
		&detailedInvClaim.BookID,
		&detailedInvClaim.LocationID,
		&detailedInvClaim.UserBookIDOriginal,
	)

	if err != nil {
		return nil, err
	}

	u, err := user.FindByID(ctx, int(detailedInvClaim.UserID.Int))
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("cannot find user")
	}
	detailedInvClaim.User = u

	b, err := book.FindByID(ctx, detailedInvClaim.BookID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("cannot find book")
	}
	detailedInvClaim.Book = b

	return &detailedInvClaim, nil
}

func DoInventoryCheckoutWithToken(ctx context.Context, tokenString string) (bool, error) {
	cl, err := VerifyClaimToken(ctx, tokenString)
	if err != nil {
		return false, err
	}
	if cl == nil {
		return false, errors.New("cannot find claim")
	}
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return false, err
	}

	defer tx.Rollback(ctx)

	tNow := time.Now()

	var oldUserBookID int

	err = tx.QueryRow(ctx, "UPDATE public.inventory SET removed_at = $2 WHERE id = $1 RETURNING user_book_id", cl.InventoryID, tNow).Scan(&oldUserBookID)
	if err != nil {
		return false, err
	}

	if cl.UserBookIDOriginal == nil {
		cl.UserBookIDOriginal = &oldUserBookID
	}

	var newUserBookID int
	tx.QueryRow(ctx, "INSERT INTO public.user_book(user_id, book_id, id_original) VALUES ($1, $2, $3) RETURNING id", cl.UserID, cl.BookID, cl.UserBookIDOriginal).Scan(&newUserBookID)

	rows, err := tx.Query(ctx, "INSERT INTO public.exchange(user_book_id_old, user_book_id_new, user_book_id_original, exchanged_at) VALUES ($1, $2, $3, $4)", oldUserBookID, newUserBookID, cl.UserBookIDOriginal, tNow)
	if err != nil {
		return false, err
	}

	rows.Close()

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}
