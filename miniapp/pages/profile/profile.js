var app = getApp();

Page({
  data: {
    nickname: '点击设置昵称', avatarUrl: '',
    totalDecks: 0, totalCards: 0, reviewedTotal: 0,
    reviewReminder: false,
    achievements: [],
    menus: [
      { title: '我的订阅', action: 'subscribe' },
      { title: '学习记录', action: 'history' },
      { title: '邀请好友', action: 'invite' },
      { title: '导出数据', action: 'export' },
      { title: '关于闪卡记忆', action: 'about' },
      { title: '意见反馈', action: 'feedback' }
    ]
  },

  onShow: function() {
    this.loadFromStorage();
    this.loadStats();
    this.loadAchievements();
    var rr = wx.getStorageSync('reviewReminder');
    this.setData({ reviewReminder: !!rr });
    wx.showShareMenu({ withShareTicket: true, menus: ['shareAppMessage'] });
  },

  loadFromStorage: function() {
    var ui = wx.getStorageSync('userInfo');
    if (ui && ui.nickName) {
      this.setData({ nickname: ui.nickName, avatarUrl: ui.avatarUrl || '' });
    }
  },

  loadStats: function() {
    var that = this;
    var t = app.globalData.token || wx.getStorageSync('token');
    wx.request({ url: app.globalData.apiBase + '/api/decks', header: { Authorization: 'Bearer ' + t },
      success: function(r) {
        var d = r.data || [];
        var cards = 0;
        for (var i = 0; i < d.length; i++) { cards += d[i].card_count || 0; }
        that.setData({ totalDecks: d.length, totalCards: cards });
      } });
    wx.request({ url: app.globalData.apiBase + '/api/review/today', header: { Authorization: 'Bearer ' + t },
      success: function(r) { that.setData({ reviewedTotal: (r.data && r.data.reviewed_today) || 0 }); } });
  },

  loadAchievements: function() {
    var that = this;
    // Hardcoded defaults
    var def = [
      {key:'first_review',title:'初次记忆',description:'完成第一次复习',icon:'star',earned:false},
      {key:'cards_10',title:'十卡入门',description:'累计复习10张',icon:'diamond',earned:false},
      {key:'cards_50',title:'勤学不辍',description:'累计复习50张',icon:'fire',earned:false},
      {key:'streak_3',title:'三日坚持',description:'连续3天复习',icon:'streak',earned:false},
      {key:'streak_7',title:'一周成习',description:'连续7天复习',icon:'crown',earned:false},
      {key:'cards_100',title:'百卡达人',description:'累计复习100张',icon:'trophy',earned:false}
    ];
    // Always show defaults immediately
    that.setData({ achievements: def });
    // Try to load real data
    wx.request({ url: app.globalData.apiBase + '/api/achievements',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: function(r) {
        var data = r.data && r.data.achievements;
        if (data && data.length) { that.setData({ achievements: data }); }
      }
    });
  },

  onChooseAvatar: function(e) {
    var url = e.detail.avatarUrl;
    this.setData({ avatarUrl: url });
    this.syncProfile({ avatarUrl: url });
  },

  onNicknameBlur: function(e) {
    var name = e.detail.value;
    if (!name) return;
    this.setData({ nickname: name });
    this.syncProfile({ nickName: name });
  },

  syncProfile: function(updates) {
    var ui = wx.getStorageSync('userInfo') || {};
    if (updates.nickName) ui.nickName = updates.nickName;
    if (updates.avatarUrl) ui.avatarUrl = updates.avatarUrl;
    wx.setStorageSync('userInfo', ui);
    wx.request({
      url: app.globalData.apiBase + '/api/user/profile', method: 'PUT',
      data: { nickname: ui.nickName || '', avatar_url: ui.avatarUrl || '' },
      header: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) }
    });
  },

  toggleReviewReminder: function(e) {
    var v = e.detail.value;
    this.setData({ reviewReminder: v });
    wx.setStorageSync('reviewReminder', v);
  },

  onMenuTap: function(e) {
    var action = e.currentTarget.dataset.action;
    var that = this;
    if (action === 'about') { wx.showModal({ title: '关于', content: 'AI 智能闪卡，高效记忆', showCancel: false, confirmText: '好的' }); }
    else if (action === 'feedback') { wx.showModal({ title: '反馈', content: 'feedback@flashcard.app', showCancel: false, confirmText: '好的' }); }
    else if (action === 'invite') {
      wx.request({ url: app.globalData.apiBase + '/api/invite/my-code',
        header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
        success: function(r) {
          var code = r.data && r.data.invite_code || 'FLASH';
          wx.setClipboardData({ data: '邀请码：' + code, success: function() { wx.showToast({ title: '已复制', icon: 'success' }); } });
        } });
    }
    else if (action === 'export') {
      wx.showLoading({ title: '导出中' });
      wx.request({ url: app.globalData.apiBase + '/api/export',
        header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
        success: function(r) { wx.hideLoading(); wx.setClipboardData({ data: JSON.stringify(r.data) }); wx.showToast({ title: '已复制', icon: 'success' }); },
        fail: function() { wx.hideLoading(); wx.showToast({ title: '失败', icon: 'none' }); } });
    }
    else { wx.showToast({ title: '即将上线', icon: 'none' }); }
  },

  onShareAppMessage: function() {
    return { title: '来闪卡记忆一起高效学习吧', path: '/pages/index/index' };
  }
});
