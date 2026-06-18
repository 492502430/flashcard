const app = getApp();
Page({
  data: { deck: {}, cards: [] },
  onLoad(opts) {
    wx.request({
      url: app.globalData.apiBase + '/api/decks/' + opts.id,
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => this.setData({ deck: res.data.deck, cards: res.data.cards || [] })
    });
  }
});