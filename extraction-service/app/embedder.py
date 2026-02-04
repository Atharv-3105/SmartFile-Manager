from sentence_transformers import SentenceTransformer


class Embedder:
    def __init__(self):
        self.model_name = "all-MiniLM-L6-v2"
        self.model = SentenceTransformer(self.model_name, device = "cpu")
    
    def embed(self, text: str):
        return self.model.encode(text).tolist()
    