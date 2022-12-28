package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cr00z/goSimpleChat/internal/domain"
)

type ResponseOK struct {
	ID int64 `json:"id"`
}

type ResponseToken struct {
	Token string `json:"token"`
}

// @summary Register
// @description Create account
// @tags Auth
// @accept json
// @produce json
// @param input body domain.User true "account info"
// @router /register [post]
// @success 200 {object} ResponseOK
// @failure 400 {object} nil
// @failure 500 {object} nil
func (h Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input domain.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateUser(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	renderJSON(w, ResponseOK{id})
}

// @summary Login
// @description Login to account
// @tags Auth
// @accept json
// @produce json
// @param input body domain.User true "account info"
// @router /login [post]
// @success 200 {object} ResponseToken
// @failure 400,401 {object} nil
// @failure 500 {object} nil
func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input domain.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.service.Authorization.GenerateJWT(input)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	renderJSON(w, ResponseToken{token})
}
