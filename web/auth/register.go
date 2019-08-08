package auth

import (
	"github.com/gin-gonic/gin"
	builder "github.com/pufferpanel/apufferi/response"
	"github.com/pufferpanel/pufferpanel/database"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterPost(c *gin.Context) {
	response := builder.From(c)
	response.Fail()
	response.Message("unknown error occurred")

	request := &registerRequest{}
	err := c.BindJSON(request)

	if err != nil {
		response.Fail().Status(400).Data(err)
		return
	}

	validate := validator.New()
	err = validate.Struct(request.Data)
	if err != nil {
		response.Fail().Status(400).Data(err)
		return
	}

	db, err := database.GetConnection()
	if shared.HandleError(response, err) {
		return
	}

	us := &services.User{DB: db}
	os := services.GetOAuth(db)

	user := &models.User{Username: request.Data.Username, Email: request.Data.Email}
	err = user.SetPassword(request.Data.Password)
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	err = us.Create(user)
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	client, _, err := os.GetByUser(user)
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	err = os.AddScope(client, nil, "servers.view")
	if err != nil {
		response.Fail().Status(400).Error(err)
		return
	}

	response.Success().Message("")
}

type registerRequest struct {
	Data registerRequestData `json:"data"`
}

type registerRequestData struct {
	Username string `json:"username" validate:"min=3,printascii,required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
