package repository

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

func InitializeDatabase(v *viper.Viper) (*gorm.DB, error) {

	server := v.Get("database.host")
	database := v.Get("database.database_name")
	port := v.Get("database.port")
	username := v.Get("database.username")
	password := v.Get("database.password")

	dbString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", username, password, server, port, database)

	db, err := gorm.Open("postgres", dbString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db connection failed: %s", err.Error())
		return nil, err
	}
	db.LogMode(true)
	return db, nil
}
