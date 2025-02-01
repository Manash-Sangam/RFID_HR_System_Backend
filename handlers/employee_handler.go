package handlers

import (
	"encoding/json"
	"net/http"
	"rfid_backend/db"
	"rfid_backend/websocket"
)

func AddEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	rfidTag, err := db.AddEmployee(request.Name)
	if err != nil {
		http.Error(w, "Failed to add employee", http.StatusInternalServerError)
		return
	}

	response := struct {
		RFIDTag string `json:"rfid_tag"`
	}{
		RFIDTag: rfidTag,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Send RFID tag to the ESP device
	message := struct {
		Command string `json:"command"`
		RFIDTag string `json:"rfid_tag"`
	}{
		Command: "Write Tag",
		RFIDTag: rfidTag,
	}
	messageBytes, _ := json.Marshal(message)
	websocket.HubInstance.SendToClient("ESP8266-01", messageBytes)
}
