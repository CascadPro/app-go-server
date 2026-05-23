package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"cascade/internal/models"
	"cascade/pkg/filter"
	"cascade/pkg/utils"
	"cascade/pkg/utils/authutils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")

		if bearer == "" {
			filter.Error(c, filter.ErrorParams{
				Status:  http.StatusUnauthorized,
				Message: "Ключ авторизации не найден в запросе!"})
			return
		}

		cfg, err := utils.LoadConfig()
		if err != nil {
			filter.Error(c, filter.ErrorParams{
				Status:  http.StatusInternalServerError,
				Message: "Something went wrong!",
				Cause:   err.Error()})
			return
		}

		tokenString := strings.Split(bearer, " ")[1]

		parsedToken, v_err := authutils.ValidateToken(tokenString, cfg)
		if v_err != nil {
			filter.Error(c, filter.ErrorParams{
				Status:  http.StatusUnauthorized,
				Message: "Ошибка во время проверки ключа!",
				Cause:   v_err.Error()})
			return
		}

		if len(roles) > 0 && !slices.Contains(roles, parsedToken.Role) {
			filter.Error(c, filter.ErrorParams{
				Status:  http.StatusForbidden,
				Message: "У вас недостаточно прав!"})
			return
		}

		c.Set("userID", parsedToken.UserID)
		c.Set("role", parsedToken.Role)
		c.Set("sessionID", parsedToken.SessionID)

		c.Next()
	}
}
