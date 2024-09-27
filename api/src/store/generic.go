package store

import (
	"context"
)

// InsertIntoTable inserts a given model into the database.
func InsertIntoTable(model interface{}) (interface{}, error) {
	db := GetBunConnection()

	// Insert the model into the table.
	res, err := db.NewInsert().Model(model).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}
