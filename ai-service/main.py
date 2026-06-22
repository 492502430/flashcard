"""
FlashCard AI Service — Card Generation via DeepSeek API.
POST /generate  →  text, deck_id  →  structured flashcard array
POST /extract   →  file upload     →  extracted text
POST /optimize  →  cards array     →  optimized cards
"""
import json, os, io
from fastapi import FastAPI, HTTPException, UploadFile, File
from pydantic import BaseModel
from openai import OpenAI

app = FastAPI(title="FlashCard AI Service", version="1.0.0")

def _load_key():
    env_file = os.path.join(os.path.dirname(__file__), ".env")
    if os.path.exists(env_file):
        for line in open(env_file):
            if line.startswith("DEEPSEEK_API_KEY="):
                return line.strip().split("=", 1)[1]
    return os.environ.get("DEEPSEEK_API_KEY", "")

client = OpenAI(api_key=_load_key(), base_url="https://api.deepseek.com/v1")

# ── Generate ──

class GenerateRequest(BaseModel):
    text: str
    deck_id: str
    card_count: int = 25

class Card(BaseModel):
    q: str
    a: str
    tags: list[str] = []

class GenerateResponse(BaseModel):
    deck_id: str
    cards: list[Card]
    count: int

PROMPT = """你是一个严谨的学习资料整理专家。请根据以下文本生成 {card_count} 张高质量闪卡（问答题卡片）。

规则：
1. 必须只依据原文内容生成，不要补充原文没有的信息，不要编造例子、数字、结论或定义。
2. 优先选择真正值得记忆的重点：核心概念、定义、公式、步骤、条件、对比关系、原因结果、易混点、考试/工作高频点。
3. 每张卡只考一个知识点；问题要具体，避免“介绍一下/谈谈/是什么内容”这类泛问题。
4. 答案要简洁但足够完整，适合主动回忆；长答案请压缩成 2-4 个关键要点。
5. 避免重复卡片；相近知识点要合并或改成对比题。
6. 尽量接近用户要求的数量：目标是 {card_count} 张；如果原文有效信息不足，可以少量减少，但不要自行限制为 25 张。
7. tags 应提取自原文主题或章节名，优先使用 1-3 个短标签。
8. 返回 JSON 数组：每项必须包含 q（问题）、a（答案）、tags（标签数组）三个字段，字段名必须是英文 q、a、tags。

文本：
{text}

请返回纯 JSON 数组，不要加 markdown 代码块标记。"""

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/generate", response_model=GenerateResponse)
def generate(req: GenerateRequest):
    if not req.text or len(req.text) < 50:
        raise HTTPException(400, f"text too short (min 50 chars, got {len(req.text)})")
    card_count = max(5, min(req.card_count or 25, 300))

    try:
        response = client.chat.completions.create(
            model="deepseek-chat",
            messages=[
                {"role": "system", "content": PROMPT.format(card_count=card_count, text=req.text)},
                {"role": "user", "content": req.text},
            ],
            temperature=0.3,
            max_tokens=max(3000, min(12000, card_count * 120)),
        )
        content = response.choices[0].message.content
        cards = json.loads(content)
        return GenerateResponse(deck_id=req.deck_id, cards=cards, count=len(cards))
    except json.JSONDecodeError:
        raise HTTPException(500, "AI returned invalid JSON")
    except Exception as e:
        raise HTTPException(500, f"AI generation failed: {str(e)}")

# ── Extract ──

@app.post("/extract")
async def extract_text(file: UploadFile = File(...)):
    content = await file.read()
    filename = file.filename or ""

    if filename.lower().endswith(".txt"):
        text = content.decode("utf-8", errors="ignore")
        return {"text": text, "filename": filename, "size": len(text)}

    if filename.lower().endswith(".pdf"):
        try:
            import fitz
            doc = fitz.open(stream=content, filetype="pdf")
            text = "".join(page.get_text() for page in doc)
            doc.close()
            return {"text": text, "filename": filename, "size": len(text)}
        except ImportError:
            raise HTTPException(500, "PDF requires: pip install PyMuPDF")

    if filename.lower().endswith(".docx"):
        try:
            from docx import Document
            doc = Document(io.BytesIO(content))
            text = "\n".join(p.text for p in doc.paragraphs)
            return {"text": text, "filename": filename, "size": len(text)}
        except ImportError:
            raise HTTPException(500, "DOCX requires: pip install python-docx")

    if filename.lower().endswith((".png", ".jpg", ".jpeg")):
        try:
            from paddleocr import PaddleOCR
            ocr = PaddleOCR(lang='ch')
            from PIL import Image
            img = Image.open(io.BytesIO(content))
            import tempfile
            with tempfile.NamedTemporaryFile(suffix='.png', delete=False) as tmp:
                img.save(tmp.name)
            result = ocr.predict(tmp.name)
            os.unlink(tmp.name)
            text = ""
            if result:
                for r in result:
                    j = r.json() if callable(r.json) else r.json
                    rec = j.get("res", {}).get("rec_texts", "")
                    if rec and rec != "00":
                        text += rec + "\n"
            return {"text": text, "filename": filename, "size": len(text)}
        except ImportError:
            raise HTTPException(500, "Image OCR requires: pip install paddleocr paddlepaddle pillow")
        except Exception as e:
            raise HTTPException(500, f"OCR failed: {str(e)}")

    raise HTTPException(400, f"unsupported format: {filename}")

# ── Optimize ──

class OptimizeRequest(BaseModel):
    cards: list[dict]

OPTIMIZE_PROMPT = """你是闪卡优化专家。以下卡片需要优化，原因已标注。请改进每张卡：

规则：
1. 如果问题过长，拆成单一知识点
2. 如果答案过短，补充关键信息（但保持简洁）
3. 如果答案过长，提炼核心要点
4. 如果与其他卡片重复，合并或差异化
5. 保持原卡片的核心知识点不变

原卡片：
{cards_json}

返回纯 JSON 数组，每项包含 q（问题）、a（答案），字段名必须是英文 q、a。"""

@app.post("/optimize")
def optimize(req: OptimizeRequest):
    if not req.cards:
        raise HTTPException(400, "cards is required")

    try:
        response = client.chat.completions.create(
            model="deepseek-chat",
            messages=[{"role": "user", "content": OPTIMIZE_PROMPT.format(
                cards_json=json.dumps(req.cards, ensure_ascii=False)
            )}],
            temperature=0.3,
            max_tokens=2000,
        )
        content = response.choices[0].message.content
        cards = json.loads(content)
        return {"cards": cards, "count": len(cards)}
    except json.JSONDecodeError:
        raise HTTPException(500, "AI returned invalid JSON")
    except Exception as e:
        raise HTTPException(500, f"AI optimization failed: {str(e)}")
