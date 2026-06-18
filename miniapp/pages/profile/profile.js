Page({
  data: {
    menus: [
      { title: '我的订阅', action: 'subscribe' },
      { title: '关于闪卡记忆', action: 'about' },
    ]
  },
  onTap(e) {
    const action = e.currentTarget.dataset.action;
    if (action === 'about') {
      wx.showModal({ title: '闪卡记忆', content: 'AI 驱动的智能记忆工具，上传文本自动生成闪卡。' });
    }
  }
});