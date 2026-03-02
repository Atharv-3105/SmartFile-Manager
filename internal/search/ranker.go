package search 

import "sort"


//Function to compute CosineSimilarity between the query vector
// and all embeddings, then returns the top-k results.
func RankTopK(query []float32,records []EmbeddingRecord, k int) []SearchResult {

	results := make([]SearchResult, 0, len(records))

	for _, rec := range records {
		score := CosineSimilarity(query, rec.Vector)

		results = append(results, SearchResult{
			Path: rec.Path,
			Score: score,
		})
	}

	//Sort by Similarity Score(In Descending order)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	//Handle k > len(results)
	if k > len(results) {
		k = len(results)
	}

	//Retur the top-k results
	return results[:k]
}