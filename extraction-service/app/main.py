from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from app.extractor import TextExtractor
from app.embedder import Embedder

app = FastAPI(title= "Extraction-Service")

extractor = TextExtractor()
embedder = Embedder()

class ExtractRequest(BaseModel):
    file_path : str

class ExtractResponse(BaseModel):
    file_path: str
    text: str
    embedding: list[float]
    model: str
    status: str 
    
@app.post("/extract", response_model=ExtractResponse)
def extract(req: ExtractRequest):
    try:
        text = extractor.extract(req.file_path)
        
        if not text.strip():
            raise ValueError("[ES] No text extracted")

        embedding = embedder.embed(text)
        
        return {
            "file_path": req.file_path,
            "text": text,
            "embedding": embedding,
            "model":embedder.model_name,
            "status": "ok",
        }
        
    except Exception as e:
        raise HTTPException(status_code=400, detail = str(e))