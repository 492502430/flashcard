const app = getApp();

Page({
  data: {
    current: 0,
    steps: [
      {
        title: 'AI 智能生成闪卡',
        desc: '上传 PDF、Word 或粘贴文本，AI 自动提取关键知识点并生成问答卡片，告别手动录入。',
        icon: 'spark'
      },
      {
        title: '科学间隔重复',
        desc: '基于遗忘曲线的智能调度算法，在最佳时机推送复习，让记忆更牢固、更持久。',
        icon: 'repeat'
      },
      {
        title: '开始高效记忆',
        desc: '每天只需几分钟，系统自动安排复习计划，轻松掌握任何学科内容。',
        icon: 'check'
      }
    ]
  },

  onSwiperChange(e) {
    this.setData({ current: e.detail.current });
  },

  finishOnboard() {
    wx.setStorageSync('onboarded', true);
    wx.switchTab({ url: '/pages/index/index' });
  },

  skipOnboard() {
    wx.setStorageSync('onboarded', true);
    wx.switchTab({ url: '/pages/index/index' });
  }
});
