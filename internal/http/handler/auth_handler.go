package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/auth"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	httpx "github.com/adityanuriskandar17/HRIS-BE/internal/http"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret string
}

type loginReq struct{ Email, Password string }

type loginRes struct{ Token, Role string }

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	var u model.UserAccount
	if err := h.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		http.Error(w, "invalid", 401)
		return
	}
	// TODO: check password hash (use util)
	tok, _ := auth.SignJWT(u.ID, string(u.Role), h.JWTSecret, 24*time.Hour)
	httpx.OK(w, loginRes{Token: tok, Role: string(u.Role)})
}
