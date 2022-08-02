package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenConnection(t *testing.T) {
	t.Run("Invalid credentials", func(t *testing.T) {
		credentials := Credentials{}

		db, err := GetConnection(credentials)

		require.Nil(t, db)
		require.NotNil(t, err)
	})
	t.Run("Valid credentials", func(t *testing.T) {
		credentials := Credentials{
			Username: "root",
			Password: "",
			Host:     "localhost",
			Port:     "3306",
			Database: "db_learngo",
		}

		db, err := GetConnection(credentials)

		require.NotNil(t, db)
		require.Nil(t, err)
	})
}

func TestGetDefaultDatabaseConfig(t *testing.T) {
	t.Run("Invalid path", func(t *testing.T) {
		credentials, err := GetDefaultDatabaseConfig("randompathsomewhere")

		require.NotNil(t, err)
		require.Equal(t, Credentials{}, credentials)
	})

	t.Run("Correct path", func(t *testing.T) {
		credentials, err := GetDefaultDatabaseConfig("../.env")

		require.Nil(t, err)
		require.Equal(t, Credentials{
			Username: "root",
			Password: "",
			Host:     "localhost",
			Port:     "3306",
			Database: "db_learngo",
		}, credentials)
	})
}
