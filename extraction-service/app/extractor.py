import os
import easyocr
import pdfplumber

class TextExtractor:
    def __init__(self):
        self.ocr_reader = easyocr.Reader(["en"], gpu = False)
        
    def extract(self, file_path: str) -> str:
        if not os.path.exists(file_path):
            raise FileNotFoundError("[ES] File does not exist")
        
        #Check the extension of the file
        ext = os.path.splitext(file_path)[1].lower()
        
        if ext == ".txt":
            return self._from_txt(file_path)

        if ext in [".png", ".jpg", ".jpeg"]:
            return self._from_image(file_path)
        
        if ext == ".pdf":
            return self._from_pdf(file_path)
        
        raise ValueError("[ES] unsupported file type")
    
    def _from_txt(self, path: str) -> str:
        with open(path, "r", encoding = "utf-8", errors = "ignore") as file:
            return file.read()
    
    def _from_image(self, path: str) -> str:
        results = self.ocr_reader.readtext(path, detail = 0)
        return " ".join(results)
    
    def _from_pdf(self, path: str) -> str:
        txt_chunks = []
        with pdfplumber.open(path) as pdf:
            for page in pdf.pages:
                txt = page.extract_text()
                if txt:
                    txt_chunks.append(txt)
        
        return "\n".join(txt_chunks)