package database

import (
	"context"
	"database/sql"
	"fmt"
)

func (d *Database) Exec(query string, args ...any) (sql.Result, error) {
	return d.Db.Exec(query, args...)
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.Db.ExecContext(ctx, query, args...)
}

func (d *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return d.Db.Query(query, args...)
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.Db.QueryContext(ctx, query, args...)
}

func (d *Database) QueryRow(query string, args ...any) *sql.Row {
	return d.Db.QueryRow(query, args...)
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.Db.QueryRowContext(ctx, query, args...)
}

func (d *Database) Begin() (*Transaction, error) {
	tx, err := d.Db.Begin()
	if err != nil {
		return nil, err
	}

	return &Transaction{
		tx: tx,
		db: d,
	}, nil
}

func (d *Database) Transaction(transactFunc func(tx *Transaction) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	err = transactFunc(tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("transaction error: %s, rollback error: %s", err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func (d *Database) getTableColumns() map[string][]string {
	return d.tableColumns
}
