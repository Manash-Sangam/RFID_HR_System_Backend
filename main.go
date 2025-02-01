package main

import (
	"fmt"
	"net"
	"net/http"
	"rfid_backend/db"
	"rfid_backend/handlers"
	"rfid_backend/websocket"
)

func main() {
	// Initialize the database connection
	db.InitDB("root:manash@tcp(127.0.0.1:3306)/rfid_backend")

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
		fmt.Println("Error getting local IP address:", err)
	} else {
		fmt.Printf("Starting server on http://%s:9191\n", ip)
	}

	http.ListenAndServe(":9191", nil)
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
