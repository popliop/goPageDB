package db

import (
	"fmt"
	"strconv"
	"time"
)

const timeLayout = "2006-01-02 15:04"

// ExportRows represents a slice of ExportRow
type ExportRows []ExportRow

// ExportRow maps to your database table.
type ExportRow struct {
	SpedNr        int        `db:"SpedNr"`        // INTEGER PRIMARY KEY
	AWB           string     `db:"AWB"`           // VARCHAR(20)
	RegDate       *time.Time `db:"RegDate"`       // DATE, nullable
	CreatedDate   *time.Time `db:"CreatedDate"`   // TIMESTAMP, nullable
	ArrivalScan   *time.Time `db:"ArrivalScan"`   // TIMESTAMP, nullable
	GTW           string     `db:"GTW"`           // VARCHAR(10)
	ShipperName   string     `db:"ShipperName"`   // VARCHAR(255)
	LastSign      string     `db:"LastSign"`      // VARCHAR(50)
	ProductCode   string     `db:"ProductCode"`   // VARCHAR(10)
	LineItems     int        `db:"LineItems"`     // INTEGER
	HoldCode      string     `db:"HoldCode"`      // VARCHAR(50)
	HoldCodeDate  *time.Time `db:"HoldCodeDate"`  // TIMESTAMP, nullable
	TullStatus    string     `db:"TullStatus"`    // VARCHAR(10)
	TullStatusDT  *time.Time `db:"TullStatusDT"`  // TIMESTAMP, nullable
	ControllCheck string     `db:"ControllCheck"` // VARCHAR(10)
	ControllDate  *time.Time `db:"ControllDate"`  // TIMESTAMP, nullable
	BPOCheck      string     `db:"BPOCheck"`      // VARCHAR(10)
	BPODate       *time.Time `db:"BPODate"`       // TIMESTAMP, nullable
	ErrorCheck    string     `db:"ErrorCheck"`    // VARCHAR(10)
	ErrorDate     *time.Time `db:"ErrorDate"`     // TIMESTAMP, nullable
	Image         string     `db:"Image"`         // VARCHAR(255)
	ImageDate     *time.Time `db:"ImageDate"`     // TIMESTAMP, nullable
}

// parseTime parses a time string using the predefined layout.
// If the input is empty, it returns nil.
func parseTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(timeLayout, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ConvertCsvToExportRows converts a slice of CSV rows into ExportRows.
func (exRows *ExportRows) ConvertCsvToExportRows(rows [][]string) error {
	for _, row := range rows {
		er, err := createExportRowFromCsvRow(row)
		if err != nil {
			return fmt.Errorf("CreateExportRowFromCsvRow: %w", err)
		}
		*exRows = append(*exRows, er)
	}
	return nil
}

// createExportRowFromCsvRow converts a single CSV row to an ExportRow.
// It assumes that the row has at least 23 columns.
func createExportRowFromCsvRow(row []string) (ExportRow, error) {
	if len(row) < 23 {
		return ExportRow{}, fmt.Errorf("not enough columns: expected at least 23, got %d", len(row))
	}

	var exportRow ExportRow

	// Simple string assignments
	exportRow.AWB = row[2]
	exportRow.GTW = row[6]
	exportRow.ShipperName = row[7]
	exportRow.LastSign = row[8]
	exportRow.ProductCode = row[9]
	exportRow.HoldCode = row[11]
	exportRow.TullStatus = row[13]
	exportRow.ControllCheck = row[15]
	exportRow.BPOCheck = row[17]
	exportRow.ErrorCheck = row[19]
	exportRow.Image = row[21]

	// Integer conversion
	spedNr, err := strconv.Atoi(row[1])
	if err != nil {
		return exportRow, fmt.Errorf("unable to convert SpedNr to int: %w", err)
	}
	exportRow.SpedNr = spedNr

	// Time parsing
	if regDate, err := parseTime(row[3]); err != nil {
		return exportRow, fmt.Errorf("unable to convert RegDate to time: %w", err)
	} else {
		exportRow.RegDate = regDate
	}

	if createDate, err := parseTime(row[4]); err != nil {
		return exportRow, fmt.Errorf("unable to convert CreatedDate to time: %w", err)
	} else {
		exportRow.CreatedDate = createDate
	}

	if arrivalScan, err := parseTime(row[5]); err != nil {
		return exportRow, fmt.Errorf("unable to convert ArrivalScan to time: %w", err)
	} else {
		exportRow.ArrivalScan = arrivalScan
	}

	lineItems, err := strconv.Atoi(row[10])
	if err != nil {
		return exportRow, fmt.Errorf("unable to convert LineItems to int: %w", err)
	}
	exportRow.LineItems = lineItems

	if holdCodeDate, err := parseTime(row[12]); err != nil {
		return exportRow, fmt.Errorf("unable to convert HoldCodeDate to time: %w", err)
	} else {
		exportRow.HoldCodeDate = holdCodeDate
	}

	if tullStatusDT, err := parseTime(row[14]); err != nil {
		return exportRow, fmt.Errorf("unable to convert TullStatusDT to time: %w", err)
	} else {
		exportRow.TullStatusDT = tullStatusDT
	}

	if controllDate, err := parseTime(row[16]); err != nil {
		return exportRow, fmt.Errorf("unable to convert ControllDate to time: %w", err)
	} else {
		exportRow.ControllDate = controllDate
	}

	if bpoDate, err := parseTime(row[18]); err != nil {
		return exportRow, fmt.Errorf("unable to convert BPODate to time: %w", err)
	} else {
		exportRow.BPODate = bpoDate
	}

	if errorDate, err := parseTime(row[20]); err != nil {
		return exportRow, fmt.Errorf("unable to convert ErrorDate to time: %w", err)
	} else {
		exportRow.ErrorDate = errorDate
	}

	if imageDate, err := parseTime(row[22]); err != nil {
		return exportRow, fmt.Errorf("unable to convert ImageDate to time: %w", err)
	} else {
		exportRow.ImageDate = imageDate
	}

	return exportRow, nil
}
