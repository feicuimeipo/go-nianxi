import request from '@/utils/request'

// 获取菜单树
export function getMenuTree(appId) {
  return request({
    url: '/menu/tree/' + appId,
    method: 'get'
  })
}

// 获取菜单列表
export function getMenus(appId) {
  return request({
    url: '/menu/list/' + appId,
    method: 'get'
  })
}

// 创建菜单
export function createMenu(data) {
  return request({
    url: '/menu/create',
    method: 'post',
    data
  })
}

// 更新菜单
export function updateMenuById(Id, data) {
  return request({
    url: '/menu/update/' + Id,
    method: 'patch',
    data
  })
}

// 批量删除菜单
export function batchDeleteMenuByIds(data) {
  return request({
    url: '/menu/delete/batch',
    method: 'delete',
    data
  })
}

// 获取用户的可访问菜单列表
export function getUserMenusByUserId(Id) {
  return request({
    url: '/menu/access/list/' + Id,
    method: 'get'
  })
}

// 获取用户的可访问菜单树
export function getUserMenuTreeByUserId(Id, appId) {
  return request({
    url: '/menu/access/tree/' + Id + '/' + appId,
    method: 'get'
  })
}
