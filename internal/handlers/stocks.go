package handlers

import (
	"bufio"
	"fmt"
	goquery2 "github.com/PuerkitoBio/goquery"
	"github.com/vibin18/go-shares/internal/config"
	"github.com/vibin18/go-shares/internal/models"
	"log"
	"net/http"
	"strings"
)

type repo config.AppConfig

func UpdateStrCodes(shareNames []string) ([]string, error) {
	var codelist []string

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://static.quandl.com/BSE+Descriptions/stocks.txt", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var lines []string
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	log.Println("Updating code lists using UpdateStrCodes")
	for _, n := range shareNames { // my list range needed
		log.Println("Checking share :" + n)
		for _, share := range lines { // every line ranged each line is called share
			shares := strings.Split(share, "EOD Prices|BOM") // share and code splited

			if strings.Contains(share, n) {
				log.Printf("Adding share code for %v with %v", shares[0], shares[1])

				codelist = append(codelist, shares[1])
			}
		}
	}
	log.Printf("Code list generated inside functions UpdateStrCodes", codelist)

	return codelist, nil
}

func getStockQuote(URL string) (*goquery2.Document, error) {

	client := http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("Failed to create http client" + err.Error())
	}

	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36 Edg/83.0.478.45"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to create request" + err.Error())
	}

	defer res.Body.Close()

	doc, err := goquery2.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GetStock(scripcd string) models.Stock {
	url := fmt.Sprintf("https://m.bseindia.com/StockReach.aspx?scripcd=%v", scripcd)
	doc, err := getStockQuote(url)
	if err != nil {
		log.Println(err)
	}
	Share := models.Stock{}

	price := doc.Find(
		".srcovalue strong").Text()
	companyName := doc.Find(

		".companydetail .companyname ").Text()
	GetTime := doc.Find(
		".companydetail span#strongDate").Text()
	Stime := strings.Split(GetTime, "-")

	UpdateTime := strings.TrimSpace(Stime[1])

	PreCloOpen := doc.Find(
		".menuarea td#tdpcloseopen").Text()
	PreCloOpenVars := strings.Split(PreCloOpen, "/")
	DayHL := doc.Find(
		".menuarea td#tdDHL").Text()
	DayHLVars := strings.Split(
		DayHL, "/")
	WeekAvg := doc.Find(
		".menuarea td#tdWAp").Text()

	Share.CompanyName = companyName
	Share.CurrentValue = price
	Share.PreviousOpen = PreCloOpenVars[1]
	Share.PreviousClose = PreCloOpenVars[0]
	Share.DayHigh = DayHLVars[0]
	Share.DayLow = DayHLVars[1]
	Share.WeekAverage = WeekAvg
	Share.UpdateTime = UpdateTime

	return Share
}
