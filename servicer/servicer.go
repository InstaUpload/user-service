package servicer

import (
	"context"
	"log/slog"
	"fmt"

	"github.com/instaUpload/user-service/database"
	"github.com/instaUpload/user-service/types"
	u "github.com/instaUpload/user-service/utils"
	"github.com/jackc/pgx/v5"
)

type Servicer interface {
	GetVersion() string
	CreateUser(context.Context, CreateUserInput) (CreateUserOutput, error)
	Close(ctx context.Context, cancel context.CancelFunc)
}

type service struct {
	conn *pgx.Conn
	db *database.Queries
}

func NewService(ctx context.Context) *service {
	ctx, cancel := context.WithCancel(ctx)
	dbConfig := types.NewDatabaseConfig()
	slog.Info("connecting to database", slog.String("connectionString", dbConfig.ConnectionString()))
	conn, err := pgx.Connect(ctx, dbConfig.ConnectionString())
	if err != nil {
		slog.Error("failed to connect to database", "slog", err)
		cancel()
	}
	serv := &service{
		db: database.New(conn),
		conn: conn,
	}
	go serv.Close(ctx, cancel)
	return serv
}

// Implement servicer methods for service here
func (s *service) GetVersion() string {
	return "v0.0.1"
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	// Step 1: Check if user with same email already exists
	user, err := s.db.GetUserByEmail(ctx, input.Email)
	if err != nil && err != pgx.ErrNoRows {
	return CreateUserOutput{}, err
	}
	if user.Email != "" {
	return CreateUserOutput{}, fmt.Errorf("user with email %s already exists", input.Email)
	}
	// Step 2: Hash the password
	hashedPassword, err := u.HashPassword(input.Password)
	if err != nil {
		return CreateUserOutput{}, err
	}
	// Step 3: Insert user into database
	userToBeCreated := database.CreateUserParams{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: hashedPassword,
	}
	createdUser, err := s.db.CreateUser(ctx, userToBeCreated)
	if err != nil {
		return CreateUserOutput{}, err
	}
	// Step 4: Return user ID and success message
	output := CreateUserOutput{
		UserID: createdUser.ID,
	}
	return output, nil
}

func (s *service) Close(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	<-ctx.Done()
	slog.Warn("server shoutdown...")
	s.conn.Close(ctx)
	slog.Info("database connection closed")
	slog.Info("server stopped")
}
