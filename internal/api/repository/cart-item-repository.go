package repository

import (
	"ecommerce-app/internal/api/model"

	"github.com/doug-martin/goqu/v9"
)

type CartItemRepository struct {
	db        *goqu.Database
	tableName string
}

func NewCartItemRepository(db *goqu.Database) *CartItemRepository {
	return &CartItemRepository{
		db:        db,
		tableName: "cart_items",
	}
}

func (r CartItemRepository) FindByCartIdAndProductVariantIdTx(tx *goqu.TxDatabase, cartId, productVariantId int64) (*model.CartItemModel, error) {
	result := new(model.CartItemModel)

	found, err := tx.
		From(goqu.T(r.tableName)).
		Where(
			goqu.C("cart_id").Eq(cartId),
			goqu.C("product_variant_id").Eq(productVariantId),
			goqu.C("deleted_at").IsNull(),
		).
		ForUpdate(goqu.Wait).
		ScanStruct(result)
	return handleNullAndError(result, found, err)
}

func (r CartItemRepository) InsertTx(tx *goqu.TxDatabase, itemModel *model.CartItemModel) (id int64, err error) {
	insert := tx.
		Insert(goqu.T(r.tableName)).
		Rows(itemModel).
		Returning("id").
		Executor()

	_, err = insert.ScanVal(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r CartItemRepository) UpdateByIdTx(tx *goqu.TxDatabase, itemModel *model.CartItemModel) error {
	_, err := tx.
		Update(goqu.T(r.tableName)).
		Set(itemModel).
		Where(goqu.C("id").Eq(itemModel.ID.Int64)).
		Executor().
		Exec()

	return err
}

func (r CartItemRepository) UpdateQuantityByIdTx(tx *goqu.TxDatabase, itemModel *model.CartItemModel) error {
	_, err := tx.
		Update(goqu.T(r.tableName)).
		Set(goqu.Record{
			"qty":        itemModel.Qty.Int32,
			"updated_at": itemModel.UpdatedAt.Time,
			"updated_by": itemModel.UpdatedBy.Int64,
		}).
		Where(goqu.C("id").Eq(itemModel.ID.Int64)).
		Executor().
		Exec()

	return err
}
