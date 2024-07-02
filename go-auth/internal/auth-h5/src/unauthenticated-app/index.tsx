import React, {useEffect, useState} from "react";
import {Button, Card, Divider, Input} from "antd";
import styled from "@emotion/styled";
import {useDocumentTitle} from "@/utils";
import left from "@/assets/left.svg";
import right from "@/assets/right.svg";
import logo from "@/assets/logo.png";
import {Outlet, useLocation, } from "react-router";
import {Link } from "react-router-dom";

export interface TabProps{
    key: string
    mobile: string
    userId: number
    multiRecord?: boolean
}


export interface PanelProps{
    title: string,
    routes: {
        path: string,
        label: string
    }[],

}

const UnAuthenticatedApp = () => {
    const [panelProps, setPanelProps] = useState<PanelProps>({
        title: "请登录",
        routes: [
            {path: "register",label: "没有账号？注册新账号"},
            {path: "forgetPassword",label: "忘记密码"}
        ]
    });

    useDocumentTitle("请登录注册以继续");

    const location = useLocation ()
    useEffect(()=>{
        if (location.pathname.startsWith("/auth/register")){
            setPanelProps({
                title: "请注册",
                routes: [
                    {path: "login",label: "已经有账号了？直接登录"},
                ],

            })
        }else if (location.pathname.startsWith("/auth/forgetPassword")){
            setPanelProps({
                title: "找回密码",
                routes: [
                    {path: "login",label: "返回登录页"}
                ]
            })
        }else if (location.pathname.startsWith("/auth/login")){
            setPanelProps({
                title: "请登录",
                routes: [
                    {path: "register",label: "没有账号？注册新账号"},
                    {path: "forgetPassword",label: "忘记密码"}
                ]
            })
        }
    },[location])

    return (<Container>
            <Header />
            <Background />
            <ShadowCard>
                <Title>{panelProps.title}</Title>
                <Outlet />
                <Divider />
                {
                    panelProps.routes.map((v,k) => (
                        <p key={k}><Link  to={v.path} key={v.path} > {v.label} </Link> </p>
                    ))
                }
            </ShadowCard>

        </Container>
    );
}

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  width: 100vw;
`;


const Header = styled.header`
  background: url(${logo}) no-repeat center;
  padding: 5rem 0;
  background-size: 8rem;
  width: 100%;
`;


const Title = styled.h2`
  margin-bottom: 1.2rem;
  color: rgb(94, 108, 132);
`;

const Background = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  background-repeat: no-repeat;
  background-attachment: fixed;
  background-position: left bottom, right bottom;
  background-size: calc(((100vw - 40rem) / 2) - 3.2rem),
    calc(((100vw - 40rem) / 2) - 3.2rem), cover;
  background-image: url(${left}), url(${right});
`;


const ShadowCard = styled(Card)`
  width: 22%;
  min-height: 25rem;
  padding: 1.5rem 1.5rem;
  border-radius: 0.3rem;
  box-sizing: border-box;
  box-shadow: rgba(0, 0, 0, 0.1) 0 0 10px;
  text-align: center;
`;


export const LongButton = styled(Button)`
  width: 100%;
`;

export const VerifyInput = styled(Input)`
   width: 20em;
`


export default UnAuthenticatedApp
