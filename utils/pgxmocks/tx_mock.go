package pgxmocks

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

type TxMock struct {
	mock.Mock
}

// Begin starts a pseudo nested transaction.
func (tx *TxMock) Begin(ctx context.Context) (pgx.Tx, error) {
	panic("not impl")
}

// BeginFunc starts a pseudo nested transaction and executes f. If f does not return an err the pseudo nested
// transaction will be committed. If it does then it will be rolled back.
func (tx *TxMock) BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error) {
	panic("not impl")
}

// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
// transaction. Commit will return ErrTxClosed if the Tx is already closed, but is otherwise safe to call multiple
// times. If the commit fails with a rollback status (e.g. the transaction was already in a broken state) then
// ErrTxCommitRollback will be returned.
func (tx *TxMock) Commit(ctx context.Context) error {
	args := tx.Called(ctx)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}

// Rollback rolls back the transaction if this is a real transaction or rolls back to the savepoint if this is a
// pseudo nested transaction. Rollback will return ErrTxClosed if the Tx is already closed, but is otherwise safe to
// call multiple times. Hence, a defer tx.Rollback() is safe even if tx.Commit() will be called first in a non-error
// condition. Any other failure of a real transaction will result in the connection being closed.
func (tx *TxMock) Rollback(ctx context.Context) error {
	args := tx.Called(ctx)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}

func (tx *TxMock) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	panic("not impl")
}
func (tx *TxMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	panic("not impl")
}
func (tx *TxMock) LargeObjects() pgx.LargeObjects {
	panic("not impl")
}

func (tx *TxMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	panic("not impl")
}

func (tx *TxMock) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	args := tx.Called(ctx, sql, arguments)

	var r0 pgconn.CommandTag
	if args.Get(0) != nil {
		r0 = args.Get(0).(pgconn.CommandTag)
	}
	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}
	return r0, r1
}
func (tx *TxMock) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	panic("not impl")
}
func (tx *TxMock) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	panic("not impl")
}
func (tx *TxMock) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	panic("not impl")
}

// Conn returns the underlying *Conn that on which this transaction is executing.
func (tx *TxMock) Conn() *pgx.Conn {
	panic("not impl")
}
