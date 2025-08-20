package repository

import (
	"ecommerce-app/internal/api/model"

	"github.com/doug-martin/goqu/v9"
)

type CartRepository struct {
	db        *goqu.Database
	tableName string
}

func NewCartRepository(db *goqu.Database) *CartRepository {
	return &CartRepository{
		db:        db,
		tableName: "carts",
	}
}

func (r CartRepository) FindByUserIdTx(tx *goqu.TxDatabase, userId int64) (*model.CartModel, error) {
	query, args, err := tx.Select(
		goqu.C("id"),
		goqu.C("user_id"),
	).
		From(goqu.T(r.tableName)).
		Where(goqu.Ex{
			"user_id":    userId,
			"deleted_at": goqu.Op{"isnot": nil},
		}).
		Prepared(true).
		ToSQL()

	if err != nil {
		return nil, err
	}

	result := new(model.CartModel)
	err = tx.
		QueryRow(query, args).
		Scan()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r CartRepository) InsertTx(tx *goqu.TxDatabase, cartModel *model.CartModel) (id int64, err error) {
	record := goqu.Record{
		"user_id":    cartModel.UserId.Int64,
		"created_at": cartModel.CreatedAt.Time,
	}

	query, args, err := tx.Insert(goqu.T(r.tableName)).
		Rows(record).
		Prepared(true).
		ToSQL()

	if err != nil {
		return 0, err
	}

	err = tx.
		QueryRow(query, args).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
