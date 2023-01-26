package driver

import (
	"fmt"
	"log"
	"os"
	"react-go/internal/models"

	"github.com/jinzhu/gorm"
)

type Config struct {
	host     string
	port     string
	password string
	user     string
	dbname   string
	sslomde  string
}

func Connect() (*gorm.DB, error) {
	var cfg = &Config{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		password: os.Getenv("DB_PASSWORD"),
		user:     os.Getenv("DB_USER"),
		dbname:   os.Getenv("DB_NAME"),
		sslomde:  os.Getenv("SSLMODE"),
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.host, cfg.user, cfg.password, cfg.dbname, cfg.port, cfg.sslomde)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return db, nil
}


func CloseDatabase(connection *gorm.DB) {
	sql := connection.DB()
	sql.Close()
}

func InitialMigration() {
	connection, err := Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer CloseDatabase(connection)
	connection.AutoMigrate(&models.Book{})
}
