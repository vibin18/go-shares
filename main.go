package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/vibin18/go-shares/internal/config"
	"github.com/vibin18/go-shares/internal/driver"
	"github.com/vibin18/go-shares/internal/handlers"
	"github.com/vibin18/go-shares/internal/models"
	opts "github.com/vibin18/go-shares/internal/ops"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	port = "0.0.0.0:8080"
)

var app config.AppConfig

var (
	argparser *flags.Parser
	arg       opts.Params
)

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {
	initArgparser()

	db, err := run()
	if err != nil {
		log.Panic(err)
	}

	err = db.SQL.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Connected to database successfully")

	defer db.SQL.Close()

	log.Println("Starting http Server..")
	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	log.Printf("Starting HTTP server on port %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	log.Println("Connecting to DB...")
	dsn := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v", arg.DbServer, arg.DbPort, arg.DbName, arg.DbUser, arg.DbPass)
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		return nil, err
		log.Fatal(err)
	}

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)
	Updater()
	go func() {
		for range time.Tick(time.Second * 10) {
			Updater()
		}
	}()
	return db, nil
}

func Updater() {

	//shareList := []string{"500469", "500209", "543441", "521070", "532155", "500183", "542724"}
	log.Printf("Updating stock for CodeList %v", app.DashShareCodeList)
	app.ShareCache = &[]models.Stock{}
	if len(app.DashShareCodeList) <= 0 {
		log.Println("Sharelist is empty, waiting to populate again")
		return
	}
	for _, share := range app.DashShareCodeList {
		s := handlers.GetStock(share)
		*app.ShareCache = append(*app.ShareCache, s)
	}
}
