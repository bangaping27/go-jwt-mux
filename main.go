package main

import (
	"go-jwt-mux/controllers/authcontroller"
	"go-jwt-mux/controllers/productcontroller"
	"go-jwt-mux/middlewares"
	"go-jwt-mux/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	models.ConnectDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productcontroller.Index).Methods("GET")
	api.Use(middlewares.JWTMiddleware)
	log.Fatal(http.ListenAndServe(":3000", r))
}
