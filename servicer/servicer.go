package servicer

import (
	"context"
	"log/slog"

	"github.com/instaUpload/user-service/database"
	"github.com/instaUpload/user-service/types"
	"github.com/jackc/pgx/v5"
)

type Servicer interface {
	GetVersion() string
	CreateUser(context.Context, CreateUserInput) (CreateUserOutput, error)
	Close(ctx context.Context, cancel context.CancelFunc)
}

type Service struct {
	conn *pgx.Conn
	db   *database.Queries
}

func NewService(ctx context.Context) *Service {
	ctx, cancel := context.WithCancel(ctx)
	dbConfig := types.NewDatabaseConfig()
	slog.Info("connecting to database", slog.String("connectionString", dbConfig.ConnectionString()))
	conn, err := pgx.Connect(ctx, dbConfig.ConnectionString())
	if err != nil {
		slog.Error("failed to connect to database", "slog", err)
		cancel()
	}
	serv := &Service{
		db:   database.New(conn),
		conn: conn,
	}
	go serv.Close(ctx, cancel)
	return serv
}

// Implement servicer methods for service here
func (s *Service) GetVersion() string {
	return "v0.0.1"
}

func (s *Service) Close(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	<-ctx.Done()
	slog.Warn("server shoutdown...")
	s.conn.Close(ctx)
	slog.Info("database connection closed")
	slog.Info("server stopped")
}
