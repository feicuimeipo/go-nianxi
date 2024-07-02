import React, { useState} from "react";
import {Button, Form, Input,  Space} from "antd";
import {LongButton, TabProps, VerifyInput} from "../index";
import {useAsync} from "@/utils/use-async";
import {useTimer} from "@/utils/useTimer";

import VerifyPointFixed from "@/components/verify/verifyPointFixed";

import {RegisterStep1Form, } from "@/types/auth-types";
import {register1, sendSmsVerifyCode} from "@/api/auth/auth-provider";
import {CaptchaResult} from "@/types/captcha-type";
import Agreement from "@/unauthenticated-app/agreement";
import {ErrorBox} from "@/components/lib";


export const Step1MobileRegister = ({setTabProps}: {setTabProps: (tabProps: TabProps) => void;}) => {
    const { run, isLoading } = useAsync(undefined, { throwOnError: true });
    const [form] = Form.useForm();
    const [error, setError] = useState<Error | null>(null);

    const handleSubmit = async (values) => {
        const formData = {
            mobile: values.mobile,
            smsCode: values.smsCode
        } as RegisterStep1Form

        try{
            await run(
                register1(formData).then(res => {
                    setTabProps({key: "3",mobile: formData.mobile,userId: res.data.id})
                }).catch(reason => {
                    console.error("register1.cache:",reason)
                    setError(new Error(reason.msg));
                })
            )
        }catch (e:any){
            console.error("try.cache:",e)
            setError(e)
        }
    }

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
                const mobile = form.getFieldValue("mobile")
                sendSmsVerifyCode({mobile: mobile,captcha: data}).then(res=>{
                    setIsShow(false)
                    start()
                }).catch(reason => {
                    console.error(reason)
                    setError(new Error(reason.msg));
                });
                setVerifyShow(false)
            }else{
                setVerify(null)
                setVerifyShow(true)
            }
        }
    }

    //手机短信验证码
    const [IsShow, setIsShow] = useState(true)
    const {num, start } = useTimer(60, () => {
        setIsShow(true)
    })
    const onSendSms = () => {
        form.validateFields(['mobile']).then(()=> {
            if (verify){
                verifyPointFixedChild(verify)
            }else{
                //console.log("2......")
                setVerify( null)
                setVerifyShow(true)
            }
        })
    }


    return  <>
            <ErrorBox error={error} />
            <Form form={form} onFinish={handleSubmit} autoComplete="on">
                    <Form.Item
                        name={"mobile"}
                        rules={[{ required: true, message: '请输入正确的手机号',pattern:new RegExp(/^1(3|4|5|6|7|8|9)\d{9}$/, "g")}]}
                    >
                        <Input addonBefore="86"  placeholder={"手机号"} id={"mobile"} />
                    </Form.Item>
                    <Space>
                        <Form.Item name={"smsCode"} id={"smsCode"}  rules={[{required: true, message: "请输入验证码" }]}  >
                            <VerifyInput placeholder={"验证码"}  id={"smsCode"}   />
                        </Form.Item>
                        <Form.Item>
                            <Button htmlType="button" disabled={!IsShow}  onClick={onSendSms}> {IsShow? '发送' : ''+num} </Button>
                        </Form.Item>
                    </Space>

                    <Form.Item>
                        <LongButton loading={isLoading} htmlType={"submit"} type={"primary"}  >
                            下一步
                        </LongButton>

                    </Form.Item>
                    <VerifyPointFixed isPointShow={verifyShow}  verifyPointFixedChild={verifyPointFixedChild}/>
                    <Agreement />
            </Form>
        </>

}
