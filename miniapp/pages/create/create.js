const app = getApp();
Page({
  data: { title: '', text: '' },
  onTitle(e) { this.setData({ title: e.detail.value }); },
  onText(e) { this.setData({ text: e.detail.value }); },
  submit() {
    if (this.data.text.length < 50) {
      wx.showToast({ title: '文本至少50字', icon: 'none' }); return;
    }
    const t = this.data.title || '未命名牌组';
    wx.showLoading({ title: 'AI 生成中...' });
    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      method: 'POST',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      data: { title: t, text: this.data.text },
      success: (res) => {
        wx.hideLoading();
        wx.showToast({ title: '创建成功！' });
        setTimeout(() => wx.navigateBack(), 1000);
      },
      fail: () => { wx.hideLoading(); wx.showToast({ title: '网络错误', icon: 'none' }); }
    });
  }
});