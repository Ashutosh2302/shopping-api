package types

type ShoppingListWithoutItems struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type ShoppingListsData struct {
	TotalCount int                         `json:"totalCount"`
	Data       []*ShoppingListWithoutItems `json:"data"`
}

type ShoppingListItem struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Picked    bool    `json:"picked"`
	Price     float32 `json:"price"`
	CreatedAt int64   `json:"createdAt"`
}

type ShoppingListItemsData struct {
	TotalCount  int                 `json:"totalCount"`
	TotalPicked int                 `json:"totalPicked"`
	Data        []*ShoppingListItem `json:"items"`
}
