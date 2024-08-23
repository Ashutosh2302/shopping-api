package controller

import (
	"fmt"
	"net/http"
	"shopping_api/pkg/dto"
	"shopping_api/pkg/service"

	"github.com/gin-gonic/gin"
)

type ShoppingController struct {
	service *service.ShoppingService
}

func NewShoppingController(service *service.ShoppingService) *ShoppingController {
	return &ShoppingController{service: service}
}

func (c *ShoppingController) CreateShoppingList(ctx *gin.Context) {
	loggedInuserId := ctx.GetString("userId")
	if loggedInuserId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No userId found"})
		return
	}

	var createListReq dto.CreateShoppingListRequest
	err := ctx.ShouldBindJSON(&createListReq)

	if err != nil {
		fmt.Print("Error parsing create list equest body:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request body"})
		return
	}

	list, err := c.service.CreateShoppingList(ctx, loggedInuserId, createListReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not create shopping list"})
		return

	}
	ctx.JSON(http.StatusOK, list)
}

func (c *ShoppingController) CreateShoppingListItem(ctx *gin.Context) {
	loggedInUserId := ctx.GetString("userId")
	if loggedInUserId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No userId found"})
		return
	}

	var createListItemReq dto.CreateShoppingListItemRequest
	err := ctx.ShouldBindJSON(&createListItemReq)

	if err != nil {
		fmt.Print("Error parsing create list item request body:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request body"})
		return
	}
	listId := ctx.Param("id")
	item, err := c.service.CreateShoppingListItem(ctx, listId, createListItemReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not create shopping list item"})
		return

	}
	ctx.JSON(http.StatusOK, item)
}

func (c *ShoppingController) PickupListItem(ctx *gin.Context) {
	loggedInUserId := ctx.GetString("userId")
	if loggedInUserId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No userId found"})
		return
	}
	listId := ctx.Param("id")
	itemId := ctx.Param("itemId")
	err := c.service.PickupListItem(ctx, itemId, listId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not pickup item"})
		return

	}
	ctx.Status(http.StatusOK)
}

func (c *ShoppingController) DeleteListItem(ctx *gin.Context) {
	loggedInUserId := ctx.GetString("userId")
	if loggedInUserId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No userId found"})
		return
	}
	listId := ctx.Param("id")
	itemId := ctx.Param("itemId")
	err := c.service.DeleteListItem(ctx, itemId, listId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not pickup item"})
		return

	}
	ctx.Status(http.StatusOK)
}

func (c *ShoppingController) GetShoppingLists(ctx *gin.Context) {
	loggedInuserId := ctx.GetString("userId")
	if loggedInuserId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No userId found"})
		return
	}

	lists, err := c.service.GetShoppingLists(ctx, loggedInuserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch shopping lists"})
		return

	}
	ctx.JSON(http.StatusOK, lists)
}

func (c *ShoppingController) GetShoppingListItems(ctx *gin.Context) {
	loggedInuserId := ctx.GetString("userId")
	if loggedInuserId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "No userId found"})
		return
	}

	listId := ctx.Param("id")
	list, err := c.service.GetShoppingListItems(ctx, loggedInuserId, listId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch shopping list"})
		return

	}

	ctx.JSON(http.StatusOK, list)
}
