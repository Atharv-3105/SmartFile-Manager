package search

import "math"


//CosineSimilarity for computing similarity between 2 vectors
//Will return value in Range [-1,1]

func CosineSimilarity(a, b []float32) float32{
	if len(a) != len(b) || len(a) == 0{
		return 0
	}

	var dot, normA,normB float64

	for i := 0; i < len(a); i++{
		av := float64(a[i])
		bv := float64(b[i])

		dot += av * bv
		normA += av * av
		normB += bv * bv
	}

	if normA == 0 || normB == 0{
		return 0
	}

	return float32(dot /(math.Sqrt(normA) * math.Sqrt(normB)))
}