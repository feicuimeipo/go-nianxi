<template>
  <div class="login-container" >
    <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" autocomplete="on" label-position="left">

      <div class="title-container">
        <h3 class="title">Go Nianxi </h3>
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <svg-icon icon-class="user" />
        </span>
        <el-input
          ref="username"
          v-model="loginForm.username"
          placeholder="用户名"
          name="username"
          type="text"
          tabindex="1"
          autocomplete="on"
        />
      </el-form-item>

      <el-tooltip v-model="capsTooltip" content="Caps lock is On" placement="right" manual>
        <el-form-item prop="password">
          <span class="svg-container">
            <svg-icon icon-class="password" />
          </span>
          <el-input
            :key="passwordType"
            ref="password"
            v-model="loginForm.password"
            :type="passwordType"
            placeholder="密码"
            name="password"
            tabindex="2"
            autocomplete="on"
            @keyup.native="checkCapslock"
            @blur="capsTooltip = false"
            @keyup.enter.native="handleLogin"
          />
          <span class="show-pwd" @click="showPwd">
            <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>
      </el-tooltip>

      <el-button :loading="loading" type="primary" class="login-btn" @click.native.prevent="handleLogin">Login</el-button>

    </el-form>

  </div>
</template>

<script>
import JSEncrypt from 'jsencrypt'

export default {
  name: 'Login',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error('The password can not be less than 6 digits'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: '',
        password: ''
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur' }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }]
      },
      passwordType: 'password',
      publicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1CkCsyd0rKlJXGCWhA3v
Rlyy3CEKvqUIj4M2O50ed/PhViF5/KcjSG7m+SjIkIz8ouD5A6Y3023Rt5ICsykr
eTajOp3xMRPbkxo+37Nyi9B5WwrBRg1c4IfvBjq4zc2Ab276dx+yWhuW3AYKStRo
rXLRGQdu3FZbvwTXEKKb5DL1ZlfPrwOjbRXYmYxbF0EIyhMJU+DtabKosc8Va/19
1wBcEn7pqOfNmHav2fEvPVo5BByDVXl/94SBdQTrRIEa31XA8m9u6t36ji0kr78T
zY7JHOJ6XTH6TzDB3M/CzPuyWdGB3XYRtJKGSQlNGhIw3tC91w9VgfBSmaoBfoWR
LwIDAQAB
-----END PUBLIC KEY-----`,
      capsTooltip: false,
      loading: false,
      redirect: undefined,
      otherQuery: {}
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        const query = route.query
        if (query) {
          this.redirect = query.redirect
          this.otherQuery = this.getOtherQuery(query)
        }
      },
      immediate: true
    }
  },
  created() {
    // window.addEventListener('storage', this.afterQRScan)
  },
  mounted() {
    if (this.loginForm.username === '') {
      this.$refs.username.focus()
    } else if (this.loginForm.password === '') {
      this.$refs.password.focus()
    }
  },
  destroyed() {
    // window.removeEventListener('storage', this.afterQRScan)
  },
  methods: {
    checkCapslock(e) {
      const { key } = e
      this.capsTooltip = key && key.length === 1 && (key >= 'A' && key <= 'Z')
    },
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          // 密码RSA加密处理
          const encryptor = new JSEncrypt()
          // 设置公钥
          encryptor.setPublicKey(this.publicKey)
          // 加密密码
          const encPassword = encryptor.encrypt(this.loginForm.password)
          const encLoginForm = { username: this.loginForm.username, password: encPassword }
          this.$store.dispatch('user/login', encLoginForm)
            .then(() => {
              console.log('login success')
              this.$router.push({ path: this.redirect || '/', query: this.otherQuery })
              this.loading = false
            })
            .catch(() => {
              this.loading = false
            })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    getOtherQuery(query) {
      return Object.keys(query).reduce((acc, cur) => {
        if (cur !== 'redirect') {
          acc[cur] = query[cur]
        }
        return acc
      }, {})
    }
  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #edf1f4;
$light_gray:#fff;
$cursor: #fff;
$bg: #edf1f4;
$dark_gray: #889aa4;
$form_bg: rgba(38, 142, 159, 0.3);
$form_input_color: rgba(0, 0, 0, 0.6);

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: $bg;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $form_input_color;
      height: 47px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $form_input_color !important;
      }
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: $bg;
    border-radius: 5px;
    color: $form_input_color;
    margin-left: 20px;
    margin-right: 20px;
  }

  .login-btn{
    width: calc(100% - 40px);
    height: 45px;
    margin-left: 20px;
    margin-right: 20px;
    margin-bottom: 30px;
    border: 0;
    font-size: 16px;
    font-weight: bold;
    // opacity: 0.9;
    border-radius: 20px;
    background-image: linear-gradient(to right, #eea2a2 0%, #7ac5d8 100%);
  }

  .login-btn:hover{
    color: greenyellow;
    background-image: linear-gradient(to right, #eea2a2 0%, #bbc1bf 19%, #57c6e1 42%, #b49fda 79%, #7ac5d8 100%);
  }
}
</style>

<style lang="scss" scoped>
$bg: #edf1f4;;
$dark_gray:#889aa4;
$light_gray:#fff;
$form_bg: rgba(38, 142, 159, 0.3);

.login-container {
  min-height: 100%;
  width: 100%;
  height: 100%;
  background-color: $bg;
  // overflow: hidden;
  background-size: cover;
  display: flex;
  justify-content: center;
  align-items: center;

  .login-form {
    // position: relative;
    width: 520px;
    max-width: 100%;
    // padding: 160px 35px 0;
    // margin: 0 auto;
    background-color: $form_bg;
    // overflow: hidden;
    border-radius: 8px;
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: rgba(0,0,0, 0.6);
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: rgba(0,0,0, 0.6);
      margin: 20px auto 20px auto;
      text-align: center;
      font-weight: bold;
      // letter-spacing: 2px;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: rgba(0,0,0,0.6);
    cursor: pointer;
    user-select: none;
  }

  .thirdparty-button {
    position: absolute;
    right: 0;
    bottom: 6px;
  }

  @media only screen and (max-width: 470px) {
    .thirdparty-button {
      display: none;
    }
  }
}
</style>
