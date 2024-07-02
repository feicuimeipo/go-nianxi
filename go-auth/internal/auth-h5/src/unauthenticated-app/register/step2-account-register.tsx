import React, {useState} from "react";
import {Button, Form, Input} from "antd";
import {LongButton, TabProps} from "../index";
import { useAsync } from "@/utils/use-async";
import {RegisterStep2Form} from "@/types/auth-types";
import { register2} from "@/api/auth/auth-provider";
import {ErrorBox} from "@/components/lib";
import styled from "@emotion/styled";


export const Step2AccountRegister = ({tabProps}: {tabProps:TabProps}) => {
    const [form] = Form.useForm();
    const [error, setError] = useState<Error | null>(null);

    const { run, isLoading } = useAsync(undefined, { throwOnError: true });
    const [finish,setFinish] = useState<{status: boolean,username: string ,redirectUrl:string}>({status: false,username: '',redirectUrl: ""})


    const handleSubmit = async ({rePassword, ...values }: {
        username: string;
        password: string;
        rePassword: string;
        verificationCode: string
    }) => {
        if (tabProps.userId<=0 || tabProps.mobile===""){
            setError(new Error("用户编号或手机号不可以为空！"));
            return;
        }
        if (rePassword !== values.password) {
            setError(new Error("请确认两次输入的密码相同"));
            return;
        }

        const formData = {
            username: values.username,
            password: values.password,
            userId: tabProps.userId,
            rePassword: rePassword,
            mobile: tabProps.mobile
        } as RegisterStep2Form

        try{
            await run(
                register2(formData).then(async res => {

                    setFinish(
                        {
                            status: true,
                            username: values.username.length > 0 ? values.username : tabProps.mobile.substring(0, 3) + "***" + tabProps.mobile.substring(5, 7),
                            redirectUrl: res.data.redirectUrl
                        })
                }).catch(reason => {
                    console.error("register2.cache:", reason)
                    setError(new Error(reason.msg));
                })
            )
        }catch (e:any){
            console.error("try.cache:",e)
            setError(e)
        }
    };

    return  <>
        {finish.status? <><SuccessMessage> {finish.username} 恭喜您，注册成功！ </SuccessMessage></>: <>
                    <ErrorBox error={error} />
                    <Form form={form} onFinish={handleSubmit}>
                        <Form.Item
                            name={"username"}
                            rules={[{ required: true, message: "请输入用户名" }]} >
                            <Input placeholder={"用户名"} type="text" id={"username"} />
                        </Form.Item>
                        <Form.Item
                            name={"password"}
                            rules={[{ required: true, message: "请输入密码" }]}
                        >
                            <Input placeholder={"密码"} type="password" id={"password"} />
                        </Form.Item>
                        <Form.Item
                            name={"rePassword"}
                            rules={[{ required: true, message: "请确认密码" }]}
                        >
                            <Input placeholder={"确认密码"} type="password" id={"rePassword"} />
                        </Form.Item>
                        <Form.Item>
                            <LongButton loading={isLoading} htmlType={"submit"} type={"primary"}>
                                保存
                            </LongButton>
                        </Form.Item>
                    </Form>
                </>
            }
        </>
}


const SuccessMessage = styled.h2`
  margin-bottom: 1.2rem;
  color: rgb(94, 108, 132);
`;
