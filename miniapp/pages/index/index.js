const app = getApp();

Page({
  data: {
    count: 0,
    totalCards: 0,
    totalDecks: 0,
    reviewedToday: 0,
    streak: 0
  },

  onShow() {
    this.loadData();
  },

  loadData() {
    // Fetch today's review count
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const data = res.data || {};
        this.setData({
          count: data.total || 0,
          reviewedToday: data.reviewed_today || 0,
          streak: data.streak || 0
        });
      },
      fail: (err) => {
        console.error('Failed to load review:', err);
      }
    });

    // Fetch decks for totals
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const decks = Array.isArray(res.data) ? res.data : [];
        this.setData({
          totalDecks: decks.length,
          totalCards: decks.reduce((s, d) => s + (d.card_count || 0), 0)
        });
      },
      fail: (err) => {
        console.error('Failed to load decks:', err);
      }
    });
  },

  startReview() {
    wx.navigateTo({ url: '/pages/review/review' });
  },

  goCreate() {
    wx.navigateTo({ url: '/pages/create/create' });
  },

  goDecks() {
    wx.switchTab({ url: '/pages/decks/decks' });
  }
});
