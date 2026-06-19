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
    generating: false,
    submitFailed: false
  },

  onTitle(e) { this.setData({ title: e.detail.value }); },
  onText(e) { this.setData({ text: e.detail.value }); },

  chooseFile() {
    wx.chooseMessageFile({
      count: 1,
      type: 'file',
      extension: ['pdf', 'doc', 'docx', 'txt', 'png', 'jpg', 'jpeg'],
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
    let token = app.globalData.token || wx.getStorageSync('token');
    console.log('[upload] token exists:', !!token);

    if (!token) {
      // Token missing — re-login and retry
      console.log('[upload] No token, re-logging...');
      app.login();
      wx.showLoading({ title: '登录中...' });
      setTimeout(() => {
        wx.hideLoading();
        let retryToken = app.globalData.token || wx.getStorageSync('token');
        if (retryToken) {
          app.globalData.token = retryToken;
          this.setData({ uploading: true, uploadProgress: 20 });
          this.doUpload(filePath, retryToken);
        } else {
          wx.showToast({ title: '登录失败，请重启小程序', icon: 'none' });
          this.setData({ uploading: false, filePath: '', fileName: '', fileSize: '' });
        }
      }, 2500);
      return;
    }

    app.globalData.token = token;
    this.setData({ uploading: true, uploadProgress: 20 });
    this.doUpload(filePath, token);
  },

  doUpload(filePath, token) {
    const fs = wx.getFileSystemManager();
    const ext = (this.data.fileName || '').split('.').pop().toLowerCase();

    // Read file as base64
    fs.readFile({
      filePath: filePath,
      encoding: 'base64',
      success: (res) => {
        console.log('[upload] File read, sending via wx.request...');
        wx.request({
          url: app.globalData.apiBase + '/api/upload',
          method: 'POST',
          header: { 
            'Authorization': 'Bearer ' + token,
            'Content-Type': 'application/json'
          },
          data: {
            filename: this.data.fileName,
            content: res.data,
            encoding: 'base64'
          },
          success: (resp) => {
            this.setData({ uploadProgress: 100 });
            const data = resp.data || {};
            const extracted = (data.text || '').trim();
            if (extracted) {
              this.setData({ text: extracted.substring(0, 5000), uploading: false });
              wx.showToast({ title: `已解析 ${extracted.length} 字`, icon: 'success' });
            } else {
              wx.showToast({ title: '文件没有可识别的文字', icon: 'none' });
              this.setData({ uploading: false, filePath: '', fileName: '', fileSize: '' });
            }
          },
          fail: (err) => {
            console.error('[upload] HTTP error:', JSON.stringify(err));
            wx.showToast({ title: '上传失败: ' + (err.errMsg || '网络错误'), icon: 'none' });
            this.setData({ uploading: false, filePath: '', fileName: '', fileSize: '' });
          }
        });
      },
      fail: (err) => {
        console.error('[upload] ReadFile error:', JSON.stringify(err));
        wx.showToast({ title: '读取文件失败', icon: 'none' });
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
    this.setData({ generating: true, submitFailed: false });

    wx.request({
      url: app.globalData.apiBase + '/api/decks',
      method: 'POST',
      header: { Authorization: 'Bearer ' + (app.globalData.token || wx.getStorageSync('token')) },
      data: { title: deckTitle, text: this.data.text },
      success: (res) => {
        this.setData({ generating: false, submitFailed: false });
        if (res.data && res.data.id) {
          wx.showToast({ title: 'AI 正在生成卡片...', icon: 'none' });
          setTimeout(() => {
            wx.navigateTo({ url: '/pages/deck-detail/deck-detail?id=' + res.data.id });
          }, 500);
        } else {
          wx.showToast({ title: '创建成功！', icon: 'success' });
          setTimeout(() => wx.navigateBack(), 800);
        }
      },
      fail: (err) => {
        console.error('Create deck failed:', err);
        this.setData({ generating: false, submitFailed: true });
        wx.showToast({ title: '生成失败，请点击重试', icon: 'none', duration: 3000 });
      }
    });
  },

  retrySubmit() {
    this.setData({ submitFailed: false });
    this.submit();
  },

  goBack() {
    wx.navigateBack();
  }
});
