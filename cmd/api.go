package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	repo "github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/adapters/postgresql/sqlc"
	"github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/orders"
	"github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// user -> handler GET /products -> service e getProducts -> repo SELECT * FROM PRODUCTS

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)	// important for rate limiting and tracings
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)  // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})
		
	// http.ListenAndServe(":3333", r)
	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Println("invalid product id")
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		productHandler.GetProduct(w, r, id)
	})

	ordersService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(ordersService)
	r.Post("/orders", ordersHandler.CreateOrder)

	return r
	}

// run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Second * 60,
		}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}


type application struct {
	config config
	db *pgx.Conn
	// logger
	// db driver
}

type config struct {
	addr string
	db dbConfig
}

type dbConfig struct {
	dsn string
}
