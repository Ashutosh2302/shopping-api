package routes

import (
	"database/sql"
	"shopping_api/middlewares"
	"shopping_api/pkg/controller"
	"shopping_api/pkg/service"

	"github.com/gin-gonic/gin"
)

func AddShoppingEndpoints(db *sql.DB, rg *gin.RouterGroup) {
	shoppingService := service.NewShoppingService(db)
	controller := controller.NewShoppingController(shoppingService)
	shoppingGroup := rg.Group("/shopping")
	shoppingGroup.Use(middlewares.Authenticate)
	shoppingGroup.GET("/:id", controller.GetShoppingListItems)
	shoppingGroup.GET("", controller.GetShoppingLists)
	shoppingGroup.POST("", controller.CreateShoppingList)
	shoppingGroup.POST("/:id/item", controller.CreateShoppingListItem)
	shoppingGroup.PATCH("/:id/pickup/:itemId", controller.PickupListItem)
	shoppingGroup.DELETE("/:id/remove/:itemId", controller.DeleteListItem)
}
