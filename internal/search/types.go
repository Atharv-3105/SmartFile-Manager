package search

type EmbeddingRecord  struct {
	FileID 	int64
	Path    string 
	Vector  []float32
}

type SearchResult struct {
	Path    string 
	Score   float32
}