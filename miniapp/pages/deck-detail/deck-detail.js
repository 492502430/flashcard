const app = getApp();

Page({
  data: {
    deck: {},
    cards: [],
    newCount: 0,
    dueCount: 0
  },

  onLoad(opts) {
    const id = opts.id;
    this.loadDeck(id);
  },

  loadDeck(id) {
    wx.request({
      url: app.globalData.apiBase + '/api/decks/' + id,
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const data = res.data;
        const deck = data.deck || {};
        const cards = data.cards || [];

        const newCount = cards.filter(c => c.state === 'new').length;
        const dueCount = cards.filter(c => c.state === 'learning' || c.state === 'review').length;

        this.setData({
          deck,
          cards,
          newCount,
          dueCount
        });
      }
    });
  },

  startReview() {
    const id = this.data.deck.id;
    wx.navigateTo({ url: '/pages/review/review?deck_id=' + id });
  },

  goBack() {
    wx.navigateBack();
  }
});
