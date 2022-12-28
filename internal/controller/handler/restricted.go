package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cr00z/goSimpleChat/internal/domain"
	"github.com/go-chi/chi/v5"
)

type ResponseStatus struct {
	Status string `json:"status"`
}

// @summary Post message
// @description Post message from user to chat
// @tags Message
// @security ApiKeyAuth
// @accept json
// @produce json
// @param input body domain.Message true "text only"
// @router /api/messages [post]
// @success 200 {object} ResponseStatus
// @failure 400,401 {object} nil
// @failure 500 {object} nil
func (h Handler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var message domain.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message.FromUserID = userID

	if err := h.service.CreateMessage(message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	renderJSON(w, ResponseStatus{"ok"})
}

// @summary Get messages
// @description Get messages from chat
// @tags Message
// @security ApiKeyAuth
// @produce json
// @router /api/messages [get]
// @success 200 {array} domain.Message
// @failure 401 {object} nil
// @failure 500 {object} nil
func (h Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := getUserID(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	renderJSON(w, h.service.GetMessages(0))
}

// @summary Send private message
// @description Send private message from one user to another
// @tags Message
// @security ApiKeyAuth
// @accept json
// @produce json
// @param id path int true "user id"
// @param input body domain.Message true "to_user and text only"
// @router /api/users/{id}/messages [post]
// @success 200 {object} ResponseStatus
// @failure 400,401 {object} nil
// @failure 500 {object} nil
func (h Handler) PostPrivateMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	toUserID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || toUserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message domain.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message.FromUserID = userID
	message.ToUserID = toUserID

	if err := h.service.CreateMessage(message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	renderJSON(w, map[string]string{"status": "ok"})
}

// @summary Get private messages
// @description Get your private messages
// @tags Message
// @security ApiKeyAuth
// @produce json
// @router /api/users/me/messages [get]
// @success 200 {array} domain.Message
// @failure 401 {object} nil
// @failure 500 {object} nil
func (h Handler) GetPrivateMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	renderJSON(w, h.service.GetMessages(userID))
}

// service

func getUserID(r *http.Request) (int64, error) {
	id := r.Context().Value(contextKey("ID"))
	if id == nil {
		return 0, domain.ErrorUserIDNotFound
	}

	userID, ok := id.(int64)
	if !ok {
		return 0, domain.ErrorUserIDInvalidType
	}

	return userID, nil
}
