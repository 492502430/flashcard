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

        this.setData({
          filePath: file.path,
          fileName: file.name,
          fileSize: size,
          text: ''
        });

        // Auto-set deck name from filename
        if (!this.data.title) {
          const name = file.name.replace(/\.(pdf|docx?|txt)$/i, '');
          this.setData({ title: name });
        }

        // Auto-upload and extract text
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
          if (data.text) {
            this.setData({
              text: data.text.substring(0, 5000),
              uploading: false
            });
            wx.showToast({ title: '文件解析成功', icon: 'success' });
          }
        } catch (e) {
          wx.showToast({ title: '文件解析失败', icon: 'none' });
          this.setData({ uploading: false });
        }
      },
      fail: () => {
        wx.showToast({ title: '上传失败', icon: 'none' });
        this.setData({ uploading: false });
      }
    });
  },

  removeFile() {
    this.setData({
      filePath: '',
      fileName: '',
      fileSize: '',
      text: ''
    });
  },

  submit() {
    if (this.data.text.length < 50 && !this.data.uploading) {
      wx.showToast({ title: '文本至少50字', icon: 'none' });
      return;
    }

    if (this.data.uploading) {
      wx.showToast({ title: '文件正在上传中...', icon: 'none' });
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
        wx.showToast({ title: '创建成功！', icon: 'success' });
        setTimeout(() => wx.navigateBack(), 1000);
      },
      fail: () => {
        this.setData({ generating: false });
        wx.showToast({ title: '创建失败，请重试', icon: 'none' });
      }
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
