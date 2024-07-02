import React, {CSSProperties, ReactNode} from "react";
import Icon from "@/components/icon";
import {Button, Dropdown, MenuProps} from "antd";
import styled from "@emotion/styled";



type MenuItem = Required<MenuProps>['items'][number];
function getItem(
    label: React.ReactNode,
    key?: React.Key | null,
    icon?: React.ReactNode,
    children?: MenuItem[]|null,
    type?: 'group',
    style?: CSSProperties
): MenuItem {
    return {
        label,
        key,
        icon,
        children,
        type,
        style,
    } as MenuItem;
}

export interface DownMenuTree {
    id: number,
    path: string,
    name: string,
    icon: string,
    type?: 'group',
    children: DownMenuTree[]
}


export interface DownMenuTreeProps {
    menuList: DownMenuTree[]
    handleBtn: ReactNode
    onClick: (item: DownMenuTree) => void
    style?: CSSProperties | undefined
    topPlusElement?: ReactNode | undefined
    bottomPlusElement?: ReactNode | undefined
    menuStyle?:  CSSProperties| undefined
}


export const DownMenuBar = ({props}:{props: DownMenuTreeProps}) => {
    const menuList:DownMenuTree[] = props.menuList;

    const result: MenuItem[] = []

    if (props.topPlusElement) {
        result.push(getItem("", "", props.topPlusElement, null,undefined))
    }

    if (!props.menuStyle){
        props.menuStyle = {
            padding: '10px'
        }
    }

    menuList.forEach(item => {
        if (item.children && item.children.length>0){
            const childrenItem = buildMenuTree(item.children, result, props.onClick);
            result.push(
                getItem(item.name, item.id, <Icon icon={item.icon}/>,childrenItem,undefined,props.menuStyle)
            )
        }else{
            result.push(getItem(<a href={"#"}  style={{paddingLeft: '5px'}} onClick={() => props.onClick(item)} > {item.name}  </a>,item.id,<Icon icon={item.icon} />,null,undefined,props.menuStyle))
        }
    })

    if (props.bottomPlusElement) {
        result.push(getItem("", "", props.bottomPlusElement, null,))
    }

    return  <DropdownContainer menu={{items: result }} placement="bottomLeft" >
                {props.handleBtn}
            </DropdownContainer>
}

const buildMenuTree = (downMenuTree:DownMenuTree[], result: MenuItem[],onClickItem: (item: DownMenuTree) => void):MenuItem[] =>{

    downMenuTree.forEach(item => {
        if (item.children && item.children.length > 0) {
            const childrenItem = buildMenuTree(item.children, result, onClickItem);
            result.push(getItem(item.name, item.id, <Icon icon={item.icon}/>,childrenItem))
        } else {
            result.push(getItem(<a href={"#"}  onClick={() => onClickItem(item)} > {item.name}  </a>, item.id,<Icon icon={item.icon} />))
        }
        return result;
    })
    return result;
}

const DropdownContainer = styled(Dropdown)`
   
`
