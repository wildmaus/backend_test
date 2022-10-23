package psql

import (
	"backend_test/internal/model"
	"backend_test/pkg/client/postgressql"
	"context"
	"fmt"
	"log"
)

type TxRepository struct {
	client postgressql.Client
}

func (r *TxRepository) FindOne(ctx context.Context, id int32) (model.TransactionDto, error) {
	q := `
		SELECT id, fromId, toId, amount, date, type, detailsId 
		FROM transactions 
		WHERE id = $1
	`
	var tx model.TransactionDto
	err := r.client.QueryRow(ctx, q, id).Scan(&tx.Id, &tx.FromId, &tx.ToId, &tx.Amount, &tx.Date, &tx.Type, &tx.DetailsId)
	if err != nil {
		log.Println(err)
		return model.TransactionDto{}, err
	}
	return tx, nil
}

func (r *TxRepository) FindTxByUser(ctx context.Context, id int32) ([]model.Transaction, error) {
	q := `
		WITH tx AS
		(SELECT * FROM transactions
		WHERE type != 4 AND (fromid = $1 or toid = $1))
		SELECT tx.id, fromId, toId, amount, date, type, orderid, serviceid
		FROM tx
		LEFT JOIN details ON detailsid = details.id
	`
	rows, err := r.client.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	txs := make([]model.Transaction, 0)
	for rows.Next() {
		var tx model.Transaction
		err = rows.Scan(&tx.Id, &tx.FromId, &tx.ToId, &tx.Amount, &tx.Date, &tx.Type, &tx.OrderId, &tx.ServiceId)
		txs = append(txs, tx)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *TxRepository) FindByDetails(ctx context.Context, userId int32, detailsId int32) (int32, error) {
	q := `
		SELECT amount FROM transactions
		WHERE fromId = $1 AND detailsId = $2 AND type = 2
	`
	var amount int32
	err := r.client.QueryRow(ctx, q, userId, detailsId).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (r *TxRepository) Transfer(ctx context.Context, from *model.User, to *model.User, tx *model.TransactionDto, create bool) error {
	q := `
		WITH updto AS
	`
	if create {
		q += fmt.Sprintf("(%s)", CrtUser)
	} else {
		q += fmt.Sprintf("(%s)", UpdUser)
	}
	q += `
		, updfrom AS
		(UPDATE users 
		SET balance = $5, reserved = $6
		WHERE id = $4)
		INSERT INTO transactions (fromId, toId, amount, date, type, detailsId) 
		VALUES ($7, $8, $9, $10, $11, $12) 
		RETURNING id
	`
	_, err := r.client.Exec(ctx, q, to.Id, to.Balance, to.Reserved, from.Id, from.Balance, from.Reserved, tx.FromId, tx.ToId, tx.Amount, tx.Date, tx.Type, tx.DetailsId)
	if err != nil {
		return err
	}
	return nil
}

func (r *TxRepository) GetReport(ctx context.Context, month int32, year int32) ([]model.Report, error) {
	q := `
		WITH tx AS 
		(SELECT amount, detailsid 
			FROM transactions 
			WHERE EXTRACT (MONTH FROM date) = $1 
			AND EXTRACT (YEAR FROM date) = $2 and type = 4) 
		SELECT serviceid, sum(amount) 
		FROM tx 
		JOIN details ON detailsid=details.id 
		GROUP BY serviceId
	`
	rows, err := r.client.Query(ctx, q, month, year)
	if err != nil {
		return nil, err
	}
	report := make([]model.Report, 0)
	for rows.Next() {
		var rptRow model.Report
		err = rows.Scan(&rptRow.ServiseId, &rptRow.Amount)
		report = append(report, rptRow)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return report, nil
}
