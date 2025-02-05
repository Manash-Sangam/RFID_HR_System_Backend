package models

type RFIDData struct {
    TagID     string `json:"tag_id"`
    DeviceID  string `json:"device_id"`
    Timestamp string `json:"timestamp"`
}
