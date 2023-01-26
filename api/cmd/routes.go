package main

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/chi"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"https://*", "http://*"},
		AllowOriginFunc:    nil,
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "From"},
		ExposedHeaders:     nil,
		AllowCredentials:   false,
		MaxAge:             0,
		OptionsPassthrough: false,
		Debug:              false,
	}))

	// admin page routes
	router.Route("/admin", func(mux chi.Router) {
		mux.Get("/", app.AdminMainPage)
		mux.Get("/books", app.ManageBooks)

		mux.Get("/users", app.ManageUsers)

		mux.Get("/admin/books/create", app.CreateBook)
		mux.Post("/admin/books/create", app.ProceedCreateBook)

		mux.Get("/admin/books/update", app.UpdateBook)
		mux.Post("/admin/books/update", app.ProceedUpdateBook)

		mux.Get("/admin/books/details", app.DetailedBook)

		mux.Get("/admin/books/delete", app.DeleteBook)

		mux.Get("/admin/users/create", app.CreateUser)
		mux.Post("/admin/users/create", app.ProceedCreateUser)

		mux.Get("/admin/users/update", app.UpdateUser)
		mux.Post("/admin/users/update", app.ProceedUpdateUser)

		mux.Get("/admin/users/details", app.DetailedUser)

		mux.Get("/admin/users/delete", app.DeleteUser)
	})

	// for interacting with front-end
	router.Post("/api/authorized", app.Authorize)
	router.Post("/api/products", app.ListProducts)
	router.Get("/api/get-crypto-info", app.GetCryptoInfo)
	router.Post("/api/signup", app.SignUp)
	router.Post("/api/authenticate", app.Login)
	router.Get("/api/categories", app.Categories)

	return router
}
