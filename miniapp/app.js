App({
  globalData: {
    token: '',
    userInfo: null,
    apiBase: 'http://localhost:8080'
  },

  onLaunch() {
    const token = wx.getStorageSync('token');
    if (!token) {
      this.login();
    } else {
      this.globalData.token = token;
    }
  },

  login() {
    wx.login({
      success: (res) => {
        wx.request({
          url: `${this.globalData.apiBase}/api/auth/login`,
          method: 'POST',
          data: { code: res.code },
          success: (resp) => {
            this.globalData.token = resp.data.token;
            wx.setStorageSync('token', resp.data.token);
          }
        });
      }
    });
  }
});
