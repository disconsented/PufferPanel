package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/middleware"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"net/http"
)

func Reauth(c *gin.Context) {
	db := middleware.GetDatabase(c)
	ps := &services.Permission{DB: db}

	user, _ := c.MustGet("user").(*models.User)

	perms, err := ps.GetForUserAndServer(user.ID, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	session, err := services.GenerateSession(user.ID)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	data := &LoginResponse{}
	data.Session = session
	data.Scopes = perms.ToScopes()

	c.JSON(http.StatusOK, data)
}
