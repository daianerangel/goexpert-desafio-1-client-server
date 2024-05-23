CREATE DATABASE currency_exchange;
USE currency_exchange;

CREATE TABLE quotations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    quotation VARCHAR(255) NOT NULL
);