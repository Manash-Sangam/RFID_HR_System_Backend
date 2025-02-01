package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected!")
}

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AddEmployee(name string) (string, error) {
	rfidTag := GenerateRandomString(16)
	_, err := db.Exec("INSERT INTO employees (name, rfid_tag) VALUES (?, ?)", name, rfidTag)
	if err != nil {
		return "", err
	}
	return rfidTag, nil
}

func LogRFIDData(tagID, deviceID string) error {
	_, err := db.Exec("INSERT INTO rfid_logs (tag_id, device_id, timestamp) VALUES (?, ?, ?)", tagID, deviceID, time.Now())
	return err
}

func VerifyPerson(tagID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM employees WHERE rfid_tag = ?)", tagID).Scan(&exists)
	if err != nil {
		log.Println("Error verifying person:", err)
		return false
	}
	return exists
}
