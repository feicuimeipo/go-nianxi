import {CaptchaResult} from "@/types/captcha-type";


export interface AuthUser {
    ID: number
    token: string | undefined
    username: string
    nickname: string
    avatar: string
    mobile: string
    email: string
    type: string
}

export interface AccountLoginForm{
    loginName: string;
    password: string;
    captcha: CaptchaResult | null
}

export interface WechatQRLoginForm {
    openId: string;
    qrcode: string;
}

export interface SMSLoginForm {
    mobile: string
    smsCode: string
}


/**
 * 注册信息
 */
export interface RegisterStep1Form {
    mobile: string;
    smsCode: string;
}

/**
 * 注册信息，配置用户与密码
 */
export interface RegisterStep2Form {
    userId: number;
    username: string;
    password: string;
    rePassword: string;
    mobile: string;
}



/**
 * 找回密码
 */
export interface ForgetPasswordFindUserForm {
    mobile: string;
    smsCode: string;
    username: string;
}


/**
 * 找回密码-重置密码
 */
export interface ForgetPasswordResetPasswordForm {
    reNewPassWord: string;
    newPassWord: string;
    username: string;
    mobile: String;
}


/**
 * 修改用户信息
 */
export interface UpdateUserNameForm {
    userId: number;
    userName: string;
}

export interface UpdatePasswordForm {
    userId: number;
    oldPassword: string;
    newPassword: string;
    reNewPassword: string;
}

export interface UpdateNickNameForm {
    userId: number;
    nickName: string;
}

export interface UpdateMobileForm {
    userId: number;
    newMobile: string;
    smsCode: string;
}


export interface UpdateEmailForm {
    userId: number;
    email: string;
    verifyCode: string;
}

export interface SendEmailVerifyCodeForm {
    email: number
    captcha: CaptchaResult|null
    use?: string
}

export interface SendSmsVerifyCodeForm {
    mobile: number
    captcha: CaptchaResult|null
    use?: string
}
