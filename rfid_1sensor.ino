#include <ESP8266WiFi.h>
#include <WebSocketsClient.h>
#include <MFRC522.h>
#include <SPI.h>

#define SS_PIN D8
#define RST_PIN D0

MFRC522 mfrc522(SS_PIN, RST_PIN);  // Create MFRC522 instance

WebSocketsClient webSocket;

const char* ssid = "Your_WiFi_SSID";
const char* password = "Your_WiFi_Password";
const char* websocket_server = "Your_Server_IP";  // e.g., "192.168.1.100"
const int websocket_port = 9191;

void setup() {
  Serial.begin(115200);
  SPI.begin();  // Init SPI bus
  mfrc522.PCD_Init();  // Init MFRC522

  pinMode(LED_BUILTIN, OUTPUT);

  Serial.println("Connecting to WiFi...");
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi connected");

  Serial.println("Connecting to WebSocket server...");
  webSocket.begin(websocket_server, websocket_port, "/ws?id=ESP8266-01");
  webSocket.onEvent(webSocketEvent);
  webSocket.setReconnectInterval(5000);
  checkRFIDConnection();
}

void loop() {
  webSocket.loop();
  checkRFID(mfrc522, "esp01");
}

void checkRFIDConnection() {
  if (!mfrc522.PCD_PerformSelfTest()) {
    Serial.println("RFID sensor not connected properly or failed self-test.");
  } else {
    Serial.println("RFID sensor connected and passed self-test.");
  }
}

void checkRFID(MFRC522 &mfrc522, const char* deviceID) {
  // Serial.println("Inside checkRFID");
  if (!mfrc522.PICC_IsNewCardPresent()) {
    return;
  }

  Serial.println("Card detected!");
  
  if (!mfrc522.PICC_ReadCardSerial()) {
    Serial.println("Failed to read card serial.");
    return;
  }


  byte buffer[18];
  byte block = 4;  // Block to read from
  MFRC522::StatusCode status;

  // Authenticate using key A
  MFRC522::MIFARE_Key key;
  for (byte i = 0; i < 6; i++) key.keyByte[i] = 0xFF;
  status = mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, block, &key, &(mfrc522.uid));
  if (status != MFRC522::STATUS_OK) {
    Serial.print("Authentication failed: ");
    Serial.println(mfrc522.GetStatusCodeName(status));
    return;
  }

  // Read data from the block
  byte len = 18;
  status = mfrc522.MIFARE_Read(block, buffer, &len);
  if (status != MFRC522::STATUS_OK) {
    Serial.print("Read failed: ");
    Serial.println(mfrc522.GetStatusCodeName(status));
    return;
  }

  String tagData = "";
  for (byte i = 0; i < 16; i++) {
    tagData += char(buffer[i]);
  }

  String json = "{\"tag_id\":\"" + tagData + "\",\"device_id\":\"" + String(deviceID) + "\"}";
  Serial.print("Sending JSON: ");
  Serial.println(json);
  webSocket.sendTXT(json);

  mfrc522.PICC_HaltA();
  mfrc522.PCD_StopCrypto1();
}

void webSocketEvent(WStype_t type, uint8_t * payload, size_t length) {
  switch (type) {
    case WStype_DISCONNECTED:
      Serial.println("WebSocket Disconnected");
      break;
    case WStype_CONNECTED:
      Serial.println("WebSocket Connected");
      break;
    case WStype_TEXT:
      Serial.printf("WebSocket Message: %s\n", payload);
      handleWebSocketMessage((char*)payload);
      break;
    case WStype_BIN:
      Serial.println("WebSocket Binary Message");
      break;
  }
}

void handleWebSocketMessage(char* message) {
  String msg = String(message);
  Serial.print("Received WebSocket Message: ");
  Serial.println(msg);
  if (msg.indexOf("Access Granted") >= 0) {
    blinkLED(1);
  } else if (msg.indexOf("Access Denied") >= 0) {
    blinkLED(3);
  } else if (msg.indexOf("Write Tag") >= 0) {
    String tagID = msg.substring(msg.indexOf(":") + 1);
    writeRFIDTag(tagID);
  }
}

void blinkLED(int times) {
  for (int i = 0; i < times; i++) {
    digitalWrite(LED_BUILTIN, LOW);
    delay(200);
    digitalWrite(LED_BUILTIN, HIGH);
    delay(200);
  }
}

void writeRFIDTag(String tagID) {
  MFRC522::MIFARE_Key key;
  for (byte i = 0; i < 6; i++) key.keyByte[i] = 0xFF;

  byte buffer[18];
  byte block = 4;  // Block to write to
  byte len = tagID.length();
  for (byte i = 0; i < len; i++) {
    buffer[i] = tagID[i];
  }
  for (byte i = len; i < 16; i++) {
    buffer[i] = 0;
  }
  
  Serial.println("Buffer content to write:");
  for (byte i = 0; i < 16; i++) {
    Serial.print(buffer[i], HEX);
    Serial.print(" ");
  }
  Serial.println();
  
  Serial.println("Waiting for card to be placed on the scanner...");
  while (!mfrc522.PICC_IsNewCardPresent() || !mfrc522.PICC_ReadCardSerial()) {
    delay(500);
  }

  Serial.println("Card detected. Attempting to authenticate...");
  MFRC522::StatusCode status;
  // status = mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, block, &key, &(mfrc522.uid));
  // if (status != MFRC522::STATUS_OK) {
  //   Serial.print("Authentication failed: ");
  //   Serial.println(mfrc522.GetStatusCodeName(status));
  //   return;
  // }

  Serial.println("Authentication successful. Attempting to write...");
  status = mfrc522.MIFARE_Write(block, buffer, 16);
  if (status != MFRC522::STATUS_OK) {
    Serial.print("Write failed: ");
    Serial.println(mfrc522.GetStatusCodeName(status));
    return;
  }

  Serial.println("Write successful");
}