const app = getApp();

Page({
  data: {
    nickname: '闪卡用户',
    userInitial: '闪',
    totalDecks: 0,
    totalCards: 0,
    reviewedTotal: 0,
    reviewReminder: false,
    achievements: []
  },

  onShow() {
    this.loadUserData();
    this.loadAchievements();
    // Load reminder preference
    const reviewReminder = wx.getStorageSync('reviewReminder') || false;
    this.setData({ reviewReminder });
  },

  loadUserData() {
    // Load from local storage
    const userInfo = wx.getStorageSync('userInfo');
    if (userInfo && userInfo.nickName) {
      const name = userInfo.nickName;
      this.setData({
        nickname: name,
        userInitial: name[0] || '闪'
      });
    }

    // Load deck stats
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const decks = res.data || [];
        const cards = decks.reduce((s, d) => s + (d.card_count || 0), 0);
        this.setData({ totalDecks: decks.length, totalCards: cards });
      },
      fail: () => {
        // Offline — show cached data
        const cachedDecks = wx.getStorageSync('cachedDecks') || [];
        const cachedCards = cachedDecks.reduce((s, d) => s + (d.card_count || 0), 0);
        this.setData({ totalDecks: cachedDecks.length, totalCards: cachedCards });
      }
    });

    // Load review stats
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const total = res.data.reviewed_total || res.data.total || 0;
        this.setData({ reviewedTotal: total });
      }
    });
  },

  loadAchievements() {
    wx.request({
      url: app.globalData.apiBase + '/api/achievements',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const achievements = (res.data && res.data.achievements) || [];
        this.setData({ achievements });
      },
      fail: () => {
        // Show cached or empty
      }
    });
  },

  toggleReviewReminder(e) {
    const enabled = e.detail.value;
    this.setData({ reviewReminder: enabled });
    wx.setStorageSync('reviewReminder', enabled);

    if (enabled) {
      // Mock notification — show a modal since real push needs backend
      wx.showModal({
        title: '复习提醒已开启',
        content: '每天 20:00 将通过微信服务通知提醒你复习（模拟提示：实际推送需后端支持）。开启后，每日定时收到复习提醒。',
        showCancel: false,
        confirmText: '知道了'
      });
    } else {
      wx.showToast({ title: '复习提醒已关闭', icon: 'none' });
    }
  },

  onMenuTap(e) {
    const action = e.currentTarget.dataset.action;
    switch (action) {
      case 'about':
        wx.showModal({
          title: '关于闪卡记忆',
          content: 'AI 驱动的智能记忆工具。上传文本，AI 自动生成问答闪卡，基于间隔重复算法安排每日复习，帮助你高效记忆任何内容。',
          showCancel: false,
          confirmText: '知道了'
        });
        break;
      case 'feedback':
        wx.showModal({
          title: '意见反馈',
          content: '感谢你的反馈！请发送邮件至：feedback@flashcard.app',
          showCancel: false,
          confirmText: '好的'
        });
        break;
      default:
        wx.showToast({ title: '即将上线', icon: 'none' });
    }
  },

  exportData() {
    wx.showLoading({ title: '导出中...' });
    wx.request({
      url: app.globalData.apiBase + '/api/export',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        wx.hideLoading();
        const data = JSON.stringify(res.data, null, 2);
        wx.setClipboardData({
          data,
          success: () => {
            wx.showToast({ title: '已复制到剪贴板', icon: 'success' });
          }
        });
      },
      fail: () => {
        wx.hideLoading();
        wx.showToast({ title: '导出失败', icon: 'none' });
      }
    });
  },

  showInvite() {
    wx.showLoading({ title: '加载中...' });
    wx.request({
      url: app.globalData.apiBase + '/api/invite/my-code',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        wx.hideLoading();
        const inviteCode = (res.data && res.data.invite_code) || '';
        if (!inviteCode) {
          wx.showToast({ title: '获取邀请码失败', icon: 'none' });
          return;
        }

        // Load invite stats
        wx.request({
          url: app.globalData.apiBase + '/api/invite/stats',
          header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
          success: (statsRes) => {
            const invitedCount = (statsRes.data && statsRes.data.invited_count) || 0;
            const shareText = '来闪卡记忆一起高效学习吧！\n我的邀请码：' + inviteCode + '\n已邀请 ' + invitedCount + ' 位好友加入';

            wx.showModal({
              title: '邀请好友',
              content: shareText,
              confirmText: '复制邀请语',
              cancelText: '关闭',
              success: (modalRes) => {
                if (modalRes.confirm) {
                  wx.setClipboardData({
                    data: shareText,
                    success: () => {
                      wx.showToast({ title: '已复制到剪贴板', icon: 'success' });
                    }
                  });
                }
              }
            });
          },
          fail: () => {
            wx.showToast({ title: '加载邀请统计失败', icon: 'none' });
          }
        });
      },
      fail: () => {
        wx.hideLoading();
        wx.showToast({ title: '网络请求失败', icon: 'none' });
      }
    });
  }
});
