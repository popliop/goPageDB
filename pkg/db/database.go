package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

const upsertQuery = `
INSERT INTO export_shipments (
    SpedNr, AWB, RegDate, CreatedDate, ArrivalScan, GTW, ShipperName,
    LastSign, ProductCode, LineItems, HoldCode, HoldCodeDate, TullStatus,
    TullStatusDT, ControllCheck, ControllDate, BPOCheck, BPODate, ErrorCheck,
    ErrorDate, Image, ImageDate
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13,
    $14, $15, $16, $17, $18, $19,
    $20, $21, $22
)
ON CONFLICT (SpedNr) DO UPDATE SET
    AWB = EXCLUDED.AWB,
    RegDate = EXCLUDED.RegDate,
    CreatedDate = EXCLUDED.CreatedDate,
    ArrivalScan = EXCLUDED.ArrivalScan,
    GTW = EXCLUDED.GTW,
    ShipperName = EXCLUDED.ShipperName,
    LastSign = EXCLUDED.LastSign,
    ProductCode = EXCLUDED.ProductCode,
    LineItems = EXCLUDED.LineItems,
    HoldCode = EXCLUDED.HoldCode,
    HoldCodeDate = EXCLUDED.HoldCodeDate,
    TullStatus = EXCLUDED.TullStatus,
    TullStatusDT = EXCLUDED.TullStatusDT,
    ControllCheck = EXCLUDED.ControllCheck,
    ControllDate = EXCLUDED.ControllDate,
    BPOCheck = EXCLUDED.BPOCheck,
    BPODate = EXCLUDED.BPODate,
    ErrorCheck = EXCLUDED.ErrorCheck,
    ErrorDate = EXCLUDED.ErrorDate,
    Image = EXCLUDED.Image,
    ImageDate = EXCLUDED.ImageDate;
`

// nullIfEmpty returns nil if the provided string is empty after ensuring valid UTF-8.
func nullIfEmpty(s string) interface{} {
	s = strings.ToValidUTF8(s, "")
	if s == "" {
		return nil
	}
	return s
}

// parseNullableInt converts a non-empty string to an integer, or returns nil.
func parseNullableInt(s string) (interface{}, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// parseRow converts a CSV row (slice of strings) into a slice of interface{}
// suitable for stmt.Exec. Row index is used in error messages.
func parseRow(i int, row []string) ([]interface{}, error) {
	if len(row) != 22 {
		return nil, fmt.Errorf("row %d does not have enough columns, got %d, expected 22", i, len(row))
	}

	// SpedNr is the primary key and must be a valid integer.
	spedNr, err := strconv.Atoi(row[0])
	if err != nil {
		return nil, fmt.Errorf("invalid SpedNr in row %d: %w", i, err)
	}

	awb := nullIfEmpty(row[1])
	regDate := nullIfEmpty(row[2])
	createdDate := nullIfEmpty(row[3])
	arrivalScan := nullIfEmpty(row[4])
	gtw := nullIfEmpty(row[5])
	shipperName := nullIfEmpty(row[6])
	lastSign := nullIfEmpty(row[7])
	productCode := nullIfEmpty(row[8])

	lineItems, err := parseNullableInt(row[9])
	if err != nil {
		return nil, fmt.Errorf("invalid LineItems in row %d: %w", i, err)
	}

	holdCode := nullIfEmpty(row[10])
	holdCodeDate := nullIfEmpty(row[11])
	tullStatus := nullIfEmpty(row[12])
	tullStatusDT := nullIfEmpty(row[13])

	// For booleans, no need to check if they are empty/null.
	controllCheck := strings.ToLower(row[14]) == "true"
	controllDate := nullIfEmpty(row[15])
	bpoCheck := strings.ToLower(row[16]) == "true"
	bpoDate := nullIfEmpty(row[17])
	errorCheck := strings.ToLower(row[18]) == "true"
	errorDate := nullIfEmpty(row[19])
	image := nullIfEmpty(row[20])
	imageDate := nullIfEmpty(row[21])

	return []interface{}{
		spedNr,
		awb,
		regDate,
		createdDate,
		arrivalScan,
		gtw,
		shipperName,
		lastSign,
		productCode,
		lineItems,
		holdCode,
		holdCodeDate,
		tullStatus,
		tullStatusDT,
		controllCheck,
		controllDate,
		bpoCheck,
		bpoDate,
		errorCheck,
		errorDate,
		image,
		imageDate,
	}, nil
}

// UpsertCSVData inserts or updates CSV data into the database.
func UpsertCSVData(db *sql.DB, data [][]string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Rollback if any error occurs.
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(upsertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Skip header row.
	for i, row := range data {
		if i == 0 {
			continue
		}
		params, err := parseRow(i, row)
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(params...); err != nil {
			return fmt.Errorf("failed to execute statement for row %d: %w", i, err)
		}
	}
	return tx.Commit()
}

// loadDBConfig loads database connection details from environment variables.
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

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode), nil
}

// NewPostgresStorage returns a new database connection.
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

// CreateAccountTable creates the export_shipments table if it doesn't exist.
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
		ImageDate TIMESTAMP
	)`
	result, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("problem creating DB table: %w", err)
	}
	return result, nil
}
