const app = getApp();

Page({
  data: {
    title: '',
    text: '',
    filePath: '',
    fileName: '',
    fileSize: '',
    uploading: false,
    uploadProgress: 0,
    generating: false
  },

  onTitle(e) { this.setData({ title: e.detail.value }); },
  onText(e) { this.setData({ text: e.detail.value }); },

  chooseFile() {
    wx.chooseMessageFile({
      count: 1,
      type: 'file',
      extension: ['pdf', 'doc', 'docx', 'txt'],
      success: (res) => {
        const file = res.tempFiles[0];
        const size = file.size > 1024 * 1024
          ? (file.size / 1024 / 1024).toFixed(1) + ' MB'
          : (file.size / 1024).toFixed(0) + ' KB';

        if (!this.data.title) {
          const name = file.name.replace(/\.(pdf|docx?|txt)$/i, '');
          this.setData({ title: name });
        }

        this.setData({
          filePath: file.path,
          fileName: file.name,
          fileSize: size,
          text: ''
        });

        this.uploadAndExtract(file.path);
      }
    });
  },

  uploadAndExtract(filePath) {
    this.setData({ uploading: true, uploadProgress: 20 });

    wx.uploadFile({
      url: app.globalData.apiBase + '/api/upload',
      filePath: filePath,
      name: 'file',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      success: (res) => {
        this.setData({ uploadProgress: 100 });
        try {
          const data = JSON.parse(res.data);
          const extracted = (data.text || '').trim();
          if (extracted) {
            this.setData({
              text: extracted.substring(0, 5000),
              uploading: false
            });
            wx.showToast({ title: `已解析 ${extracted.length} 字`, icon: 'success' });
          } else {
            wx.showToast({ title: '文件没有可识别的文字', icon: 'none' });
            this.setData({ uploading: false, filePath: '', fileName: '', fileSize: '' });
          }
        } catch (e) {
          wx.showToast({ title: '文件解析失败', icon: 'none' });
          this.setData({ uploading: false, filePath: '', fileName: '', fileSize: '' });
        }
      },
      fail: (err) => {
        console.error('Upload failed:', err);
        wx.showToast({ title: '上传失败，请检查网络', icon: 'none' });
        this.setData({ uploading: false, filePath: '', fileName: '', fileSize: '' });
      }
    });
  },

  removeFile() {
    this.setData({ filePath: '', fileName: '', fileSize: '', text: '' });
  },

  submit() {
    // Still uploading — wait
    if (this.data.uploading) {
      wx.showToast({ title: '文件正在解析中...', icon: 'none' });
      return;
    }

    // No content at all
    if (!this.data.text || this.data.text.trim().length < 50) {
      if (this.data.filePath) {
        wx.showToast({ title: '文件解析结果太短，请手动粘贴文本', icon: 'none' });
      } else {
        wx.showToast({ title: '请输入至少 50 字内容', icon: 'none' });
      }
      return;
    }

    const deckTitle = this.data.title || '未命名牌组';
    this.setData({ generating: true });

    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      method: 'POST',
      header: { Authorization: 'Bearer ' + app.globalData.token },
      data: { title: deckTitle, text: this.data.text },
      success: (res) => {
        this.setData({ generating: false });
        if (res.data && res.data.id) {
          wx.showToast({ title: `已生成 ${res.data.card_count || ''} 张卡片`, icon: 'success' });
        } else {
          wx.showToast({ title: '创建成功！', icon: 'success' });
        }
        setTimeout(() => wx.navigateBack(), 1200);
      },
      fail: (err) => {
        console.error('Create deck failed:', err);
        this.setData({ generating: false });
        wx.showToast({ title: '创建失败，请重试', icon: 'none' });
      }
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
