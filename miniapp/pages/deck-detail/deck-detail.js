const app = getApp();

function parseTags(tags) {
  if (!tags) return [];
  if (Array.isArray(tags)) return tags.filter(Boolean);
  try {
    const parsed = JSON.parse(tags);
    return Array.isArray(parsed) ? parsed.filter(Boolean) : [];
  } catch (e) {
    return String(tags).split(',').map(t => t.trim()).filter(Boolean);
  }
}

function assessCards(cards) {
  const questionMap = {};
  cards.forEach((card) => {
    const q = (card.question || '').trim();
    if (!q) return;
    questionMap[q] = (questionMap[q] || 0) + 1;
  });

  let optimizeCount = 0;
  let duplicateCount = 0;
  const tagCount = {};
  const enriched = cards.map((card) => {
    const q = (card.question || '').trim();
    const a = (card.answer || '').trim();
    const tags = parseTags(card.tags);
    tags.forEach(tag => { tagCount[tag] = (tagCount[tag] || 0) + 1; });

    let qualityStatus = 'good';
    let qualityLabel = '优秀';
    let qualityClass = 'good';
    let qualityScore = 100;
    const qualityReasons = [];

    if (!q) {
      qualityReasons.push('问题为空');
      qualityScore -= 40;
    } else if (q.length > 90) {
      qualityReasons.push('问题过长，建议拆成单一知识点');
      qualityScore -= 20;
    }

    if (!a) {
      qualityReasons.push('答案为空');
      qualityScore -= 45;
    } else if (a.length < 6) {
      qualityReasons.push('答案过短，信息不足');
      qualityScore -= 25;
    } else if (a.length > 220) {
      qualityReasons.push('答案过长，建议压缩成可回忆要点');
      qualityScore -= 20;
    }

    if (q && questionMap[q] > 1) {
      duplicateCount += 1;
      qualityReasons.push('问题与其他卡片重复');
      qualityScore -= 35;
    }

    if (qualityReasons.length > 0) {
      qualityStatus = q && questionMap[q] > 1 ? 'duplicate' : 'optimize';
      qualityLabel = qualityStatus === 'duplicate' ? '重复疑似' : '待优化';
      qualityClass = qualityStatus === 'duplicate' ? 'bad' : 'warn';
      if (qualityStatus === 'optimize') optimizeCount += 1;
    }

    const documentName = card.document_name || card.documentName || '未分组文档';
    return {
      ...card,
      tags,
      documentName,
      qualityScore: Math.max(0, qualityScore),
      qualityReasons,
      qualityReasonText: qualityReasons.join('；'),
      qualityStatus,
      qualityLabel,
      qualityClass
    };
  });

  const clearCount = Math.max(0, cards.length - optimizeCount - duplicateCount);
  const qualityScore = cards.length
    ? Math.round(enriched.reduce((sum, card) => sum + card.qualityScore, 0) / cards.length)
    : 0;
  const qualityLabel = cards.length === 0 ? '暂无' : qualityScore >= 90 ? '优秀' : qualityScore >= 80 ? '良好' : qualityScore >= 65 ? '可优化' : '待优化';
  const deckTagLabel = Object.keys(tagCount).sort((a, b) => tagCount[b] - tagCount[a])[0] || '未设置';

  return { enriched, clearCount, optimizeCount, duplicateCount, qualityScore, qualityLabel, deckTagLabel };
}

function groupCardsByDocument(cards) {
  const map = {};
  const groups = [];
  cards.forEach((card) => {
    const name = card.documentName || '未分组文档';
    if (!map[name]) {
      map[name] = { name, cards: [], cardCount: 0, optimizeCount: 0 };
      groups.push(map[name]);
    }
    map[name].cards.push(card);
    map[name].cardCount += 1;
    if (card.qualityStatus !== 'good') map[name].optimizeCount += 1;
  });
  return groups;
}

