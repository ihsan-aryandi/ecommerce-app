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
	result := new(model.CartModel)

	found, err := tx.
		From(goqu.T(r.tableName)).
		Where(
			goqu.C("user_id").Eq(userId),
			goqu.C("deleted_at").IsNull(),
		).
		ScanStruct(result)

	return handleNullAndError(result, found, err)
}

func (r CartRepository) InsertTx(tx *goqu.TxDatabase, cartModel *model.CartModel) (id int64, err error) {
	insert := tx.Insert(goqu.T(r.tableName)).
		Returning("id").
		Rows(cartModel).
		Executor()

	_, err = insert.ScanVal(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
