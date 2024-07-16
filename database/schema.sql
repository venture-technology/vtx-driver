-- Table drivers
DROP TABLE drivers;

CREATE TABLE IF NOT EXISTS drivers (
    id SERIAL,
    name VARCHAR(100) NOT NULL,
    cpf VARCHAR(14) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    cnh VARCHAR(20) PRIMARY KEY NOT NULL,
    qrcode VARCHAR(100) NOT NULL,
    amount BIGINT NOT NULL,
    street VARCHAR(100) NOT NULL,
    number VARCHAR(10) NOT NULL,
    complement VARCHAR(10),
    zip VARCHAR(8) NOT NULL
); 