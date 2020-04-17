package manager

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"
	"user-service/model"
)

type UserHandler struct {
	UserPath    string
	userService UserService
}

func NewUserHandler() UserHandler {

	return UserHandler{
		UserPath:    "/user/",
		userService: NewUserService(),
	}

}

type IUserHandler interface {
	ManageUser(w http.ResponseWriter, r *http.Request)
	CheckPassword(w http.ResponseWriter, r *http.Request)
	ReadUser(w http.ResponseWriter, r *http.Request)
	ReadRolesFromUser(w http.ResponseWriter, r *http.Request)
}

func (h *UserHandler) ManageUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var userDTO model.UserDTO
		err := json.NewDecoder(r.Body).Decode(&userDTO)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		err = h.createUser(userDTO)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
	if r.Method == "PUT" {

		path := html.EscapeString(r.URL.Path)
		userID, err := strconv.ParseUint(strings.ReplaceAll(path, h.UserPath, ""), 10, 64)
		if userID == 0 || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Id of user must be specified and not be null if put is http Method"))
		}
		var userDTO model.UserDTO

		err = json.NewDecoder(r.Body).Decode(&userDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userDTO.Password != "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not allowed to set Password with this Method"))
		}

		h.userService.UpdateUser(uint(userID), userDTO)

	}
	if r.Method == "GET" {

		path := html.EscapeString(r.URL.Path)
		userID, err := strconv.ParseUint(strings.ReplaceAll(path, h.UserPath, ""), 10, 64)
		if userID != 0 {
			user, err := h.userService.FindUser(uint(userID))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
		users, err := h.userService.FindAllUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}

}

func (h *UserHandler) createUser(userDTO model.UserDTO) error {
	return h.userService.CreateUser(userDTO)
}
