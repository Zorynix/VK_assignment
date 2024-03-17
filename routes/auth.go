package routes

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"vk.com/m/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles user login requests
// @Summary User login
// @Description handles login requests by checking username and password
// @Accept  json
// @Produce  json
// @Param   LoginRequest  body      LoginRequest  true  "Login Credentials"
// @Success 200 {object} LoginResponse "Returns login token"
// @Failure 400,401 "Invalid request or Unauthorized"
// @Router /login [post]
func (router *Router) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Invalid login request payload")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var role string
	if req.Username == "admin" && req.Password == "password" {
		role = "admin"
		log.Info().Str("username", req.Username).Str("role", role).Msg("User logged in successfully")
	} else if req.Username == "user" && req.Password == "password" {
		role = "user"
		log.Info().Str("username", req.Username).Str("role", role).Msg("User logged in successfully")
	} else {
		log.Warn().Str("username", req.Username).Msg("Unauthorized login attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(1, role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token: token,
	}
	json.NewEncoder(w).Encode(resp)
}
