const app = getApp();

Page({
  data: {
    nickname: 'User',
    email: '',
    userInitial: 'U',
    totalDecks: 0,
    totalCards: 0,
    reviewedTotal: 0,
    accountMenus: [
      {
        title: 'My Subscription',
        desc: 'Manage plan',
        geoStyle: 'geo-diamond',
        color: 'cyan',
        action: 'subscribe'
      },
      {
        title: 'Study History',
        desc: 'View review log',
        geoStyle: 'geo-bars',
        color: 'green',
        action: 'history'
      }
    ],
    generalMenus: [
      {
        title: 'Notifications',
        desc: 'Review reminders',
        geoStyle: 'geo-bell',
        color: 'purple',
        action: 'notifications'
      },
      {
        title: 'About FlashCard',
        desc: 'v1.0.0',
        geoStyle: 'geo-info',
        color: 'blue',
        action: 'about'
      },
      {
        title: 'Feedback',
        desc: 'Help us improve',
        geoStyle: 'geo-chat',
        color: 'orange',
        action: 'feedback'
      }
    ]
  },

  onShow() {
    this.loadUserData();
  },

  loadUserData() {
    const userInfo = wx.getStorageSync('userInfo');
    if (userInfo) {
      this.setData({
        nickname: userInfo.nickName || 'User',
        email: userInfo.email || '',
        userInitial: (userInfo.nickName || 'U')[0]
      });
    }

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
          title: 'FlashCard',
          content: 'AI-powered smart memory tool. Upload text, and AI analyzes and generates flashcards based on spaced repetition algorithms.',
          showCancel: false
        });
        break;
      case 'feedback':
        wx.showModal({
          title: 'Feedback',
          content: 'Email us at feedback@flashcard.app',
          showCancel: false
        });
        break;
      default:
        wx.showToast({ title: 'Coming soon', icon: 'none' });
    }
  }
});
