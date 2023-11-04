package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DbName   string
	SSLMode  string
}

func NewConnection(c *Config) (*gorm.DB, error) {

	dsn :=
		fmt.Sprintf("host=%s port=%s password=%s  dbname=%s port=%v sslmode=%s",
			c.Host, c.Port, c.Password, c.DbName, c.Port, c.SSLMode,
		)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Postgres connection established")
	return db, nil
}
