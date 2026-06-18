const app = getApp();
Page({
  data: { decks: [] },
  onShow() {
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => this.setData({ decks: res.data || [] })
    });
  },
  openDeck(e) {
    wx.navigateTo({ url: '/pages/deck-detail/deck-detail?id=' + e.currentTarget.dataset.id });
  },
  goCreate() { wx.navigateTo({ url: '/pages/create/create' }); }
});