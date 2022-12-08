package admin

import (
	"bookstop/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Profit struct {
	BookId          int64   `json:"book_id"`
	InvoiceEntrySum float64 `json:"invoice_entry_sum"`
	TradeInSum      float64 `json:"trade_in_sum"`
	Profit          float64 `json:"profit"`
	InvoiceCount    int64   `json:"invoice_count"`
	TradeInCount    int64   `json:"trade_in_count"`
}

func apiProfit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := getDbConn(r)
	ctx := r.Context()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer conn.Close(ctx)

	profits, err := models.CalculateProfit(ctx, conn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profits)
}

func apiTradeinPreview(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := getDbConn(r)
	ctx := r.Context()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer conn.Close(ctx)
	ubID := p.ByName("userBookID")
	ubIDNum, _ := strconv.Atoi(ubID)
	preview, err := models.CalculateTradeInPreview(ctx, conn, ubIDNum, models.BookCondition(r.URL.Query().Get("condition")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(preview)
}

type TradeinInput struct {
	Condition  string `json:"condition"`
	LocationID int    `json:"location_id"`
}

func apiTradeinDo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := getDbConn(r)
	ctx := r.Context()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer conn.Close(ctx)
	ubID := p.ByName("userBookID")
	ubIDNum, _ := strconv.Atoi(ubID)

	var input TradeinInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tradein, err := models.DoTradeIn(ctx, conn, ubIDNum, models.BookCondition(input.Condition), input.LocationID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tradein)
}

func apiGetLocations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	locations, err := models.FindAllLocations(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(locations)
}

func apiCreateLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var location *models.Location
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	location, err = models.CreateLocation(ctx, location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}

func apiUpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	var location *models.Location
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	locId := p.ByName("locationID")
	locIdNum, _ := strconv.Atoi(locId)
	location.ID = locIdNum

	location, err = models.UpdateLocation(ctx, location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}

func apiGetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	conn, err := getDbConn(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	users, err := models.FindAllUsers(ctx, conn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func apiCreateBrowse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var browse *models.Browse
	err := json.NewDecoder(r.Body).Decode(&browse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	browse, err = models.CreateBrowse(ctx, browse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(browse)
}

func apiUpdateBrowse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	browseId := p.ByName("browseID")
	browseIdNum, _ := strconv.Atoi(browseId)

	var browse *models.Browse
	err := json.NewDecoder(r.Body).Decode(&browse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	browse.ID = browseIdNum

	browse, err = models.UpdateBrowse(ctx, browse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(browse)
}

func apiDeleteBrowse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	browseId := p.ByName("browseID")
	browseIdNum, _ := strconv.Atoi(browseId)

	err := models.DeleteBrowse(ctx, browseIdNum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func apiAddBookToBrowse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	browseId := p.ByName("browseID")
	browseIdNum, _ := strconv.Atoi(browseId)

	bookId := p.ByName("bookID")
	bookIdNum, _ := strconv.Atoi(bookId)

	err := models.AddBookToBrowse(ctx, browseIdNum, bookIdNum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func apiRemoveBookFromBrowse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	browseId := p.ByName("browseID")
	browseIdNum, _ := strconv.Atoi(browseId)

	bookId := p.ByName("bookID")
	bookIdNum, _ := strconv.Atoi(bookId)

	err := models.RemoveBookFromBrowse(ctx, browseIdNum, bookIdNum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func apiGetAllBrowses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	browses, err := models.FindAllBrowses(ctx)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(browses)
}

func apiGetBrowseBooks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	browseId := p.ByName("browseID")
	browseIdNum, _ := strconv.Atoi(browseId)

	books, err := models.FindBrowseBooks(ctx, browseIdNum)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func Router(router *httprouter.Router) {
	router.GET("/admin/profit", apiProfit)
	router.GET("/admin/tradein/:userBookID", apiTradeinPreview)
	router.POST("/admin/tradein/:userBookID", apiTradeinDo)
	router.GET("/admin/locations", apiGetLocations)
	router.PUT("/admin/locations/:locationID", apiUpdateLocation)
	router.POST("/admin/locations", apiCreateLocation)
	router.GET("/admin/users", apiGetAllUsers)
	router.GET("/admin/browses", apiGetAllBrowses)
	router.POST("/admin/browses", apiCreateBrowse)
	router.PUT("/admin/browses/:browseID", apiUpdateBrowse)
	router.DELETE("/admin/browses/:browseID", apiDeleteBrowse)
	router.GET("/admin/browses/:browseID/books", apiGetBrowseBooks)
	router.POST("/admin/browses/:browseID/books/:bookID", apiAddBookToBrowse)
	router.DELETE("/admin/browses/:browseID/books/:bookID", apiRemoveBookFromBrowse)
}
