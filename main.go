package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/emavillamayorpsh/rest-ws/handlers"
	"github.com/emavillamayorpsh/rest-ws/middleware"
	"github.com/emavillamayorpsh/rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// go get file and loaded so that go can access to the values in the .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get values from .env file variables
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")


	// create a new server
	s, err := server.NewServer(context.Background(), &server.Config{
		JWTSecret: JWT_SECRET,
		Port: PORT,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	// FOR EACH ROUTE WE WILL APPLY THIS MIDDLEWARE
	r.Use(middleware.CheckAuthMiddleware(s))


	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.InsertPostHandler((s))).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler((s))).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.UpdatePostHandler((s))).Methods(http.MethodPut)
	r.HandleFunc("/posts/{id}", handlers.DeletePostHandler((s))).Methods(http.MethodDelete)
}