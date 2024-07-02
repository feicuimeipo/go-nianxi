import React, {useCallback, useRef, useState} from "react";
import {bootstrapUser, useAuth} from "@/context/auth-context";
import {Button, Descriptions, Form, Input, message, Space} from "antd";
import styled from "@emotion/styled";
import {
    updatePassword,

    updateUsername
} from "@/api/auth/auth-provider";
import {
    UpdatePasswordForm,

    UpdateUserNameForm
} from "@/types/auth-types";
import {useMount} from "@/utils";
import {useAsync} from "@/utils/use-async";

export const AccountInfoUserName = (props:{userId:number|undefined,username:string|undefined,setError: (error:Error)=>void}) =>{
    const [usernameForm] = Form.useForm();
    const auth = useAuth();
    const [editUserName,setEditUserName] = useState<boolean>(false)
    const [editPassword,setEditPassword] = useState<boolean>(false)
    const { run, isLoading } = useAsync(undefined, { throwOnError: true });


    useMount(
        useCallback(() => {
            usernameForm.setFieldValue("username",props.username)
        }, [])
    );

    const onChangeUsername =  async ({...values})=>{
        const username = values.username
        if (username== props.username){
            props.setError(new Error("用户名相同!"))
        }
        const dataForm = {
            userId: props.userId,
            userName: values.username
        } as UpdateUserNameForm

        try{
          await run(
                updateUsername(dataForm)
                    .then(async res=>{
                        if (auth.user!=null) {
                            let newUser = auth.user
                            newUser.username = res.data.username
                            await auth.setUserData(newUser)
                        }
                        setEditUserName(false)
                    }).catch(reason => {
                        props.setError(new Error(reason.msg))
                    })
            )
        }catch (e:any){
            console.error("try.cache:",e)
            props.setError(e)
        }

    }

    const onUpdateChangePassword = async ({...values}) =>{
        if (values.newPassword ==""){
            props.setError(new Error("密码中不可空！"))
        }

        if (values.newPassword!==values.reNewPassword){
            props.setError(new Error("新密码与确认密码不匹配！"))
        }

        const dataForm = {
            userId: props.userId,
            oldPassword: values.oldPassword,
            newPassword: values.newPassword,
            reNewPassword: values.reNewPassword
        } as UpdatePasswordForm


        try{
            await run(
                updatePassword(dataForm)
                    .then(async res=>{
                        message.success("重置密码成功",100)
                        //await auth.setUserData(res.data)
                        setEditPassword(false)
                    }).catch(reason => {
                    props.setError(new Error(reason.msg))
                })
            )
        }catch (e:any){
            console.error("try.cache:",e)
            props.setError(e)
        }
    }

    return <>
        <Space>{props.username}
                <a href={"#"} onClick={()=> setEditUserName(true)}>修改用户名</a>
                <a href={"#"} onClick={()=> setEditPassword(true)}>重置密码</a>
            </Space>
            <Form form={usernameForm} onFinish={onChangeUsername}>
            {editUserName ? <EditorContainer>
                    <Space>
                        <Form.Item  name={"username"}  rules={[{ required: true, message: "请输入用户名" }]} >
                            <Input placeholder={"用户名"} type="text" id={"username"} />
                        </Form.Item>
                        <Form.Item>  <Button  loading={isLoading} htmlType={"submit"} type={"primary"} >保存</Button>   </Form.Item>
                        <Form.Item>  <Button  onClick={()=>setEditUserName(false)}>   取消 </Button>  </Form.Item>
                    </Space>
            </EditorContainer>:null}
            </Form>


            {editPassword ?
                <EditorContainer>
                    <Form  onFinish={onUpdateChangePassword} autoComplete={"on"}>

                            <Form.Item
                                name={"oldPassword"}
                                rules={[{ required: true, message: "请输入原密码" }]}
                            >
                                <Input placeholder={"密码"} type="password" id={"oldPassword"} />
                            </Form.Item>
                            <Form.Item
                                name={"newPassword"}
                                rules={[{ required: true, message: "请输入密码" }]}
                            >
                                <Input placeholder={"密码"} type="password" id={"newPassword"} />
                            </Form.Item>
                            <Form.Item
                                name={"reNewPassword"}
                                rules={[{ required: true, message: "重置新密码" }]}
                            >
                                <Input placeholder={"确认新密码"} type="password" id={"reNewPassword"} />
                            </Form.Item>
                         <Space>
                            <Form.Item>
                                <Button  loading={isLoading} htmlType={"submit"} type={"primary"} >
                                    保存
                                </Button>
                            </Form.Item>
                            <Form.Item>
                                <Button  onClick={()=>setEditPassword(false)}>
                                    取消
                                </Button>
                            </Form.Item>
                        </Space>
                    </Form>
                </EditorContainer>:null}
    </>
}

const Container = styled.div`
 padding: 10px;
 font-size: 1.5rem;
`

const EditorSpace = styled(Space)`
  padding-top: 5px;
`


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
