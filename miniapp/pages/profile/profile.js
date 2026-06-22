const app = getApp();

Page({
  data: {
    nickname: '闪卡用户',
    avatarUrl: '',
    userInitial: '闪',
    totalDecks: 0,
    totalCards: 0,
    reviewedTotal: 0,
    reviewReminder: false,
    achievements: [],
    inviteCode: '',
    invitedCount: 0
  },

  onShow() {
    this.loadUserData();
    this.loadAchievements();
    this.loadInviteInfo(false);
    wx.showShareMenu({ withShareTicket: true, menus: ['shareAppMessage'] });
    const reviewReminder = wx.getStorageSync('reviewReminder') || false;
    this.setData({ reviewReminder });
  },

  loadUserData() {
    const token = app.globalData.token || wx.getStorageSync('token');
    const userInfo = wx.getStorageSync('userInfo');
    if (userInfo && userInfo.nickName) {
      this.setData({
        nickname: userInfo.nickName,
        userInitial: (userInfo.nickName || '闪')[0],
        avatarUrl: userInfo.avatarUrl || ''
      });
    }

    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const decks = res.data || [];
        const cards = decks.reduce((s, d) => s + (d.card_count || 0), 0);
        this.setData({ totalDecks: decks.length, totalCards: cards });
      }
    });

    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const data = res.data || {};
        this.setData({ reviewedTotal: data.reviewed_today || 0 });
      }
    });
  },

  // WeChat native avatar picker
  onChooseAvatar(e) {
    const avatarUrl = e.detail.avatarUrl;
    this.setData({ avatarUrl });
    this.saveProfile({ avatarUrl });
  },

  // WeChat native nickname input
  onNicknameBlur(e) {
    const nickName = e.detail.value;
    if (!nickName) return;
    this.setData({
      nickname: nickName,
      userInitial: nickName[0]
    });
    this.saveProfile({ nickName });
  },

  saveProfile(updates) {
    const userInfo = wx.getStorageSync('userInfo') || {};
    Object.assign(userInfo, updates);
    if (updates.nickName) {
      userInfo.nickName = updates.nickName;
    }
    if (updates.avatarUrl) {
      userInfo.avatarUrl = updates.avatarUrl;
    }
    wx.setStorageSync('userInfo', userInfo);

    wx.request({
      url: app.globalData.apiBase + '/api/user/profile',
      method: 'PUT',
      header: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token'))
      },
      data: { nickname: userInfo.nickName || '闪卡用户', avatar_url: userInfo.avatarUrl || '' }
    });
  },

  loadAchievements() {
    wx.request({
      url: app.globalData.apiBase + '/api/achievements',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const achievements = (res.data && res.data.achievements) || [];
        this.setData({ achievements });
      }
    });
  },

  toggleReviewReminder(e) {
    const enabled = e.detail.value;
    this.setData({ reviewReminder: enabled });
    wx.setStorageSync('reviewReminder', enabled);
    if (enabled) {
      wx.showModal({
        title: '复习提醒已开启',
        content: '每天 20:00 将通过微信服务通知提醒你复习',
        showCancel: false, confirmText: '知道了'
      });
    } else {
      wx.showToast({ title: '复习提醒已关闭', icon: 'none' });
    }
  },

  onMenuTap(e) {
    const action = e.currentTarget.dataset.action;
    if (action === 'about') {
      wx.showModal({
        title: '关于闪卡记忆',
        content: 'AI 驱动的智能记忆工具。上传文本，AI 自动生成问答闪卡，基于间隔重复算法安排每日复习。',
        showCancel: false, confirmText: '知道了'
      });
    } else if (action === 'feedback') {
      wx.showModal({
        title: '意见反馈',
        content: '感谢你的反馈！请发送邮件至：feedback@flashcard.app',
        showCancel: false, confirmText: '好的'
      });
    } else {
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
        wx.setClipboardData({
          data: JSON.stringify(res.data, null, 2),
          success: () => wx.showToast({ title: '已复制到剪贴板', icon: 'success' })
        });
      },
      fail: () => {
        wx.hideLoading();
        wx.showToast({ title: '导出失败', icon: 'none' });
      }
    });
  },

  loadInviteInfo(showErrors) {
    wx.request({
      url: app.globalData.apiBase + '/api/invite/my-code',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const inviteCode = (res.data && res.data.invite_code) || '';
        if (!inviteCode) { if (showErrors) wx.showToast({ title: '获取邀请码失败', icon: 'none' }); return; }
        this.setData({ inviteCode });
        wx.request({
          url: app.globalData.apiBase + '/api/invite/stats',
          header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
          success: (sr) => this.setData({ invitedCount: (sr.data && sr.data.invited_count) || 0 })
        });
      },
      fail: (e) => { if (showErrors) wx.showToast({ title: '网络异常', icon: 'none' }); }
    });
  },

  showInvite() {
    if (!this.data.inviteCode) { wx.showLoading({ title: '加载中...' }); this.loadInviteInfo(true); setTimeout(() => wx.hideLoading(), 800); return; }
    const text = '来闪卡记忆一起高效学习吧！\n我的邀请码：' + this.data.inviteCode;
    wx.setClipboardData({ data: text, success: () => wx.showToast({ title: '已复制', icon: 'success' }) });
  },

  onShareAppMessage() {
    return { title: '来闪卡记忆一起高效学习吧', path: '/pages/onboard/onboard?invite_code=' + (this.data.inviteCode || 'FLASH') };
  }
});
