const app = getApp();

Page({
  data: {
    deck: {},
    cards: [],
    newCount: 0,
    dueCount: 0,
    generating: true   // Assume generating until cards appear
  },

  onLoad(opts) {
    this.deckId = opts.id;
    this.loadDeck();
    this.startPolling();
  },

  onUnload() {
    this.stopPolling();
  },

  loadDeck() {
    wx.request({
      url: app.globalData.apiBase + '/api/decks/' + this.deckId,
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const data = res.data;
        const deck = data.deck || {};
        const cards = data.cards || [];

        const newCount = cards.filter(c => c.state === 'new').length;
        const dueCount = cards.filter(c => c.state !== 'new').length;
        const hasCards = cards.length > 0;

        this.setData({ deck, cards, newCount, dueCount, generating: !hasCards });
        
        if (hasCards) this.stopPolling();
      }
    });
  },

  startPolling() {
    this.stopPolling();
    this._pollTimer = setInterval(() => {
      this.loadDeck();
    }, 3000); // Poll every 3 seconds
  },

  stopPolling() {
    if (this._pollTimer) {
      clearInterval(this._pollTimer);
      this._pollTimer = null;
    }
  },

  startReview() {
    wx.navigateTo({ url: '/pages/review/review?deck_id=' + this.deckId });
  },

  goBack() {
    wx.switchTab({ url: '/pages/index/index' });
  },

  confirmDelete() {
    wx.showModal({
      title: '删除牌组',
      content: `确定要删除「${this.data.deck.title}」吗？所有卡片也会被删除。`,
      confirmText: '删除',
      confirmColor: '#DC2626',
      success: (res) => {
        if (res.confirm) {
          this.deleteDeck();
        }
      }
    });
  },

  deleteDeck() {
    wx.request({
      url: app.globalData.apiBase + '/api/decks/' + this.deckId,
      method: 'DELETE',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: () => {
        wx.showToast({ title: '已删除', icon: 'success' });
        setTimeout(() => wx.navigateBack(), 500);
      },
      fail: () => {
        wx.showToast({ title: '删除失败', icon: 'none' });
      }
    });
  }
});
