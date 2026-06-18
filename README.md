# 闪卡记忆引擎

上传文本/PDF → AI 自动生成闪卡 → FSRS 算法智能调度复习

## 技术栈

- **后端：** Go (gin) — API + FSRS 引擎
- **AI 服务：** Python (FastAPI) — DeepSeek 卡片生成
- **前端：** 微信小程序
- **数据库：** PostgreSQL + Redis

## 项目结构

```
flashcard/
├── backend/           # Go 后端
│   └── internal/
│       ├── ai/        # AI 服务客户端
│       ├── config/    # 配置
│       ├── fsrs/      # FSRS 间隔重复算法
│       ├── handler/   # HTTP handlers
│       ├── middleware/ # 中间件
│       └── model/     # 数据模型
├── ai-service/        # Python AI 服务
├── miniapp/           # 微信小程序
│   └── pages/
│       ├── index/     # 首页
│       ├── decks/     # 牌组列表
│       ├── deck-detail/ # 牌组详情
│       ├── review/    # 复习界面
│       ├── create/    # 创建牌组
│       └── profile/   # 个人中心
└── docs/              # 文档
```

## 开发状态

- [x] Spike 001：FSRS 算法 + AI 卡片生成验证通过
- [ ] Phase 0：项目初始化
- [ ] Phase 1：Go 后端
- [ ] Phase 2：Python AI 服务
- [ ] Phase 3：小程序前端
- [ ] Phase 4：支付 + 集成
- [ ] Phase 5：部署上线
