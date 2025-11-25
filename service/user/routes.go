package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jangidRkt08/go-Ecom/config"
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

	// admin route
		router.HandleFunc("/users/{userID}", auth.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)

}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
		var payload types.LoginUserPayload

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

	u, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found valid email or password"))
		return
	}

	if !auth.ComparePassword([]byte(u.Password), []byte(payload.Password)){
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not Found, invalid email or password"))
		return
	}

	// utils.WriteJSON(w,http.StatusOK,map[string]string{"token": ""})

	secret := []byte(config.Envs.JWTSecret)
	token, err:= auth.CreateJWT(secret, u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK,map[string]string{"token":token})
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

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}