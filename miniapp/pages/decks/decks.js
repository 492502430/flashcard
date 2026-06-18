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
        // Add placeholder icons based on index
        const icons = ['📝', '🧠', '📚', '💡', '🎯', '🔬', '🌍', '📖', '✍️'];
        const enriched = decks.map((d, i) => ({
          ...d,
          icon: d.icon || icons[i % icons.length]
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
