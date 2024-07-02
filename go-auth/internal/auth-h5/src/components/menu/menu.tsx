import React, { useState } from "react";
import Icon from "@/components/icon";
import { Menu, MenuProps} from 'antd';
import { useNavigate} from "react-router-dom";


type MenuItem = Required<MenuProps>['items'][number];
export interface NavMenuTree {
    id: number,
    path: string,
    name: string,
    icon: string,
    children: NavMenuTree[]
}



export const NavMenuBar = ({menuList}:{menuList:NavMenuTree[]}) => {
    const items: MenuItem[] = []


    menuList.forEach(item => {
        if (item.children && item.children.length>0){
            const childrenItem = buildMenuTree(item.children, items);
            items.push(getItem(<><Icon icon={item.icon}/>{item.name}</>, item.id, null,childrenItem))
        }else{
            items.push(getItem(item.name,item.path,<Icon icon={item.icon} />))
        }
    })


    const [openKeys, setOpenKeys] = useState(['sub1']);
    const rootSubmenuKeys = ['sub1', 'sub2', 'sub4'];
    const onOpenChange: MenuProps['onOpenChange'] = (keys) => {
        const latestOpenKey = keys.find((key) => openKeys.indexOf(key) === -1);
        if (rootSubmenuKeys.indexOf(latestOpenKey!) === -1) {
            setOpenKeys(keys);
        } else {
            setOpenKeys(latestOpenKey ? [latestOpenKey] : []);
        }
    };

    let navigate = useNavigate()
    const onClick = (e) =>{
        let path =e.key.startsWith("/")? e.key.substring(1): e.key;
        navigate(`/`+ path)
    }



    return (
        <>
        <Menu
                mode="inline"
                openKeys={openKeys}
                onOpenChange={onOpenChange}
                items={items}
                onClick={onClick}
            />
        </>
    )

}

function getItem(
    label: React.ReactNode,
    key?: React.Key | null,
    icon?: React.ReactNode,
    children?: MenuItem[]|undefined,
    type?: 'group',
): MenuItem {
    return {
        key,
        icon,
        children,
        label,
        type,
    } as MenuItem;
}

const buildMenuTree = (menuTrees:NavMenuTree[], items: MenuItem[]):MenuItem[] =>{
    menuTrees.forEach( item => {
        if (item.children && item.children.length > 0) {
            const childrenItem = buildMenuTree(item.children, items);
            items.push(getItem(item.name, item.path, <Icon icon={item.icon} />,childrenItem))
        } else {
            items.push(getItem(item.name,item.path,<Icon icon={item.icon} />))
        }
        return items;
    })
    return items;
}
