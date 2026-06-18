const app = getApp();

Page({
  data: {
    nickname: '闪卡用户',
    email: '',
    userInitial: '闪',
    totalDecks: 0,
    totalCards: 0,
    reviewedTotal: 0,
    accountMenus: [
      {
        title: '我的订阅',
        desc: '管理订阅计划',
        icon: '💎',
        color: 'cyan',
        action: 'subscribe'
      },
      {
        title: '学习记录',
        desc: '查看复习历史',
        icon: '📊',
        color: 'green',
        action: 'history'
      }
    ],
    generalMenus: [
      {
        title: '通知设置',
        desc: '复习提醒频率',
        icon: '🔔',
        color: 'purple',
        action: 'notifications'
      },
      {
        title: '关于闪卡记忆',
        desc: 'v1.0.0',
        icon: 'ℹ️',
        color: 'blue',
        action: 'about'
      },
      {
        title: '反馈建议',
        desc: '帮助我们改进',
        icon: '💬',
        color: 'orange',
        action: 'feedback'
      }
    ]
  },

  onShow() {
    this.loadUserData();
  },

  loadUserData() {
    // Try to load user info from storage
    const userInfo = wx.getStorageSync('userInfo');
    if (userInfo) {
      this.setData({
        nickname: userInfo.nickName || '闪卡用户',
        email: userInfo.email || '',
        userInitial: (userInfo.nickName || '闪')[0]
      });
    }

    // Fetch stats
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const decks = res.data || [];
        this.setData({
          totalDecks: decks.length,
          totalCards: decks.reduce((s, d) => s + (d.card_count || 0), 0)
        });
      }
    });

    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        this.setData({
          reviewedTotal: res.data.reviewed_total || 0
        });
      }
    });
  },

  onMenuTap(e) {
    const action = e.currentTarget.dataset.action;

    switch (action) {
      case 'about':
        wx.showModal({
          title: '闪卡记忆',
          content: 'AI 驱动的智能记忆工具。上传文本，AI 自动分析并生成闪卡，基于间隔重复算法帮你高效记忆。',
          showCancel: false
        });
        break;
      case 'feedback':
        wx.showModal({
          title: '反馈建议',
          content: '请发送邮件至 feedback@flashcard.app',
          showCancel: false
        });
        break;
      default:
        wx.showToast({ title: '功能开发中', icon: 'none' });
    }
  }
});