Page({
  data: {
    deck: {},
    cards: [],
    newCount: 0,
    dueCount: 0,
    learningCount: 0,
    masteredCount: 0,
    masteryRate: 0,
    clearCount: 0,
    optimizeCount: 0,
    duplicateCount: 0,
    qualityScore: 0,
    qualityScoreText: '--',
    qualityLabel: '暂无',
    deckTagLabel: '未设置',
    activeFilter: 'all',
    filteredCards: [],
    documentGroups: [],
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
        const cards = (data.cards || []).map(card => ({
          ...card,
          document_name: card.document_name || deck.source_name || deck.title || '未分组文档'
        }));

        const newCount = cards.filter(c => c.state === 'new').length;
        const learningCount = cards.filter(c => c.state === 'learning').length;
        const masteredCount = cards.filter(c => c.state === 'review').length;
        const dueCount = cards.filter(c => c.state !== 'review').length;
        const masteryRate = cards.length ? Math.round(masteredCount / cards.length * 100) : 0;
        const quality = assessCards(cards);
        const hasCards = cards.length > 0;
        const activeFilter = this.data.activeFilter || 'all';
        const filteredCards = this.filterCards(quality.enriched, activeFilter);
        const documentGroups = groupCardsByDocument(filteredCards);

        this.setData({
          deck,
          cards: quality.enriched,
          newCount,
          dueCount,
          learningCount,
          masteredCount,
          masteryRate,
          clearCount: quality.clearCount,
          optimizeCount: quality.optimizeCount,
          duplicateCount: quality.duplicateCount,
          qualityScore: quality.qualityScore,
          qualityScoreText: hasCards ? String(quality.qualityScore) : '--',
          qualityLabel: quality.qualityLabel,
          deckTagLabel: quality.deckTagLabel,
          filteredCards,
          documentGroups,
          generating: !hasCards
        });
        
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

  filterCards(cards, filter) {
    if (filter === 'good') return cards.filter(c => c.qualityStatus === 'good');
    if (filter === 'optimize') return cards.filter(c => c.qualityStatus === 'optimize');
    if (filter === 'duplicate') return cards.filter(c => c.qualityStatus === 'duplicate');
    return cards;
  },

  setFilter(e) {
    const filter = e.currentTarget.dataset.filter || 'all';
    const filteredCards = this.filterCards(this.data.cards, filter);
    this.setData({
      activeFilter: filter,
      filteredCards,
      documentGroups: groupCardsByDocument(filteredCards)
    });
  },

  showLowQualityCards() {
    const hasOptimize = this.data.cards.some(card => card.qualityStatus === 'optimize');
    const hasDuplicate = this.data.cards.some(card => card.qualityStatus === 'duplicate');
    const filter = hasOptimize ? 'optimize' : hasDuplicate ? 'duplicate' : 'all';
    const filteredCards = this.filterCards(this.data.cards, filter);
    this.setData({
      activeFilter: filter,
      filteredCards,
      documentGroups: groupCardsByDocument(filteredCards)
    });
    wx.showToast({
      title: filter === 'all' ? '暂无待优化卡片' : '已筛选待优化卡片',
      icon: 'none'
    });
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
        this.loadDeck();
      },
      fail: () => {
        wx.showToast({ title: '保存失败', icon: 'none' });
        this.setData({ editSaving: false });
      }
    });
  },

  // ── AI Batch Optimize ──
  aiOptimizeCards() {
    const cardsToOptimize = this.data.cards.filter(
      c => c.qualityStatus === 'optimize' || c.qualityStatus === 'duplicate'
    );
    if (cardsToOptimize.length === 0) {
      wx.showToast({ title: '没有需要优化的卡片', icon: 'none' });
      return;
    }

    wx.showLoading({ title: 'AI 正在优化...' });

    const payload = cardsToOptimize.map(c => ({
      q: c.question,
      a: c.answer,
      id: c.id,
      reason: c.qualityReasonText || ''
    }));

    wx.request({
      url: app.globalData.apiBase + '/api/cards/optimize',
      method: 'POST',
      header: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token'))
      },
      data: { cards: payload },
      success: (res) => {
        wx.hideLoading();
        const optimized = res.data.cards || [];
        if (optimized.length === 0) {
          wx.showToast({ title: '优化完成', icon: 'success' });
          this.loadDeck();
          return;
        }
        // Update each card with optimized content
        let updated = 0;
        optimized.forEach((opt, i) => {
          const orig = payload[i];
          if (!orig) return;
          wx.request({
            url: app.globalData.apiBase + '/api/cards/' + orig.id,
            method: 'PUT',
            header: {
              'Content-Type': 'application/json',
              Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token'))
            },
            data: { question: opt.q, answer: opt.a },
            success: () => { updated++; if (updated >= optimized.length) this.loadDeck(); },
            fail: () => { updated++; if (updated >= optimized.length) this.loadDeck(); }
          });
        });
        wx.showToast({ title: `已优化 ${optimized.length} 张卡片`, icon: 'success' });
      },
      fail: () => {
        wx.hideLoading();
        wx.showToast({ title: '优化失败', icon: 'none' });
      }
    });
  }
});
