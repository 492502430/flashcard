package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TemplateDeck represents a preset deck available in the template library.
type TemplateDeck struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	CardCount   int            `json:"card_count"`
	Icon        string         `json:"icon"`
	Cards       []TemplateCard `json:"cards,omitempty"`
}

// TemplateCard is a single card inside a template deck.
type TemplateCard struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Tags     []string `json:"tags"`
}

// templateDecks is the seed data for preset template decks.
var templateDecks = []TemplateDeck{
	{
		ID:          "tpl-ma-yuan",
		Title:       "马原重点",
		Description: "马克思主义基本原理概论核心考点，涵盖唯物辩证法、认识论、唯物史观等重点内容",
		Category:    "考研政治",
		CardCount:   12,
		Icon:        "diamond",
		Cards: []TemplateCard{
			{Question: "马克思主义的基本立场是什么？", Answer: "以无产阶级的解放和全人类的解放为己任，以人的自由全面发展为美好目标，以人民为中心，一切为了人民，一切依靠人民。", Tags: []string{"马原", "总论"}},
			{Question: "马克思主义的鲜明特征有哪些？", Answer: "科学性、革命性、实践性、人民性、发展性。", Tags: []string{"马原", "总论"}},
			{Question: "哲学基本问题是什么？", Answer: "思维和存在的关系问题。包括两个方面：①思维和存在何者为第一性（划分唯物主义和唯心主义）；②思维能否正确认识存在（划分可知论和不可知论）。", Tags: []string{"马原", "唯物论"}},
			{Question: "物质的唯一特性是什么？", Answer: "客观实在性。物质是不依赖于人的意识并能为人的意识所反映的客观实在。", Tags: []string{"马原", "唯物论"}},
			{Question: "对立统一规律（矛盾规律）的内容是什么？", Answer: "对立统一规律是唯物辩证法的实质和核心。矛盾是事物发展的根本动力，矛盾具有同一性和斗争性两种基本属性，矛盾的普遍性和特殊性的关系是矛盾问题的精髓。", Tags: []string{"马原", "辩证法"}},
			{Question: "量变和质变的辩证关系是什么？", Answer: "量变是质变的必要准备，质变是量变的必然结果。量变和质变相互渗透，在总的量变过程中有阶段性和局部性的部分质变，在质变过程中有量的扩张。", Tags: []string{"马原", "辩证法"}},
			{Question: "辩证否定观的基本内容是什么？", Answer: "否定是事物的自我否定；否定是事物发展的环节；否定是事物联系的环节；辩证否定的实质是「扬弃」。", Tags: []string{"马原", "辩证法"}},
			{Question: "实践是认识的基础，具体表现在哪些方面？", Answer: "①实践是认识的来源；②实践是认识发展的动力；③实践是检验认识真理性的唯一标准；④实践是认识的目的。", Tags: []string{"马原", "认识论"}},
			{Question: "感性认识和理性认识的辩证关系是什么？", Answer: "感性认识是认识的初级阶段，理性认识是认识的高级阶段。感性认识有待于发展和深化为理性认识，理性认识依赖于感性认识。二者相互渗透、相互包含。", Tags: []string{"马原", "认识论"}},
			{Question: "真理的绝对性和相对性的关系是什么？", Answer: "真理的绝对性（绝对真理）和相对性（相对真理）是辩证统一的。绝对真理寓于相对真理之中，无数相对真理的总和构成绝对真理。", Tags: []string{"马原", "认识论"}},
			{Question: "社会存在和社会意识的辩证关系是什么？", Answer: "社会存在决定社会意识，社会意识是社会存在的反映并反作用于社会存在。社会意识具有相对独立性。", Tags: []string{"马原", "唯物史观"}},
			{Question: "生产力和生产关系的辩证关系是什么？", Answer: "生产力决定生产关系，生产关系反作用于生产力。当生产关系适合生产力发展状况时就会促进生产力发展，否则就会阻碍生产力发展。", Tags: []string{"马原", "唯物史观"}},
		},
	},
	{
		ID:          "tpl-mao-zhong-te",
		Title:       "毛中特核心",
		Description: "毛泽东思想和中国特色社会主义理论体系概论重点考点，涵盖新民主主义革命理论、社会主义改造、邓小平理论等",
		Category:    "考研政治",
		CardCount:   10,
		Icon:        "star",
		Cards: []TemplateCard{
			{Question: "毛泽东思想活的灵魂是什么？", Answer: "实事求是、群众路线、独立自主。", Tags: []string{"毛中特", "毛泽东思想"}},
			{Question: "新民主主义革命的总路线是什么？", Answer: "无产阶级领导的，人民大众的，反对帝国主义、封建主义和官僚资本主义的革命。", Tags: []string{"毛中特", "新民主主义革命"}},
			{Question: "新民主主义革命的三大法宝是什么？", Answer: "统一战线、武装斗争、党的建设。", Tags: []string{"毛中特", "新民主主义革命"}},
			{Question: "党在过渡时期的总路线是什么？", Answer: "要在一个相当长的时期内，逐步实现国家的社会主义工业化，并逐步实现国家对农业、手工业和资本主义工商业的社会主义改造。（「一化三改」）", Tags: []string{"毛中特", "社会主义改造"}},
			{Question: "邓小平理论回答的根本问题是什么？", Answer: "什么是社会主义、怎样建设社会主义。", Tags: []string{"毛中特", "邓小平理论"}},
			{Question: "社会主义的本质是什么？", Answer: "解放生产力，发展生产力，消灭剥削，消除两极分化，最终达到共同富裕。", Tags: []string{"毛中特", "邓小平理论"}},
			{Question: "「三个代表」重要思想的核心内容是什么？", Answer: "中国共产党始终代表中国先进生产力的发展要求，代表中国先进文化的前进方向，代表中国最广大人民的根本利益。", Tags: []string{"毛中特", "三个代表"}},
			{Question: "科学发展观的科学内涵是什么？", Answer: "第一要义是发展，核心是以人为本，基本要求是全面协调可持续，根本方法是统筹兼顾。", Tags: []string{"毛中特", "科学发展观"}},
			{Question: "习近平新时代中国特色社会主义思想的核心要义是什么？", Answer: "坚持和发展中国特色社会主义。", Tags: []string{"毛中特", "新时代"}},
			{Question: "「四个全面」战略布局是什么？", Answer: "全面建设社会主义现代化国家、全面深化改革、全面依法治国、全面从严治党。", Tags: []string{"毛中特", "新时代"}},
		},
	},
	{
		ID:          "tpl-shi-gang",
		Title:       "史纲时间线",
		Description: "中国近现代史纲要重要历史事件时间线，从鸦片战争到新时代，掌握历史脉络和重大转折点",
		Category:    "考研政治",
		CardCount:   10,
		Icon:        "fire",
		Cards: []TemplateCard{
			{Question: "中国近代史的开端是什么事件？", Answer: "1840年鸦片战争。清政府战败，被迫签订《南京条约》，中国开始沦为半殖民地半封建社会。", Tags: []string{"史纲", "近代开端"}},
			{Question: "太平天国运动的历史意义是什么？", Answer: "沉重打击了封建统治阶级和外国侵略势力，是中国旧式农民战争的最高峰，但农民阶级的局限性导致其最终失败。", Tags: []string{"史纲", "农民运动"}},
			{Question: "洋务运动的指导思想是什么？", Answer: "「中学为体，西学为用」。以曾国藩、李鸿章、左宗棠、张之洞为代表，兴办近代军事工业和民用工业。", Tags: []string{"史纲", "洋务运动"}},
			{Question: "戊戌变法的历史意义是什么？", Answer: "是一次爱国救亡运动、政治改革运动、思想启蒙运动。维新派试图通过改良实现君主立宪，虽失败但促进了思想解放。", Tags: []string{"史纲", "戊戌变法"}},
			{Question: "辛亥革命的历史意义是什么？", Answer: "推翻了清王朝的封建统治，结束了中国两千多年的封建君主专制制度，建立了中华民国，使民主共和观念深入人心。", Tags: []string{"史纲", "辛亥革命"}},
			{Question: "五四运动的历史意义是什么？", Answer: "是一次彻底的反帝反封建的爱国运动，标志着中国新民主主义革命的开端，促进了马克思主义在中国的传播。", Tags: []string{"史纲", "五四运动"}},
			{Question: "遵义会议的历史意义是什么？", Answer: "1935年1月召开。确立了毛泽东在党和红军中的领导地位，挽救了党、挽救了红军、挽救了中国革命，是党的历史上生死攸关的转折点。", Tags: []string{"史纲", "长征"}},
			{Question: "抗日战争胜利的意义是什么？", Answer: "是中国人民一百多年来第一次取得反对帝国主义侵略的完全胜利，为世界反法西斯战争作出了重大贡献，为新民主主义革命的胜利奠定了基础。", Tags: []string{"史纲", "抗日战争"}},
			{Question: "中共七届二中全会的主要内容是什么？", Answer: "1949年3月召开。提出党的工作重心由乡村转移到城市，提出「两个务必」——务必保持谦虚谨慎不骄不躁的作风，务必保持艰苦奋斗的作风。", Tags: []string{"史纲", "解放战争"}},
			{Question: "十一届三中全会的历史意义是什么？", Answer: "1978年12月召开。彻底否定了「两个凡是」的方针，重新确立了解放思想、实事求是的思想路线，作出了改革开放的伟大决策，是新中国成立以来具有深远意义的伟大转折。", Tags: []string{"史纲", "改革开放"}},
		},
	},
	{
		ID:          "tpl-cet4",
		Title:       "四级核心词汇",
		Description: "大学英语四级考试高频核心词汇，涵盖听力、阅读、写作中常见的关键词汇和短语",
		Category:    "英语四级",
		CardCount:   15,
		Icon:        "trophy",
		Cards: []TemplateCard{
			{Question: "abandon", Answer: "v. 放弃；抛弃\n例：He abandoned his plan to go abroad.", Tags: []string{"四级", "高频"}},
			{Question: "available", Answer: "adj. 可获得的；有空的\n例：Is this seat available?", Tags: []string{"四级", "高频"}},
			{Question: "benefit", Answer: "n. 利益；好处  v. 有益于；受益\n短语：benefit from 从中受益", Tags: []string{"四级", "高频"}},
			{Question: "challenge", Answer: "n. 挑战  v. 向……挑战\n短语：face a challenge 面对挑战 / meet the challenge 迎接挑战", Tags: []string{"四级", "高频"}},
			{Question: "consequence", Answer: "n. 结果；后果\n短语：as a consequence 因此 / in consequence of 由于……", Tags: []string{"四级", "高频"}},
			{Question: "decline", Answer: "v. 下降；拒绝  n. 下降\n例：The company's profits declined by 20%.", Tags: []string{"四级", "高频"}},
			{Question: "environment", Answer: "n. 环境\n短语：environmental protection 环境保护", Tags: []string{"四级", "高频"}},
			{Question: "essential", Answer: "adj. 必要的；本质的  n. 必需品\n短语：be essential to/for 对……是必要的", Tags: []string{"四级", "高频"}},
			{Question: "flexible", Answer: "adj. 灵活的；可变通的\n例：We need a more flexible approach.", Tags: []string{"四级", "高频"}},
			{Question: "guarantee", Answer: "v. 保证；担保  n. 保证书\n例：We guarantee the quality of our products.", Tags: []string{"四级", "高频"}},
			{Question: "illustrate", Answer: "v. 阐述；举例说明\n例：Let me illustrate my point with an example.", Tags: []string{"四级", "高频"}},
			{Question: "phenomenon", Answer: "n. 现象（复数：phenomena）\n例：This is a common social phenomenon.", Tags: []string{"四级", "高频"}},
			{Question: "significant", Answer: "adj. 重要的；显著的\n短语：play a significant role in 在……中扮演重要角色", Tags: []string{"四级", "高频"}},
			{Question: "traditional", Answer: "adj. 传统的\n短语：traditional culture 传统文化 / traditional values 传统价值观", Tags: []string{"四级", "高频"}},
			{Question: "various", Answer: "adj. 各种各样的\n同义词：diverse, a variety of", Tags: []string{"四级", "高频"}},
		},
	},
	{
		ID:          "tpl-cs-basics",
		Title:       "计算机基础",
		Description: "计算机基础知识要点，涵盖数据结构、操作系统、计算机网络、数据库等核心概念",
		Category:    "计算机基础",
		CardCount:   12,
		Icon:        "streak",
		Cards: []TemplateCard{
			{Question: "栈（Stack）和队列（Queue）的区别是什么？", Answer: "栈是后进先出（LIFO），只允许在一端进行操作；队列是先进先出（FIFO），一端入队、另一端出队。", Tags: []string{"数据结构"}},
			{Question: "数组和链表的区别是什么？", Answer: "数组：内存连续，支持随机访问O(1)，插入/删除O(n)；链表：内存不连续，不支持随机访问，查找O(n)，插入/删除O(1)。", Tags: []string{"数据结构"}},
			{Question: "哈希表（Hash Table）的基本原理是什么？", Answer: "通过哈希函数将键（Key）映射到表中一个位置来访问记录，以加快查找速度。理想情况下查找、插入、删除的时间复杂度都是O(1)。冲突解决方法：链地址法、开放地址法。", Tags: []string{"数据结构"}},
			{Question: "进程（Process）和线程（Thread）的区别是什么？", Answer: "进程是资源分配的基本单位，拥有独立的地址空间；线程是CPU调度的基本单位，同一进程的线程共享地址空间。线程切换开销更小，进程间通信（IPC）更复杂。", Tags: []string{"操作系统"}},
			{Question: "死锁（Deadlock）的四个必要条件是什么？", Answer: "①互斥条件：资源不能被共享；②请求与保持：已持有资源的进程可以再请求新资源；③不可剥夺：已分配的资源不能被强制抢占；④循环等待：形成进程-资源的循环等待链。", Tags: []string{"操作系统"}},
			{Question: "TCP和UDP的区别是什么？", Answer: "TCP：面向连接、可靠传输、有序、有流量控制和拥塞控制、开销大（三次握手四次挥手）；UDP：无连接、不可靠、无序、开销小、适合实时应用（视频通话、游戏）。", Tags: []string{"计算机网络"}},
			{Question: "HTTP的常用状态码有哪些？", Answer: "200 OK（成功），301 Moved Permanently（永久重定向），302 Found（临时重定向），400 Bad Request（请求错误），401 Unauthorized（未授权），403 Forbidden（禁止访问），404 Not Found（未找到），500 Internal Server Error（服务器内部错误），502 Bad Gateway（网关错误）。", Tags: []string{"计算机网络"}},
			{Question: "三次握手（Three-way Handshake）的过程是什么？", Answer: "①客户端发送SYN包，进入SYN-SENT状态；②服务器回复SYN+ACK包，进入SYN-RCVD状态；③客户端回复ACK包，双方进入ESTABLISHED状态。", Tags: []string{"计算机网络"}},
			{Question: "数据库索引的作用和原理是什么？", Answer: "索引是为了加速查询而创建的数据结构（通常使用B+树）。通过维护索引结构，避免全表扫描，将查询时间复杂度从O(n)降为O(log n)。但索引会增加写入开销和存储空间。", Tags: []string{"数据库"}},
			{Question: "SQL中INNER JOIN、LEFT JOIN、RIGHT JOIN的区别是什么？", Answer: "INNER JOIN：返回两表匹配的行；LEFT JOIN：返回左表所有行，右表无匹配时填充NULL；RIGHT JOIN：返回右表所有行，左表无匹配时填充NULL。", Tags: []string{"数据库"}},
			{Question: "面向对象编程（OOP）的四大特性是什么？", Answer: "①封装（Encapsulation）：隐藏内部实现，仅暴露接口；②继承（Inheritance）：子类继承父类的属性和方法；③多态（Polymorphism）：同一接口在不同对象上有不同表现；④抽象（Abstraction）：提取共性，忽略细节。", Tags: []string{"编程基础"}},
			{Question: "时间复杂度和空间复杂度分别表示什么？", Answer: "时间复杂度衡量算法执行所需的时间随输入规模增长的趋势，用大O表示法（如O(1), O(n), O(n²)）；空间复杂度衡量算法执行所需的存储空间随输入规模增长的趋势。", Tags: []string{"算法"}},
		},
	},
}

