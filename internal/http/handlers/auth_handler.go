package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/book_API/internal/config"
	"github.com/qullDev/book_API/internal/domain/user"
	appauth "github.com/qullDev/book_API/internal/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	ts  *appauth.TokenStore
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, ts *appauth.TokenStore, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, ts: ts, cfg: cfg}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenPairResp struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	UserID           string `json:"user_id,omitempty"`
	Username         string `json:"username,omitempty"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload tidak valid", "error": err.Error()})
		return
	}

	var u user.User
	if err := h.db.Where("username = ?", req.Username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "username atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses login"})
		return
	}

	// Verifikasi password: coba bcrypt, jika gagal coba plain match sebagai fallback dev
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		if u.Password != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "username atau password salah"})
			return
		}
	}

	// Buat token
	at, err := appauth.GenerateAccessToken(h.cfg, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat access token"})
		return
	}
	rt, jti, err := appauth.GenerateRefreshToken(h.cfg, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat refresh token"})
		return
	}

	// Simpan RT ke Redis
	if err := h.ts.SaveRefreshToken(c.Request.Context(), u.ID, jti, h.cfg.RefreshTokenTTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menyimpan refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokenPairResp{
		AccessToken:      at,
		RefreshToken:     rt,
		TokenType:        "Bearer",
		ExpiresIn:        int64(h.cfg.AccessTokenTTL.Seconds()),
		RefreshExpiresIn: int64(h.cfg.RefreshTokenTTL.Seconds()),
		UserID:           u.ID.String(),
		Username:         u.Username,
	})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload tidak valid", "error": err.Error()})
		return
	}

	// Parse dan validasi refresh token
	claims, err := appauth.ParseToken(h.cfg, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token tidak valid"})
		return
	}

	// Pastikan RT ini masih valid di Redis
	valid, err := h.ts.VerifyRefreshToken(c.Request.Context(), claims.UserID, claims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memverifikasi refresh token"})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token sudah tidak berlaku"})
		return
	}

	// Rotasi RT: hapus yang lama, buat yang baru
	if err := h.ts.RevokeRefreshToken(c.Request.Context(), claims.UserID, claims.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mencabut refresh token"})
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token tidak valid"})
		return
	}

	newAT, err := appauth.GenerateAccessToken(h.cfg, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat access token"})
		return
	}
	newRT, newJTI, err := appauth.GenerateRefreshToken(h.cfg, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat refresh token"})
		return
	}
	if err := h.ts.SaveRefreshToken(c.Request.Context(), userID, newJTI, h.cfg.RefreshTokenTTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menyimpan refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokenPairResp{
		AccessToken:      newAT,
		RefreshToken:     newRT,
		TokenType:        "Bearer",
		ExpiresIn:        int64(h.cfg.AccessTokenTTL.Seconds()),
		RefreshExpiresIn: int64(h.cfg.RefreshTokenTTL.Seconds()),
	})
}

type logoutReq struct {
	RefreshToken string `json:"refresh_token"` // optional: jika kosong, revoke semua RT user
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// userID dari middleware JWT
	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	var req logoutReq
	_ = c.ShouldBindJSON(&req) // optional body

	if req.RefreshToken == "" {
		// Revoke semua token user
		if err := h.ts.RevokeAllRefreshTokens(c.Request.Context(), userIDStr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal logout"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
		return
	}

	// Jika disediakan refresh_token tertentu, revoke token tersebut
	claims, err := appauth.ParseToken(h.cfg, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token tidak valid"})
		return
	}
	// Pastikan token milik user yang sama
	if claims.UserID != userIDStr {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token tidak sesuai"})
		return
	}
	if err := h.ts.RevokeRefreshToken(c.Request.Context(), claims.UserID, claims.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
