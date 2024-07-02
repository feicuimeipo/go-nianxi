import { constantRoutes } from '@/router'
import { getUserMenuTreeByUserId } from '@/api/system/menu'
import { GetAppsByUserId } from '@/api/system/app'
import Layout from '@/layout'
import { isExternal } from '@/utils/validate'
import path from 'path'
import { getToken } from '@/utils/auth'

const resolvePath = (routePath, basePath) => {
  if (isExternal(routePath)) {
    return routePath
  }

  if (isExternal(basePath)) {
    if (basePath !== routePath && routePath && routePath !== '' && routePath !== '#') {
      if (basePath.endsWith('/')) {
        basePath = basePath.substring(0, basePath.length - 1)
      }
      if (routePath.startsWith('/')) {
        routePath = routePath.substring(1)
      }
      return basePath + '/' + routePath
    } else {
      return basePath
    }
  }
  return path.resolve(basePath, routePath)
}

export const getRoutesFromMenuTree = (menuTree) => {
  const routes = []
  menuTree.forEach(menu => {
    if (menu.children && menu.children.length > 0) {
      menu.children = getRoutesFromMenuTree(menu.children)
    }

    let path = menu.path
    let component = menu.component
    let iframeUrl = resolvePath(menu.basePath, menu.path)
    if (!isExternal(iframeUrl)) {
      iframeUrl = ''
    } else {
      if (path.startsWith('#')) {
        path = path.substring(1)
      }
      if (path.startsWith('/')) {
        path = path.substring(1)
      }
      if (menu.children && menu.children.length > 0) {
        component = 'Layout'
      } else {
        component = '/external/index'
      }
      path = '/external/p_' + menu.ID

      if (iframeUrl.indexOf('?') > -1) {
        iframeUrl += '&authToken=' + getToken()
      } else {
        iframeUrl += '?authToken=' + getToken()
      }
    }

    routes.push({
      path: path,
      name: menu.name,
      component: loadComponent(component),
      hidden: menu.hidden === 1,
      redirect: menu.redirect,
      alwaysShow: menu.alwaysShow === 1,
      children: menu.children,
      meta: {
        name: menu.name,
        title: menu.title,
        icon: menu.icon,
        noCache: menu.noCache === 1,
        breadcrumb: menu.breadcrumb === 1,
        activeMenu: menu.activeMenu,
        iframeUrl: iframeUrl
      }
    })
  })
  return routes
}

export const loadComponent = (component) => {
  if (component === '' || component === 'Layout') {
    // 组件不存在使用默认布局
    return Layout
  }
  // 动态获取组件
  return (resolve) => require([`@/views${component}`], resolve)
}

const state = {
  routes: [],
  addRoutes: [],
  apps: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  },
  SET_APPS: (state, apps) => {
    state.apps = apps
  }
}

const actions = {
  generateRoutes({ commit }, userinfo) {
    return new Promise((resolve, reject) => {
      let accessedRoutes = []
      // 获取菜单树
      console.log('userinfo.currentAppId=', userinfo.currentAppId)
      getUserMenuTreeByUserId(userinfo.id, userinfo.currentAppId).then(res => {
        const { data } = res
        const menuTree = data.menuTree

        accessedRoutes = getRoutesFromMenuTree(menuTree)
        commit('SET_ROUTES', accessedRoutes)
        resolve(accessedRoutes)
      }).catch(err => {
        reject(err)
      })
    })
  },
  getApps({ commit }, userinfo) {
    return new Promise((resolve, reject) => {
      // 获取应用权
      GetAppsByUserId(userinfo.id).then(res => {
        const { data } = res
        const apps = data
        console.log(JSON.stringify(apps))
        commit('SET_APPS', apps)
        resolve(apps)
      }).catch(err => {
        reject(err)
      })
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
