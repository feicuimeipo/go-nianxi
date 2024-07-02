
export interface CaptchaResult {
    captchaType: string,
    pointJson: string,
    token: string,
    clientUid: string,
    ts: number
}


// export interface VerifyProps {
//     isPointShow: boolean
//     verifyResult: boolean
//     verifyData: CaptchaCheckResult|null
// }
//
// export const verifyPropsDefaultValue:VerifyProps = {
//     isPointShow: false,
//     verifyResult: false,
//     verifyData: null
// }
