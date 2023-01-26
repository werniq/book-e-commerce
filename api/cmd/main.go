package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"react-go/internal/driver"
	models "react-go/internal/models"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"time"
)

type config struct {
	port    int
	env     string
	version string
	stripe  struct {
		key    string
		secret string
	}
	db struct {
		dsn string
	}
}

type application struct {
	cfg           config
	infoLog       *log.Logger
	errorLog      *log.Logger
	database      models.DatabaseModel
	templateCache map[string]*template.Template
}

func (app *application) serve() error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":4001"),
		Handler:           app.routes(),
		IdleTimeout:       15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println("Running development server on :4001")

	return server.ListenAndServe()
}

func main() {

	var cfg config
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	cfg.port, _ = strconv.Atoi(os.Getenv("API_PORT"))
	cfg.env = "development"
	cfg.version = "1.0.0"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := driver.Connect()
	db.AutoMigrate(&models.Book{}, &models.User{}, &models.Token{})

	if err != nil {
		errorLog.Println(err)
	}
	defer db.Close()

	// page name -> page
	tc := make(map[string]*template.Template)

	app := &application{
		cfg:           cfg,
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: tc,
		database: models.DatabaseModel{
			DBModel: db,
		},
	}

	if err = app.serve(); err != nil {
		app.errorLog.Println(err)
		return
	}

	db.Close()
}
