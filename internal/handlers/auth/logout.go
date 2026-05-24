package auth

import (
	"fmt"

	"cascade/pkg/filter"
	"cascade/pkg/utils/authutils/sessions"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	sessionID, isExists := c.Get("sessionID")

	if isExists {
		sessions.DeleteSession(fmt.Sprintf("%s", sessionID))
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, false)

	filter.Success(c, "Вы успешно вышли из аккаунта")
}
