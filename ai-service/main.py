"""
FlashCard AI Service — Card Generation via DeepSeek API.
POST /generate  →  text, deck_id  →  structured flashcard array
"""
import json
import os
import io
from fastapi import FastAPI, HTTPException, UploadFile, File
from pydantic import BaseModel
from openai import OpenAI

app = FastAPI(title="FlashCard AI Service", version="1.0.0")

# Load key from .env file or environment
def _load_key():
    env_file = os.path.join(os.path.dirname(__file__), ".env")
    if os.path.exists(env_file):
        for line in open(env_file):
            if line.startswith("DEEPSEEK_API_KEY="):
                return line.strip().split("=", 1)[1]
    return os.environ.get("DEEPSEEK_API_KEY", "")

client = OpenAI(
    api_key=_load_key(),
    base_url="https://api.deepseek.com/v1",
)

PROMPT = """你是一个教育专家。请根据以下文本生成闪卡（问答题卡片）。

规则：
1. 只提取文本中的核心知识点和关键概念，忽略铺垫和废话
2. 每张卡一个问题，答案简洁但完整
3. 如果知识点之间有逻辑关联，可以适当分组
4. 一般生成 8-20 张卡片即可，质量优先，不要凑数

文本：
{text}

请返回纯 JSON 数组，不要加 markdown 代码块标记。"""


class GenerateRequest(BaseModel):
    text: str
    deck_id: str


class Card(BaseModel):
    q: str
    a: str
    tags: list[str] = []


class GenerateResponse(BaseModel):
    deck_id: str
    cards: list[Card]
    count: int


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/generate", response_model=GenerateResponse)
def generate(req: GenerateRequest):
    if len(req.text) < 50:
        raise HTTPException(400, "text too short (min 50 chars)")

    try:
        response = client.chat.completions.create(
            model="deepseek-chat",
            messages=[
                {"role": "system", "content": "你是教育专家。只返回 JSON 数组，不要加 markdown 标记。"},
                {"role": "user", "content": PROMPT.format(text=req.text)},
            ],
            temperature=0.3,
            max_tokens=2000,
        )
        content = response.choices[0].message.content
        cards = json.loads(content)
        return GenerateResponse(deck_id=req.deck_id, cards=cards, count=len(cards))
    except json.JSONDecodeError:
        raise HTTPException(500, "AI returned invalid JSON")
    except Exception as e:
        raise HTTPException(500, f"AI generation failed: {str(e)}")


@app.post("/extract")
async def extract_text(file: UploadFile = File(...)):
    """Extract text from uploaded file (PDF, TXT, DOCX)."""
    content = await file.read()
    filename = file.filename or ""

    # TXT — direct read
    if filename.lower().endswith(".txt"):
        text = content.decode("utf-8", errors="ignore")
        return {"text": text, "filename": filename, "size": len(text)}

    # PDF — PyMuPDF
    if filename.lower().endswith(".pdf"):
        try:
            import fitz  # PyMuPDF
            doc = fitz.open(stream=content, filetype="pdf")
            text = ""
            for page in doc:
                text += page.get_text()
            doc.close()
            return {"text": text, "filename": filename, "size": len(text)}
        except ImportError:
            raise HTTPException(500, "PDF support requires PyMuPDF: pip install PyMuPDF")

    # DOCX — python-docx
    if filename.lower().endswith(".docx"):
        try:
            from docx import Document
            doc = Document(io.BytesIO(content))
            text = "\n".join(p.text for p in doc.paragraphs)
            return {"text": text, "filename": filename, "size": len(text)}
        except ImportError:
            raise HTTPException(500, "DOCX support requires python-docx: pip install python-docx")

    # Image — OCR via pytesseract
    if filename.lower().endswith((".png", ".jpg", ".jpeg")):
        try:
            from PIL import Image
            import pytesseract
            img = Image.open(io.BytesIO(content))
            text = pytesseract.image_to_string(img, lang="chi_sim+eng")
            return {"text": text, "filename": filename, "size": len(text)}
        except ImportError:
            raise HTTPException(500, "Image OCR requires: pip install pytesseract pillow")
        except Exception as e:
            raise HTTPException(500, f"OCR failed: {str(e)}")

    raise HTTPException(400, f"unsupported format: {filename}")
