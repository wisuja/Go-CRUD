package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Credentials struct {
	Username, Password, Host, Port, Database string
}

func GetConnection(credentials Credentials) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", credentials.Username, credentials.Password, credentials.Host, credentials.Port, credentials.Database)
	db, _ := sql.Open("mysql", dataSourceName)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}

func GetDefaultDatabaseConfig(path string) (Credentials, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return Credentials{}, err
	}

	username := viper.Get("DB_USERNAME").(string)
	password := viper.Get("DB_PASSWORD").(string)
	host := viper.Get("DB_HOST").(string)
	port := viper.Get("DB_PORT").(string)
	database := viper.Get("DB_NAME").(string)

	return Credentials{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}, nil
}
