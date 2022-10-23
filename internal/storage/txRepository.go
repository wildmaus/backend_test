package storage

import (
	"backend_test/internal/model"
	"context"
)

type TxRepository interface {
	FindOne(ctx context.Context, id int32) (model.TransactionDto, error)
	GetReport(ctx context.Context, month int32, year int32) ([]model.Report, error)
	FindTxByUser(ctx context.Context, id int32, sort_param string) ([]model.Transaction, error)
	FindByDetails(ctx context.Context, userId int32, detailsId int32) (int32, error)
	Transfer(ctx context.Context, from *model.User, to *model.User, tx *model.TransactionDto, create bool) error
}
