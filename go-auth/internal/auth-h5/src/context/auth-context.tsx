import React, {ReactNode, useCallback} from "react";
import * as authProvider from "@/api/auth/auth-provider";
import { useMount } from "@/utils";
import { useAsync } from "@/utils/use-async";
import { FullPageErrorFallback, FullPageLoading } from "@/components/lib";
import { useQueryClient } from "react-query";
import {
    AccountLoginForm,
    AuthUser,
    SMSLoginForm
} from "@/types/auth-types";
import {message} from "antd";
import {getToken, removeToken, setToken} from "@/api/auth/token";
import request from "@/utils/request";


export const handleUserResponse = (data) =>{
    if (data && data.token) {
        setToken(data.token || "")
    }

    let redirectUrl  = data.redirectUrl
    if (redirectUrl!=="") {
        const sign = data.sign
        const encodeToken = data.encodeToken
        const flag = data.flag
        try {
            redirectUrl += redirectUrl.indexOf("?") > -1 ? "&token=" + encodeToken : "?token=" + encodeToken
            redirectUrl = redirectUrl + "&sign=" + sign + "&flag=" + flag
            console.log("redirectUrl=" + redirectUrl)
            //window.location.href = redirectUrl
        } catch (e) {
            //window.location.reload()
        }
    }else{
        window.location.reload()
    }
}

export const bootstrapUser = async () => {
    let user = null;
    const token = getToken();
    if (token) {
        const res = await request("/auth/me", {
            method: "GET"
        })
        user = res.data
    }
    return user
};

const AuthContext = React.createContext<
  | {
        user: AuthUser | null;

        me: () => void;

        login: (form: AccountLoginForm) => Promise<void>;//帐号登录

        smsLogin: (form: SMSLoginForm) => Promise<void>;//帐号登录

        setUserData: (user: AuthUser | null) => Promise<void>; //返回userId

        logout: () => Promise<void>;  //出

    }
  | undefined
>(undefined);



AuthContext.displayName = "AuthContext";
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const {
    data: user,
    error,
    isLoading,
    isIdle,
    isError,
    run,
    setData: setUser,
  } = useAsync<AuthUser | null>();

    const queryClient = useQueryClient();



    useMount(
        useCallback(() => {
            run(bootstrapUser());
        }, [])
    );

    const me = () => {
        run(bootstrapUser())
    }


    // point free
    const login = async (form:AccountLoginForm) => {
        console.log("帐号登录",JSON.stringify(form))
        authProvider.login(form)
            .then(res=>{
                console.log("成功登录：",JSON.stringify(res))
                setUser(res.data.user)
                handleUserResponse(res.data)
                return Promise.resolve(res)
            }).catch(reason => {
                console.log(JSON.stringify(reason))
                message.error("2."+reason.msg ,1000,()=>{})
                return Promise.reject(reason)
         })
    }

    const smsLogin = async (form:SMSLoginForm) => {
        authProvider.smsLogin(form)
            .then(res=>{
                console.log("成功登录：",JSON.stringify(res))
                setUser(res.data.user)
                handleUserResponse(res.data)
                return Promise.resolve(res.data)
            })
            .catch(reason => {
                console.log(JSON.stringify(reason))
                message.error("2."+reason.msg ,1000,()=>{})
                return Promise.reject(reason)
           })
    }

    const logout = async () => {
        authProvider.logout().then(() => {
        }).catch(reason => {
            message.error(reason.msg ,1000,()=>{})
            return Promise.reject(reason)
        }).finally(()=>{
            setUser(null);
            queryClient.clear();
            removeToken()
            window.location.reload()
        })
    }

    const setUserData = async (user: AuthUser |null) => {
        await setUser(user)
    }

  if (isIdle || isLoading) {
    return <FullPageLoading />;
  }

  if (isError) {
    return <FullPageErrorFallback error={error} />;
  }

  return (
    <AuthContext.Provider
      children={children}
      value={{user, me, login,smsLogin, setUserData, logout }}
    />
  );
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth必须在AuthProvider中使用");
  }
  return context;
};
