package core_service

import (
	"backend/internal/db"
	"backend/internal/runtime_errors"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
)


func UploadCsvService(file multipart.File, filename string, uploadedBy int) error {
	reader := csv.NewReader(file)

	// Insert into csv_table
	var fileID int64
	err := db.DB.QueryRow(`
		INSERT INTO csv_table (file_name, uploaded_by) 
		VALUES ($1, $2) 
		RETURNING id
	`, filename, uploadedBy).Scan(&fileID)
	if err != nil {
		return &runtime_errors.InternalServerError{
			Message: fmt.Sprintf("failed to insert file record: %v", err),
		}
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return &runtime_errors.InternalServerError{
			Message: fmt.Sprintf("failed to begin transaction: %v", err),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO csv_rows (csv_file_id, position, input_text) 
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return &runtime_errors.InternalServerError{
			Message: fmt.Sprintf("failed to prepare statement: %v", err),
		}
	}
	defer stmt.Close()

	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		return &runtime_errors.BadRequestError{
			Message: fmt.Sprintf("error reading CSV header: %v", err),
		}
	}

	pos := 10.0
	step := 10.0
	rowCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &runtime_errors.BadRequestError{
				Message: fmt.Sprintf("invalid CSV: %v", err),
			}
		}

		inputText := ""
		if len(record) > 2 { // ensure at least 3 columns
			inputText = record[2] // take the third column for input_text
		}

		_, err = stmt.Exec(fileID, pos, inputText)
		if err != nil {
			return &runtime_errors.InternalServerError{
				Message: fmt.Sprintf("failed to insert row: %v", err),
			}
		}

		pos += step
		rowCount++
	}

	fmt.Printf("Uploaded %d rows for file %d\n", rowCount, fileID)
	return nil
}
