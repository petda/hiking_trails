package models

import (
	"database/sql"
	"fmt"
)

type DatabaseHandle interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SQLQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type SQLExecer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SQLTransaction interface {
	Commit() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Rollback() error
}

type Model interface {
	Save(execer SQLExecer) error
	Load(queryer SQLQueryer) error
	Update(execer SQLExecer) error
	// Delete(execer SQLExecer) error
	RequireTransaction() bool
	DatabaseTable() string
	Type() string
}

func Save(model Model, db *sql.DB) error {
	if model.RequireTransaction() {
		return saveUsingTransaction(model, db)
	}

	return model.Save(db)
}

func saveUsingTransaction(model Model, db *sql.DB) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	err = model.Save(transaction)

	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit()
}

func Load(model Model, db *sql.DB) error {
	if model.RequireTransaction() {
		return loadUsingTransaction(model, db)
	}

	return model.Load(db)
}

func loadUsingTransaction(model Model, db *sql.DB) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	err = model.Load(transaction)

	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit()
}

func Update(model Model, db *sql.DB) error {
	return model.Update(db)
}

func Delete(model Model, id int64, execer SQLExecer) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", model.DatabaseTable())
	result, err := execer.Exec(query, id)

	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to delete %s with id %d", model.Type(), id), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to get rows affected when deleting %s with id %d",
			model.Type(), id), err)
	}

	if rowsAffected == 0 {
		return NewAPIError(404, fmt.Sprintf("No %s with id %d exist", model.Type(), id), nil)
	}

	return nil
}
