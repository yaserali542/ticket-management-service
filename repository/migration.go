package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/yaserali542/ticket-management-service/models"
)

func DropBookingTable(db *gorm.DB) error {
	return db.DropTable(&models.Booking{}).Error
}
func MigrateBookingTable(db *gorm.DB) error {
	if !db.HasTable(&models.Booking{}) {
		sqlStr := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
		tx := db.Exec(sqlStr)

		if tx.Error != nil {
			fmt.Println(tx.Error)
		} else {
			fmt.Println("create extension ran successfully")
		}
	}

	if err := db.AutoMigrate(&models.Booking{}); err.Error != nil {
		return err.Error
	}
	if !db.HasTable(&models.Booking{}) {
		return errors.New("table doesn't exist")
	}
	return nil
}

func MigrateEventTable(db *gorm.DB) error {
	if !db.HasTable(&models.Booking{}) {
		sqlStr := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
		tx := db.Exec(sqlStr)

		if tx.Error != nil {
			fmt.Println(tx.Error)
		} else {
			fmt.Println("create extension ran successfully")
		}
	}

	if err := db.AutoMigrate(&models.Events{}); err.Error != nil {
		return err.Error
	}
	if !db.HasTable(&models.Events{}) {
		return errors.New("table doesn't exist")
	}
	return SeedEventRecords(db)
}

func SeedEventRecords(db *gorm.DB) error {
	if err := db.First(&models.Events{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		testObject := models.Events{
			EventName:  "Quay walk",
			EventImage: []byte(seedMilliniumBridge),
			StartDate:  time.Now(),
			EndDate:    time.Now().Add(4 * time.Hour),
			Venue:      "Quay Side Newcastle",
		}
		db.Create(&testObject)

		testObject = models.Events{
			EventName:  "USB building",
			EventImage: []byte(seedUSBBuilding),
			StartDate:  time.Now(),
			EndDate:    time.Now().Add(4 * time.Hour),
			Venue:      "Urban side Building",
		}
		db.Create(&testObject)
	}
	return nil
}

func MigrateBookingDataTable(db *gorm.DB) error {
	if !db.HasTable(&models.BookingData{}) {
		sqlStr := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
		tx := db.Exec(sqlStr)

		if tx.Error != nil {
			fmt.Println(tx.Error)
		} else {
			fmt.Println("create extension ran successfully")
		}
	}

	if err := db.AutoMigrate(&models.BookingData{}); err.Error != nil {
		return err.Error
	}
	if !db.HasTable(&models.BookingData{}) {
		return errors.New("table doesn't exist")
	}
	return nil
}
