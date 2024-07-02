import request from "@/utils/request";


export function getPicture(params) {
  return request('/captcha/get',{
    method: 'POST',
    data: params
  })
}

export function reqCheck(params) {
  return request('/captcha/check',{
    method: 'POST',
    data: params,
  })
}


