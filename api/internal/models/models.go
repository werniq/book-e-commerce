package models

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type DatabaseModel struct {
	DBModel *gorm.DB
}

type Book struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      Author    `json:"author'"`
	Amount      int       `json:"quantity"`
	Price       int       `json:"price"`
	Image       string    `json:"image"`
	IsAvaliable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Author struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Surname       string      `json:"surname"`
	DateOfBirth   time.Time   `json:"date_of_birth"`
	Books         []Book      `json:"books"`
	Description   string      `json:"description"`
	ActivityYears []time.Time `json:"activity_years"`
	//For instance [1921, 1948]
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	gorm.Model
	Username string `json:"varchar"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
	IsAdmin  bool   `json:"admin"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string    `json:"role"`
	Email       string    `json:"email"`
	TokenString string    `json:"token"`
	Expiry      time.Time `json:"expiry"`
}

func (m *DatabaseModel) ReturnBooks() ([]Book, error) {
	var books []Book
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// stmt := `select * from books`

	res := m.DBModel.Table("books").Order("id").Find(&books)

	if res.Error != nil {
		fmt.Println(res.Error)
		return nil, res.Error
	}

	return books, nil
}

func (m *DatabaseModel) DeleteBook(book Book) error {
	row := m.DBModel.Model(&Book{}).Where("id = $1", book.Id).Delete(book)
	if row.Error != nil {
		log.Fatalf("%v", row.Error)
		return row.Error
	}
	return nil
}

func (m *DatabaseModel) UpdateBook(book Book) error {
	row := m.DBModel.Model(&Book{}).Where("id = $1", book.Id).Updates(&book)
	if row.Error != nil {
		log.Fatalf("%v", row.Error)
		return row.Error
	}

	if row.RowsAffected == 0 {
		m.DBModel.Create(&book)
	}

	return nil
}

func (m *DatabaseModel) InsertBook(book Book) error {
	row := m.DBModel.Model(&Book{}).Select("title", "description", "price", "is_avaliable", "created_at", "updated_at", "amount", "image").Create(&book)
	if row.Error != nil {
		log.Fatalf("%v", row.Error)
		return row.Error
	}

	return nil
}

func (m *DatabaseModel) InsertUser(user User) error {
	row := m.DBModel.Model(&Book{}).Select("Username", "Email", "Password", "IsAdmin").Create(&user)
	if row.Error != nil {
		log.Fatalf("%v", row.Error)
		return row.Error
	}

	return nil
}

func (m *DatabaseModel) GetUserForToken(token string) (User, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	var user User

	/*
		Username string `json:"varchar"`
		Email    string `json:"email"`
		Password []byte `json:"password"`
		IsAdmin  bool   `json:"admin"`
	*/

	query := `
		select
			user.Id, user.Username, user.Password, user.IsAdmin
		from 
		    users u
			inner join tokens t on (user.Id = t.userID)
		where 
		    t.hash = $1
			and t.expiry > $2
	`

	m.DBModel.Raw(query, tokenHash, time.Now()).Scan(&user)

	return user, nil
}

func (m *DatabaseModel) GetUserByToken(id int) (User, error) {
	var u User
	m.DBModel.Where("id = $1", id).Scan(&u)
	if u.ID != 0 {
		return u, nil
	}
	return u, errors.New("user not found")
}

func (m *DatabaseModel) GetUserByEmail(email string) (User, error) {
	var u User
	m.DBModel.Where("email = $1", email).Scan(&u)
	if u.ID != 0 {
		return u, nil
	}
	return u, errors.New("user not found")
}

func (m *DatabaseModel) GetAllBookRecords() ([]Book, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	var books []Book

	defer cancel()

	res := m.DBModel.Find(&books)
	if res.Error != nil {
		log.Fatalf("%v", res.Error)
		return books, res.Error
	}

	fmt.Println("Rows affected: ", res.RowsAffected)

	return books, nil
}

func (m *DatabaseModel) GetAllUsers() ([]User, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	var users []User

	defer cancel()

	res := m.DBModel.Find(&users)
	if res.Error != nil {
		log.Fatalf("%v", res.Error)
		return users, res.Error
	}

	fmt.Println("Rows affected: ", res.RowsAffected)

	return users, nil
}

func (m *DatabaseModel) GetBookById(id int) (Book, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var book Book

	row := m.DBModel.First(&book, id)

	if row.Error != nil {
		log.Fatalln(row.Error)
		return book, row.Error
	}

	return book, nil
}
