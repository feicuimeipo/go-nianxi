import React, {useCallback, useState} from "react";
import {useAuth} from "@/context/auth-context";
import { Descriptions} from "antd";
import {ErrorBox} from "@/components/lib";
import styled from "@emotion/styled";

import {UserInfoMobile} from "@/page/workstation/view/update-mobile";
import {UserInfoEmail} from "@/page/workstation/view/update-email";
import {AccountInfoUserName} from "@/page/workstation/view/user-username";
import {AccountInfoNickName} from "@/page/workstation/view/user-nickname";
import {useMount} from "@/utils";

export const AccountInfoScreen = () =>{
    const {user} = useAuth();
    const [error, setError] = useState<Error | null>(null);


    return <>
        <Container>
              <h3>帐户信息</h3>
              <ErrorBox error={error} />
                <Descriptions  column={2} bordered>
                <Descriptions.Item label="用户名" >
                      <AccountInfoUserName username={user?.username} userId={user?.ID} setError={setError}></AccountInfoUserName>
                </Descriptions.Item>
                <Descriptions.Item label="昵称">
                    <AccountInfoNickName nickname={user?.nickname} userId={user?.ID} setError={setError}></AccountInfoNickName>
                 </Descriptions.Item>
                    <Descriptions.Item label="手机号">
                        <UserInfoMobile mobile={user?.mobile} userId={user?.ID} setError={setError}></UserInfoMobile>
                    </Descriptions.Item>
                    <Descriptions.Item label="邮箱">
                        <UserInfoEmail email={user?.email} userId={user?.ID} setError={setError}></UserInfoEmail>
                    </Descriptions.Item>

            </Descriptions>
        </Container>
    </>
}

const Container = styled.div`
 padding: 10px;
 font-size: 1.5rem;
`


export const EditorContainer = styled.div`
  padding-top: 5px;
`
