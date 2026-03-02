package search

import(
	"database/sql"
	"fmt"
)

//Function to LoadEmbeddings for a given model from the DB.
func LoadEmbeddings(db *sql.DB, model string) ([]EmbeddingRecord, error) {
	rows, err := db.Query(
		`SELECT e.file_id, f.path, e.vector
		FROM embeddings e
		JOIN files f on e.file_id = f.id
		WHERE e.model = ?`,
		model,
	)
	if err != nil {
		return nil, fmt.Errorf("[SEARCH] query embeddings: %w", err)
	}
	defer rows.Close()

	var results []EmbeddingRecord

	for rows.Next() {
		var (
			fileID  int64
			path   string
			blob   []byte
		)

		if err := rows.Scan(&fileID, &path, &blob); err != nil{
			return nil, fmt.Errorf("[SEARCH] scan embedding row: %w", err)
		}

		vec, err := BlobToFloat32Slice(blob)
		if err != nil {
			return nil, fmt.Errorf("[SEARCH] decode vector for %s: %w", path, err)
		}

		results = append(results, EmbeddingRecord{
			FileID: fileID,
			Path: path,
			Vector: vec,
		})
	}

	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("[SEARCH] iterate embeddings: %w", err)
	}

	return results, nil

}