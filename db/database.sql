CREATE DATABASE rfid_backend;

USE rfid_backend;

CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rfid_tag VARCHAR(16) UNIQUE NOT NULL
);

CREATE TABLE rfid_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tag_id VARCHAR(16) NOT NULL,
    device_id VARCHAR(50) NOT NULL,
    timestamp DATETIME NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES employees(rfid_tag)
);