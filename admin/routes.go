package admin

import (
	"bookstop/book"
	"bookstop/browse"
	"bookstop/db"
	"bookstop/event"
	"bookstop/graph/model"
	"bookstop/inventory"
	"bookstop/location"
	"bookstop/user"
	"bookstop/userbook"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var errNotFound = errors.New("404 page not found")

func writeErr(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

// /admin/home
var tmplHome = template.Must(template.ParseFiles("admin/index.html", "admin/base.html"))

func apiHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := tmplHome.Execute(w, nil); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
	}
}

// /admin/browse
var tmplBrowse = template.Must(template.ParseFiles("admin/browse.html", "admin/base.html"))

type DataBrowse struct {
	Browses []*browse.Browse
}

func apiBrowse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	results, err := browse.FindAll(r.Context(), nil)
	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	if err := tmplBrowse.Execute(w, DataBrowse{
		Browses: results,
	}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func apiBrowseCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	_, err := browse.Create(r.Context(),
		r.PostForm.Get("name"),
		r.PostForm.Get("description"),
		r.PostForm.Get("started_at"),
		r.PostForm.Get("ended_at"))

	if err != nil {
		http.Redirect(w, r, "/admin/browse?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/browse", http.StatusSeeOther)
}

// /admin/browse/:id
type DataBrowseEach struct {
	Browse *browse.Browse
	Books  []*book.Book
}

var tmplBrowseEach = template.Must(template.ParseFiles("admin/browse_each.html", "admin/base.html"))

func apiBrowseEach(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	intID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	result, err := browse.FindByID(ctx, intID)
	resultBooks, errs := browse.FindBooksByBrowseID(ctx, intID)

	var books []*book.Book
	for idx, b := range resultBooks {
		if errs[idx] == nil {
			books = append(books, b)
		}
	}

	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	if result == nil {
		writeErr(w, errNotFound, 404)
		return
	}
	if err := tmplBrowseEach.Execute(w, DataBrowseEach{
		Browse: result,
		Books:  books,
	}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func apiBrowseEachEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	strID := ps.ByName("id")
	intID, err := strconv.Atoi(strID)
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	_, err = browse.UpdateByID(r.Context(), intID,
		r.PostForm.Get("name"),
		r.PostForm.Get("description"),
		r.PostForm.Get("started_at"),
		r.PostForm.Get("ended_at"))

	redirectUrl := "/admin/browse/" + strID

	if err != nil {
		redirectUrl += "?error=" + err.Error()
		return
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func apiBrowseEachDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	intID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	_, err = browse.DeleteByID(r.Context(), intID)
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func apiBrowseEachAddBooks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	strID := ps.ByName("id")
	redirectUrl := "/admin/browse/" + strID

	intID, err := strconv.Atoi(strID)
	if err != nil {
		http.Redirect(w, r, redirectUrl+err.Error(), http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, redirectUrl+err.Error(), http.StatusFound)
		return
	}

	_, err = browse.AddBooksByIDs(r.Context(), intID, strings.Split(r.PostForm.Get("book_ids"), ","))
	if err != nil {
		http.Redirect(w, r, redirectUrl+"?error="+err.Error(), http.StatusFound)
		return
	}

	http.Redirect(w, r, redirectUrl+"?book_added=1", http.StatusSeeOther)
}

func apiBrowseEachDeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	strID := ps.ByName("id")
	intID, err := strconv.Atoi(strID)
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	_, err = browse.DeleteBooksByIDs(r.Context(), intID, strings.Split(r.URL.Query().Get("book_ids"), ","))
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin/browse/"+strID+"?deleted=1", http.StatusSeeOther)
}

// /admin/inventory
var tmplInventory = template.Must(template.ParseFiles("admin/inventory.html", "admin/base.html"))

type DataInventory struct {
	Inventories []*DataInventoryItem
}

type DataInventoryItem struct {
	ID         int
	BookID     string
	CreatedAt  *time.Time
	RemovedAt  *time.Time
	UserBookID int
	LocationID int
}

func apiInventory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := db.Conn.Query(r.Context(), `SELECT inventory.id, user_book.book_id, created_at, removed_at, user_book_id, location_id
	FROM public.inventory
	INNER JOIN public.user_book ON public.inventory.user_book_id = public.user_book.id`)
	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	data := DataInventory{}

	for rows.Next() {
		iv := DataInventoryItem{}
		err := rows.Scan(&iv.ID, &iv.BookID, &iv.CreatedAt, &iv.RemovedAt, &iv.UserBookID, &iv.LocationID)
		if err != nil {
			writeErr(w, err, http.StatusInternalServerError)
			return
		}
		data.Inventories = append(data.Inventories, &iv)
	}

	if err := tmplInventory.Execute(w, data); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

// admin/check-out
var tmplCheckOut = template.Must(template.ParseFiles("admin/checkout.html", "admin/base.html"))

func apiCheckOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := tmplCheckOut.Execute(w, nil); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func apiCheckOutAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		writeErr(w, errors.New("empty token"), http.StatusBadRequest)
		return
	}
	token, err := inventory.VerifyClaimToken(r.Context(), tokenString)
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func apiCheckOutActionCommit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		writeErr(w, errors.New("empty token"), http.StatusBadRequest)
		return
	}
	ok, err := inventory.DoInventoryCheckoutWithToken(r.Context(), tokenString)
	if err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}
	if !ok {
		writeErr(w, errors.New("unsuccessful"), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// admin/check-in
var tmplCheckIn = template.Must(template.ParseFiles("admin/checkin.html", "admin/base.html"))

type DataCheckIn struct {
	Locations       []*model.Location
	CheckedUserBook *userbook.UserBook
}

func apiCheckIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	locations, err := location.FindAllLocations(r.Context())

	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
	}

	if err := tmplCheckIn.Execute(w, DataCheckIn{Locations: locations}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func apiCheckInCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	ubID, err := strconv.Atoi(r.PostForm.Get("user_book_id"))
	if err != nil {
		http.Redirect(w, r, "/admin/check-in?error="+err.Error(), http.StatusSeeOther)
		return
	}
	locID, err := strconv.Atoi(r.PostForm.Get("location_id"))
	if err != nil {
		http.Redirect(w, r, "/admin/check-in?error="+err.Error(), http.StatusSeeOther)
		return
	}

	ctx := r.Context()

	inv, err := inventory.InsertInventoryAndReward(ctx, ubID, locID)

	if err != nil {
		http.Redirect(w, r, "/admin/check-in?error="+err.Error(), http.StatusSeeOther)
		return
	}

	ub, _ := userbook.FindByID(ctx, int(inv.UserBookID.Int))
	if ub == nil {
		writeErr(w, errors.New("cannot load user book"), http.StatusInternalServerError)
		return
	}

	locations, err := location.FindAllLocations(r.Context())

	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
	}

	if err := tmplCheckIn.Execute(w, DataCheckIn{Locations: locations, CheckedUserBook: ub}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

// /admin/location
var tmplLocation = template.Must(template.ParseFiles("admin/location.html", "admin/base.html"))

type DataLocation struct {
	Locations []*model.Location
}

func apiLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	locations, err := location.FindAllLocations(r.Context())

	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
	}

	if err := tmplLocation.Execute(w, DataLocation{Locations: locations}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func apiLocationCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	_, err := location.Create(r.Context(),
		r.PostForm.Get("name"),
		r.PostForm.Get("parent_name"),
		r.PostForm.Get("address_line"))

	if err != nil {
		http.Redirect(w, r, "/admin/location?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/location", http.StatusSeeOther)
}

// admin/user
type DataUser struct {
	Users []*user.User
}

var tmplUser = template.Must(template.ParseFiles("admin/user.html", "admin/base.html"))

func apiUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, _ := user.FindAll(r.Context())
	if err := tmplUser.Execute(w, DataUser{Users: users}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

// admin/event
type DataEvent struct {
	Events []*event.Event
}

var tmplEvent = template.Must(template.ParseFiles("admin/event.html", "admin/base.html"))

func apiEvent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events, err := event.FindAll(r.Context())
	if err != nil {
		writeErr(w, err, http.StatusInternalServerError)
	}
	if err := tmplEvent.Execute(w, DataEvent{Events: events}); err != nil {
		writeErr(w, err, http.StatusInternalServerError)
		return
	}
}

func apiEventCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		writeErr(w, err, http.StatusBadRequest)
		return
	}

	_, err := event.Create(r.Context(),
		r.PostForm.Get("title"),
		r.PostForm.Get("description"),
		r.PostForm.Get("href"),
		r.PostForm.Get("user_id"),
		r.PostForm.Get("started_at"),
		r.PostForm.Get("ended_at"),
	)

	if err != nil {
		http.Redirect(w, r, "/admin/event?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/event", http.StatusSeeOther)
}

// router

func basicAuth(next httprouter.Handle) httprouter.Handle {
	adminAuth := strings.Split(os.Getenv("ADMIN_AUTH"), ":")
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		username, password, ok := r.BasicAuth()
		if ok {
			if username == adminAuth[0] && password == adminAuth[1] {
				next(w, r, ps)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func Router(router *httprouter.Router) {
	(*router).GET("/admin", basicAuth(apiHome))
	(*router).GET("/admin/browse", basicAuth(apiBrowse))
	(*router).POST("/admin/browse", basicAuth(apiBrowseCreate))
	(*router).GET("/admin/browse/:id", basicAuth(apiBrowseEach))
	(*router).POST("/admin/browse/:id", basicAuth(apiBrowseEachEdit))
	(*router).DELETE("/admin/browse/:id", basicAuth(apiBrowseEachDelete))
	(*router).POST("/admin/browse/:id/books", basicAuth(apiBrowseEachAddBooks))
	(*router).DELETE("/admin/browse/:id/books", basicAuth(apiBrowseEachDeleteBook))
	(*router).GET("/admin/inventory", basicAuth(apiInventory))
	(*router).GET("/admin/check-out", basicAuth(apiCheckOut))
	(*router).GET("/admin/check-out/action", basicAuth(apiCheckOutAction))
	(*router).POST("/admin/check-out/action", basicAuth(apiCheckOutActionCommit))
	(*router).GET("/admin/check-in", basicAuth(apiCheckIn))
	(*router).POST("/admin/check-in", basicAuth(apiCheckInCreate))
	(*router).GET("/admin/location", basicAuth(apiLocation))
	(*router).POST("/admin/location", basicAuth(apiLocationCreate))
	(*router).GET("/admin/user", basicAuth(apiUser))
	(*router).GET("/admin/event", basicAuth(apiEvent))
	(*router).POST("/admin/event", basicAuth(apiEventCreate))
}
