package psql

import (
	"backend_test/internal/model"
	"backend_test/pkg/client/postgressql"
	"context"
	"log"
)

type DetailsRepository struct {
	client postgressql.Client
}

func (r *DetailsRepository) FindOne(ctx context.Context, orderId int32, serviceId int32) (model.Details, error) {
	q := `
		SELECT id, status
		FROM details
		WHERE orderId = $1 and serviceId = $2
	`
	var dtl model.Details
	err := r.client.QueryRow(ctx, q, orderId, serviceId).Scan(&dtl.Id, &dtl.Status)
	if err != nil {
		log.Println(err)
		return model.Details{}, err
	}
	return dtl, nil
}

func (r *DetailsRepository) Update(ctx context.Context, details *model.Details) error {
	q := `
		UPDATE details
		SET orderId = $2, serviceId = $3, status = $4
		WHERE id = $1
	`
	_, err := r.client.Exec(ctx, q, details.Id, details.OrderId, details.ServiceId, details.Status)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *DetailsRepository) Reserve(ctx context.Context, user *model.User, tx *model.TransactionDto, detail *model.Details) error {
	q := `
		WITH upd AS
		(UPDATE users
		SET balance = $2, reserved = $3
		WHERE id = $1),
		crt AS
		(INSERT INTO details (orderid, serviceid, status)
		VALUES ($4, $5, $6) 
		RETURNING id)
		INSERT INTO transactions (fromid, toid, amount, date, type, detailsid)
		VALUES ($7, $8, $9, $10, $11, (SELECT id FROM crt))
	`
	_, err := r.client.Exec(ctx, q,
		user.Id, user.Balance, user.Reserved,
		detail.OrderId, detail.ServiceId, detail.Status,
		tx.FromId, tx.ToId, tx.Amount, tx.Date, tx.Type)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *DetailsRepository) SolveReserve(ctx context.Context, user *model.User, tx *model.TransactionDto) error {
	q := `
		WITH upd AS
		(UPDATE users
		SET balance = $2, reserved = $3
		WHERE id = $1),
		dtl AS
		(UPDATE details
		SET status = true
		WHERE id = $4)
		INSERT INTO transactions (fromid, toid, amount, date, type, detailsid)
		VALUES ($5, $6, $7, $8, $9, $4)
	`
	_, err := r.client.Exec(ctx, q,
		user.Id, user.Balance, user.Reserved, tx.DetailsId,
		tx.FromId, tx.ToId, tx.Amount, tx.Date, tx.Type)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
