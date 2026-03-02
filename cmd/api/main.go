package main

import (
	// "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Atharv-3105/File-Manager/internal/client"
	"github.com/Atharv-3105/File-Manager/internal/search"
	"github.com/Atharv-3105/File-Manager/internal/storage"
)

type SearchRequest struct {
	Query	string   `json:"query"`
	TopK    int      `json:"top_k"`
}

type SearchResponse struct {
	Results []search.SearchResult  `json:"results"`
}

func main() {
	log.Println("[API] starting search service")

	//DB initialization
	_ = os.MkdirAll("data", os.ModePerm)

	db, err := storage.Open("data/index.db")
	if err != nil {
		log.Fatal("[API] db oepn:", err)
	}
	defer db.Close()

	if err := storage.InitSchema(db); err != nil {
		log.Fatal("[API] schema init:", err)
	}

	extractor := client.NewExtractorClient("http://127.0.0.1:8001")

	//Routes    
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return 
		}

		var req SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		//If top_k is <= 0; by default consider top_k as 5
		if req.TopK <= 0{
			req.TopK = 5
		}

		//Embed Natural Language Query into Vector
		queryVec, model, err := extractor.Embed(req.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		//Load Embeddings from the DB
		records, err := search.LoadEmbeddings(db, model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		//return the TopK results
		results := search.RankTopK(queryVec, records, req.TopK)

		resp := SearchResponse{Results: results}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("[API] listening on Port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}