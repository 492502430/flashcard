App({
  globalData: {
    token: '',
    userInfo: null,
    apiBase: 'http://192.168.0.103:8088'
  },

  onLaunch(options) {
    this.captureInviteCode(options);
    const token = wx.getStorageSync('token');
    if (token) {
      this.globalData.token = token;
      console.log('[app] Found stored token');
    }
    this.login();

    // Onboarding check — redirect first-time users
    const onboarded = wx.getStorageSync('onboarded');
    if (!onboarded) {
      console.log('[app] First launch — redirecting to onboard');
      wx.reLaunch({ url: '/pages/onboard/onboard' });
    }
  },

  onShow(options) {
    this.captureInviteCode(options);
  },

  captureInviteCode(options) {
    const inviteCode = options && options.query && options.query.invite_code;
    if (inviteCode) {
      wx.setStorageSync('pendingInviteCode', decodeURIComponent(inviteCode));
    }
  },

  login() {
    wx.login({
      success: (res) => {
        console.log('[app] wx.login code:', res.code ? res.code.substring(0,10) + '...' : 'EMPTY');
        const code = res.code || 'dev-fallback-' + Date.now();
        
        console.log('[app] POST to:', `${this.globalData.apiBase}/api/auth/login`);
        const inviteCode = wx.getStorageSync('pendingInviteCode') || '';
        wx.request({
          url: `${this.globalData.apiBase}/api/auth/login`,
          method: 'POST',
          data: { code: code, invite_code: inviteCode },
          success: (resp) => {
            console.log('[app] Login response status:', resp.statusCode);
            console.log('[app] Login response data:', JSON.stringify(resp.data).substring(0, 150));
            if (resp.data && resp.data.token) {
              this.globalData.token = resp.data.token;
              wx.setStorageSync('token', resp.data.token);
              console.log('[app] Login OK — token saved');
            } else {
              console.error('[app] Login response missing token');
            }
          },
          fail: (err) => {
            console.error('[app] Login HTTP fail:', JSON.stringify(err));
          }
        });
      },
      fail: (err) => {
        console.error('[app] wx.login failed:', JSON.stringify(err));
        this.fallbackLogin();
      }
    });
  },

  fallbackLogin() {
    console.log('[app] Using fallback login...');
    const inviteCode = wx.getStorageSync('pendingInviteCode') || '';
    wx.request({
      url: `${this.globalData.apiBase}/api/auth/login`,
      method: 'POST',
      data: { code: 'dev-' + Date.now(), invite_code: inviteCode },
      success: (resp) => {
        if (resp.data && resp.data.token) {
          this.globalData.token = resp.data.token;
          wx.setStorageSync('token', resp.data.token);
          console.log('[app] Fallback login OK');
        }
      }
    });
  }
});
