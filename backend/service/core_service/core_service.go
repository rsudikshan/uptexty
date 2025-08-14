package core_service

import (
	"backend/internal/db"
	"backend/internal/runtime_errors"
	"backend/payloads/response"
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

func GetUploadedFilesService(uploadedBy int) ([]response.GetFilesResponse,error){
	var err error
	var responseList []response.GetFilesResponse 


	queryStr := "SELECT id,file_name,uploaded_at FROM csv_table WHERE uploaded_by = $1 "

	resultSet,err := db.DB.Query(queryStr,uploadedBy)

	if err!=nil{
		return nil,&runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	for resultSet.Next() {
		var responseVar response.GetFilesResponse
		err = resultSet.Scan(&responseVar.ID,&responseVar.Filename,&responseVar.UploadedAt)
		if err!=nil {
			return nil,&runtime_errors.InternalServerError{
				Message: err.Error(),
			};
		}

		responseList = append(responseList, responseVar)
	}


	return responseList,nil

}

func GetRows(userId int,fileId int)([]response.GetRowsResponse,error){
	var err error
	var responseList []response.GetRowsResponse

	//need to check authenticated user specific files

	queryStr := "SELECT id,position,input_text FROM csv_rows WHERE csv_file_id = $1 "

	resultSet,err := db.DB.Query(queryStr,fileId)

	if err!=nil{
		return nil,&runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	for resultSet.Next() {
		var responseVar response.GetRowsResponse
		err = resultSet.Scan(&responseVar.Id,&responseVar.Position,&responseVar.InputText)
		if err!=nil {
			return nil,&runtime_errors.InternalServerError{
				Message: err.Error(),
			};
		}

		responseList = append(responseList, responseVar)
	}


	return responseList,nil
}

func CreateRowService(userID, fileID int, position float64, inputText string) (*response.GetRowsResponse, error) {
	var fileExists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM csv_table WHERE id = $1 AND uploaded_by = $2
		)`, fileID, userID).Scan(&fileExists)
	if err != nil {
		return nil, &runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}
	if !fileExists {
		return nil, &runtime_errors.BadRequestError{
			Message: "File not found or access denied",
		}
	}

	var newRow response.GetRowsResponse
	err = db.DB.QueryRow(`
		INSERT INTO csv_rows (csv_file_id, position, input_text, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, position, input_text`,
		fileID, position, inputText,
	).Scan(&newRow.Id, &newRow.Position, &newRow.InputText)
	if err != nil {
		return nil, &runtime_errors.InternalServerError{
			Message: err.Error(),
		}
	}

	return &newRow, nil
}

func UpdateRowService(userID, fileID, rowID int, position float64, inputText string) (*response.GetRowsResponse, error) {
	var rowExists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM csv_rows r
			JOIN csv_table f ON r.csv_file_id = f.id
			WHERE r.id = $1 AND r.csv_file_id = $2 AND f.uploaded_by = $3
		)`, rowID, fileID, userID).Scan(&rowExists)
	if err != nil {
		return nil, &runtime_errors.InternalServerError{Message: err.Error()}
	}
	if !rowExists {
		return nil, &runtime_errors.BadRequestError{Message: "Row doesnt exist"}
	}

	var updatedRow response.GetRowsResponse
	err = db.DB.QueryRow(`
		UPDATE csv_rows
		SET position = $1, input_text = $2
		WHERE id = $3 AND csv_file_id = $4
		RETURNING id, position, input_text`,
		position, inputText, rowID, fileID,
	).Scan(&updatedRow.Id, &updatedRow.Position, &updatedRow.InputText)
	if err != nil {
		return nil, &runtime_errors.InternalServerError{Message: err.Error()}
	}

	return &updatedRow, nil
}

func DeleteRowService(userID, fileID, rowID int) error {
	result, err := db.DB.Exec(`
		DELETE FROM csv_rows
		WHERE id = $1 AND csv_file_id = $2 AND EXISTS(
			SELECT 1 FROM csv_table WHERE id = $2 AND uploaded_by = $3
		)`, rowID, fileID, userID)
	if err != nil {
		return &runtime_errors.InternalServerError{Message: err.Error()}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &runtime_errors.InternalServerError{Message: err.Error()}
	}
	if rowsAffected == 0 {
		return &runtime_errors.BadRequestError{Message: "Row not found"}
	}

	return nil
}


//
// // Enhanced GetRows with optional pagination and search
// func GetRowsWithPagination(userID, fileID int, limit, offset int, searchTerm string) ([]*response.GetRowsResponse, int, error) {
// 	// Verify file ownership
// 	var fileExists bool
// 	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM files WHERE id = $1 AND user_id = $2)", fileID, userID).Scan(&fileExists)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("error checking file ownership: %v", err)
// 	}
// 	if !fileExists {
// 		return nil, 0, fmt.Errorf("file not found or access denied")
// 	}

// 	// Build query with search
// 	baseQuery := `
// 		FROM csv_rows 
// 		WHERE file_id = $1`
	
// 	args := []interface{}{fileID}
// 	argCount := 1

// 	if searchTerm != "" {
// 		argCount++
// 		baseQuery += fmt.Sprintf(" AND input_text ILIKE $%d", argCount)
// 		args = append(args, "%"+searchTerm+"%")
// 	}

// 	// Get total count
// 	var totalCount int
// 	countQuery := "SELECT COUNT(*) " + baseQuery
// 	err = db.DB.QueryRow(countQuery, args...).Scan(&totalCount)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("error getting total count: %v", err)
// 	}

// 	// Get paginated results
// 	dataQuery := `
// 		SELECT id, position, input_text ` + baseQuery + `
// 		ORDER BY position ASC 
// 		LIMIT $%d OFFSET $%d`
	
// 	argCount++
// 	args = append(args, limit)
// 	argCount++
// 	args = append(args, offset)
	
// 	dataQuery = fmt.Sprintf(dataQuery, argCount-1, argCount)
	
// 	rows, err := db.DB.Query(dataQuery, args...)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("error fetching rows: %v", err)
// 	}
// 	defer rows.Close()

// 	var result []*response.GetRowsResponse
// 	for rows.Next() {
// 		var row response.GetRowsResponse
// 		err := rows.Scan(&row.Id, &row.Position, &row.InputText)
// 		if err != nil {
// 			return nil, 0, fmt.Errorf("error scanning row: %v", err)
// 		}
// 		result = append(result, &row)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, 0, fmt.Errorf("error iterating rows: %v", err)
// 	}

// 	return result, totalCount, nil
// }