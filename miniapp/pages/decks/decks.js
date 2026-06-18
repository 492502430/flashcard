const app = getApp();

Page({
  data: {
    decks: [],
    refreshing: false
  },

  onShow() {
    this.loadDecks();
    // Auto-refresh for 30s after creating deck (cards generate async)
    this.startAutoRefresh();
  },

  onHide() {
    this.stopAutoRefresh();
  },

  loadDecks() {
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const decks = res.data || [];
        const geoStyles = ['geo-a','geo-b','geo-c','geo-d','geo-e','geo-f','geo-g','geo-h','geo-i'];
        const enriched = decks.map((d, i) => ({
          ...d,
          geoStyle: d.geoStyle || geoStyles[i % geoStyles.length]
        }));
        this.setData({ decks: enriched });
      }
    });
  },

  startAutoRefresh() {
    this.stopAutoRefresh();
    let count = 0;
    this._refreshTimer = setInterval(() => {
      this.loadDecks();
      count++;
      if (count > 10) this.stopAutoRefresh(); // Stop after 30s
    }, 3000);
  },

  stopAutoRefresh() {
    if (this._refreshTimer) {
      clearInterval(this._refreshTimer);
      this._refreshTimer = null;
    }
  },

  openDeck(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: '/pages/deck-detail/deck-detail?id=' + id });
  },

  goCreate() {
    wx.navigateTo({ url: '/pages/create/create' });
  },

  confirmDelete(e) {
    const id = e.currentTarget.dataset.id;
    const title = e.currentTarget.dataset.title || '这个牌组';
    wx.showModal({
      title: '删除牌组',
      content: `确定要删除「${title}」吗？`,
      confirmText: '删除',
      confirmColor: '#DC2626',
      success: (res) => {
        if (res.confirm) this.deleteDeck(id);
      }
    });
  },

  deleteDeck(id) {
    wx.request({
      url: app.globalData.apiBase + '/api/decks/' + id,
      method: 'DELETE',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: () => {
        wx.showToast({ title: '已删除', icon: 'success' });
        this.loadDecks();
      }
    });
  }
});
