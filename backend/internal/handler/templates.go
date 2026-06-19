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
		CardCount:   25,
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
			{Question: "经济基础和上层建筑的辩证关系是什么？", Answer: "经济基础决定上层建筑，上层建筑反作用于经济基础。上层建筑一定要适合经济基础发展的要求。", Tags: []string{"马原", "唯物史观"}},
			{Question: "意识的本质是什么？", Answer: "意识是物质世界长期发展的产物，是人脑的机能和属性，是客观世界在人脑中的主观映像。", Tags: []string{"马原", "唯物论"}},
			{Question: "运动和静止的辩证关系是什么？", Answer: "运动是绝对的、无条件的，静止是相对的、有条件的。动中有静，静中有动，物质世界是绝对运动和相对静止的统一。", Tags: []string{"马原", "唯物论"}},
			{Question: "矛盾的普遍性和特殊性的关系是什么？", Answer: "普遍性寓于特殊性之中并通过特殊性表现出来；特殊性包含普遍性。二者相互联结、相互转化。矛盾的普遍性和特殊性辩证关系是矛盾问题的精髓。", Tags: []string{"马原", "辩证法"}},
			{Question: "原因和结果的辩证关系是什么？", Answer: "原因和结果是相互依存、相互作用、相互转化的。原因引起结果，结果又反作用于原因。二者的区分既是确定的又是不确定的。", Tags: []string{"马原", "辩证法"}},
			{Question: "必然性和偶然性的辩证关系是什么？", Answer: "必然性居于支配地位决定事物发展方向，偶然性居于从属地位加速或延缓发展。必然性通过偶然性为自己开辟道路，偶然性是必然性的表现和补充。", Tags: []string{"马原", "辩证法"}},
			{Question: "实践的基本特征是什么？", Answer: "①实践是客观的物质性活动（直接现实性）；②实践是自觉的能动性活动（主观能动性）；③实践是社会历史性活动（社会历史性）。", Tags: []string{"马原", "认识论"}},
			{Question: "真理的客观性指什么？", Answer: "真理的内容是客观的；检验真理的标准——实践——也是客观的。真理中不包含同客观实际相违背的主观成分。", Tags: []string{"马原", "认识论"}},
			{Question: "实践标准的确定性和不确定性是什么？", Answer: "确定性（绝对性）：实践是检验真理的唯一标准，一切认识最终都要经过实践检验。不确定性（相对性）：一定历史阶段的实践具有局限性，不能完全证实或驳倒一切认识。", Tags: []string{"马原", "认识论"}},
			{Question: "社会基本矛盾是什么？", Answer: "生产力和生产关系的矛盾、经济基础和上层建筑的矛盾是社会基本矛盾，是社会发展的根本动力。这两对矛盾贯穿人类社会发展始终。", Tags: []string{"马原", "唯物史观"}},
			{Question: "人民群众在历史中的作用是什么？", Answer: "人民群众是历史的创造者：是社会物质财富的创造者，是社会精神财富的创造者，是社会变革的决定力量。", Tags: []string{"马原", "唯物史观"}},
			{Question: "商品二因素是什么？", Answer: "使用价值和价值。使用价值是商品的自然属性，构成社会财富的物质内容；价值是商品的社会属性，是凝结在商品中的一般人类劳动。二者对立统一。", Tags: []string{"马原", "政治经济学"}},
			{Question: "剩余价值是什么？", Answer: "剩余价值是由雇佣工人创造的、被资本家无偿占有的、超过劳动力价值以上的那部分价值。它是资本主义生产关系的本质体现，反映了资本家对工人的剥削关系。", Tags: []string{"马原", "政治经济学"}},
		},
	},
	{
		ID:          "tpl-mao-zhong-te",
		Title:       "毛中特核心",
		Description: "毛泽东思想和中国特色社会主义理论体系概论重点考点，涵盖新民主主义革命理论、社会主义改造、邓小平理论等",
		Category:    "考研政治",
		CardCount:   20,
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
			{Question: "新民主主义革命的性质是什么？", Answer: "新民主主义革命是无产阶级领导的资产阶级民主革命。其革命对象是帝国主义、封建主义和官僚资本主义，革命动力包括工人、农民、小资产阶级和民族资产阶级。", Tags: []string{"毛中特", "新民主主义革命"}},
			{Question: "农村包围城市、武装夺取政权道路的必要性是什么？", Answer: "①中国是半殖民地半封建国家，内部没有民主制度，外部没有民族独立；②农民占人口大多数，是中国革命的主力军；③敌人长期占据中心城市，农村是其统治的薄弱环节。", Tags: []string{"毛中特", "新民主主义革命"}},
			{Question: "社会主义改造的历史经验是什么？", Answer: "①坚持社会主义工业化与社会主义改造同时并举；②采取积极引导、逐步过渡的方式；③用和平方法进行改造。", Tags: []string{"毛中特", "社会主义改造"}},
			{Question: "社会主义初级阶段的基本路线是什么？", Answer: "领导和团结全国各族人民，以经济建设为中心，坚持四项基本原则，坚持改革开放，自力更生，艰苦创业，为把我国建设成为富强民主文明和谐美丽的社会主义现代化强国而奋斗。（「一个中心、两个基本点」）", Tags: []string{"毛中特", "初级阶段"}},
			{Question: "「一国两制」的基本内容是什么？", Answer: "在一个中国的前提下，国家主体坚持社会主义制度，香港、澳门、台湾保持原有的资本主义制度长期不变。核心是实现祖国统一。", Tags: []string{"毛中特", "邓小平理论"}},
			{Question: "「五位一体」总体布局是什么？", Answer: "经济建设、政治建设、文化建设、社会建设、生态文明建设五位一体，全面推进。", Tags: []string{"毛中特", "新时代"}},
			{Question: "全面深化改革的总目标是什么？", Answer: "完善和发展中国特色社会主义制度，推进国家治理体系和治理能力现代化。", Tags: []string{"毛中特", "新时代"}},
			{Question: "新发展理念（五大发展理念）是什么？", Answer: "创新、协调、绿色、开放、共享。创新是引领发展的第一动力；协调是持续健康发展的内在要求；绿色是永续发展的必要条件；开放是国家繁荣发展的必由之路；共享是中国特色社会主义的本质要求。", Tags: []string{"毛中特", "新时代"}},
			{Question: "人类命运共同体的核心内涵是什么？", Answer: "建设持久和平、普遍安全、共同繁荣、开放包容、清洁美丽的世界。坚持对话协商、共建共享、合作共赢、交流互鉴、绿色低碳。", Tags: []string{"毛中特", "新时代"}},
			{Question: "中国特色社会主义最本质的特征是什么？", Answer: "中国共产党的领导是中国特色社会主义最本质的特征，是中国特色社会主义制度的最大优势。党是最高政治领导力量。", Tags: []string{"毛中特", "新时代"}},
		},
	},
	{
		ID:          "tpl-shi-gang",
		Title:       "史纲时间线",
		Description: "中国近现代史纲要重要历史事件时间线，从鸦片战争到新时代，掌握历史脉络和重大转折点",
		Category:    "考研政治",
		CardCount:   20,
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
			{Question: "《南京条约》的主要内容是什么？", Answer: "1842年签订。①割让香港岛给英国；②赔款2100万银元；③开放广州、厦门、福州、宁波、上海五处为通商口岸；④协定关税。中国开始丧失领土和主权。", Tags: []string{"史纲", "近代开端"}},
			{Question: "义和团运动的历史意义是什么？", Answer: "1900年爆发。显示了中国人民反抗外来侵略的坚强意志，沉重打击了帝国主义瓜分中国的野心。但由于农民阶级的局限性，最终在中外联合镇压下失败。", Tags: []string{"史纲", "农民运动"}},
			{Question: "新文化运动的主要内容和意义是什么？", Answer: "1915年陈独秀创办《青年杂志》为开端。提倡民主与科学（「德先生」和「赛先生」），反对封建专制和迷信盲从。是一次空前的思想解放运动，为马克思主义传播创造了条件。", Tags: []string{"史纲", "新文化运动"}},
			{Question: "中国共产党的成立及其意义是什么？", Answer: "1921年7月，中共一大在上海召开。中国共产党的成立是中国历史上「开天辟地的大事变」，从此中国革命有了坚强的领导核心、科学的指导思想和崭新的奋斗目标。", Tags: []string{"史纲", "建党"}},
			{Question: "南昌起义的历史意义是什么？", Answer: "1927年8月1日，周恩来、贺龙、叶挺、朱德、刘伯承等领导。打响了武装反抗国民党反动派的第一枪，标志着中国共产党独立领导革命战争、创建人民军队和武装夺取政权的开始。", Tags: []string{"史纲", "武装革命"}},
			{Question: "长征胜利的历史意义是什么？", Answer: "1934年10月至1936年10月。粉碎了国民党反动派的围追堵截，保存了党和红军的基干力量，使中国革命转危为安。铸就了伟大的长征精神。", Tags: []string{"史纲", "长征"}},
			{Question: "开国大典（新中国成立）的历史意义是什么？", Answer: "1949年10月1日。结束了帝国主义、封建主义和官僚资本主义在中国的统治，中国人民从此站起来了。改变了世界政治力量的对比，鼓舞了世界被压迫民族的解放斗争。", Tags: []string{"史纲", "建国"}},
			{Question: "抗美援朝战争的意义是什么？", Answer: "1950-1953年。保卫了国家安全，捍卫了世界和平，打破了美军不可战胜的神话，极大地提高了新中国的国际威望，为国内经济建设赢得了相对稳定的和平环境。", Tags: []string{"史纲", "抗美援朝"}},
			{Question: "南方谈话（邓小平南方谈话）的核心内容是什么？", Answer: "1992年春。深刻回答了「什么是社会主义、怎样建设社会主义」的问题。提出：①发展才是硬道理；②三个有利于标准；③坚持党的基本路线一百年不动摇。标志着改革开放进入新阶段。", Tags: []string{"史纲", "改革开放"}},
			{Question: "中国特色社会主义进入新时代的标志是什么？", Answer: "2017年党的十九大召开，将习近平新时代中国特色社会主义思想确立为党的指导思想。我国社会主要矛盾转化为人民日益增长的美好生活需要和不平衡不充分的发展之间的矛盾。", Tags: []string{"史纲", "新时代"}},
		},
	},
	{
		ID:          "tpl-cet4",
		Title:       "四级核心词汇",
		Description: "大学英语四级考试高频核心词汇，涵盖听力、阅读、写作中常见的关键词汇和短语",
		Category:    "英语四级",
		CardCount:   40,
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
			{Question: "approach", Answer: "n. 方法；途径  v. 接近；处理\n短语：take a new approach 采用新方法 / approach the problem 处理问题", Tags: []string{"四级", "高频"}},
			{Question: "assume", Answer: "v. 假设；承担\n例：We cannot assume the result. / She assumed responsibility for the project.", Tags: []string{"四级", "高频"}},
			{Question: "circumstance", Answer: "n. 环境；情况（常用复数）\n短语：under no circumstances 无论如何都不 / under the circumstances 在这种情况下", Tags: []string{"四级", "高频"}},
			{Question: "contribute", Answer: "v. 贡献；捐献；促成\n短语：contribute to 有助于；捐献 / make a contribution to 为……做贡献", Tags: []string{"四级", "高频"}},
			{Question: "demonstrate", Answer: "v. 展示；证明；示威游行\n例：The experiment demonstrates the theory.", Tags: []string{"四级", "高频"}},
			{Question: "establish", Answer: "v. 建立；确立\n短语：establish a company 创办公司 / establish a relationship 建立关系", Tags: []string{"四级", "高频"}},
			{Question: "identify", Answer: "v. 确认；识别；认同\n短语：identify with 认同 / identify the cause 找出原因", Tags: []string{"四级", "高频"}},
			{Question: "maintain", Answer: "v. 维持；保持；维修\n短语：maintain order 维持秩序 / maintain good health 保持健康", Tags: []string{"四级", "高频"}},
			{Question: "obvious", Answer: "adj. 明显的；显然的\n同义词：apparent, evident, clear\n例：It is obvious that he is lying.", Tags: []string{"四级", "高频"}},
			{Question: "perspective", Answer: "n. 观点；视角；前景\n短语：from a different perspective 从不同角度 / in perspective 正确地", Tags: []string{"四级", "高频"}},
			{Question: "relevant", Answer: "adj. 相关的；切题的\n短语：be relevant to 与……相关\n反义词：irrelevant", Tags: []string{"四级", "高频"}},
			{Question: "strategy", Answer: "n. 策略；战略\n短语：develop a strategy 制定策略 / marketing strategy 营销策略", Tags: []string{"四级", "高频"}},
			{Question: "concentrate", Answer: "v. 集中；专注\n短语：concentrate on 专注于 / concentration camp 集中营", Tags: []string{"四级", "高频"}},
			{Question: "determine", Answer: "v. 决定；确定；下决心\n短语：be determined to do 决心做…… / determine the cause 确定原因", Tags: []string{"四级", "高频"}},
			{Question: "efficient", Answer: "adj. 高效的\n短语：energy efficient 节能的 / an efficient way 高效的方式\n区别：effective 有效的≠efficient 高效的", Tags: []string{"四级", "高频"}},
			{Question: "fundamental", Answer: "adj. 基本的；根本的  n. 基本原则\n同义词：basic, essential\n例：Freedom of speech is a fundamental right.", Tags: []string{"四级", "高频"}},
			{Question: "independent", Answer: "adj. 独立的；自主的\n短语：be independent of 不依赖…… / independent thinking 独立思考\n反义词：dependent", Tags: []string{"四级", "高频"}},
			{Question: "motivate", Answer: "v. 激励；激发\n短语：motivate sb. to do 激励某人做…… / stay motivated 保持积极性", Tags: []string{"四级", "高频"}},
			{Question: "participate", Answer: "v. 参与；参加\n短语：participate in 参加 / active participation 积极参与", Tags: []string{"四级", "高频"}},
			{Question: "recommend", Answer: "v. 推荐；建议\n短语：highly recommend 强烈推荐 / recommend doing 建议做……", Tags: []string{"四级", "高频"}},
			{Question: "sufficient", Answer: "adj. 足够的；充分的\n短语：sufficient evidence 充分证据 / be sufficient for 对……足够\n反义词：insufficient", Tags: []string{"四级", "高频"}},
			{Question: "accomplish", Answer: "v. 完成；达成\n短语：accomplish a task 完成任务 / mission accomplished 任务完成", Tags: []string{"四级", "高频"}},
			{Question: "confident", Answer: "adj. 自信的；有信心的\n短语：be confident of/about 对……有信心 / overconfident 过分自信的", Tags: []string{"四级", "高频"}},
			{Question: "essential", Answer: "adj. 必要的；本质的  n. 必需品\n短语：be essential to/for 对……是必要的", Tags: []string{"四级", "高频"}},
			{Question: "opportunity", Answer: "n. 机会；时机\n短语：seize the opportunity 抓住机会 / equal opportunity 机会均等", Tags: []string{"四级", "高频"}},
		},
	},
	{
		ID:          "tpl-cs-basics",
		Title:       "计算机基础",
		Description: "计算机基础知识要点，涵盖数据结构、操作系统、计算机网络、数据库等核心概念",
		Category:    "计算机基础",
		CardCount:   25,
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
			{Question: "二叉树（Binary Tree）的三种遍历方式是什么？", Answer: "①前序遍历：根→左→右；②中序遍历：左→根→右（对于二叉搜索树可得到有序序列）；③后序遍历：左→右→根。每种遍历都可以用递归或迭代方式实现。", Tags: []string{"数据结构"}},
			{Question: "堆（Heap）和二叉搜索树（BST）的区别是什么？", Answer: "堆：父节点 ≤（或≥）子节点，不保证左右子树有序，主要用于优先队列和堆排序，查找最值O(1)；BST：左子树 < 根 < 右子树，中序遍历得到有序序列，查找/插入/删除平均O(log n)。", Tags: []string{"数据结构"}},
			{Question: "虚拟内存（Virtual Memory）的作用是什么？", Answer: "①为每个进程提供独立的地址空间，实现进程隔离；②使得程序可以使用比物理内存更大的地址空间；③通过页面置换算法（LRU、FIFO等）管理物理内存和磁盘之间的页面交换。", Tags: []string{"操作系统"}},
			{Question: "OSI七层模型是什么？", Answer: "物理层→数据链路层→网络层→传输层→会话层→表示层→应用层。自上而下：应用层（HTTP、FTP），传输层（TCP、UDP），网络层（IP），数据链路层（MAC），物理层（电缆、光纤）。", Tags: []string{"计算机网络"}},
			{Question: "DNS的作用和工作原理是什么？", Answer: "DNS（域名系统）将域名转换为IP地址。工作原理：客户端查询本地DNS缓存→递归查询本地DNS服务器→迭代查询根域名服务器→顶级域服务器→权威DNS服务器，逐级解析获取IP地址。", Tags: []string{"计算机网络"}},
			{Question: "数据库事务（Transaction）的ACID特性是什么？", Answer: "①原子性（Atomicity）：事务中的操作要么全做要么全不做；②一致性（Consistency）：事务执行前后数据库保持一致状态；③隔离性（Isolation）：并发事务之间互不干扰；④持久性（Durability）：已提交的事务结果永久保存。", Tags: []string{"数据库"}},
			{Question: "SQL注入（SQL Injection）是什么？如何防范？", Answer: "SQL注入是攻击者通过构造恶意SQL语句来操控数据库的攻击方式。防范措施：①使用参数化查询（PreparedStatement）；②输入验证和过滤；③最小权限原则；④使用ORM框架。", Tags: []string{"数据库"}},
			{Question: "快速排序（Quick Sort）的基本原理和时间复杂度？", Answer: "选取基准元素（pivot），将数组分为小于和大于基准的两部分，递归排序子数组。平均时间O(n log n)，最坏O(n²)（当每次选到最小或最大元素），空间O(log n)。通过随机化选择pivot可避免最坏情况。", Tags: []string{"算法"}},
			{Question: "动态规划（Dynamic Programming）的核心思想是什么？", Answer: "将复杂问题分解为重叠子问题，保存子问题的解避免重复计算。两个核心要素：①最优子结构（问题最优解包含子问题最优解）；②重叠子问题（子问题被多次计算）。经典应用：背包问题、最长公共子序列。", Tags: []string{"算法"}},
			{Question: "RESTful API的设计原则是什么？", Answer: "①资源（Resource）通过URL唯一标识；②使用标准HTTP方法（GET查询、POST创建、PUT全量更新、PATCH部分更新、DELETE删除）；③无状态（每个请求包含全部所需信息）；④使用标准HTTP状态码；⑤支持多种表示格式（JSON、XML）。", Tags: []string{"编程基础"}},
			{Question: "什么是缓存穿透、缓存击穿和缓存雪崩？", Answer: "缓存穿透：查询不存在的数据，请求直达数据库（解决方案：布隆过滤器、缓存空值）；缓存击穿：热点key过期瞬间大量请求到数据库（解决方案：互斥锁、永不过期）；缓存雪崩：大量key同时过期导致数据库崩溃（解决方案：过期时间加随机值、多级缓存）。", Tags: []string{"编程基础"}},
			{Question: "HTTPS的工作原理是什么？", Answer: "HTTPS = HTTP + SSL/TLS。流程：①客户端发送支持的加密算法；②服务器返回证书（包含公钥）；③客户端验证证书，生成对称密钥并用公钥加密发送；④服务器用私钥解密获得对称密钥；⑤双方使用对称密钥加密通信。实现了机密性、完整性和身份认证。", Tags: []string{"计算机网络"}},
			{Question: "什么是CAP定理？", Answer: "CAP定理（Brewer定理）：分布式系统不能同时满足一致性（Consistency）、可用性（Availability）和分区容错性（Partition Tolerance）三个特性，最多只能同时满足两个。实际系统中必须选择P，因此是在CP和AP之间取舍。", Tags: []string{"编程基础"}},
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
