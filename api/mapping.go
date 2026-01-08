package api

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/db"
)

func (a *Api) GetItemMappings(c echo.Context) error {
	store := db.NewItemAccountMappingStore(db.GetInstance())
	res, err := store.GetAll()
	if err != nil {
		return a.ReturnError(c, 500, "retreiving data failed")
	}

	return c.JSON(200, map[string]any{
		"mapping": res,
	})
}

func (a *Api) CreateItemMapping(c echo.Context) error {
	mapping := db.NewItemAccountMapping()
	err := json.NewDecoder(c.Request().Body).Decode(mapping)
	if err != nil {
		return a.ReturnError(c, 500, "request parsing failed")
	}

	mapping.Id = ""
	store := db.NewItemAccountMappingStore(db.GetInstance())
	item, err := store.GetByItemName(mapping.ItemName)
	if err == nil && item != nil {
		item.AccountId = mapping.AccountId
		mapping = item
		err = store.Update(mapping)
		if err != nil {
			return a.ReturnError(c, 500, "saving data failed")
		}
	} else {
		err = store.Insert(mapping)
		if err != nil {
			return a.ReturnError(c, 500, "saving data failed")
		}
	}

	return c.JSON(200, mapping)
}
