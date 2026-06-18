const app = getApp();

Page({
  data: {
    deck: {},
    cards: [],
    newCount: 0,
    dueCount: 0,
    generating: true,   // Assume generating until cards appear
    // Edit modal state
    showEditModal: false,
    editCardId: '',
    editQuestion: '',
    editAnswer: '',
    editSaving: false
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
    // Standard back — goes to wherever user came from
    const pages = getCurrentPages();
    if (pages.length > 1) {
      wx.navigateBack();
    } else {
      wx.switchTab({ url: '/pages/index/index' });
    }
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
  },

  confirmDeleteCard(e) {
    const cardId = e.currentTarget.dataset.id;
    wx.showModal({
      title: '删除卡片',
      content: '确定要删除这张卡片吗？',
      confirmText: '删除',
      confirmColor: '#DC2626',
      success: (res) => {
        if (res.confirm) this.deleteCard(cardId);
      }
    });
  },

  deleteCard(cardId) {
    wx.request({
      url: app.globalData.apiBase + '/api/cards/' + cardId,
      method: 'DELETE',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: () => {
        wx.showToast({ title: '已删除', icon: 'success' });
        this.loadDeck();  // Refresh list
      }
    });
  },

  // ── Edit modal handlers ──

  openEditModal(e) {
    const { id, question, answer } = e.currentTarget.dataset;
    this.setData({
      showEditModal: true,
      editCardId: id,
      editQuestion: question,
      editAnswer: answer,
      editSaving: false
    });
  },

  closeEditModal() {
    this.setData({ showEditModal: false });
  },

  onEditQuestionInput(e) {
    this.setData({ editQuestion: e.detail.value });
  },

  onEditAnswerInput(e) {
    this.setData({ editAnswer: e.detail.value });
  },

  saveEditCard() {
    const { editCardId, editQuestion, editAnswer } = this.data;
    if (!editQuestion.trim() || !editAnswer.trim()) {
      wx.showToast({ title: '问题和答案不能为空', icon: 'none' });
      return;
    }

    this.setData({ editSaving: true });

    wx.request({
      url: app.globalData.apiBase + '/api/cards/' + editCardId,
      method: 'PUT',
      header: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token'))
      },
      data: {
        question: editQuestion.trim(),
        answer: editAnswer.trim()
      },
      success: () => {
        wx.showToast({ title: '已保存', icon: 'success' });
        this.setData({ showEditModal: false, editSaving: false });
        this.loadDeck();  // Refresh list
      },
      fail: () => {
        wx.showToast({ title: '保存失败', icon: 'none' });
        this.setData({ editSaving: false });
      }
    });
  }
});