// ListTemplates returns all available template decks (no auth required for listing).
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	// Return summaries (without card details to keep response thin)
	type TemplateSummary struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		CardCount   int    `json:"card_count"`
		Icon        string `json:"icon"`
	}
	summaries := make([]TemplateSummary, len(templateDecks))
	for i, t := range templateDecks {
		summaries[i] = TemplateSummary{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Category:    t.Category,
			CardCount:   t.CardCount,
			Icon:        t.Icon,
		}
	}
	writeJSON(w, 200, summaries)
}

// GetTemplate returns a single template with its cards.
func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for _, t := range templateDecks {
		if t.ID == id {
			writeJSON(w, 200, t)
			return
		}
	}
	writeError(w, 404, "template not found")
}

// ImportTemplate copies a template deck (with all its cards) into the user's own decks.
func (h *Handler) ImportTemplate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	templateID := chi.URLParam(r, "id")

	var tmpl *TemplateDeck
	for i := range templateDecks {
		if templateDecks[i].ID == templateID {
			tmpl = &templateDecks[i]
			break
		}
	}
	if tmpl == nil {
		writeError(w, 404, "template not found")
		return
	}

	// Create a new deck for the user
	var deck struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		CardCount int    `json:"card_count"`
	}

	err := h.DB.Raw(`
		INSERT INTO decks (user_id, title, source) VALUES (?, ?, 'template')
		RETURNING id, title, card_count
	`, userID, tmpl.Title).Scan(&deck).Error

	if err != nil {
		writeError(w, 500, "failed to create deck from template")
		return
	}

	// Insert all template cards into the new deck
	for _, card := range tmpl.Cards {
		tagsJSON := "[]"
		if len(card.Tags) > 0 {
			b, _ := json.Marshal(card.Tags)
			tagsJSON = string(b)
		}
		h.DB.Exec(`
			INSERT INTO cards (deck_id, question, answer, tags, next_review_at)
			VALUES (?, ?, ?, ?, NOW())
		`, deck.ID, card.Question, card.Answer, tagsJSON)
	}

	// Update card count
	h.DB.Exec(`UPDATE decks SET card_count = ?, updated_at = NOW() WHERE id = ?`,
		len(tmpl.Cards), deck.ID)

	deck.CardCount = len(tmpl.Cards)
	writeJSON(w, 201, deck)
}
