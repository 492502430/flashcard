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

PROMPT = """你是一位资深教育内容设计师。请根据以下文本生成高质量的学习闪卡。

## 题型要求（混合使用）
1. **概念问答**（占 40%）：针对核心概念、定理、定义提问
2. **填空题**（占 30%）：将关键术语或数字挖空，用「______」表示。q 字段填完整句子（含空），a 字段填答案
3. **应用题/场景题**（占 30%）：设计一个具体场景，让学习者运用知识点解决问题

## 内容筛选规则
- 只选取考试常考点、易错点、核心定义，忽略叙述性铺垫
- 如果原文描述了因果关系或对比关系，必须出题
- 如果原文有具体数字、日期、公式，优先出题
- 如果原文案例有实际应用价值，转化为场景题

## ⚠️ 输出格式（极其重要）
你必须且只能返回一个 JSON 数组，不要任何其他文字。不要加 ```json 标记。不要加任何解释。直接返回数组本身。

格式示例：
[{"q":"什么是辩证唯物主义？","a":"认为物质决定意识的哲学观点","tags":["哲学","唯物主义"],"type":"qa"}, {"q":"唯物辩证法的三大规律是______、______和否定之否定规律","a":"对立统一、量变质变","tags":["辩证法"],"type":"fill"}]

type 必须是 "qa"、"fill"、"scenario" 之一。"""

GENERATE_USER = """请根据以下文本生成 {card_count} 张闪卡：

{text}"""


class GenerateRequest(BaseModel):
    text: str
    deck_id: str


class Card(BaseModel):
    q: str
    a: str
    tags: list[str] = []
    type: str = "qa"


class GenerateResponse(BaseModel):
    deck_id: str
    cards: list[Card]
    count: int
    tokens_used: int = 0


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/generate", response_model=GenerateResponse)
def generate(req: GenerateRequest):
    if len(req.text) < 50:
        raise HTTPException(400, "text too short (min 50 chars)")

    card_count = max(8, min(25, len(req.text) // 80))

    try:
        response = client.chat.completions.create(
            model="deepseek-chat",
            messages=[
                {"role": "system", "content": PROMPT},
                {"role": "user", "content": GENERATE_USER.format(card_count=card_count, text=req.text)},
            ],
            temperature=0.4,
            max_tokens=3000,
        )
        content = response.choices[0].message.content
        cards = json.loads(content)
        usage = response.usage.total_tokens if response.usage else 0
        return GenerateResponse(deck_id=req.deck_id, cards=cards, count=len(cards), tokens_used=usage)
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
