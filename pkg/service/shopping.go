package service

import (
	"database/sql"
	"errors"
	"fmt"
	"shopping_api/pkg/dto"
	"shopping_api/pkg/types"
	"shopping_api/utils"

	"github.com/gin-gonic/gin"
)

type ShoppingService struct {
	db *sql.DB
}

func NewShoppingService(db *sql.DB) *ShoppingService {
	return &ShoppingService{db}
}

func (s *ShoppingService) updateListUpdatedAt(listId string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE shopping_list set updatedAt = now() WHERE id = $1", listId)
	if err != nil {
		fmt.Print("Failed to update list", err)
		return err
	}
	return nil
}

func (s *ShoppingService) CreateShoppingList(ctx *gin.Context, userId string, l dto.CreateShoppingListRequest) (*types.ShoppingListWithoutItems, error) {

	query := `
	INSERT INTO shopping_list(userId, name)
	VALUES($1, $2)
	RETURNING id, name, createdAt, updatedAt
	`
	var createdAt string
	var updatedAt string
	var list types.ShoppingListWithoutItems
	err := s.db.QueryRow(query, userId, l.Name).Scan(&list.Id, &list.Name, &createdAt, &updatedAt)

	if err != nil {
		fmt.Print("Error creating shopping list", err)
		return nil, errors.New("failed to create shopping list")
	}

	time, err := utils.GetEpochTime(createdAt)
	if err != nil {
		return nil, errors.New("failed to get epoch time")
	}
	list.CreatedAt = time

	time, err = utils.GetEpochTime(updatedAt)
	if err != nil {
		return nil, errors.New("failed to get epoch time")
	}
	list.UpdatedAt = time
	return &list, nil
}

func (s *ShoppingService) CreateShoppingListItem(ctx *gin.Context, listId string, i dto.CreateShoppingListItemRequest) (*types.ShoppingListItem, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	query := `
	INSERT INTO shopping_list_item(shoppingListId, name, price)
	VALUES($1, $2, $3)
	RETURNING id, name, picked, price, createdAt
	`
	var createdAt string

	var item types.ShoppingListItem

	err = tx.QueryRow(query, listId, i.Name, i.Price).Scan(&item.Id, &item.Name, &item.Picked, &item.Price, &createdAt)

	if err != nil {
		fmt.Print("Error creating shopping list item", err)
		return nil, errors.New("failed to create shopping list item")
	}

	err = s.updateListUpdatedAt(listId, tx)
	if err != nil {
		return nil, errors.New("failed to update list")
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	time, err := utils.GetEpochTime(createdAt)
	if err != nil {
		return nil, errors.New("failed to get epoch time")
	}
	item.CreatedAt = time

	return &item, nil
}

func (s *ShoppingService) PickupListItem(ctx *gin.Context, itemId string, listId string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
	UPDATE shopping_list_item
	SET 
		picked = true
	WHERE
		id = $1 AND shoppingListId = $2

	`
	_, err = tx.Exec(query, itemId, listId)

	if err != nil {
		fmt.Print("Error creating shopping list item", err)
		return errors.New("failed to create shopping list item")
	}

	err = s.updateListUpdatedAt(listId, tx)
	if err != nil {
		return errors.New("failed to update list")
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ShoppingService) DeleteListItem(ctx *gin.Context, itemId string, listId string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
	DELETE FROM shopping_list_item
	WHERE id = $1 and shoppingListId = $2
	`

	_, err = tx.Exec(query, itemId, listId)

	if err != nil {
		fmt.Print("Error deleting shopping list item", err)
		return errors.New("failed to delete shopping list item")
	}

	err = s.updateListUpdatedAt(listId, tx)
	if err != nil {
		return errors.New("failed to update list")
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ShoppingService) GetShoppingLists(ctx *gin.Context, loggedInuserId string) (*types.ShoppingListsData, error) {
	query := `
	WITH lists AS (
		SELECT id, name, createdAt, updatedAt
		FROM shopping_list
		WHERE userId = $1
	)
	SELECT *, (SELECT COUNT(*) FROM lists) AS total_count
	FROM lists
	ORDER BY
	lists.updatedAt DESC
	`
	rows, err := s.db.Query(query, loggedInuserId)

	if err != nil {
		fmt.Print("Error scanning lists:", err)
		return nil, errors.New("failed to fetch shopping lists")
	}

	var lists []*types.ShoppingListWithoutItems
	c := 0
	var count *int = &c
	for rows.Next() {
		var list types.ShoppingListWithoutItems
		var createdAt string
		var updatedAt string
		err := rows.Scan(
			&list.Id,
			&list.Name,
			&createdAt,
			&updatedAt,
			&count,
		)

		if err != nil {
			fmt.Print("Error scanning lists:", err)
			return nil, errors.New("failed to fetch shopping lists")
		}

		time, err := utils.GetEpochTime(createdAt)
		if err != nil {
			return nil, errors.New("failed to get epoch time")
		}
		list.CreatedAt = time

		time, err = utils.GetEpochTime(updatedAt)
		if err != nil {
			return nil, errors.New("failed to get epoch time")
		}

		list.UpdatedAt = time

		lists = append(lists, &list)
	}

	return &types.ShoppingListsData{
		TotalCount: *count,
		Data:       lists,
	}, nil
}

func (s *ShoppingService) GetShoppingListItems(ctx *gin.Context, loggedInuserId string, listId string) (*types.ShoppingListItemsData, error) {

	query := `
	WITH items AS (
		SELECT 
			sli.id,
			sli.name, 
			sli.picked,
			sli.price,
			sli.createdAt
		FROM 
			shopping_list_item sli
		JOIN 
			shopping_list sl on sli.shoppingListId = sl.id
		WHERE
			sl.id = $1
	)
	select *, (SELECT COUNT(*) FROM items) AS total_items_count, (SELECT COUNT(*) FROM items WHERE picked is True) AS total_picked_count
	FROM items
	`
	rows, err := s.db.Query(query, listId)

	if err != nil {
		fmt.Print("Error fetching items:", err)
		return nil, errors.New("failed to fetch shopping list items")
	}

	var items []*types.ShoppingListItem
	c := 0
	var count *int = &c
	p := 0
	var pickedCount *int = &p
	for rows.Next() {
		var item types.ShoppingListItem
		var createdAt string

		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Picked,
			&item.Price,
			&createdAt,
			&count,
			&pickedCount,
		)

		if err != nil {
			fmt.Print("Error scanning list items:", err)
			return nil, errors.New("failed to fetch shopping list items")
		}

		time, err := utils.GetEpochTime(createdAt)
		if err != nil {
			return nil, errors.New("failed to get epoch time")
		}
		item.CreatedAt = time

		items = append(items, &item)
	}

	return &types.ShoppingListItemsData{
		TotalCount:  *count,
		TotalPicked: *pickedCount,
		Data:        items,
	}, nil
}
