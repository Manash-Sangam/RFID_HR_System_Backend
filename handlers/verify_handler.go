package handlers

import (
	"encoding/json"
	"net/http"
	"rfid_backend/db"
	"rfid_backend/models"
	"time"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var data models.RFIDData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	data.Timestamp = time.Now().Format(time.RFC3339)

	if db.VerifyPerson(data.TagID) {
		db.LogRFIDData(data.TagID, data.DeviceID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access Granted"))
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Access Denied"))
	}
}
