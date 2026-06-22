const app = getApp();

Page({
  data: {
    nickname: '点击设置昵称',
    avatarUrl: '',
    totalDecks: 0,
    totalCards: 0,
    reviewedTotal: 0,
    reviewReminder: false,
    achievements: []
  },

  onShow() {
    this.loadFromStorage();
    this.loadStats();
    this.loadAchievements();
  },

  loadFromStorage() {
    const ui = wx.getStorageSync('userInfo');
    if (ui && ui.nickName) {
      this.setData({ nickname: ui.nickName, avatarUrl: ui.avatarUrl || '' });
    }
    const rr = wx.getStorageSync('reviewReminder');
    this.setData({ reviewReminder: !!rr });
  },

  loadStats() {
    const token = app.globalData.token || wx.getStorageSync('token');
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const decks = (res.data || []);
        const cards = decks.reduce((s, d) => s + (d.card_count || 0), 0);
        this.setData({ totalDecks: decks.length, totalCards: cards });
      }
    });
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const d = res.data || {};
        this.setData({ reviewedTotal: d.reviewed_today || 0 });
      }
    });
  },

  onChooseAvatar(e) {
    const avatarUrl = e.detail.avatarUrl;
    this.setData({ avatarUrl });
    this.syncProfile({ avatarUrl });
  },

  onNicknameBlur(e) {
    const nickName = e.detail.value;
    if (!nickName || nickName === '点击设置昵称') return;
    this.setData({ nickname: nickName });
    this.syncProfile({ nickName });
  },

  syncProfile(updates) {
    const ui = wx.getStorageSync('userInfo') || {};
    if (updates.nickName) ui.nickName = updates.nickName;
    if (updates.avatarUrl) ui.avatarUrl = updates.avatarUrl;
    wx.setStorageSync('userInfo', ui);
    wx.request({
      url: app.globalData.apiBase + '/api/user/profile',
      method: 'PUT',
      data: { nickname: ui.nickName || '', avatar_url: ui.avatarUrl || '' },
      header: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token'))
      }
    });
  },

  loadAchievements() {
    wx.request({
      url: app.globalData.apiBase + '/api/achievements',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        this.setData({ achievements: (res.data && res.data.achievements) || [] });
      }
    });
  },

  toggleReviewReminder() {
    const v = !this.data.reviewReminder;
    this.setData({ reviewReminder: v });
    wx.setStorageSync('reviewReminder', v);
  },

  showInvite() {
    wx.request({
      url: app.globalData.apiBase + '/api/invite/my-code',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const code = (res.data && res.data.invite_code) || 'FLASH';
        wx.setClipboardData({ data: '来闪卡记忆一起高效学习吧！我的邀请码：' + code,
          success: () => wx.showToast({ title: '已复制邀请码', icon: 'success' }) });
      }
    });
  },

  exportData() {
    wx.showLoading({ title: '导出中...' });
    wx.request({
      url: app.globalData.apiBase + '/api/export',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        wx.hideLoading();
        wx.setClipboardData({ data: JSON.stringify(res.data),
          success: () => wx.showToast({ title: '已复制到剪贴板', icon: 'success' }) });
      },
      fail: () => { wx.hideLoading(); wx.showToast({ title: '导出失败', icon: 'none' }); }
    });
  },

  onMenuTap(e) {
    if (e.currentTarget.dataset.action === 'about') {
      wx.showModal({ title: '关于', content: 'AI 驱动智能闪卡，高效记忆', showCancel: false, confirmText: '好的' });
    }
  }
});
