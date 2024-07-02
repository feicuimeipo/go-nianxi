import React, {useCallback, useState} from "react";
import {useAuth} from "@/context/auth-context";
import {Button,  Form, Input, Space} from "antd";

import {
    updateNickname,
} from "@/api/auth/auth-provider";
import {

    UpdateNickNameForm,


} from "@/types/auth-types";

import {useMount} from "@/utils";
import {EditorContainer} from "@/page/workstation/view/user-info";
import {useAsync} from "@/utils/use-async";


export const AccountInfoNickName = (props:{userId:number|undefined,nickname:string|undefined,setError: (error:Error)=>void}) =>{
    const [nickNameForm] = Form.useForm();
    const auth = useAuth();
    const [editNickName,setNickName] = useState<boolean>(false)
    const { run, isLoading } = useAsync(undefined, { throwOnError: true });


    useMount(
        useCallback(() => {
            nickNameForm.setFieldValue("nickname",props.nickname)
        }, [])
    );

    const onChangeNickname =  async ({...values})=>{
        if (values.nickname==props.nickname){
            props.setError(new Error("昵称名与原昵称名相同!"))
        }
        const dataForm = {
            userId: props.userId,
            nickName: values.nickname
        } as UpdateNickNameForm
        updateNickname(dataForm)
            .then(async res=>{
                if (auth.user!=null) {
                    let newUser = auth.user
                    newUser.nickname = res.data.nickname
                    await auth.setUserData(newUser)
                }
                setNickName(false)
            }).catch(reason => {
            props.setError(new Error(reason.msg))
        })
    }



    return <> <Space>{props.nickname}
                <a href={"#"} onClick={()=> setNickName(true)}>修改昵称</a>
            </Space>
            <Form form={nickNameForm} onFinish={onChangeNickname}>
            {editNickName ? <EditorContainer>
                    <Space>
                        <Form.Item  name={"nickname"}  rules={[{ required: true, message: "请输入昵称" }]} >
                            <Input placeholder={"昵称"} type="text" id={"nickname"}     />
                        </Form.Item>
                        <Form.Item>  <Button  loading={isLoading} htmlType={"submit"} type={"primary"}>保存</Button>   </Form.Item>
                        <Form.Item>  <Button  onClick={()=>setNickName(false)}>   取消 </Button>  </Form.Item>
                    </Space>
            </EditorContainer>:null}
            </Form>
    </>
}

