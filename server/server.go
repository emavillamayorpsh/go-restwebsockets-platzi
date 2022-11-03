package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/emavillamayorpsh/rest-ws/database"
	"github.com/emavillamayorpsh/rest-ws/repository"
	"github.com/gorilla/mux"
)

// config of the server in order to be executed

type Config struct {
	Port string
	JWTSecret string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker , error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database is required")
	}

	broker := &Broker{
		config: config,
		router: *mux.NewRouter(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func (s Server, r *mux.Router)) {
	b.router = *mux.NewRouter()
	binder(b, &b.router)

	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	// call the abstract interface in order to instance the "postgres" repo
	// if need to change to another db then we pass the other db here
	repository.SetRepository(repo)

	log.Println("Starting server on port ", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, &b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}