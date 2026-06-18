"""
FlashCard AI Service — Card Generation via DeepSeek API.
POST /generate  →  text, deck_id  →  structured flashcard array
"""
import json
import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from openai import OpenAI

app = FastAPI(title="FlashCard AI Service", version="1.0.0")

client = OpenAI(
    api_key=os.environ["DEEPSEEK_API_KEY"],
    base_url="https://api.deepseek.com/v1",
)

PROMPT = """你是一个教育专家。请根据以下文本生成闪卡（问答题卡片）。

规则：
1. 每张卡一个知识点
2. 问题简洁明确，答案准确完整
3. 只基于给定文本，不添加文本中没有的内容
4. 返回 JSON 数组格式：每项包含 q（问题）、a（答案）、tags（标签数组）三个字段
5. 生成 10-15 张卡片

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
