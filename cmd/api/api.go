package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jangidRkt08/go-Ecom/service/user"

)

type APIserver struct{
	addr string
	db *sql.DB
}


func NewAPIserver(addr string, db *sql.DB) *APIserver{
	return &APIserver{
		addr: addr,
		db: db,
	}
}

func (s *APIserver) Run() error{
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	
	return http.ListenAndServe(s.addr, router)
}
