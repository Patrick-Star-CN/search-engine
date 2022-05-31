package sessionService

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"search-engine/app/models"
	"search-engine/app/services/userService"
)

func SetUserSession(c *gin.Context, user *models.User) error {
	webSession := sessions.Default(c)
	webSession.Options(sessions.Options{
		Domain: c.Request.Header.Get("Origin"),
		MaxAge: 3600 * 24 * 7})
	webSession.Set("id", user.ID)
	return webSession.Save()
}

func GetUserSession(c *gin.Context) (*models.User, error) {
	webSession := sessions.Default(c)
	id := webSession.Get("id")
	if id == nil {
		return nil, errors.New("")
	}
	user, _ := userService.GetUserID(id.(int))
	if user == nil {
		ClearUserSession(c)
		return nil, errors.New("")
	}
	return user, nil
}

func CheckUserSession(c *gin.Context) bool {
	webSession := sessions.Default(c)
	id := webSession.Get("id")
	if id == nil {
		return false
	}
	return true
}

func ClearUserSession(c *gin.Context) {
	webSession := sessions.Default(c)
	webSession.Delete("id")
	webSession.Save()
	return
}
