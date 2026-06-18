const app = getApp();

Page({
  data: {
    decks: []
  },

  onShow() {
    this.loadDecks();
  },

  loadDecks() {
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const decks = res.data || [];
        /* Geometric icon style classes — NO emoji */
        const geoStyles = ['geo-a','geo-b','geo-c','geo-d','geo-e','geo-f','geo-g','geo-h','geo-i'];
        const enriched = decks.map((d, i) => ({
          ...d,
          geoStyle: d.geoStyle || geoStyles[i % geoStyles.length]
        }));
        this.setData({ decks: enriched });
      }
    });
  },

  openDeck(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: '/pages/deck-detail/deck-detail?id=' + id });
  },

  goCreate() {
    wx.navigateTo({ url: '/pages/create/create' });
  }
});
