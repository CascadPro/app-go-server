package auth

import (
	"fmt"
	"net/http"

	"cascade/pkg/filter"
	"cascade/pkg/utils/authutils/sessions"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	sessionID, isSIDExists := c.Get("sessionID")
	userID, isIDExists := c.Get("userID")

	if isSIDExists && isIDExists {
		err := sessions.DeleteUserSession(fmt.Sprint(userID), fmt.Sprint(sessionID))
		if err != nil {
			filter.Error(c, filter.ErrorParams{Status: http.StatusInternalServerError, Cause: err.Error()})
			return
		}
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, false)

	filter.Success(c, "Вы успешно вышли из аккаунта")
}
