package repository

import (
	"context"
	"ecommerce-app/internal/api/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CheckoutSessionRepository struct {
	rdb    *redis.Client
	prefix string
}

func NewCheckoutSessionRepository(rdb *redis.Client) *CheckoutSessionRepository {
	return &CheckoutSessionRepository{
		rdb:    rdb,
		prefix: "checkout:",
	}
}

func (r CheckoutSessionRepository) generateKey(userId int64) string {
	return fmt.Sprintf("%s%d", r.prefix, userId)
}

func (r CheckoutSessionRepository) GetByUserId(userId int64) (*model.CheckoutSession, error) {
	key := r.generateKey(userId)

	cmd := r.rdb.Get(context.Background(), key)
	result, err := cmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	sessionResult := new(model.CheckoutSession)
	if err = json.Unmarshal([]byte(result), &sessionResult); err != nil {
		return nil, err
	}

	return sessionResult, cmd.Err()
}

func (r CheckoutSessionRepository) Save(checkoutSession *model.CheckoutSession, expiration time.Duration) (string, error) {
	key := r.generateKey(checkoutSession.UserId)

	uid := uuid.New()
	checkoutSession.CheckoutID = uid.String()

	data, _ := json.Marshal(checkoutSession)
	cmd := r.rdb.Set(context.Background(), key, string(data), expiration)

	return checkoutSession.CheckoutID, cmd.Err()
}
