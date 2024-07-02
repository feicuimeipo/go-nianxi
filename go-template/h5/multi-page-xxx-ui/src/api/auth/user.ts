import request from '@/common/utils/request'

export function login(data:{}) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

export function refreshToken() {
  return request({
    url: '/auth/refreshToken',
    method: 'post'
  })
}

export function logout(token:string|null) {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

// 获取当前登录用户信息
export function getInfo(token:string|null) {
  return request({
    url: '/auth/info',
    method: 'post'
  })
}


// 更新用户登录密码
export function changePwd(data:{}) {
  return request({
    url: '/auth/changePwd',
    method: 'put',
    data
  })
}


// 更新用户
export function updateUserById(id:string, data:{}) {
  return request({
    url: '/auth/update/' + id,
    method: 'patch',
    data
  })
}


