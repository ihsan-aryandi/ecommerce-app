package repository

import (
	"github.com/redis/go-redis/v9"
)

type CheckoutSessionRepository struct {
	rdb *redis.Client
}

func NewCheckoutSessionRepository(rdb *redis.Client) *CheckoutSessionRepository {
	return &CheckoutSessionRepository{
		rdb: rdb,
	}
}
