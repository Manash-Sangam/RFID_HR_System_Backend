# RFID Backend System

This project is an RFID backend system that allows you to manage employees and log RFID scans using an ESP8266 platform. The backend is built with Go and uses a MySQL database to store employee and RFID log data. The ESP8266 communicates with the backend server via WebSocket.

## Project Structure

- `db/`: Contains database-related files.
- `handlers/`: Contains HTTP handlers for the backend server.
- `models/`: Contains data models.
- `websocket/`: Contains WebSocket-related files.
- `main.go`: Entry point for the backend server.
- `.env`: Environment variables file.
- `.gitignore`: Git ignore file.

## Prerequisites

- Go 1.16 or later
- MySQL
- ESP8266 platform
- Arduino IDE

## Setup

### Step 1: Clone the Repository

```sh
git clone https://github.com/yourusername/rfid-backend.git
cd rfid-backend
```

### Step 2: Set Up the Database

1. Create the database and tables by running the SQL script in `db/database.sql`.

```sh
mysql -u your_db_user -p < db/database.sql
```

2. Replace `your_db_user` with your MySQL username and enter your MySQL password when prompted.

### Step 3: Configure Environment Variables

1. Create a `.env` file in the root directory with the following content:

```plaintext
DB_CONNECTION_STRING=your_db_user:your_db_password@tcp(127.0.0.1:3306)/rfid_backend
```

2. Replace `your_db_user` and `your_db_password` with your MySQL username and password.

### Step 4: Install Dependencies

1. Install Go dependencies:

```sh
go mod tidy
```

2. Install the `godotenv` package:

```sh
go get github.com/joho/godotenv
```

### Step 5: Run the Backend Server

1. Run the backend server:

```sh
go run main.go
```

2. The server will start and print the local IP address. Note this IP address as you will need it for the ESP8266 configuration.

### Step 6: Configure and Upload Code to ESP8266

1. Open the `rfid_1sensor.ino` file in the Arduino IDE.
2. Edit the WiFi SSID and password in the `rfid_1sensor.ino` file:

```cpp
const char* ssid = "your_SSID";
const char* password = "your_PASSWORD";
```

3. Edit the WebSocket server IP address in the `rfid_1sensor.ino` file:

```cpp
const char* websocket_server = "your_server_ip";  // e.g., "192.168.1.100"
```

4. Upload the code to the ESP8266 platform.

### Usage

1. **Add an Employee**:
   - Use a tool like Thunder Client or Postman to send a POST request to add a new employee.
   - URL: `http://<your_ip>:9191/add_employee`
   - Method: `POST`
   - Body: JSON

     ```json
     {
       "name": "John Doe"
     }
     ```

   - The backend will respond with the generated RFID tag for the new employee and send a "Write Tag" command to the ESP8266 device to write the new RFID tag to a blank card.

2. **Verify RFID Tags**:
   - Scan an RFID card with the RFID reader connected to the ESP8266.
   - The ESP8266 will send the tag data to the backend via WebSocket.
   - The backend will respond with "Access Granted" or "Access Denied" based on the verification result.

### License

This project is licensed under the MIT License.