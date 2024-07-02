import React, {useState} from "react";
import { Tabs ,TabsProps } from "antd";
import {WechatQr} from "../wechat-qr";
import {MobileLogin} from "./mobile-login";
import {AccountLogin} from "./account-login";


export const LoginScreen = () => {

    const onChange = (key: string) => {
        //console.log(key);
    };

    const items: TabsProps['items'] = [
        {
            key: '1',
            label: `微信扫码`,
            children: <WechatQr />,
            disabled: true,
        },
        {
            key: '2',
            label: `手机号登录`,
            children: <MobileLogin  />,
        },
        {
            key: '3',
            label: `帐号登录`,
            children: <AccountLogin  />,
        },
    ];

    return (<Tabs defaultActiveKey="2"  items={items} onChange={onChange} />);
};
