package rest_api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mikhatanu/dating-test/auth"
)

type Response struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// http handler for login
func Login(w http.ResponseWriter, r *http.Request) {
	response := &Response{}
	request := &signupRequest{}

	// allow only POST method
	if r.Method == http.MethodPost {
		// Check for content-type application/json
		ct := r.Header.Get("Content-Type")
		if ct != "" {
			mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
			if mediaType != "application/json" {
				response.Message = "Content-Type header is not application/json"
				JSONError(w, response, http.StatusUnsupportedMediaType)
				return
			}
		}

		// decode to struct
		decodeErr := json.NewDecoder(r.Body).Decode(request)
		if decodeErr != nil {
			log.Println(decodeErr.Error())
			response.Message = "error in post body"
			JSONError(w, response, http.StatusBadRequest)
			return
		}

		// get user
		u, err := auth.GetUser(request.Username)
		if err != nil {
			log.Println(err.Error())
			response.Message = err.Error()
			JSONError(w, response, http.StatusBadRequest)
			return
		}

		// check password
		err2 := auth.CheckPassword(request.Password, u.Password)
		if err2 != nil {
			log.Println(err2.Error())
			response.Message = err2.Error()
			JSONError(w, response, http.StatusBadRequest)
			return
		}

		response.Message = "success"
		response.Data = map[string]any{
			"username": u.Username,
			"id":       u.Id,
		}

		JSONError(w, response, http.StatusOK)
	} else {
		JSONError(w, "", http.StatusMethodNotAllowed)
		return
	}
}

// http handler for signup
func Signup(w http.ResponseWriter, r *http.Request) {
	response := &Response{}
	request := &signupRequest{}

	// only allow post method
	if r.Method == http.MethodPost {
		// Check for content-type application/json
		ct := r.Header.Get("Content-Type")
		if ct != "" {
			mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
			if mediaType != "application/json" {
				response.Message = "Content-Type header is not application/json"
				JSONError(w, response, http.StatusUnsupportedMediaType)
				return
			}
		}

		// decode body to struct
		decodeErr := json.NewDecoder(r.Body).Decode(request)
		if request.Username == "" || request.Password == "" {
			log.Println("empty password or username")
			response.Message = "empty password or username"
			JSONError(w, response, http.StatusBadRequest)
			return
		}
		if decodeErr != nil {
			log.Println(decodeErr.Error())
			response.Message = "error in post body"
			JSONError(w, response, http.StatusBadRequest)
			return
		}

		// sign up user
		err := auth.Signup(request.Username, request.Password)
		if err != nil {
			response.Message = err.Error()
			JSONError(w, response, http.StatusInternalServerError)
			return
		}

		response.Message = "success"
		JSONError(w, response, http.StatusOK)
		return
	} else {
		JSONError(w, "", http.StatusMethodNotAllowed)
		return
	}

}

// helper to return json
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; ")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
