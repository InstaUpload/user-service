package types

import (
	"strconv"

	"fmt"

	u "github.com/instaUpload/user-service/utils"
)

type ApplicationConfig struct {
	Host string
	Port int
}

func NewApplicationConfig() *ApplicationConfig {
	host := u.GetEnvAsString("APP_HOST", "localhost")
	port := u.GetEnvAsInt("APP_PORT", 8001)
	return &ApplicationConfig{
		Host: host,
		Port: port,
	}

}

func (c *ApplicationConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, strconv.Itoa(c.Port))
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	Server   string
}

func NewDatabaseConfig() *DatabaseConfig {
	host := u.GetEnvAsString("DB_HOST", "localhost")
	port := u.GetEnvAsInt("DB_PORT", 5432)
	user := u.GetEnvAsString("DB_USER", "postgres")
	password := u.GetEnvAsString("DB_PASSWORD", "password")
	name := u.GetEnvAsString("DB_NAME", "userdb")
	sslmode := u.GetEnvAsString("DB_SSLMODE", "disable")
	server := u.GetEnvAsString("DB_SERVER", "postgresql")
	return &DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslmode,
		Server:   server,
	}
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", c.Server, c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

type TokenizerConfig struct {
	SecretKey       string
	ExpirationHours int
}

func NewTokenizerConfig() *TokenizerConfig {
	secretKey := u.GetEnvAsString("TOKEN_SECRET_KEY", "mysecretkey")
	expirationHours := u.GetEnvAsInt("TOKEN_EXPIRATION_HOURS", 72)
	return &TokenizerConfig{
		SecretKey:       secretKey,
		ExpirationHours: expirationHours,
	}
}
