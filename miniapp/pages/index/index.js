const app = getApp();
Page({
  data: { count: 0, totalCards: 0, totalDecks: 0 },
  onShow() { this.loadData(); },
  loadData() {
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => this.setData({ count: res.data.total })
    });
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const decks = res.data;
        this.setData({
          totalDecks: decks.length,
          totalCards: decks.reduce((s, d) => s + d.card_count, 0)
        });
      }
    });
  },
  startReview() { wx.navigateTo({ url: '/pages/review/review' }); },
  goCreate() { wx.navigateTo({ url: '/pages/create/create' }); }
});