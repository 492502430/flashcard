const app = getApp();

Page({
  data: {
    nickname: '点击设置昵称',
    avatarUrl: '',
    totalDecks: 0, totalCards: 0, reviewedTotal: 0,
    reviewReminder: false,
    achievements: [],
    inviteCode: '', invitedCount: 0
  },

  onShow() {
    this.loadFromStorage();
    this.loadStats();
    this.loadAchievements();
    this.loadInviteInfo();
    const rr = wx.getStorageSync('reviewReminder');
    this.setData({ reviewReminder: !!rr });
    wx.showShareMenu({ withShareTicket: true, menus: ['shareAppMessage'] });
  },

  loadFromStorage() {
    const ui = wx.getStorageSync('userInfo');
    if (ui && ui.nickName) {
      this.setData({ nickname: ui.nickName, avatarUrl: ui.avatarUrl || '' });
    }
  },

  loadStats() {
    const t = app.globalData.token || wx.getStorageSync('token');
    wx.request({ url: app.globalData.apiBase + '/api/decks', header: { Authorization: 'Bearer ' + t },
      success: (r) => { const d = r.data || []; this.setData({ totalDecks: d.length, totalCards: d.reduce((s,x)=>s+(x.card_count||0),0) }); } });
    wx.request({ url: app.globalData.apiBase + '/api/review/today', header: { Authorization: 'Bearer ' + t },
      success: (r) => { const d = r.data || {}; this.setData({ reviewedTotal: d.reviewed_today || 0 }); } });
  },

  onChooseAvatar(e) {
    const avatarUrl = e.detail.avatarUrl;
    this.setData({ avatarUrl });
    this.syncProfile({ avatarUrl });
  },

  onNicknameBlur(e) {
    const nickName = e.detail.value;
    if (!nickName) return;
    this.setData({ nickname: nickName });
    this.syncProfile({ nickName });
  },

  syncProfile(updates) {
    const ui = wx.getStorageSync('userInfo') || {};
    if (updates.nickName) ui.nickName = updates.nickName;
    if (updates.avatarUrl) ui.avatarUrl = updates.avatarUrl;
    wx.setStorageSync('userInfo', ui);
    wx.request({
      url: app.globalData.apiBase + '/api/user/profile', method: 'PUT',
      data: { nickname: ui.nickName || '', avatar_url: ui.avatarUrl || '' },
      header: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) }
    });
  },

  loadAchievements() {
    wx.request({ url: app.globalData.apiBase + '/api/achievements',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (r) => this.setData({ achievements: (r.data && r.data.achievements) || [] }) });
  },

  loadInviteInfo() {
    wx.request({ url: app.globalData.apiBase + '/api/invite/my-code',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (r) => {
        const code = (r.data && r.data.invite_code) || '';
        if (!code) return;
        this.setData({ inviteCode: code });
        wx.request({ url: app.globalData.apiBase + '/api/invite/stats',
          header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
          success: (sr) => this.setData({ invitedCount: (sr.data && sr.data.invited_count) || 0 }) });
      } });
  },

  toggleReviewReminder(e) {
    const v = e.detail.value;
    this.setData({ reviewReminder: v });
    wx.setStorageSync('reviewReminder', v);
  },

  onMenuTap(e) {
    const action = e.currentTarget.dataset.action;
    if (action === 'about') wx.showModal({ title: '关于闪卡记忆', content: 'AI 驱动智能闪卡，高效记忆', showCancel: false, confirmText: '好的' });
    else if (action === 'feedback') wx.showModal({ title: '意见反馈', content: '感谢反馈！feedback@flashcard.app', showCancel: false, confirmText: '好的' });
    else wx.showToast({ title: '即将上线', icon: 'none' });
  },

  exportData() {
    wx.showLoading({ title: '导出中...' });
    wx.request({ url: app.globalData.apiBase + '/api/export',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (r) => { wx.hideLoading(); wx.setClipboardData({ data: JSON.stringify(r.data), success: ()=>wx.showToast({ title: '已复制', icon: 'success' }) }); },
      fail: () => { wx.hideLoading(); wx.showToast({ title: '失败', icon: 'none' }); } });
  },

  showInvite() {
    if (!this.data.inviteCode) return;
    wx.setClipboardData({ data: '来闪卡记忆一起高效学习！邀请码：' + this.data.inviteCode, success: ()=>wx.showToast({ title: '已复制', icon: 'success' }) });
  },

  onShareAppMessage() {
    return { title: '来闪卡记忆一起高效学习吧', path: '/pages/onboard/onboard?invite_code=' + (this.data.inviteCode || 'FLASH') };
  }
});
