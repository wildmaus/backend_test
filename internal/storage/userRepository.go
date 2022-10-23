package storage

import (
	"backend_test/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindOne(ctx context.Context, id int32) (model.User, error)
	Update(ctx context.Context, user *model.User) error
	CreateWithTx(ctx context.Context, user *model.User, tx *model.TransactionDto, create bool) error
}
