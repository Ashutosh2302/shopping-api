package routes

import (
	"database/sql"
	"shopping_api/pkg/controller"
	"shopping_api/pkg/service"

	"github.com/gin-gonic/gin"
)

func AddLoginEndpoints(db *sql.DB, rg *gin.RouterGroup) {
	loginService := service.NewLoginService(db)
	controller := controller.NewLoginController(loginService)
	loginGroup := rg.Group("")
	loginGroup.POST("/login", controller.Login)
	loginGroup.POST("/signup", controller.SignUp)
}
