package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vibin18/go-shares/internal/config"
	"github.com/vibin18/go-shares/internal/driver"
	"github.com/vibin18/go-shares/internal/models"
	"github.com/vibin18/go-shares/internal/repository"
	"github.com/vibin18/go-shares/internal/repository/dbrepo"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var tmpl *template.Template
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}

	ps, err := Repo.DB.GetAllPurchaseReport()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}
	Repo.App.DashShareList = []string{}
	for _, s := range ps {
		Repo.App.DashShareList = append(Repo.App.DashShareList, s.Name)
	}

	Repo.App.DashShareCodeList, err = UpdateStrCodes(Repo.App.DashShareList)
	if err != nil {
		log.Println("Failed to update shares" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "index.html", Repo.App.ShareCache)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	buff.WriteTo(w)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	tmpl.ExecuteTemplate(w, "error_main.html", nil)
}

func AddSharesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	tmpl.ExecuteTemplate(w, "add_shares_main.html", nil)
}

func ListSharesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []models.Share{}

	shares, err = Repo.DB.GetAllShares()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	tmpl.ExecuteTemplate(w, "list_shares_main.html", shares)
}

func random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func AddSharesPostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to parse template" + err.Error())
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	name := r.Form.Get("company")

	id := random(10, 9999999)

	//log.Printf("Adding ID %v", id)

	myShare := models.Share{
		Name: name,
		Id:   id,
	}
	err = Repo.DB.InsertNewShare(myShare)
	if err != nil {
		log.Println("DB execution failed to add shares " + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	buff := bytes.Buffer{}

	log.Println(myShare.Name + " added")

	tmpl.ExecuteTemplate(&buff, "add_shares_success_main.html", myShare)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	buff.WriteTo(w)
}

func UpdateShareHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to parse template" + err.Error())
	}

	shares := []models.Share{}
	shares, err = Repo.DB.GetAllShares()

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "update_shares_main.html", shares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	buff.WriteTo(w)
}

func UpdateSharePostHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	name := r.Form.Get("name")
	cdate := r.Form.Get("date")
	datetime, err := time.Parse("2006-01-02", cdate)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	c := r.Form.Get("count")
	count, err := strconv.Atoi(c)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	log.Printf("Date entered : %v", cdate)
	shareIdName := strings.Split(name, "---")

	idn, err := strconv.Atoi(shareIdName[1])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	p := r.Form.Get("price")
	price, err := strconv.ParseFloat(p, 32)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	pf := float32(price)
	um := r.Form.Get("update_type")
	var umType string

	buySellShare := models.SellBuyShare{
		Id:        idn,
		Name:      shareIdName[0],
		Count:     count,
		Price:     pf,
		CreatedAt: datetime,
		UpdatedAt: datetime,
		Type:      umType,
	}

	if um == "Buy" {
		umType = "bought"
		buySellShare.Type = umType
		log.Printf("%v %v shares of %v", um, count, shareIdName[0])
		err = Repo.DB.BuyShare(buySellShare)
		if err != nil {
			log.Println("DB execution failed to add shares " + err.Error())
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	} else {
		umType = "sold"
		buySellShare.Type = umType
		log.Printf("%v %v shares of %v", um, count, shareIdName[0])
		err = Repo.DB.SellShare(buySellShare)
		if err != nil {
			log.Println("DB execution failed to add shares " + err.Error())
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	}

	log.Printf("%v %v shares of %v", um, count, shareIdName[0])

	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to parse template" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "update_shares_success_main.html", buySellShare)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	buff.WriteTo(w)
}

func ListTotalSharesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []models.TotalShare{}
	ushares := []models.TotalShare{}

	shares, err = Repo.DB.GetAllSharesWithData()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	for _, share := range shares {
		share.TCount = share.PCount - share.SCount
		ushares = append(ushares, share)
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_shares_main.html", ushares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func ListAllPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []models.SellBuyShare{}
	pshares := []models.SellShare{}

	shares, err = Repo.DB.GetAllPurchases()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	for _, share := range shares {
		//share.CreatedAt.Format("02-Jan-2006")

		p := models.SellShare{
			CreatedAt: share.CreatedAt.Format("2006-Jan-02"),
			Name:      share.Name,
			Count:     share.Count,
			Price:     share.Price,
		}
		pshares = append(pshares, p)
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_purchases_main.html", pshares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func ListAllSalesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []models.SellBuyShare{}
	pshares := []models.SellShare{}

	shares, err = Repo.DB.GetAllSales()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	for _, share := range shares {
		share.CreatedAt.Format("02-Jan-2006")

		p := models.SellShare{
			CreatedAt: share.CreatedAt.Format("02-Jan-2006"),
			Name:      share.Name,
			Count:     share.Count,
			Price:     share.Price,
		}
		pshares = append(pshares, p)
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_sales_main.html", pshares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func ReportAllPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []models.ShareReport{}

	shares, err = Repo.DB.GetAllPurchaseReport()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_purchase_report_main.html", shares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func ReportAllSalesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []models.ShareReport{}

	shares, err = Repo.DB.GetAllSalesReport()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_sales_report_main.html", shares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func StockHandler(w http.ResponseWriter, r *http.Request) {
	data := Repo.App.ShareCache
	json, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed marshal json" + err.Error())
	}
	io.Writer(w).Write(json)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Share added.")
}
