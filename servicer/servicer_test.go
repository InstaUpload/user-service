package servicer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	u "github.com/instaUpload/user-service/utils"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var serv Servicer

func CreateTestServicer(ctx context.Context) Servicer {

	if serv == nil {
		serv = NewService(ctx)
	}
	return serv
}

func setupDatabaseContainer(ctx context.Context) *postgres.PostgresContainer {
	user := u.GetEnvAsString("DB_USER", "postgres")
	password := u.GetEnvAsString("DB_PASSWORD", "password")
	name := u.GetEnvAsString("DB_NAME", "userdb")
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(name),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}
	p, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		panic(err)
	}
	os.Setenv("DB_PORT", p.Port())
	// Wait for the database to be ready

	return pgContainer
}

func TestMain(m *testing.M) {
	// This is a placeholder main function.
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("No .env file found", err)
		return
	}
	migrationFolder := filepath.Join("../", "migrations")
	fmt.Println("Migration folder set to:", migrationFolder)
	ctx := context.Background()

	pgContainer := setupDatabaseContainer(ctx)
	defer pgContainer.Terminate(ctx)

	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}
	fmt.Println("Postgres connection string:", connectionString)

	migration, err := migrate.New("file://"+migrationFolder, connectionString)
	if err != nil {
		panic(err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	fmt.Println("Database migrations applied successfully")

	fmt.Println("Test application ran successfully")
	_ = CreateTestServicer(ctx)
	exitCode := m.Run()
	time.Sleep(10 * time.Minute)
	os.Exit(exitCode)

}
