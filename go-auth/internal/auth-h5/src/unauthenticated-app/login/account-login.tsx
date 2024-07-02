import React, {useState} from "react";
import {Form, Input, message} from "antd";
import {LongButton} from "../index";
import { useAsync } from "@/utils/use-async";
import VerifyPointFixed from "@/components/verify/verifyPointFixed";
import {CaptchaResult} from "@/types/captcha-type";
import {useAuth} from "@/context/auth-context";
import {AccountLoginForm} from "@/types/auth-types";
import {KeyOutlined, UserOutlined} from "@ant-design/icons";
import {ErrorBox} from "@/components/lib";


export const AccountLogin = () => {
    const [form] = Form.useForm();
    const {login} = useAuth();
    const {run, isLoading} = useAsync(undefined, {throwOnError: true});
    const [error, setError] = useState<Error | null>(null);


    //验证码
    const [verify,setVerify] = useState<CaptchaResult|null>(null);
    const [verifyShow,setVerifyShow] = useState<boolean>(false);
    const verifyPointFixedChild = (data:CaptchaResult|null|"close") => {
        if (data==="close"){
            setVerifyShow(false)
            setVerify(null)
        }else{
            setVerify(data)
            if (data!=null){
                setVerifyShow(false)
                const formData = {
                    loginName: form.getFieldValue("username"),
                    password: form.getFieldValue("password"),
                    captcha: data
                } as AccountLoginForm
                try {
                    run(login(formData));
                    //关闭对话框
                } catch (e: any) {
                    setError(e);
                }

            }else{
                setVerify(null)
                setVerifyShow(true)
            }
        }
    }

    const handleSubmit = async () => {
        if (verify){
            verifyPointFixedChild(verify)
        }else{
            setVerify(null)
            setVerifyShow(true)
        }
    };


    return  <> <ErrorBox error={error} />
            <Form form={form} onFinish={handleSubmit} autoComplete={"on"}>
            <Form.Item
                name={"username"}
                rules={[{ required: true, message: "请输入用户名" }]}
            >
                <Input prefix={<UserOutlined />} placeholder={"用户名"} type="text" id={"username"}  />
            </Form.Item>
            <Form.Item
                name={"password"}
                rules={[{ required: true, message: "请输入密码" }]}
            >
                <Input prefix={<KeyOutlined />} placeholder={"密码"} type="password" id={"password"}  />
            </Form.Item>
            <Form.Item>
                <LongButton loading={isLoading} htmlType={"submit"} type={"primary"}>登录</LongButton>
            </Form.Item>
            <VerifyPointFixed isPointShow={verifyShow}  verifyPointFixedChild={verifyPointFixedChild}/>
         </Form>
    </>
}
