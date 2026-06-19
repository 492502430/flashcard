const app = getApp();

Page({
  data: {
    templates: [],
    importing: null  // template id currently importing
  },

  onShow() {
    this.loadTemplates();
  },

  loadTemplates() {
    wx.request({
      url: app.globalData.apiBase + '/api/templates',
      success: (res) => {
        const templates = res.data || [];
        // Group by category
        const categories = [];
        const catMap = {};
        templates.forEach(t => {
          if (!catMap[t.category]) {
            catMap[t.category] = { category: t.category, items: [] };
            categories.push(catMap[t.category]);
          }
          catMap[t.category].items.push(t);
        });
        this.setData({ templates, categories });
      },
      fail: (err) => {
        console.error('加载模板列表失败:', err);
        wx.showToast({ title: '加载模板失败: ' + (err.errMsg || '网络异常'), icon: 'none', duration: 2500 });
      }
    });
  },

  importTemplate(e) {
    const id = e.currentTarget.dataset.id;
    const title = e.currentTarget.dataset.title;
    const token = app.globalData.token || wx.getStorageSync('token');

    if (!token) {
      wx.showToast({ title: '请先登录', icon: 'none' });
      return;
    }

    this.setData({ importing: id });

    wx.request({
      url: app.globalData.apiBase + '/api/templates/' + id + '/import',
      method: 'POST',
      header: { Authorization: 'Bearer ' + token },
      success: (res) => {
        if (res.statusCode === 201 && res.data && res.data.id) {
          wx.showToast({ title: '「' + title + '」已导入', icon: 'success' });
          // Navigate to the new deck
          setTimeout(() => {
            wx.navigateTo({ url: '/pages/deck-detail/deck-detail?id=' + res.data.id });
          }, 800);
        } else {
          const errMsg = (res.data && res.data.error) ? res.data.error : '导入失败';
          wx.showToast({ title: errMsg, icon: 'none' });
        }
      },
      fail: (err) => {
        console.error('导入模板失败:', err);
        wx.showToast({ title: '导入失败: ' + (err.errMsg || '网络异常，请稍后重试'), icon: 'none', duration: 2500 });
      },
      complete: () => {
        this.setData({ importing: null });
      }
    });
  },

  previewTemplate(e) {
    const id = e.currentTarget.dataset.id;
    // Show template cards in a modal-like view
    wx.request({
      url: app.globalData.apiBase + '/api/templates/' + id,
      success: (res) => {
        const tmpl = res.data;
        if (!tmpl || !tmpl.cards) {
          wx.showToast({ title: '加载模板详情失败', icon: 'none' });
          return;
        }
        // Show first few cards in a modal
        const cardList = tmpl.cards.slice(0, 5).map((c, i) =>
          (i + 1) + '. ' + c.question
        ).join('\n');
        const more = tmpl.cards.length > 5 ? '\n...共' + tmpl.cards.length + '张卡片' : '';
        wx.showModal({
          title: tmpl.title,
          content: cardList + more,
          confirmText: '导入此牌组',
          cancelText: '取消',
          success: (modalRes) => {
            if (modalRes.confirm) {
              this.importTemplate({ currentTarget: { dataset: { id: id, title: tmpl.title } } });
            }
          }
        });
      },
      fail: (err) => {
        console.error('加载模板详情失败:', err);
        wx.showToast({ title: '加载模板详情失败: ' + (err.errMsg || '网络异常'), icon: 'none', duration: 2500 });
      }
    });
  }
});
