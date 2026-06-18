const app = getApp();
Page({
  data: { cards: [], current: 1, total: 0, card: null, flipped: false },
  onLoad() {
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        const cards = res.data.cards || [];
        this.setData({ cards, total: cards.length, card: cards[0] || null });
      }
    });
  },
  flip() { this.setData({ flipped: !this.data.flipped }); },
  rate(e) {
    const rating = parseInt(e.currentTarget.dataset.r);
    wx.request({
      url: app.globalData.apiBase + '/api/review',
      method: 'POST',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      data: { card_id: this.data.card.id, rating },
      success: () => {
        const next = this.data.current + 1;
        if (next > this.data.total) {
          this.setData({ card: null, flipped: false });
        } else {
          this.setData({ current: next, card: this.data.cards[next - 1], flipped: false });
        }
      }
    });
  },
  goHome() { wx.switchTab({ url: '/pages/index/index' }); }
});