import React, {useState} from "react";
import {Tabs, TabsProps} from "antd";
import {Step1MobileRegister} from "./step1-mobile-register";
import {Step2AccountRegister} from "@/unauthenticated-app/register/step2-account-register";
import {TabProps} from "@/unauthenticated-app";

export const RegisterScreen = () => {

    const [tabProps, setTabProps] = useState<TabProps>({
        key: '2',
        mobile: '',
        userId: -1
    })


    const items: TabsProps['items'] = [
        {
            key: '2',
            label: `用户注册`,
            children: <Step1MobileRegister setTabProps={setTabProps} />,
            disabled: true,
        }
        ,
        {
            key: '3',
            label: `帐号信息`,
            children: <Step2AccountRegister tabProps={tabProps} />,
            disabled: true,
        },
    ];

    const onChange = (key: string) => {
        items.forEach(item=>{
            if (item.key !== key){
                item.disabled = true;
            }
        })
    };

    return (<>
         <Tabs activeKey={tabProps.key}  items={items} onChange={onChange} />
         </>
       )
};
