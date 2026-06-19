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
    searchTimer: null,
    checkinDays: []  // 30-day heatmap data
  },

  onShow() {
    this.loadData();
  },

  loadData() {
    // Reset to zero before fetching — prevent stale cached values
    this.setData({ count: 0, totalCards: 0, totalDecks: 0, reviewedToday: 0, streak: 0 });

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

    // Load 30-day checkin for calendar heatmap
    wx.request({
      url: app.globalData.apiBase + '/api/checkin',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const checkinData = res.data || [];
        this.buildCheckinGrid(checkinData);
      }
    });
  },

  buildCheckinGrid(apiData) {
    // Build a map from date string to count
    const countMap = {};
    apiData.forEach(d => { countMap[d.date] = d.count; });

    // Generate past 30 days (including today)
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const grid = [];
    const weekDays = ['日','一','二','三','四','五','六'];

    for (let i = 29; i >= 0; i--) {
      const d = new Date(today);
      d.setDate(d.getDate() - i);
      const dateStr = d.toISOString().slice(0, 10);
      const count = countMap[dateStr] || 0;
      // Level: 0=none, 1=1-3, 2=4-9, 3=10-19, 4=20+
      let level = 0;
      if (count > 0) level = 1;
      if (count >= 4) level = 2;
      if (count >= 10) level = 3;
      if (count >= 20) level = 4;
      grid.push({
        date: dateStr,
        count,
        level,
        dayOfWeek: weekDays[d.getDay()],
        isToday: i === 0
      });
    }
    this.setData({ checkinDays: grid });
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
