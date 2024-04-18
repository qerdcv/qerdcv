package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/qerdcv/qerdcv/pkg/domain"
	"github.com/qerdcv/qerdcv/pkg/sqlutils"
)

type UserRepo struct {
	db sqlutils.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) WithTX(ctx context.Context, cmd func(r *UserRepo) error) (err error) {
	txDB, ok := r.db.(sqlutils.TxDB)
	if !ok {
		return errors.New("begin method not implemented")
	}

	var tx *sql.Tx
	tx, err = txDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	defer func() {
		if err = errors.Join(err, tx.Rollback()); err != nil && errors.Is(err, sql.ErrTxDone) {
			err = nil
		}
	}()

	if err = cmd(&UserRepo{db: tx}); err != nil {
		return fmt.Errorf("exec cmd: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}

	return nil
}

func (r *UserRepo) CreateUser(ctx context.Context, username, passwordHash string) error {
	query := `INSERT INTO users(username, password_hash) VALUES ($1, $2)`
	if _, err := r.db.ExecContext(ctx, query, username, passwordHash); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == sqlutils.PQErrCodeUniqueConstraint {
			return ErrUniqueConstraint
		}

		return fmt.Errorf("pool exec: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `SELECT
		id,
		username,
		password_hash,
		created_at,
		updated_at
		FROM users WHERE username=$1`

	user, err := r.scanUser(r.db.QueryRowContext(ctx, query, username))
	if err != nil {
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) CreateUserSession(
	ctx context.Context,
	user domain.User,
	expiresAt time.Time,
) (domain.UserSession, error) {
	query := `INSERT INTO user_sessions(user_id, expires_at) VALUES ($1, $2) RETURNING id, user_id, created_at, expires_at`

	userSession, err := r.scanUserSession(r.db.QueryRowContext(ctx, query, user.ID, expiresAt))
	if err != nil {
		return domain.UserSession{}, fmt.Errorf("scan user session: %w", err)
	}

	return userSession, nil
}

func (r *UserRepo) GetUserSession(ctx context.Context, sID, uID int) (domain.UserSession, error) {
	query := `SELECT id, user_id, created_at, expires_at from user_sessions WHERE id=$1 AND user_id=$2`

	session, err := r.scanUserSession(r.db.QueryRowContext(ctx, query, sID, uID))
	if err != nil {
		return domain.UserSession{}, fmt.Errorf("scan user session: %w", err)
	}

	return session, nil
}

func (r *UserRepo) DeleteUserSession(ctx context.Context, sID, uID int) error {
	query := `DELETE FROM user_sessions WHERE id=$1 AND user_id=$2`

	res, err := r.db.ExecContext(ctx, query, sID, uID)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UserRepo) scanUser(row *sql.Row) (domain.User, error) {
	var u domain.User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrNotFound
		}

		return domain.User{}, fmt.Errorf("row scan: %w", err)
	}

	return u, nil
}

func (r *UserRepo) scanUserSession(row *sql.Row) (domain.UserSession, error) {
	var session domain.UserSession
	if err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.UserSession{}, ErrNotFound
		}

		return domain.UserSession{}, fmt.Errorf("row scan: %w", err)
	}

	return session, nil
}
