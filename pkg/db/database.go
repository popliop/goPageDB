package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func CreateTable() {
	fmt.Println("s")
}

func loadDBConfig() (string, error) {
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" || sslmode == "" {
		return "", fmt.Errorf("missing required database environment variables")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode), nil
}

func NewPostgresStorage() (*sql.DB, error) {
	connStr, err := loadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load DB config: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}
	return db, nil
}

func CreateAccountTable(db *sql.DB) (sql.Result, error) {
	query := `CREATE TABLE IF NOT EXISTS export_shipments (
    SpedNr INTEGER PRIMARY KEY,
    AWB VARCHAR(20),
    RegDate DATE,
    CreatedDate TIMESTAMP,
    ArrivalScan TIMESTAMP,
    GTW VARCHAR(10),
    ShipperName VARCHAR(255),
    LastSign VARCHAR(50),
    ProductCode VARCHAR(10),
    LineItems INTEGER,
    HoldCode VARCHAR(50),
    HoldCodeDate TIMESTAMP,
    TullStatus VARCHAR(10),
    TullStatusDT TIMESTAMP,
    ControllCheck BOOLEAN,
    ControllDate TIMESTAMP,
    BPOCheck BOOLEAN,
    BPODate TIMESTAMP,
    ErrorCheck BOOLEAN,
    ErrorDate TIMESTAMP,
    Image VARCHAR(255),
    ImageDate TIMESTAMP,
	Team VARCHAR(10)
	)`
	result, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("problem creating DB table %w", err)
	}

	return result, nil
}
