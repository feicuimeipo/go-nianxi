import React, {useState} from "react";
import { Tabs ,TabsProps } from "antd";
import {FindUser} from "@/unauthenticated-app/forget-password/step1-confirm-user";
import {Step2ResetPassword} from "@/unauthenticated-app/forget-password/step2-reset-password";
import {TabProps} from "@/unauthenticated-app";


export const ForgetPasswordScreen = () => {
    const [tabProps, setTabProps] = useState<TabProps>({
        key: '1',
        mobile: '',
        userId: -1
    })


    const items: TabsProps['items'] = [
        {
            key: '1',
            label: `忘记密码`,
            children: <FindUser   setTabProps={setTabProps} />,
        },
        {
            key: '2',
            label: `重置密码`,
            children: <Step2ResetPassword tabProps={tabProps} />,
            disabled: true
        },
    ];

    const onChange = (key: string) => {
        items.forEach(item=>{
            if (item.key !== key){
                item.disabled = true;
            }
        })
    };

    return (<Tabs defaultActiveKey="1"  activeKey={tabProps.key} items={items} onChange={onChange} />);
};
