const app = getApp();

Page({
  data: {
    cards: [],
    current: 1,
    total: 0,
    card: null,
    flipped: false,
    flipping: false,
    streak: 0,
    reviewedSession: 0,  // How many reviewed this session
    // Feedback state
    showFeedback: false,
    lastRatedCardId: '',
    feedbackOptions: [
      { type: 'content_error', label: '内容有误' },
      { type: 'answer_too_brief', label: '答案太简略' },
      { type: 'question_unclear', label: '问题不清晰' }
    ],
    // New achievements from the most recent review
    newAchievements: [],
    showPosterCanvas: false
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
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      success: (res) => {
        const data = res.data || {};
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
    this.setData({ flipped: !this.data.flipped });
    setTimeout(() => this.setData({ flipping: false }), 450);
  },

  rate(e) {
    const rating = parseInt(e.currentTarget.dataset.r);
    const card = this.data.card;
    if (!card) return;

    // Clear any existing feedback timer
    this.clearFeedbackTimer();

    wx.request({
      url: app.globalData.apiBase + '/api/review',
      method: 'POST',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      data: { card_id: card.id, rating },
      success: (res) => {
        // Track: one more card reviewed this session
        this.setData({ reviewedSession: this.data.reviewedSession + 1 });

        // Check for new achievements
        const newAchievements = (res.data && res.data.new_achievements) || [];
        if (newAchievements.length > 0) {
          this.setData({ newAchievements });
          this.showAchievementToasts();
        }

        // Show feedback button for this card
        this.setData({
          showFeedback: true,
          lastRatedCardId: card.id
        });
        // Auto-hide after 3 seconds
        this._feedbackTimer = setTimeout(() => {
          this.setData({ showFeedback: false });
        }, 3000);

        this.nextCard();
      },
      fail: () => {
        wx.showToast({ title: '网络错误', icon: 'none' });
      }
    });
  },

  clearFeedbackTimer() {
    if (this._feedbackTimer) {
      clearTimeout(this._feedbackTimer);
      this._feedbackTimer = null;
    }
  },

  submitFeedback(e) {
    const type = e.currentTarget.dataset.type;
    const cardId = this.data.lastRatedCardId;
    if (!cardId) return;

    wx.request({
      url: app.globalData.apiBase + '/api/cards/' + cardId + '/feedback',
      method: 'POST',
      header: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token'))
      },
      data: { type: type },
      success: () => {
        wx.showToast({ title: '感谢反馈', icon: 'success' });
      },
      fail: () => {
        wx.showToast({ title: '反馈失败', icon: 'none' });
      }
    });

    // Hide feedback options after selection
    this.setData({ showFeedback: false });
    this.clearFeedbackTimer();
  },

  dismissFeedback() {
    this.setData({ showFeedback: false });
    this.clearFeedbackTimer();
  },

  showAchievementToasts() {
    const defs = {
      first_review: '初次记忆',
      cards_10: '十卡入门',
      cards_50: '勤学不辍',
      cards_100: '百卡达人',
      streak_7: '七日坚持',
      streak_30: '月度之星'
    };
    const achievements = this.data.newAchievements;
    if (achievements.length === 0) return;

    // Show toasts sequentially with a small delay
    achievements.forEach((key, i) => {
      setTimeout(() => {
        const title = defs[key] || key;
        wx.showToast({
          title: '🏆 获得成就：「' + title + '」',
          icon: 'none',
          duration: 2500
        });
      }, i * 2800);
    });

    // Clear after showing
    setTimeout(() => {
      this.setData({ newAchievements: [] });
    }, achievements.length * 2800 + 500);
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

  generatePoster() {
    const that = this;
    this.setData({ showPosterCanvas: true });

    // Wait for canvas to render in DOM
    setTimeout(() => {
      const query = wx.createSelectorQuery();
      query.select('#posterCanvas')
        .fields({ node: true, size: true })
        .exec((res) => {
          if (!res || !res[0] || !res[0].node) {
            wx.showToast({ title: '生成失败', icon: 'none' });
            that.setData({ showPosterCanvas: false });
            return;
          }

          const canvas = res[0].node;
          const ctx = canvas.getContext('2d');
          const dpr = wx.getSystemInfoSync().pixelRatio;

          // Fixed poster dimensions (375x550 logical px)
          const width = 375;
          const height = 550;
          canvas.width = width * dpr;
          canvas.height = height * dpr;
          ctx.scale(dpr, dpr);

          // Background — purple gradient
          const grad = ctx.createLinearGradient(0, 0, 0, height);
          grad.addColorStop(0, '#5B21B6');
          grad.addColorStop(0.5, '#7C3AED');
          grad.addColorStop(1, '#6D28D9');
          ctx.fillStyle = grad;
          this.roundRect(ctx, 0, 0, width, height, 20);
          ctx.fill();

          // Decorative circles (subtle)
          ctx.globalAlpha = 0.08;
          ctx.fillStyle = '#FFFFFF';
          ctx.beginPath(); ctx.arc(320, 60, 120, 0, Math.PI * 2); ctx.fill();
          ctx.beginPath(); ctx.arc(50, 480, 80, 0, Math.PI * 2); ctx.fill();
          ctx.globalAlpha = 1;

          // Title
          ctx.fillStyle = 'rgba(255,255,255,0.7)';
          ctx.font = 'bold 18px PingFang SC, sans-serif';
          ctx.textAlign = 'center';
          ctx.fillText('今日复习完成', width / 2, 100);

          // Large review count
          ctx.fillStyle = '#FFFFFF';
          ctx.font = 'bold 80px PingFang SC, sans-serif';
          ctx.fillText(String(that.data.reviewedSession), width / 2, 200);

          // "张卡片" subtitle
          ctx.fillStyle = 'rgba(255,255,255,0.55)';
          ctx.font = '20px PingFang SC, sans-serif';
          ctx.fillText('张卡片', width / 2, 240);

          // Divider line
          ctx.strokeStyle = 'rgba(255,255,255,0.2)';
          ctx.lineWidth = 1;
          ctx.beginPath();
          ctx.moveTo(60, 275);
          ctx.lineTo(width - 60, 275);
          ctx.stroke();

          // Streak info
          ctx.fillStyle = '#FFFFFF';
          ctx.font = 'bold 36px PingFang SC, sans-serif';
          ctx.fillText('🔥 ' + String(that.data.streak), width / 2, 320);
          ctx.fillStyle = 'rgba(255,255,255,0.55)';
          ctx.font = '16px PingFang SC, sans-serif';
          ctx.fillText('连续打卡天数', width / 2, 348);

          // QR code placeholder (white rounded square)
          const qrSize = 100;
          const qrX = (width - qrSize) / 2;
          const qrY = 380;
          ctx.fillStyle = '#FFFFFF';
          this.roundRect(ctx, qrX, qrY, qrSize, qrSize, 12);
          ctx.fill();
          ctx.fillStyle = '#5B21B6';
          ctx.font = '12px PingFang SC, sans-serif';
          ctx.fillText('扫码体验', width / 2, qrY + 52);

          // Branding
          ctx.fillStyle = 'rgba(255,255,255,0.5)';
          ctx.font = '14px PingFang SC, sans-serif';
          ctx.fillText('闪卡记忆 · AI智能复习', width / 2, height - 30);

          // Export to image
          wx.canvasToTempFilePath({
            canvas,
            x: 0,
            y: 0,
            width: width * dpr,
            height: height * dpr,
            destWidth: width * dpr,
            destHeight: height * dpr,
            success(resTemp) {
              that.setData({ showPosterCanvas: false });
              wx.showShareImageMenu({
                path: resTemp.tempFilePath,
                success() {
                  wx.showToast({ title: '长按图片即可分享', icon: 'none' });
                },
                fail() {
                  // Fallback: preview image
                  wx.previewImage({
                    urls: [resTemp.tempFilePath],
                    current: resTemp.tempFilePath
                  });
                }
              });
            },
            fail() {
              that.setData({ showPosterCanvas: false });
              wx.showToast({ title: '生成失败', icon: 'none' });
            }
          });
        });
    }, 300);
  },

  // Helper: draw rounded rectangle
  roundRect(ctx, x, y, w, h, r) {
    ctx.beginPath();
    ctx.moveTo(x + r, y);
    ctx.lineTo(x + w - r, y);
    ctx.arcTo(x + w, y, x + w, y + r, r);
    ctx.lineTo(x + w, y + h - r);
    ctx.arcTo(x + w, y + h, x + w - r, y + h, r);
    ctx.lineTo(x + r, y + h);
    ctx.arcTo(x, y + h, x, y + h - r, r);
    ctx.lineTo(x, y + r);
    ctx.arcTo(x, y, x + r, y, r);
    ctx.closePath();
  },

  goHome() {
    this.clearFeedbackTimer();
    wx.switchTab({ url: '/pages/index/index' });
  },

  onUnload() {
    this.clearFeedbackTimer();
  }
});
