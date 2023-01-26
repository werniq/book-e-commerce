package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"react-go/internal/driver"
	"react-go/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UsedId   int       `json:"userID"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
	Token    string    `json:"token"`
	Expiry   time.Time `json:"expriy"`
}

type CartProduct struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId"`
	ProductId int       `json:"bookID"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Authentication struct {
	Email    string `json:"email"'`
	Password string `json:"password"`
}

type Categories struct {
	gorm.Model
	Name  string        `json:"name"`
	Books []models.Book `json:"books"`
}

type User struct {
	gorm.Model
	Id       int    `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
	IsAdmin  bool   `json:"admin"`
}

const (
	maxBytes = 1048576
)

func (app *application) ListProducts(w http.ResponseWriter, r *http.Request) {
	var books []models.Book
	var err error
	db, err := driver.Connect()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer driver.CloseDatabase(db)
	//dbmodel := models.DatabaseModel{DBModel: db}
	books, err = app.database.ReturnBooks()

	if err != nil {
		log.Fatal(err)
		return
	}
	res, err := json.Marshal(books)
	if err != nil {
		fmt.Println("Error parsing body. Controllers line 36.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (app *application) GetCryptoInfo(w http.ResponseWriter, r *http.Request) {

}

func (app *application) AddToCart(w http.ResponseWriter, r *http.Request) {

}

// GenerateToken func generates token for authentication
func (app *application) GenerateToken(email, role string) (string, error) {
	var superSecretKey = []byte(os.Getenv("SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 30).Unix()

	tokenString, err := token.SignedString(superSecretKey)
	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

// invalidCredentials is helper func which returns error if not valid email/password printed
func (app *application) invalidCredentials(w http.ResponseWriter) error {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	w.Header().Set("Content-Type", "application/json")
	res, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		log.Fatal(err)
		return err
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(res)

	return nil
}

// SignUp is function to handle "/api/sign-up"
func (app *application) SignUp(w http.ResponseWriter, r *http.Request) {
	var u User
	db, err := driver.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Fatal(err)
	}
	var dbUser User
	db.Where("email = $1", u.Email).First(&dbUser)

	if dbUser.Email != "" {
		err = errors.New("email is already using")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	u.Password, err = bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		log.Fatal(err)
		return
	}
	u.IsAdmin = false
	db.Table("users").Create(&u)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	http.Redirect(w, r, "localhost:3000/login", http.StatusOK)

	driver.CloseDatabase(db)
}

/*
	db.Raw(`
	SELECT
	    	user.id, user.email, user.firstname, user.lastname
	FROM users u
		INNER JOIN tokens t on (u.ID = t.userID)
	WHERE
	    t.hash = $1
		and t.expiry > $2
	`, t.Token, time.Now()).Scan(&u)


*/

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var validToken string
	db, err := driver.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.CloseDatabase(db)

	var authDetails Authentication
	var dbuser User

	err = json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err error
		err = errors.New("error reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	db.Where("email = $1", authDetails.Email).First(&dbuser)
	if dbuser.Email == "" {
		var err error
		err = errors.New("invalid authentication credentials")
		w.Header().Set("Content-Type", "application/json0")
		json.NewEncoder(w).Encode(err)
		return
	}

	err = bcrypt.CompareHashAndPassword(dbuser.Password, []byte(authDetails.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		err = errors.New("invalid authentication credentials")
		w.Header().Set("Content-Type", "application/json0")
		json.NewEncoder(w).Encode(err)
		return
	} else if err != nil {
		err = errors.New("error compering passwords")
		w.Header().Set("Content-Type", "application/json0")
		json.NewEncoder(w).Encode(err)
		return
	}

	if !dbuser.IsAdmin {
		validToken, err = app.GenerateToken(authDetails.Email, "user")
		var t Token
		t.Email = dbuser.Email
		t.Role = "user"
		t.Token = validToken
		t.Expiry = time.Now().Add(30 * time.Minute)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			fmt.Println("error generating json response")
			return
		}
		return
	}
	validToken, err = app.GenerateToken(authDetails.Email, "admin")

	var u models.User
	// store token in database, encode userID
	u, err = app.database.GetUserByEmail(authDetails.Email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	var t Token
	t.UsedId = int(u.ID)
	t.Username = u.Username
	t.Email = authDetails.Email
	t.Role = "admin"
	t.Token = validToken
	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		fmt.Println("error generating json response")
		return
	}
}

// if user is authenticated
func (app *application) Authorize(w http.ResponseWriter, r *http.Request) {
	var auth Authentication

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	u, err := app.database.GetUserByEmail(auth.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		out, err := json.MarshalIndent(errors.New("can not retrieve user from database"), "", "\t")
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		w.Write(out)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func (app *application) Categories(w http.ResponseWriter, r *http.Request) {
	db, err := driver.Connect()
	if err != nil {
		app.errorLog.Println(err)
	}
	defer driver.CloseDatabase(db)
	var categories []Categories

	db.Table("categories").Find(&categories)
	json.NewEncoder(w).Encode(categories)
}

func (app *application) AdminMainPage(w http.ResponseWriter, r *http.Request) {
	books, err := app.database.ReturnBooks()
	if err != nil {
		app.errorLog.Println(err)
	}
	data := make(map[string]interface{})
	data["books"] = books
	users, err := app.database.GetAllUsers()
	if err != nil {
		app.errorLog.Println(err)
	}
	data["users"] = users
	if err := app.renderTemplate(w, r, "index", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ManageUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.database.GetAllUsers()
	if err != nil {
		app.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["users"] = users

	if err := app.renderTemplate(w, r, "manage-users", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ManageBooks(w http.ResponseWriter, r *http.Request) {
	books, err := app.database.GetAllBookRecords()
	if err != nil {
		app.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["books"] = books

	if err := app.renderTemplate(w, r, "manage-books", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ProceedUpdateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
}

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	if err := app.renderTemplate(w, r, "create", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ProceedCreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user models.User
	var err error

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Password, err = bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		app.errorLog.Println(err)
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if r.FormValue("is_admin") == "true" {
		user.IsAdmin = true
	} else {
		user.IsAdmin = false
	}

	json.NewEncoder(w).Encode(user)
	app.database.InsertUser(user)
	http.Redirect(w, r, "localhost:4001/admin", http.StatusOK)
}

func (app *application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "create", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) DetailedUser(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	if err := app.renderTemplate(w, r, "create", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}
func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var id int
	json.NewDecoder(r.Body).Decode(id)
	data := make(map[string]interface{})

	if err := app.renderTemplate(w, r, "create", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) CreateBook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var book models.Book
	var err error

	book.Title = r.FormValue("title")
	book.Description = r.FormValue("description")
	err = json.Unmarshal([]byte(r.FormValue("author")), &book.Author)
	if err != nil {
		app.errorLog.Println(err)
	}
	book.Price, _ = strconv.Atoi(r.FormValue("price"))
	book.Amount, _ = strconv.Atoi(r.FormValue("amount"))
	book.Image = r.FormValue("image")

	data := make(map[string]interface{})

	if err := app.renderTemplate(w, r, "create", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ProceedCreateBook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var book models.Book
	var err error

	book.Title = r.FormValue("title")
	book.Description = r.FormValue("description")
	json.Unmarshal([]byte(r.FormValue("author")), &book.Author)
	book.Price, err = strconv.Atoi(r.FormValue("price"))
	if err != nil {
		app.errorLog.Println(err)
	}
	book.Amount, err = strconv.Atoi(r.FormValue("amount"))
	// file -> image
	if err != nil {
		app.errorLog.Println(err)
	}

	book.Image = r.FormValue("image")

	if r.FormValue("is_avaliable") == "true" {
		book.IsAvaliable = true
	} else {
		book.IsAvaliable = false
	}

	http.Redirect(w, r, "localhost:4001", http.StatusCreated)
}

func (app *application) UpdateBook(w http.ResponseWriter, r *http.Request) {
	// form action => ProceedUpdateBook
	if err := app.renderTemplate(w, r, "update-book", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ProceedUpdateBook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var err error
	var newBook models.Book
	if app.HelperParseForm(r, "title") {
		newBook.Title = r.FormValue("title")
	}
	if app.HelperParseForm(r, "description") {
		newBook.Description = r.FormValue("description")
	}
	if app.HelperParseForm(r, "price") {
		newBook.Price, err = strconv.Atoi(r.FormValue("price"))
		if err != nil {
			app.errorLog.Println(err)
		}
	}
	if app.HelperParseForm(r, "amount") {
		newBook.Amount, err = strconv.Atoi(r.FormValue("amount"))
		if err != nil {
			app.errorLog.Println(err)
		}
	}
	if app.HelperParseForm(r, "image") {
		newBook.Image = r.FormValue("image")
	}
	if r.FormValue("is_avaliable") == "true" {
		newBook.IsAvaliable = true
	} else {
		newBook.IsAvaliable = false
	}

	if err := app.database.UpdateBook(newBook); err != nil {
		app.errorLog.Println(err)
	}
	http.Redirect(w, r, "localhost:4001", http.StatusAccepted)
}

func (app *application) DetailedBook(w http.ResponseWriter, r *http.Request) {
	var id string
	json.NewDecoder(r.Body).Decode(&id)
	bookID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLog.Println(err)
	}
	book, err := app.database.GetBookById(bookID)
	if err != nil {
		app.errorLog.Println(err)
	}

	data := make(map[string]interface{})
	data["book"] = book

	if err := app.renderTemplate(w, r, "detailed-book", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// get book id from manage-books.go
	// delete book
	var id int
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		app.errorLog.Println(err)
	}

	book, err := app.database.GetBookById(id)
	if err != nil {
		app.errorLog.Println(err)
	}
	err = app.database.DeleteBook(book)
	if err != nil {
		app.errorLog.Println(err)
	}

	http.Redirect(w, r, "localhost:4001/admin", http.StatusOK)
}
