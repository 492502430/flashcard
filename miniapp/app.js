App({
  globalData: {
    token: '',
    userInfo: null,
    apiBase: 'http://192.168.0.103:8088',
    needAuth: true
  },

  onLaunch() {
    const token = wx.getStorageSync('token');
    if (token) {
      this.globalData.token = token;
    }
    const userInfo = wx.getStorageSync('userInfo');
    if (userInfo) {
      this.globalData.userInfo = userInfo;
      this.globalData.needAuth = false;
    }
    this.login();
  },

  login() {
    wx.login({
      success: (res) => {
        const code = res.code || 'dev-' + Date.now();
        wx.request({
          url: `${this.globalData.apiBase}/api/auth/login`,
          method: 'POST',
          data: { code: code },
          success: (resp) => {
            if (resp.data && resp.data.token) {
              this.globalData.token = resp.data.token;
              wx.setStorageSync('token', resp.data.token);
            }
          }
        });
      },
      fail: () => {
        wx.request({
          url: `${this.globalData.apiBase}/api/auth/login`,
          method: 'POST',
          data: { code: 'dev-' + Date.now() },
          success: (resp) => {
            if (resp.data && resp.data.token) {
              this.globalData.token = resp.data.token;
              wx.setStorageSync('token', resp.data.token);
            }
          }
        });
      }
    });
  }
});
