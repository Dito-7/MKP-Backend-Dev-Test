package handler

import (
	"errors"
	"net/http"

	"mkp-cinema-ticketing/internal/delivery/http/middleware"
	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/usecase"
	"mkp-cinema-ticketing/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format data login tidak valid", err.Error())
		return
	}

	result, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Email atau kata sandi salah", nil)
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, http.StatusOK, "Logout berhasil. Hapus bearer token dari client", nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data registrasi tidak valid", err.Error())
		return
	}

	result, err := h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			response.Error(c, http.StatusConflict, "Email sudah terdaftar", nil)
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Registrasi berhasil", result)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get(middleware.CtxUserIDKey)
	if !exists {
		response.Unauthorized(c, "Sesi tidak valid")
		return
	}

	result, err := h.authUsecase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		response.NotFound(c, "Pengguna tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Profil pengguna berhasil didapatkan", result)
}
