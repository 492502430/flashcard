package handler

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/492502430/flashcard/backend/internal/config"
)

// AdminAuth validates admin token from header or query param.
func AdminAuth(password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Admin-Token")
			if token == "" {
				token = r.URL.Query().Get("token")
			}
			if token == "" || token != password {
				writeError(w, http.StatusUnauthorized, "invalid admin token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AdminStatsResponse holds all dashboard metrics.
type AdminStatsResponse struct {
	UsersToday       int                `json:"users_today"`
	TotalUsers       int                `json:"total_users"`
	TotalDecks       int                `json:"total_decks"`
	TotalCards       int                `json:"total_cards"`
	TokensUsed       int                `json:"tokens_used"`
	DAU              int                `json:"dau"`
	ReviewsToday     int                `json:"reviews_today"`
	FeedbackTotal    int                `json:"feedback_total"`
	LowQualityCards  int                `json:"low_quality_cards"`
	DailyNewUsers    []DayStat          `json:"daily_new_users"`
	DailyReviews     []DayStat          `json:"daily_reviews"`
	RecentUsers      []AdminUserRow     `json:"recent_users"`
	RecentDecks      []AdminDeckRow     `json:"recent_decks"`
	FeedbackByType   []AdminTypeCount   `json:"feedback_by_type"`
	DecksBySource    []AdminTypeCount   `json:"decks_by_source"`
	QualityBreakdown []AdminQualityItem `json:"quality_breakdown"`
}

type DayStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type AdminUserRow struct {
	ID         string `json:"id"`
	Nickname   string `json:"nickname"`
	InviteCode string `json:"invite_code"`
	TokensUsed int    `json:"tokens_used"`
	CreatedAt  string `json:"created_at"`
}

type AdminDeckRow struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	SourceName string `json:"source_name"`
	CardCount  int    `json:"card_count"`
	CreatedAt  string `json:"created_at"`
}

type AdminTypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type AdminQualityItem struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// AdminStats returns aggregated platform statistics as JSON.
func (h *Handler) AdminStats(w http.ResponseWriter, r *http.Request) {
	var resp AdminStatsResponse

	// Today's new users
	h.DB.Raw(`SELECT COUNT(*) FROM users WHERE created_at::date = CURRENT_DATE`).Scan(&resp.UsersToday)

	// Total users
	h.DB.Raw(`SELECT COUNT(*) FROM users`).Scan(&resp.TotalUsers)

	// Total decks
	h.DB.Raw(`SELECT COUNT(*) FROM decks`).Scan(&resp.TotalDecks)

	// Total cards
	h.DB.Raw(`SELECT COUNT(*) FROM cards`).Scan(&resp.TotalCards)

	// Actual tokens used (sum across all users)
	h.DB.Raw(`SELECT COALESCE(SUM(tokens_used), 0) FROM users`).Scan(&resp.TokensUsed)

	// Daily active users (users with review records today)
	h.DB.Raw(`SELECT COUNT(DISTINCT user_id) FROM review_records WHERE created_at::date = CURRENT_DATE`).Scan(&resp.DAU)

	// Reviews submitted today
	h.DB.Raw(`SELECT COUNT(*) FROM review_records WHERE created_at::date = CURRENT_DATE`).Scan(&resp.ReviewsToday)

	// Feedback volume
	h.DB.Raw(`SELECT COUNT(*) FROM card_feedbacks`).Scan(&resp.FeedbackTotal)

	// Approximate low-quality cards using the same practical thresholds as the miniapp.
	h.DB.Raw(`
		SELECT COUNT(*) FROM cards
		WHERE question = '' OR answer = ''
			OR LENGTH(question) > 90
			OR LENGTH(answer) > 220
			OR LENGTH(answer) < 6
	`).Scan(&resp.LowQualityCards)

	// Daily new users for past 7 days
	h.DB.Raw(`
		SELECT d::date::text AS date, COALESCE(COUNT(u.id), 0) AS count
		FROM generate_series(CURRENT_DATE - INTERVAL '6 days', CURRENT_DATE, '1 day') d
		LEFT JOIN users u ON u.created_at::date = d::date
		GROUP BY d::date
		ORDER BY date
	`).Scan(&resp.DailyNewUsers)

	// Daily review activity for past 7 days
	h.DB.Raw(`
		SELECT d::date::text AS date, COALESCE(COUNT(r.id), 0) AS count
		FROM generate_series(CURRENT_DATE - INTERVAL '6 days', CURRENT_DATE, '1 day') d
		LEFT JOIN review_records r ON r.created_at::date = d::date
		GROUP BY d::date
		ORDER BY date
	`).Scan(&resp.DailyReviews)

	h.DB.Raw(`
		SELECT id, COALESCE(NULLIF(nickname, ''), '闪卡用户') AS nickname,
			COALESCE(invite_code, '') AS invite_code,
			COALESCE(tokens_used, 0) AS tokens_used,
			created_at::text AS created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT 8
	`).Scan(&resp.RecentUsers)

	h.DB.Raw(`
		SELECT id, title, COALESCE(source_name, '') AS source_name,
			COALESCE(card_count, 0) AS card_count,
			created_at::text AS created_at
		FROM decks
		ORDER BY created_at DESC
		LIMIT 8
	`).Scan(&resp.RecentDecks)

	h.DB.Raw(`
		SELECT type, COUNT(*) AS count
		FROM card_feedbacks
		GROUP BY type
		ORDER BY count DESC
	`).Scan(&resp.FeedbackByType)

	h.DB.Raw(`
		SELECT COALESCE(NULLIF(source, ''), 'text') AS type, COUNT(*) AS count
		FROM decks
		GROUP BY source
		ORDER BY count DESC
	`).Scan(&resp.DecksBySource)

	if resp.DailyNewUsers == nil {
		resp.DailyNewUsers = []DayStat{}
	}
	if resp.DailyReviews == nil {
		resp.DailyReviews = []DayStat{}
	}
	if resp.RecentUsers == nil {
		resp.RecentUsers = []AdminUserRow{}
	}
	if resp.RecentDecks == nil {
		resp.RecentDecks = []AdminDeckRow{}
	}
	if resp.FeedbackByType == nil {
		resp.FeedbackByType = []AdminTypeCount{}
	}
	if resp.DecksBySource == nil {
		resp.DecksBySource = []AdminTypeCount{}
	}
	resp.QualityBreakdown = []AdminQualityItem{
		{Label: "健康卡片", Count: resp.TotalCards - resp.LowQualityCards},
		{Label: "待优化", Count: resp.LowQualityCards},
	}

	writeJSON(w, 200, resp)
}

//go:embed admin_dashboard.html
var adminDashboardHTML string

// adminHTMLTemplate has __ADMIN_TOKEN__ placeholder for the password.
const adminHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Flashcard Admin Dashboard</title>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f3ff; color: #1e1b4b; min-height: 100vh; }
header { background: #7C3AED; color: #fff; padding: 24px 32px; display: flex; justify-content: space-between; align-items: center; }
header h1 { font-size: 22px; font-weight: 600; letter-spacing: -0.3px; }
header .info { font-size: 13px; opacity: 0.85; }
.container { max-width: 1100px; margin: 0 auto; padding: 32px 24px; }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; margin-bottom: 32px; }
.card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.06), 0 1px 2px rgba(0,0,0,0.04); border: 1px solid #ede9fe; }
.card .label { font-size: 13px; color: #6b7280; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 6px; }
.card .value { font-size: 32px; font-weight: 700; color: #7C3AED; }
.card .value em { font-style: normal; font-size: 16px; color: #a78bfa; }
.chart-card { background: #fff; border-radius: 12px; padding: 24px 24px 16px; box-shadow: 0 1px 3px rgba(0,0,0,0.06); border: 1px solid #ede9fe; margin-bottom: 32px; }
.chart-card h2 { font-size: 17px; font-weight: 600; color: #1e1b4b; margin-bottom: 20px; }
.chart-wrap { width: 100%; height: 260px; }
.chart-wrap svg { width: 100%; height: 100%; }
.refresh { font-size: 12px; color: #9ca3af; text-align: right; margin-top: 8px; }
.error-state { background: #fef2f2; border: 1px solid #fecaca; border-radius: 12px; padding: 48px 24px; text-align: center; color: #991b1b; grid-column: 1 / -1; }
</style>
</head>
<body>
<header>
  <h1>⚡ Flashcard Admin</h1>
  <span class="info" id="ts">Loading...</span>
</header>
<div class="container">
  <div class="grid" id="grid"></div>
  <div class="chart-card">
    <h2>New Users — Past 7 Days</h2>
    <div class="chart-wrap" id="chart"></div>
  </div>
  <div class="refresh">Auto-refreshes every 30s · <span id="countdown">30</span>s</div>
</div>
<script>
const TOKEN = '__ADMIN_TOKEN__';
const MAX_VAL = 200;

function fmt(n) { return n != null ? Number(n).toLocaleString() : '\u2014'; }

function drawChart(data) {
  const days = data.daily_new_users || [];
  const W = 800, H = 240;
  const pad = { top: 10, right: 10, bottom: 30, left: 50 };
  const pw = W - pad.left - pad.right;
  const ph = H - pad.top - pad.bottom;

  var maxY = Math.max(MAX_VAL, ...days.map(function(d) { return d.count; }), 1);
  var ticks = [0, Math.round(maxY / 2), maxY];

  var points = days.map(function(d, i) {
    var x = pad.left + (pw * i) / Math.max(days.length - 1, 1);
    var y = pad.top + ph - (ph * d.count) / maxY;
    return { x: x, y: y, date: d.date, count: d.count };
  });

  var pathParts = [];
  for (var i = 0; i < points.length; i++) {
    pathParts.push((i === 0 ? 'M' : 'L') + points[i].x.toFixed(1) + ',' + points[i].y.toFixed(1));
  }
  var pathD = pathParts.join(' ');

  var areaD = pathD;
  if (points.length > 0) {
    areaD += ' L' + points[points.length - 1].x.toFixed(1) + ',' + (pad.top + ph).toFixed(1);
    areaD += ' L' + points[0].x.toFixed(1) + ',' + (pad.top + ph).toFixed(1) + ' Z';
  }

  var svg = '<svg viewBox="0 0 ' + W + ' ' + H + '" xmlns="http://www.w3.org/2000/svg">';

  // Y-axis grid lines and labels
  for (var ti = 0; ti < ticks.length; ti++) {
    var t = ticks[ti];
    var gy = pad.top + ph - (ph * t) / maxY;
    svg += '<line x1="' + pad.left + '" y1="' + gy.toFixed(1) + '" x2="' + (pad.left + pw).toFixed(1) + '" y2="' + gy.toFixed(1) + '" stroke="#e5e7eb" stroke-width="1"/>';
    svg += '<text x="' + (pad.left - 8) + '" y="' + (gy + 4).toFixed(1) + '" text-anchor="end" font-size="11" fill="#9ca3af">' + t + '</text>';
  }

  // Area fill
  svg += '<path d="' + areaD + '" fill="url(#grad)" opacity="0.3"/>';

  // Line
  if (points.length > 0) {
    svg += '<path d="' + pathD + '" fill="none" stroke="#7C3AED" stroke-width="2.5" stroke-linejoin="round" stroke-linecap="round"/>';
  }

  // Dots and labels
  for (var pi = 0; pi < points.length; pi++) {
    var p = points[pi];
    svg += '<circle cx="' + p.x.toFixed(1) + '" cy="' + p.y.toFixed(1) + '" r="4" fill="#7C3AED" stroke="#fff" stroke-width="2"/>';
    svg += '<text x="' + p.x.toFixed(1) + '" y="' + (pad.top + ph + 16).toFixed(1) + '" text-anchor="middle" font-size="10" fill="#6b7280">' + p.date.slice(5) + '</text>';
    svg += '<text x="' + p.x.toFixed(1) + '" y="' + (p.y - 10).toFixed(1) + '" text-anchor="middle" font-size="11" font-weight="600" fill="#7C3AED">' + p.count + '</text>';
  }

  // Gradient definition
  svg += '<defs><linearGradient id="grad" x1="0" y1="0" x2="0" y2="1"><stop offset="0%" stop-color="#7C3AED" stop-opacity="0.4"/><stop offset="100%" stop-color="#7C3AED" stop-opacity="0.02"/></linearGradient></defs>';
  svg += '</svg>';

  document.getElementById('chart').innerHTML = svg;
}

async function load() {
  try {
    var res = await fetch('/api/admin/stats', { headers: { 'X-Admin-Token': TOKEN } });
    if (!res.ok) throw new Error('HTTP ' + res.status);
    var d = await res.json();
    document.getElementById('grid').innerHTML =
      '<div class="card"><div class="label">New Users Today</div><div class="value">' + fmt(d.users_today) + '</div></div>' +
      '<div class="card"><div class="label">Total Users</div><div class="value">' + fmt(d.total_users) + '</div></div>' +
      '<div class="card"><div class="label">Total Decks</div><div class="value">' + fmt(d.total_decks) + '</div></div>' +
      '<div class="card"><div class="label">Total Cards</div><div class="value">' + fmt(d.total_cards) + '</div></div>' +
      '<div class="card"><div class="label">Tokens Used</div><div class="value">' + fmt(d.tokens_used) + '</div></div>' +
      '<div class="card"><div class="label">DAU (Today)</div><div class="value">' + fmt(d.dau) + '</div></div>';
    drawChart(d);
    document.getElementById('ts').textContent = 'Last updated: ' + new Date().toLocaleTimeString();
  } catch (e) {
    document.getElementById('grid').innerHTML = '<div class="error-state">Failed to load stats: ' + e.message + '</div>';
  }
}

// Countdown timer
var sec = 30;
setInterval(function() {
  sec--;
  if (sec <= 0) { sec = 30; load(); }
  document.getElementById('countdown').textContent = sec;
}, 1000);

load();
</script>
</body>
</html>`

// AdminPage serves the admin dashboard HTML directly.
func (h *Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()
	html := strings.ReplaceAll(adminDashboardHTML, "__ADMIN_TOKEN__", cfg.AdminPassword)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
