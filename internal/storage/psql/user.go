package psql

import (
	"backend_test/internal/model"
	"backend_test/pkg/client/postgressql"
	"context"
	"fmt"
)

const (
	CrtUser = `
		INSERT INTO users (id, balance, reserved) 
		VALUES ($1, $2, $3)
	`
	UpdUser = `
		UPDATE users 
		SET balance = $2, reserved = $3
		WHERE id = $1
	`
)

type UserRepository struct {
	client postgressql.Client
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.client.QueryRow(ctx, CrtUser, user.Id, user.Balance, user.Reserved).Scan(); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindOne(ctx context.Context, id int32) (model.User, error) {
	q := `
		SELECT id, balance, reserved 
		FROM users 
		WHERE id = $1
	`
	var user model.User
	err := r.client.QueryRow(ctx, q, id).Scan(&user.Id, &user.Balance, &user.Reserved)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.client.Exec(ctx, UpdUser, user.Id, user.Balance, user.Reserved)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateWithTx(ctx context.Context, user *model.User, tx *model.TransactionDto, create bool) error {
	q := `
		WITH sub AS
	`
	if create {
		q += fmt.Sprintf("(%s)", CrtUser)
	} else {
		q += fmt.Sprintf("(%s)", UpdUser)
	}
	q += `
		INSERT INTO transactions (fromId, toId, amount, date, type, detailsId) 
		VALUES ($4, $5, $6, $7, $8, $9) 
		RETURNING id
	`
	_, err := r.client.Exec(ctx, q, user.Id, user.Balance, user.Reserved, tx.FromId, tx.ToId, tx.Amount, tx.Date, tx.Type, tx.DetailsId)
	if err != nil {
		return err
	}
	return nil
}
