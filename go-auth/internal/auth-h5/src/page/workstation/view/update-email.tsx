import React, {useCallback, useState} from "react";
import {useAuth} from "@/context/auth-context";
import {Button,  Form, Input, Space} from "antd";
import {ErrorBox} from "@/components/lib";
import styled from "@emotion/styled";
import {
    sendEmailVerifyCode, updateEmail,
} from "@/api/auth/auth-provider";
import {
    SendEmailVerifyCodeForm,
    UpdateEmailForm,

} from "@/types/auth-types";
import {CaptchaResult} from "@/types/captcha-type";
import {useTimer} from "@/utils/useTimer";
import {useMount} from "@/utils";
import VerifyPointFixed from "@/components/verify/verifyPointFixed";
import {useAsync} from "@/utils/use-async";

export const UserInfoEmail = (props:{userId:number|undefined,email:string|undefined,setError: (error:Error)=>void}) =>{
    const [form] = Form.useForm();
    const [error, setError] = useState<Error | null>(null);
    const [editEmail,setEditEmail] = useState<boolean>(false)
    const auth = useAuth()
    const { run, isLoading } = useAsync(undefined, { throwOnError: true });

    useMount(
        useCallback(() => {
            form.setFieldValue("email",props.email)
        }, [])
    );

    const handleSubmit = ({...values})=>{
        const verifyCode = values.verifyCode
        if (values.email==props.email){
            props.setError(new Error("邮箱号与原邮箱号上同相同!"))
        }
        const dataForm = {
            userId: props.userId,
            email: values.email,
            verifyCode: verifyCode
        } as UpdateEmailForm
        updateEmail(dataForm)
            .then(async res=>{
                if (auth.user!=null) {
                    let newUser = auth.user
                    newUser.email = res.data.email
                    await auth.setUserData(newUser)
                }
                setEditEmail(false)
            }).catch(reason => {
                props.setError(new Error(reason.msg))
        })
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
                const email = form.getFieldValue("email")
                sendEmailVerifyCode({email: email,captcha: data ,use: "update"} as SendEmailVerifyCodeForm).then(res=>{
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


    return <>
        <ErrorBox error={error} />
        <Space> {props.email}
            <a href={"#"} onClick={()=> setEditEmail(true)}>修改邮箱</a></Space>
        <Form form={form} onFinish={handleSubmit}>
        {editEmail? <EditorContainer>

                <Form.Item   name={"email"}  rules={[{ required: true, message: "请输入邮箱" }]} >
                    <MobileInput placeholder={"邮箱"} type="text" id={"email"}   />
                </Form.Item>

                <Space>
                    <Form.Item  name={"verifyCode"} id={"verifyCode"}  rules={[{required: true, message: "请输入验证码" }]}  >
                        <VerifyInput placeholder={"验证码"}  id={"verifyCode"}  />
                    </Form.Item>
                    <Form.Item>
                        <Button htmlType="button" disabled={!IsShow}  onClick={onSendSms}> {IsShow? '发送' : ''+num} </Button>
                    </Form.Item>
                </Space>

                <Form.Item>
                    <Button loading={isLoading} htmlType={"submit"} type={"primary"}>保存</Button>
                    &nbsp;<Button  onClick={()=>setEditEmail(false)}>取消 </Button>
                </Form.Item>
                 <VerifyPointFixed isPointShow={verifyShow}  verifyPointFixedChild={verifyPointFixedChild}/>
        </EditorContainer>  :null} </Form>
    </>
}



const EditorContainer = styled.div`
  padding-top: 5px;
  width: 30%;
`


export const VerifyInput = styled(Input)`
   width: 10em;
  //width: 270px;
`

export const MobileInput = styled(Input)`
   width: 15rem;
`
