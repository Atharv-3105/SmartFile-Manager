package client 	

import(
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type EmbedRequest struct {
	Text string  `json:"text"`
}

type EmbedResponse struct {
	Embedding []float32  `json:"embedding"`
	Model     string     `json:"model"`
}

type ExtractRequest struct {
	FilePath    string    `json:"file_path"`
}

type ExtractResponse struct {
	FilePath     string    `json:"file_path"`
	Text		 string    `json:"text"`
	Embedding    []float32	`json:"embedding"`
	Model       string      `json:"model"`
	Status      string      `json:"status"`
}

type ExtractorClient struct {
	baseURL   string 
	http      *http.Client
}


func NewExtractorClient(baseURL string) *ExtractorClient{
	return &ExtractorClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}


func (c *ExtractorClient) Embed(text string) ([]float32, string, error) {
	reqBody, _ := json.Marshal(EmbedRequest{Text: text})

	resp , err := c.http.Post(
		c.baseURL+"/embed",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("[CLIENT] emebed service returned %d", resp.StatusCode)
	}

	var out EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, "", err
	}

	return out.Embedding, out.Model, nil
}


func (c *ExtractorClient) Extract(filePath string) (string, []float32, string, error) {

	reqBody, _ := json.Marshal(ExtractRequest{
		FilePath: filePath,
	})

	resp, err := c.http.Post(
		c.baseURL+"/extract",
		"application/json",
		bytes.NewBuffer(reqBody),
	)

	//If Response fails
	if err != nil {
		return "", nil, "", err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, "", fmt.Errorf("[CLIENT] extract service returned %d", resp.StatusCode)
	}

	var out ExtractResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", nil, "",err
	}

	return out.Text, out.Embedding, out.Model, nil
}
