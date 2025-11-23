package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jangidRkt08/go-Ecom/types"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T){
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "123",
			Email: "user@gmail.com",
			Password: "testuser",
		}
		marshalled, _ := json.Marshal(payload)


		req, err := http.NewRequest(http.MethodPost,"/register",bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
		}
		// handler.HandleRegister(nil, req)
	})
}

type mockUserStore struct {
	
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {

	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}