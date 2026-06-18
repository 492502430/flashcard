const app = getApp();

Page({
  data: {
    title: '',
    text: '',
    generating: false
  },

  onTitle(e) {
    this.setData({ title: e.detail.value });
  },

  onText(e) {
    this.setData({ text: e.detail.value });
  },

  submit() {
    const { text, title } = this.data;

    if (text.length < 50) {
      wx.showToast({ title: '文本至少需要50字', icon: 'none' });
      return;
    }

    const deckTitle = title.trim() || '未命名牌组';

    this.setData({ generating: true });

    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      method: 'POST',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      data: { title: deckTitle, text },
      success: () => {
        this.setData({ generating: false });
        wx.showToast({ title: '创建成功！', icon: 'success' });
        setTimeout(() => {
          wx.navigateBack();
        }, 1000);
      },
      fail: () => {
        this.setData({ generating: false });
        wx.showToast({ title: '网络错误，请重试', icon: 'none' });
      }
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
