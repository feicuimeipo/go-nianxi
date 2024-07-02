// 在真实环境中，如果使用firebase这种第三方auth服务的话，本文件不需要开发者开发
import {
  AccountLoginForm,
  ForgetPasswordFindUserForm,
  ForgetPasswordResetPasswordForm,
  RegisterStep1Form,
  RegisterStep2Form, SendEmailVerifyCodeForm, SendSmsVerifyCodeForm,
  SMSLoginForm, UpdateEmailForm,
  UpdateMobileForm,
  UpdateNickNameForm,
  UpdatePasswordForm,
  UpdateUserNameForm
} from "@/types/auth-types";
import {http} from "@/utils/http";
import request from "@/utils/request";



//reducer是一个函数(state, action) => newState：接收当前应用的state和触发的动作action，计算并返回最新的state
export const me = () => {
  return request( "/auth/me",{
    method: "GET"
  })
}


export const login = (form: AccountLoginForm) => {
  return request( "/auth/login",{
    method: "POST",
    data: form
  })
};

export const smsLogin = (params: SMSLoginForm) => {
  return request("/auth/smsLogin",{
    method: "POST",
     data: params
  })
};


export const register1 = async (form: RegisterStep1Form) => {
  return request("/auth/register/1",{
    method: "POST",
    data: form
  })
};

export const register2= (form: RegisterStep2Form) => {
  return request("/auth/register/2",{
    method: "POST",
    data: form
  })
};


export const logout = () =>{
  return request("/auth/logout",{
     method: "GET",
  })
}

export function updateNickname(form: UpdateNickNameForm) {
  return request("/user/update/nickname",{
    method: "PATCH",
    data: form
  })
}

export function updateUsername(param: UpdateUserNameForm) {
  return request("/user/update/username",{
    method: "PATCH",
    data: param
  })
}

export function updateMobile(param: UpdateMobileForm) {
  return request("/user/update/mobile",{
    method: "PATCH",
    data: param
  })
}

export function updateEmail(param: UpdateEmailForm) {
  return request("/user/update/email",{
    method: "PATCH",
    data: param
  })
}

export function updatePassword(param: UpdatePasswordForm) {
  return request("/user/update/password",{
    method: "PATCH",
    data: param
  })
}


//发送短信码
export function sendSmsVerifyCode(params:SendSmsVerifyCodeForm) {
  return request("/auth/sendSmsVerifyCode",{
    method: "POST",
    data: params
  })
}

//发送邮箱
export function sendEmailVerifyCode(params:SendEmailVerifyCodeForm) {
  return request("/auth/sendEmailVerifyCode",{
    method: "POST",
    data: params
  })
}



//忘记密码
export function forgetPasswordStep1(params:ForgetPasswordFindUserForm) {
  return http("/auth/findPassword/1",{
    method: "POST",
    data: params
  })

  // return http('auth/forgetPassword',{
  //   method: 'POST',
  //   data: params,
  // })
}

//重置同密码
export function forgetPasswordStep2_resetPassword(params:ForgetPasswordResetPasswordForm) {
  return http("/auth/findPassword/2",{
    method: "POST",
    data: params
  })
}

