package auth

import (
	"net/http"
	"time"

	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/utils"
	"cascade/pkg/utils/authutils"
	"cascade/pkg/utils/authutils/sessions"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenPayload struct {
	jwt.Claims

	UserID    uuid.UUID       `json:"user_id"`
	Role      models.UserRole `json:"role"`
	SessionID string          `json:"session_id"`
	ExpiresAt int64           `json:"exp"`
}

func GetNewTokens(c *gin.Context) {
	rt_s, token_err := c.Cookie("refresh_token")
	if token_err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusUnauthorized,
			Message: "Ключа для обновления не найдено или он просрочен! Авторизуйтесь снова",
			Cause:   token_err.Error(),
		})
		return
	}

	cfg, err := utils.LoadConfig()
	if err != nil {
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusInternalServerError,
			Message: "Something went wrong!",
			Cause:   err.Error(),
		})
		return
	}

	rt, v_err := authutils.ValidateToken(rt_s, cfg)
	if v_err != nil {
		c.SetCookie("refresh_token", "", -1, "/", "", false, false)
		filter.Error(c, filter.ErrorParams{
			Status:  http.StatusUnauthorized,
			Message: "Ошибка во время валидации ключа!",
			Cause:   v_err.Error(),
		})
		return
	}

	_, s_err := sessions.GetSessionByID(rt.SessionID)
	if s_err != nil {
		if s_err == redis.Nil {
			c.SetCookie("refresh_token", "", -1, "/", "", false, false)

			filter.Error(c, filter.ErrorParams{Status: http.StatusUnauthorized, Message: "Сессия не обнаружена!"})
		} else {
			filter.Error(c, filter.ErrorParams{Status: http.StatusBadRequest, Message: "Something went wrong!"})
		}
		return
	}

	new_at, new_rt := authutils.IssueTokens(rt.UserID, rt.Role, rt.SessionID, cfg)

	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    new_rt,
		Domain:   cfg.Domain,
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(authutils.RefreshTokenLifetime),
	})

	filter.Success(c, "Ключи обновлены!", gin.H{"access_token": new_at})
}
