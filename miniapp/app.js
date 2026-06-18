App({
  globalData: {
    token: '',
    userInfo: null,
    apiBase: 'http://127.0.0.1:8080'
  },

  onLaunch() {
    const token = wx.getStorageSync('token');
    if (token) {
      this.globalData.token = token;
    }
    // Always try to login to ensure fresh token
    this.login();
  },

  login() {
    wx.login({
      success: (res) => {
        if (!res.code) {
          console.error('wx.login returned no code');
          return;
        }
        wx.request({
          url: `${this.globalData.apiBase}/api/auth/login`,
          method: 'POST',
          data: { code: res.code },
          success: (resp) => {
            if (resp.data && resp.data.token) {
              this.globalData.token = resp.data.token;
              wx.setStorageSync('token', resp.data.token);
              console.log('Login OK');
            } else {
              console.error('Login response missing token:', resp.data);
            }
          },
          fail: (err) => {
            console.error('Login request failed:', err);
          }
        });
      },
      fail: (err) => {
        console.error('wx.login failed:', err);
      }
    });
  }
});
