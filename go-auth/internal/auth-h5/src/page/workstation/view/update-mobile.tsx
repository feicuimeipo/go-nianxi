import React, {useState} from "react";
import {useAuth} from "@/context/auth-context";
import {Button,  Form, Input, Space} from "antd";
import styled from "@emotion/styled";
import {
    sendSmsVerifyCode, updateMobile,
} from "@/api/auth/auth-provider";
import {
    SendEmailVerifyCodeForm, SendSmsVerifyCodeForm,
    UpdateMobileForm,
} from "@/types/auth-types";
import {CaptchaResult} from "@/types/captcha-type";
import {useTimer} from "@/utils/useTimer";
import VerifyPointFixed from "@/components/verify/verifyPointFixed";
import {useAsync} from "@/utils/use-async";

export const UserInfoMobile = (props:{userId:number|undefined,mobile:string|undefined,setError: (error:Error)=>void}) => {
    const [form] = Form.useForm();
    const auth = useAuth()
    const [editMobile, setEditMobile] = useState<boolean>(false)
    const { run, isLoading } = useAsync(undefined, { throwOnError: true });

    const onChangeMobile = async ({...values}) => {
        if (values.mobile === props.mobile) {
            // props.setError(new Error("手机号与原手机号相同相同!"))
            // return
        }

        if (values.smsCode === "") {
            props.setError(new Error("验证码不可以为空!"))
            return
        }

        const dataForm = {
            userId: props.userId,
            newMobile: values.mobile,
            smsCode: values.smsCode
        } as UpdateMobileForm


        try{
            await run(
                updateMobile(dataForm)
                    .then(async res => {
                        if (auth.user!=null) {
                            let newUser = auth.user
                            newUser.mobile = res.data.mobile
                            await auth.setUserData(newUser)
                        }
                        setEditMobile(false)
                    }).catch(reason => {
                    props.setError(new Error(reason.msg))
                })
            )
        }catch (e:any){
            console.error("try.cache:",e)
            props.setError(e)
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
                sendSmsVerifyCode({mobile: mobile,captcha: data,use: "update"} as SendSmsVerifyCodeForm).then(res=>{
                    setIsShow(false)
                    start()
                }).catch(reason => {
                    console.error(reason)
                    props.setError(new Error(reason.msg));
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
         <Space> {props.mobile?.substring(0,4)} *** {props.mobile?.substring(7,11)}
         <a href={"#"} onClick={()=> setEditMobile(true)}>修改手机号</a></Space>
             <Form form={form} onFinish={onChangeMobile}>
                {editMobile? <EditorContainer>

                    <Form.Item   name={"mobile"}  rules={[{ required: true, message: "请输入新手机号" }]} >
                        <Input placeholder={"新手机号"} type="text" id={"mobile"}    />
                    </Form.Item>

                    <Space>
                        <Form.Item  name={"smsCode"}  rules={[{required: true, message: "请输入验证码" }]}  >
                            <VerifyInput placeholder={"验证码"}  id={"smsCode"}  />
                        </Form.Item>
                        <Form.Item>
                            <Button htmlType="button" disabled={!IsShow}  onClick={onSendSms}> {IsShow? '发送' : ''+num} </Button>
                        </Form.Item>
                    </Space>

                    <Form.Item>
                        <Button  loading={isLoading} htmlType={"submit"} type={"primary"} >保存</Button>
                        &nbsp;<Button  onClick={()=>setEditMobile(false)}>取消 </Button>
                    </Form.Item>
                    <VerifyPointFixed isPointShow={verifyShow}  verifyPointFixedChild={verifyPointFixedChild}/>
            </EditorContainer>  :null}  </Form>
    </>
}



const EditorContainer = styled.div`
  padding-top: 5px;
  width: 30%;
`


export const VerifyInput = styled(Input)`
   width: 10em;
`

