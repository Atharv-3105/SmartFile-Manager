package storage

import(
	"database/sql"
	"fmt"
)

func InsertExtraction(
	tx *sql.Tx,
	fileID    int64,
	text    string,
	status  string,
) error {
	_, err := tx.Exec(
		`INSERT INTO extractions (file_id, text, status)
		VALUES (?, ?, ?)`,
		fileID,text,status,
	)
	if err != nil {
		return fmt.Errorf("[DB] insert extraction: %w", err)
	}

	return nil
}