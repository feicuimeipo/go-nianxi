import request from '@/utils/request'

// 获取接口列表
export function getApps(params) {
  return request({
    url: '/app/list',
    method: 'get',
    params
  })
}

// 获取树状列表
export function getAppTree() {
  return request({
    url: '/app/tree',
    method: 'get'
  })
}

// 获取树状列表
export function GetAppsByUserId(uid) {
  return request({
    url: '/app/get/' + uid,
    method: 'get'
  })
}

// 获取类别列表
export function getAppTypes() {
  return request({
    url: '/app/type/list',
    method: 'get'
  })
}

// 创建接口
export function createApp(data) {
  return request({
    url: '/app/create',
    method: 'post',
    data
  })
}

// 更新接口
export function updateAppById(Id, data) {
  return request({
    url: '/app/update/' + Id,
    method: 'patch',
    data
  })
}

// 批量删除接口
export function batchDeleteAppByIds(data) {
  return request({
    url: '/app/delete/batch',
    method: 'delete',
    data
  })
}
