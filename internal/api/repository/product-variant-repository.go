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

func (r ProductVariantRepository) FindByIds(idList []int64) ([]model.ProductVariant, error) {
	var result []model.ProductVariant

	selectCols := []interface{}{
		goqu.I("pv.id"),
		goqu.I("pv.price"),
		goqu.I("pv.stock"),
		goqu.I("pv.weight"),
		goqu.I("pv.product_id"),
		goqu.I("product.name").As("product_name"),
	}
	query := r.db.
		Select(selectCols...).
		From(goqu.T(r.tableName).As("pv")).
		InnerJoin(
			goqu.T("products").As("product"),
			goqu.On(goqu.Ex{"pv.product_id": goqu.I("product.id")}),
		).
		Where(
			goqu.I("pv.id").In(idList),
			goqu.I("pv.deleted_at").IsNull(),
		)

	err := query.ScanStructs(&result)
	return result, err
}

func (r ProductVariantRepository) convertInt64ListToInterfaceList(list []int64) []interface{} {
	var result []interface{}

	for _, val := range list {
		result = append(result, val)
	}

	return result
}
