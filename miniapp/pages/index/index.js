const app = getApp();

Page({
  data: {
    count: 0,
    totalCards: 0,
    totalDecks: 0,
    reviewedToday: 0,
    streak: 0,
    showSearch: false,
    keyword: '',
    searched: false,
    searchResults: [],
    searchTimer: null
  },

  onShow() {
    this.loadData();
  },

  loadData() {
    // Fetch today's review count
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const data = res.data || {};
        this.setData({
          count: data.total || 0,
          reviewedToday: data.reviewed_today || 0,
          streak: data.streak || 0
        });
      },
      fail: (err) => {
        console.error('Failed to load review:', err);
      }
    });

    // Fetch decks for totals
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const decks = Array.isArray(res.data) ? res.data : [];
        this.setData({
          totalDecks: decks.length,
          totalCards: decks.reduce((s, d) => s + (d.card_count || 0), 0)
        });
      },
      fail: (err) => {
        console.error('Failed to load decks:', err);
      }
    });

    // Load weekly stats
    wx.request({
      url: app.globalData.apiBase + '/api/stats',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const stats = res.data || [];
        const maxCount = Math.max(1, ...stats.map(s => s.count));
        const days = ['日','一','二','三','四','五','六'];
        const weekStats = stats.map(s => ({
          ...s,
          label: days[new Date(s.date).getDay()],
          height: Math.max(4, Math.round(s.count / maxCount * 120))
        }));
        this.setData({ weekStats });
      }
    });
  },

  startReview() {
    wx.navigateTo({ url: '/pages/review/review' });
  },

  goCreate() {
    wx.navigateTo({ url: '/pages/create/create' });
  },

  goDecks() {
    wx.switchTab({ url: '/pages/decks/decks' });
  },

  // Search
  openSearch() {
    this.setData({ showSearch: true, searched: false, searchResults: [], keyword: '' });
  },

  closeSearch() {
    this.setData({ showSearch: false, keyword: '', searched: false, searchResults: [] });
  },

  onSearchInput(e) {
    const keyword = e.detail.value;
    this.setData({ keyword, searched: false });

    // Debounce search
    if (this.data.searchTimer) clearTimeout(this.data.searchTimer);
    this.data.searchTimer = setTimeout(() => {
      if (keyword.trim()) {
        this.doSearch();
      }
    }, 400);
  },

  doSearch() {
    const keyword = this.data.keyword.trim();
    if (!keyword) return;

    wx.request({
      url: app.globalData.apiBase + '/api/cards/search?q=' + encodeURIComponent(keyword),
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        this.setData({
          searchResults: Array.isArray(res.data) ? res.data : [],
          searched: true
        });
      },
      fail: (err) => {
        console.error('Search failed:', err);
        this.setData({ searched: true, searchResults: [] });
      }
    });
  },

  clearSearch() {
    this.setData({ keyword: '', searched: false, searchResults: [] });
  },

  preventTouchMove() {
    // Prevent background scroll when modal is open
  }
});
