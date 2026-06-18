const app = getApp();

Page({
  data: {
    cards: [],
    current: 1,
    total: 0,
    card: null,
    flipped: false,
    flipping: false,
    streak: 0
  },

  onLoad(opts) {
    const deckId = opts.deck_id;
    this.loadCards(deckId);
  },

  loadCards(deckId) {
    const url = deckId
      ? app.globalData.apiBase + '/api/decks/' + deckId + '/review'
      : app.globalData.apiBase + '/api/review/today';

    wx.request({
      url: url,
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const data = res.data;
        const cards = data.cards || [];
        this.setData({
          cards,
          total: cards.length,
          card: cards[0] || null,
          streak: data.streak || 0
        });
      }
    });
  },

  flipCard() {
    if (this.data.flipping) return;

    this.setData({ flipping: true });

    if (this.data.flipped) {
      this.setData({ flipped: false });
    } else {
      this.setData({ flipped: true });
    }

    setTimeout(() => {
      this.setData({ flipping: false });
    }, 450);
  },

  rate(e) {
    const rating = parseInt(e.currentTarget.dataset.r);
    const card = this.data.card;
    if (!card) return;

    wx.request({
      url: app.globalData.apiBase + '/api/review',
      method: 'POST',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      data: { card_id: card.id, rating },
      success: () => {
        this.nextCard();
      },
      fail: () => {
        wx.showToast({ title: '网络错误', icon: 'none' });
      }
    });
  },

  nextCard() {
    const next = this.data.current + 1;
    if (next > this.data.total) {
      this.setData({ card: null, flipped: false, flipping: false });
    } else {
      this.setData({
        current: next,
        card: this.data.cards[next - 1],
        flipped: false,
        flipping: false
      });
    }
  },

  goHome() {
    wx.switchTab({ url: '/pages/index/index' });
  }
});
