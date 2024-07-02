import request from '@/utils/request'

export function login(data) {
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

export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}
