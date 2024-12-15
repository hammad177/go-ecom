package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hammad177/ecom/services/cart"
	"github.com/hammad177/ecom/services/order"
	"github.com/hammad177/ecom/services/products"
	"github.com/hammad177/ecom/services/users"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	useStore := users.NewStore(s.db)
	userHandler := users.NewHandler(useStore)
	userHandler.RegisterRoutes(subRouter)

	productStore := products.NewStore(s.db)
	productHandler := products.NewHandler(productStore)
	productHandler.RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(productStore, orderStore, useStore)
	cartHandler.RegisterRoutes(subRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
