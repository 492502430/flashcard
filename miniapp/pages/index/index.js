const app = getApp();

function assessOptimizeCount(cards) {
  const questionMap = {};
  cards.forEach((card) => {
    const q = (card.question || '').trim();
    if (q) questionMap[q] = (questionMap[q] || 0) + 1;
  });
  return cards.reduce((sum, card) => {
    const q = (card.question || '').trim();
    const a = (card.answer || '').trim();
    const needsOptimize = !q || !a || (q && questionMap[q] > 1) || q.length > 90 || a.length > 220 || a.length < 6;
    return sum + (needsOptimize ? 1 : 0);
  }, 0);
}

Page({
  data: {
    count: 0, totalCards: 0, totalDecks: 0, reviewedToday: 0, streak: 0,
    showSearch: false, keyword: '', searched: false, searchResults: [],
    searchTimer: null, checkinDays: [],
    generatingDecks: 0,
    optimizeCards: 0,
    readyDecks: [],
    weeklyTotal: 0,
    showAuth: false
  },

  onShow() {
    this.loadData();
    var userInfo = wx.getStorageSync("userInfo");
    this.setData({ showAuth: !userInfo || (userInfo && userInfo.nickName === "闪卡用户") });
  },
  onPullDownRefresh() { this.loadData(); wx.stopPullDownRefresh(); },

  loadData() {
    console.log('[index] loadData START, current values:', 
      this.data.totalCards, this.data.totalDecks, this.data.reviewedToday);
    
    const token = app.globalData.token || wx.getStorageSync('token');

    // 1. Review today
    wx.request({
      url: app.globalData.apiBase + '/api/review/today',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const d = res.data || {};
        console.log('[index] review/today response:', JSON.stringify({total:d.total,rt:d.reviewed_today,streak:d.streak}));
        this.setData({
          count: d.total || 0,
          reviewedToday: d.reviewed_today || 0,
          streak: d.streak || 0
        });
      },
      fail: (e) => {
        console.error('[index] review/today FAIL:', e.errMsg);
        this.setData({ count: 0, reviewedToday: 0, streak: 0 });
      }
    });

    // 2. Decks
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const decks = Array.isArray(res.data) ? res.data : [];
        const tc = decks.reduce((s,d) => s + (d.card_count||0), 0);
        const generatingDecks = decks.filter(d => !d.card_count).length;
        const readyDecks = decks.filter(d => d.card_count > 0).slice(0, 3).map(d => ({
          ...d,
          due_count: d.due_count || 0
        }));
        console.log('[index] decks response:', decks.length, 'decks,', tc, 'cards');
        this.setData({ totalDecks: decks.length, totalCards: tc, generatingDecks, readyDecks, optimizeCards: 0 });
        this.loadOptimizeCount(decks.filter(d => d.card_count > 0));
      },
      fail: (e) => {
        console.error('[index] decks FAIL:', e.errMsg);
        this.setData({ totalDecks: 0, totalCards: 0 });
      }
    });

    // 3. Stats / Checkin
    wx.request({
      url: app.globalData.apiBase + '/api/stats',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        const stats = res.data || [];
        const maxCount = Math.max(1, ...stats.map(s => s.count));
        const days = ['日','一','二','三','四','五','六'];
        const weeklyTotal = stats.reduce((sum, item) => sum + (item.count || 0), 0);
        this.setData({ weekStats: stats.map(s => ({
          ...s, label: days[new Date(s.date).getDay()],
          height: Math.max(4, Math.round(s.count / maxCount * 120))
        })), weeklyTotal });
      }
    });

    wx.request({
      url: app.globalData.apiBase + '/api/checkin',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => { this.buildCheckinGrid(res.data || []); }
    });
  },

  loadOptimizeCount(decks) {
    const token = app.globalData.token || wx.getStorageSync('token');
    if (!decks.length) {
      this.setData({ optimizeCards: 0 });
      return;
    }
    let pending = decks.length;
    let total = 0;
    decks.forEach((deck) => {
      wx.request({
        url: app.globalData.apiBase + '/api/decks/' + deck.id,
        header: { Authorization: 'Bearer ' + token },
        success: (res) => {
          const cards = (res.data && res.data.cards) || [];
          total += assessOptimizeCount(cards);
        },
        complete: () => {
          pending -= 1;
          if (pending === 0) this.setData({ optimizeCards: total });
        }
      });
    });
  },

  buildCheckinGrid(apiData) {
    const countMap = {}; apiData.forEach(d => { countMap[d.date] = d.count; });
    const today = new Date(); today.setHours(0,0,0,0);
    const grid = []; const weekDays = ['日','一','二','三','四','五','六'];
    for (let i = 29; i >= 0; i--) {
      const d = new Date(today); d.setDate(d.getDate() - i);
      const dateStr = d.toISOString().slice(0,10);
      const count = countMap[dateStr] || 0;
      let level = 0;
      if (count > 0) level = 1; if (count >= 4) level = 2;
      if (count >= 10) level = 3; if (count >= 20) level = 4;
      grid.push({ date: dateStr, count, level, dayOfWeek: weekDays[d.getDay()], isToday: i===0 });
    }
    this.setData({ checkinDays: grid });
  },

  startReview() { wx.navigateTo({ url: '/pages/review/review' }); },
  goCreate() { wx.navigateTo({ url: '/pages/create/create' }); },
  goDecks() { wx.switchTab({ url: '/pages/decks/decks' }); },
  openDeck(e) {
    const id = e.currentTarget.dataset.id;
    if (id) wx.navigateTo({ url: '/pages/deck-detail/deck-detail?id=' + id });
  },

  openSearch() { this.setData({ showSearch: true, searched: false, searchResults: [], keyword: '' }); },
  closeSearch() { this.setData({ showSearch: false, keyword: '', searched: false, searchResults: [] }); },

  onSearchInput(e) {
    const keyword = e.detail.value; this.setData({ keyword, searched: false });
    if (this.data.searchTimer) clearTimeout(this.data.searchTimer);
    this.data.searchTimer = setTimeout(() => { if (keyword.trim()) this.doSearch(); }, 400);
  },

  doSearch() {
    const kw = this.data.keyword.trim(); if (!kw) return;
    const token = app.globalData.token || wx.getStorageSync('token');
    wx.request({
      url: app.globalData.apiBase + '/api/cards/search?q=' + encodeURIComponent(kw),
      header: { Authorization: 'Bearer ' + token },
      success: (res) => this.setData({ searchResults: Array.isArray(res.data)?res.data:[], searched: true }),
      fail: () => this.setData({ searched: true, searchResults: [] })
    });
  },

  clearSearch() { this.setData({ keyword: '', searched: false, searchResults: [] }); },


  onAuthTap() {
    var that = this;
    wx.getUserProfile({
      desc: "用于展示你的昵称和头像",
      success: function(res) {
        that.onAuthGot({ detail: res });
      },
      fail: function() {
        wx.showToast({ title: "需要授权才能使用", icon: "none" });
      }
    });
  },
  onAuthGot(e) {
    var userInfo = e.detail.userInfo;
    if (!userInfo) return;
    wx.setStorageSync("userInfo", userInfo);
    wx.request({
      url: app.globalData.apiBase + "/api/user/profile",
      method: "PUT",
      header: { "Content-Type": "application/json", Authorization: "Bearer " + (app.globalData.token || wx.getStorageSync("token")) },
      data: { nickname: userInfo.nickName, avatar_url: userInfo.avatarUrl }
    });
    this.setData({ showAuth: false });
  },

  skipAuth() {
    wx.setStorageSync("userInfo", { nickName: "闪卡用户" });
    this.setData({ showAuth: false });
  },
  preventTouchMove() {}
});
