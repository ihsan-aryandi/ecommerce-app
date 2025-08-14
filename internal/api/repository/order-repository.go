package repository

type OrderRepository struct {
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (OrderRepository) GetTotal() []string {
	return []string{"Total", "Rp. 12.000,00"}
}
