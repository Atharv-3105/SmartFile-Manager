package storage

import(
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	// "time"
)

//Function inserts a file record if it doesn't exist,
//OR updates metadata of the file if it already exists

func UpsertFile(tx *sql.Tx, path string) (int64, error ) {
	info , err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("[DB] stat file: %w", err)
	}


	filename := filepath.Base(path)
	extension := filepath.Ext(path)
	lastModified := info.ModTime()

	//Fetch existing fileID
	var fileID int64
	err = tx.QueryRow(
		`SELECT id FROM files WHERE path = ?`,
		path,	
	).Scan(&fileID)

	switch {
	//Case: When no data is present
	case err == sql.ErrNoRows:
		res, err := tx.Exec(
			`INSERT INTO files (path,filename,extension,last_modified) VALUE (?,?,?,?)`,
			path,filename,extension,lastModified,
		)

		if err != nil{
			return 0, fmt.Errorf("[DB] insert file: %w", err)
		}

		fileID, err = res.LastInsertId()
		if err != nil{
			return 0, fmt.Errorf("[DB] get inserted file id: %w", err)
		}
	
	case err != nil:
		return 0, fmt.Errorf("[DB] select file id: %w", err)
	
	//Default case will be when file already exists; Update its metadata
	default:
		_, err := tx.Exec(
			`UPDATE files
			SET filename = ?, extension = ?, last_modified = ?
			WHERE id = ?`,
			filename,extension,lastModified,fileID,
		)

		if err != nil{
			return 0, fmt.Errorf("[DB] update file: %w", err)
		}
	}

	return fileID, nil
}