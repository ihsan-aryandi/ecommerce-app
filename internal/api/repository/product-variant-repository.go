package repository

import (
	"ecommerce-app/internal/api/model"

	"github.com/doug-martin/goqu/v9"
)

type ProductVariantRepository struct {
	db        *goqu.Database
	tableName string
}

func NewProductVariantRepository(db *goqu.Database) *ProductVariantRepository {
	return &ProductVariantRepository{
		db:        db,
		tableName: "product_variants",
	}
}

func (r ProductVariantRepository) FindByIdTx(tx *goqu.TxDatabase, id int64) (*model.ProductVariant, error) {
	result := new(model.ProductVariant)
	found, err := tx.
		From(goqu.T(r.tableName)).
		Where(
			goqu.C("id").Eq(id),
			goqu.C("deleted_at").IsNull(),
		).
		ScanStruct(result)
	return handleNullAndError(result, found, err)
}
