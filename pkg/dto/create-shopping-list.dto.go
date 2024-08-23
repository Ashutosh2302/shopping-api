package dto

type CreateShoppingListRequest struct {
	Name string `json:"name" binding:"required"`
}
