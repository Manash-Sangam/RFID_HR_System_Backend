package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os" // Ensure this import statement is correct
	"rfid_backend/db"
	"rfid_backend/handlers"
	"rfid_backend/websocket"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database connection
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	if dbConnectionString == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable is required")
	}
	db.InitDB(dbConnectionString)

	// Initialize WebSocket hub
	websocket.HubInstance = websocket.NewHub()
	go websocket.HubInstance.Run()

	http.HandleFunc("/verify", handlers.VerifyHandler)
	http.HandleFunc("/add_employee", handlers.AddEmployeeHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(websocket.HubInstance, w, r)
	})

	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Error getting local IP address: %v", err)
	}
	fmt.Printf("Starting server on http://%s:9191\n", ip)

	log.Fatal(http.ListenAndServe(":9191", nil))
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to get local IP address: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
