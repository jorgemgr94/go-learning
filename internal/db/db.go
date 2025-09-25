package db

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"go-learning/internal/db/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const (
	path = "sql/"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// Db represents the database layer with pre-loaded SQL queries
type Db struct {
	db               db
	selectUserQuery  string
	insertUserExec   string
	selectUsersQuery string
}

// DBConfig represents the configuration for the database
type DBConfig struct {
	Db db
}

func (c DBConfig) Validate() error {
	var errs []error
	if c.Db == nil {
		errs = append(errs, fmt.Errorf("Db cannot be nil"))
	}
	return errors.Join(errs...)
}

//go:generate mockgen -destination=./db_mock_test.go -package=db -source=./db.go
type db interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// NewDb creates a new database instance with pre-loaded SQL queries
func NewDb(config DBConfig) (*Db, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Db{
		db:               config.Db,
		selectUserQuery:  loadSqlQueries("select_user.sql"),
		insertUserExec:   loadSqlQueries("insert_user.sql"),
		selectUsersQuery: loadSqlQueries("select_users.sql"),
	}, nil
}

// CreateUser inserts a new user into the database
func (d *Db) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.CreateUserResponse, error) {
	if err := d.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	var userID string
	err := d.db.QueryRow(ctx, d.insertUserExec, req.ID, req.Name, req.Email).Scan(&userID)
	if err != nil {
		return nil, convertPgErrorToDbError(req.ID, err)
	}

	zap.L().Info("user created successfully", zap.String("userID", userID))

	return &models.CreateUserResponse{
		ID: userID,
	}, nil
}

// GetUser retrieves a user by ID
func (d *Db) GetUser(ctx context.Context, req *models.GetUserRequest) (*models.GetUserResponse, error) {
	if req.ID == "" {
		return nil, models.NewValidationError("id", "user ID is required")
	}

	var user models.User
	err := d.db.QueryRow(ctx, d.selectUserQuery, req.ID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, convertPgErrorToDbError(req.ID, err)
	}

	return &models.GetUserResponse{
		User: &user,
	}, nil
}

// ListUsers retrieves a list of users with pagination
func (d *Db) ListUsers(ctx context.Context, req *models.ListUsersRequest) (*models.ListUsersResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10 // default limit
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	rows, err := d.db.Query(ctx, d.selectUsersQuery, req.Limit, req.Offset)
	if err != nil {
		return nil, convertPgErrorToDbError("", err)
	}
	defer closeFunc(rows)

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, convertPgErrorToDbError("", err)
		}
		users = append(users, user)
	}

	return &models.ListUsersResponse{
		Users: users,
		Total: len(users),
	}, nil
}

// validateCreateUserRequest validates the create user request
func (d *Db) validateCreateUserRequest(req *models.CreateUserRequest) error {
	var errs []error

	if req.ID == "" {
		errs = append(errs, models.NewValidationError("id", "user ID is required"))
	}
	if req.Name == "" {
		errs = append(errs, models.NewValidationError("name", "user name is required"))
	}
	if req.Email == "" {
		errs = append(errs, models.NewValidationError("email", "user email is required"))
	}

	return errors.Join(errs...)
}

// loadSqlQueries loads SQL query from a file
func loadSqlQueries(sqlFile string) string {
	content, err := sqlFiles.ReadFile(path + sqlFile)
	if err != nil {
		panic(fmt.Errorf("error reading sql file %s: %v", sqlFile, err))
	}
	return string(content)
}

// convertPgErrorToDbError converts a pgx error to a custom database error
func convertPgErrorToDbError(id string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}

	var pgErr *pgconn.PgError
	isPgErr := errors.As(err, &pgErr)
	if !isPgErr {
		return models.NewDatabaseError("unexpected error type", err)
	}

	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return models.NewDatabaseError("constraint violation", err)
	}

	return models.NewDatabaseError("unknown database error", err)
}

// closeFunc safely closes database rows and logs any errors
func closeFunc(rows pgx.Rows) {
	rows.Close()
	if err := rows.Err(); err != nil {
		zap.L().Error("error encountered iterating on rows", zap.Error(err))
	}
}
