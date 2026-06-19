package handler

import (
	"net/http"
)

// docsHTML is a self-contained Swagger-like HTML page listing all API endpoints
// with curl examples, styled like the admin page.
const docsHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Flashcard API 文档</title>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif; background: #f5f3ff; color: #1e1b4b; min-height: 100vh; }
header { background: linear-gradient(135deg, #7C3AED, #8B5CF6); color: #fff; padding: 28px 32px; }
header h1 { font-size: 22px; font-weight: 700; letter-spacing: -0.3px; margin-bottom: 4px; }
header .sub { font-size: 13px; opacity: 0.85; }
.container { max-width: 1000px; margin: 0 auto; padding: 32px 24px; }
.group { margin-bottom: 36px; }
.group-title { font-size: 18px; font-weight: 700; color: #7C3AED; margin-bottom: 16px; display: flex; align-items: center; gap: 8px; }
.group-title .icon { width: 28px; height: 28px; border-radius: 8px; background: rgba(124,58,237,0.12); display: inline-flex; align-items: center; justify-content: center; font-size: 15px; }
.endpoint { background: #fff; border-radius: 12px; padding: 20px 24px; margin-bottom: 12px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); border: 1px solid #ede9fe; }
.endpoint-header { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.method { display: inline-block; padding: 3px 10px; border-radius: 6px; font-size: 12px; font-weight: 700; text-transform: uppercase; font-family: 'SF Mono', 'Fira Code', 'Monaco', monospace; }
.method.get { background: #dcfce7; color: #16a34a; }
.method.post { background: #dbeafe; color: #2563eb; }
.method.put { background: #fef3c7; color: #d97706; }
.method.delete { background: #fce7f3; color: #db2777; }
.path { font-family: 'SF Mono', 'Fira Code', 'Monaco', monospace; font-size: 15px; font-weight: 600; color: #1e1b4b; }
.auth-badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; background: #fef3c7; color: #92400e; }
.auth-badge.no-auth { background: #f0fdf4; color: #166534; }
.desc { font-size: 13px; color: #6b7280; margin-bottom: 10px; line-height: 1.5; }
.curl-block { background: #1e1b4b; color: #e2e8f0; border-radius: 8px; padding: 12px 16px; font-family: 'SF Mono', 'Fira Code', 'Monaco', monospace; font-size: 12px; line-height: 1.8; overflow-x: auto; white-space: pre-wrap; word-break: break-all; position: relative; }
.curl-block .copy-btn { position: absolute; top: 8px; right: 8px; background: rgba(255,255,255,0.12); border: none; color: #c7d2fe; font-size: 11px; padding: 4px 10px; border-radius: 6px; cursor: pointer; font-family: inherit; }
.curl-block .copy-btn:hover { background: rgba(255,255,255,0.2); }
.curl-block .copy-btn.copied { background: #16a34a; color: #fff; }
.note { font-size: 12px; color: #9ca3af; margin-top: 8px; font-style: italic; }
.footer { text-align: center; padding: 24px; color: #9ca3af; font-size: 12px; border-top: 1px solid #ede9fe; margin-top: 20px; }
</style>
</head>
<body>
<header>
  <h1>⚡ Flashcard API 文档</h1>
  <div class="sub">所有接口列表 &amp; curl 调用示例 · BASE_URL 替换为实际地址</div>
</header>
<div class="container">

<!-- Auth -->
<div class="group">
  <div class="group-title"><span class="icon">🔑</span>认证</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method post">POST</span>
      <span class="path">/api/auth/login</span>
      <span class="auth-badge no-auth">无需认证</span>
    </div>
    <div class="desc">微信小程序登录。传入 wx.login() 返回的 code，换取 JWT token。开发环境可使用 <code>dev-</code> 前缀的假 code。</div>
    <div class="curl-block">curl -X POST BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"code":"dev-1718000000000"}'<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: { "token": "...", "user": { "id": "...", "nickname": "..." } }</div>
  </div>
</div>

<!-- Decks -->
<div class="group">
  <div class="group-title"><span class="icon">📚</span>牌组 (Decks)</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/decks</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取当前用户的所有牌组列表。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/decks \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: [{ "id": "...", "title": "马原重点", "card_count": 12, "created_at": "..." }]</div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method post">POST</span>
      <span class="path">/api/decks</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">创建新牌组。如果提供 text 字段，后台会调用 AI 自动生成闪卡（异步）。</div>
    <div class="curl-block">curl -X POST BASE_URL/api/decks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"马原重点","text":"马克思主义基本原理概论..."}'<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: { "id": "...", "title": "马原重点", "card_count": 0 }（AI 生成中）</div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/decks/{id}</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取单个牌组详情，包含所有卡片。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/decks/DECK_ID \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: { "deck": {...}, "cards": [{ "id": "...", "question": "...", "answer": "...", "state": "new" }] }</div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/decks/{id}/review</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取牌组中待复习的卡片。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/decks/DECK_ID/review \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method delete">DELETE</span>
      <span class="path">/api/decks/{id}</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">删除牌组及其所有卡片。</div>
    <div class="curl-block">curl -X DELETE BASE_URL/api/decks/DECK_ID \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<!-- Templates -->
<div class="group">
  <div class="group-title"><span class="icon">🧩</span>模板库 (Templates)</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/templates</span>
      <span class="auth-badge no-auth">无需认证</span>
    </div>
    <div class="desc">获取所有预设模板牌组列表（摘要信息，不含卡片详情）。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/templates<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: [{ "id": "tpl-ma-yuan", "title": "马原重点", "description": "...", "category": "考研政治", "card_count": 12, "icon": "diamond" }]</div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/templates/{id}</span>
      <span class="auth-badge no-auth">无需认证</span>
    </div>
    <div class="desc">获取单个模板详情，包含所有卡片。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/templates/tpl-ma-yuan<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method post">POST</span>
      <span class="path">/api/templates/{id}/import</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">将模板牌组导入当前用户的牌组（复制模板及其所有卡片到用户账户）。</div>
    <div class="curl-block">curl -X POST BASE_URL/api/templates/tpl-ma-yuan/import \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: { "id": "...", "title": "马原重点", "card_count": 12 }（已创建的新牌组）</div>
  </div>
</div>

<!-- Cards -->
<div class="group">
  <div class="group-title"><span class="icon">🃏</span>卡片 (Cards)</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method put">PUT</span>
      <span class="path">/api/cards/{id}</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">更新卡片内容（问题和答案）。</div>
    <div class="curl-block">curl -X PUT BASE_URL/api/cards/CARD_ID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"question":"新问题","answer":"新答案"}'<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method delete">DELETE</span>
      <span class="path">/api/cards/{id}</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">删除单张卡片。</div>
    <div class="curl-block">curl -X DELETE BASE_URL/api/cards/CARD_ID \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/cards/search</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">搜索当前用户的所有卡片。</div>
    <div class="curl-block">curl -X GET "BASE_URL/api/cards/search?q=马克思主义" \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<!-- Review -->
<div class="group">
  <div class="group-title"><span class="icon">🔄</span>复习 (Review)</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/review/today</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取今日待复习的卡片列表。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/review/today \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method post">POST</span>
      <span class="path">/api/review</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">提交复习评分（1=完全忘记, 2=困难, 3=良好, 4=简单），FSRS 算法计算下次复习时间。</div>
    <div class="curl-block">curl -X POST BASE_URL/api/review \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"card_id":"CARD_ID","rating":3}'<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method post">POST</span>
      <span class="path">/api/cards/{id}/feedback</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">提交卡片反馈（content_error / answer_too_brief / question_unclear）。</div>
    <div class="curl-block">curl -X POST BASE_URL/api/cards/CARD_ID/feedback \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"content_error"}'<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<!-- Stats & Achievements -->
<div class="group">
  <div class="group-title"><span class="icon">📊</span>统计 &amp; 成就</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/stats</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取用户的复习统计信息。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/stats \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/achievements</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取用户已获得的成就列表。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/achievements \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/checkin</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">获取用户的签到记录和连续签到天数。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/checkin \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<!-- Upload & Export -->
<div class="group">
  <div class="group-title"><span class="icon">📤</span>上传 &amp; 导出</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method post">POST</span>
      <span class="path">/api/upload</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">上传文件（支持文本/Pdf），从内容中 AI 生成卡片。</div>
    <div class="curl-block">curl -X POST BASE_URL/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@notes.txt" \
  -F "title=学习笔记"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/export</span>
      <span class="auth-badge">需登录</span>
    </div>
    <div class="desc">导出用户的所有牌组数据（JSON 格式）。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/export \
  -H "Authorization: Bearer YOUR_TOKEN"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<!-- Admin -->
<div class="group">
  <div class="group-title"><span class="icon">⚙️</span>管理 (Admin)</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/admin</span>
      <span class="auth-badge">X-Admin-Token</span>
    </div>
    <div class="desc">管理员仪表盘（HTML 页面），显示平台统计数据。</div>
    <div class="curl-block">curl -X GET BASE_URL/admin \
  -H "X-Admin-Token: YOUR_ADMIN_PASSWORD"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/admin/stats</span>
      <span class="auth-badge">X-Admin-Token</span>
    </div>
    <div class="desc">管理统计 JSON 接口。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/admin/stats \
  -H "X-Admin-Token: YOUR_ADMIN_PASSWORD"<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<!-- Docs & Health -->
<div class="group">
  <div class="group-title"><span class="icon">🏥</span>系统</div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/health</span>
      <span class="auth-badge no-auth">无需认证</span>
    </div>
    <div class="desc">健康检查接口。</div>
    <div class="curl-block">curl -X GET BASE_URL/health<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
    <div class="note">返回: { "status": "ok" }</div>
  </div>

  <div class="endpoint">
    <div class="endpoint-header">
      <span class="method get">GET</span>
      <span class="path">/api/docs</span>
      <span class="auth-badge no-auth">无需认证</span>
    </div>
    <div class="desc">API 文档页面（即本页）。</div>
    <div class="curl-block">curl -X GET BASE_URL/api/docs<span class="copy-btn" onclick="copyCurl(this)">复制</span></div>
  </div>
</div>

<div class="footer">Flashcard API · 闪卡记忆</div>
</div>

<script>
function copyCurl(btn) {
  var block = btn.parentElement;
  var text = block.textContent.replace(/复制$/, '').trim();
  navigator.clipboard.writeText(text).then(function() {
    btn.textContent = '已复制';
    btn.classList.add('copied');
    setTimeout(function() { btn.textContent = '复制'; btn.classList.remove('copied'); }, 2000);
  }).catch(function() {
    // Fallback for non-HTTPS
    var ta = document.createElement('textarea');
    ta.value = text;
    document.body.appendChild(ta);
    ta.select();
    document.execCommand('copy');
    document.body.removeChild(ta);
    btn.textContent = '已复制';
    btn.classList.add('copied');
    setTimeout(function() { btn.textContent = '复制'; btn.classList.remove('copied'); }, 2000);
  });
}
</script>
</body>
</html>`

// DocsPage serves the API documentation HTML page.
func (h *Handler) DocsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(docsHTML))
}
