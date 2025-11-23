package user

import (

	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jangidRkt08/go-Ecom/service/auth"
	"github.com/jangidRkt08/go-Ecom/types"
	"github.com/jangidRkt08/go-Ecom/utils"
)


type Handler struct {
	store types.UserStore
	
}

func NewHandler(store types.UserStore) *Handler{
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Get JSON Payload

	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r,&payload); err !=nil{
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v",errors))
		return
	}

	// Check user Exist
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil{
		utils.WriteError(w, http.StatusBadRequest,fmt.Errorf("user with %s email exist", payload.Email))
		return
	}

	hashedPassword,err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// if not exist create user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
