package storage

import (
	"backend_test/internal/model"
	"context"
)

type DetailsRepository interface {
	FindOne(ctx context.Context, orderId int32, serviceId int32) (model.Details, error)
	Update(ctx context.Context, details *model.Details) error
	Reserve(ctx context.Context, user *model.User, tx *model.TransactionDto, detail *model.Details) error
	SolveReserve(ctx context.Context, user *model.User, tx *model.TransactionDto) error
}
