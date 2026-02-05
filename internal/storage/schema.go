package storage

import (
	"database/sql"
	"fmt"
)

//InitSchema contains all the Schema code related to required tables & indexes.

func InitSchema(db *sql.DB) error {
	stmts := []string {

		//======Files========
		`CREATE TABLE IF NOT EXISTS files(
			id			INTEGER PRIMARY KEY AUTOINCREMENT,
			path		TEXT NOT NULL UNIQUE,
			filename	TEXT NOT NULL,
			extension	TEXT NOT NULL,
			last_modified	DATETIME NOT NULL,
			created_at		DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		//======Extraction Table Schema======
		`CREATE TABLE IF NOT EXISTS extractions(
			id			INTEGER PRIMARY KEY AUTOINCREMENT,
			file_id		INTEGER NOT NULL,
			text		TEXT NOT NULL,
			status		TEXT NOT NULL,
			extracted_at	DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);`,

		//=======Embeddings Table Schema========
		`CREATE TABLE IF NOT EXISTS embeddings(
			id			INTEGER PRIMARY KEY AUTOINCREMENT,
			file_id		INTEGER NOT NULL,
			model		TEXT NOT NULL,
			vector		BLOB NOT NULL,
			dim			INTEGER NOT NULL,
			created_at	DATETIME DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);`,

		//=======Indexes==========
		`CREATE INDEX IF NOT EXISTS idx_files_path ON files(path);`,
		`CREATE INDEX IF NOT EXISTS idx_extractions_file_id ON extractions(file_id);`,
		`CREATE INDEX IF NOT EXISTS idx_embeddings_file_id ON embeddings(file_id);`,
	}


	for _, stmt := range stmts{
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("[DB] init schema failed: %w", err)
		}
	}


	return nil
}


