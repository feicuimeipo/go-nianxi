import { NavMenuTree} from "@/components/menu/menu";
import {DownMenuTree} from "@/components/menu/down-menu";


const WorkstationData = {
    navMenuList :  [
        {id:1,name: "帐号信息",path: '/workstation/accountInfo',    icon:'UserOutlined'},
        {id:2,name: "使用详情",path: '/workstation/useLog',         icon:'InfoCircleOutlined'},
        {id:3,name: "订单信息",path: '/workstation/orderInfo',     icon:'PayCircleOutlined'}
    ] as NavMenuTree[],
    userPanelDownMenuList :  [
        {id:1,name: "帐号信息",path: '/workstation/accountInfo',    icon:'UserOutlined'},
        {id:2,name: "使用详情",path: '/workstation/useLog',         icon:'InfoCircleOutlined'},
        {id:3,name: "订单信息",path: '/workstation/orderInfo',      icon:'PayCircleOutlined'},
        {id:4,name: "退出登录",path: '/workstation/exist',          icon:'LogoutOutlined'}
    ] as DownMenuTree[]
}

export default WorkstationData
