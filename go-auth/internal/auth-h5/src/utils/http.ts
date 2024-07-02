import qs from "qs";
import { useCallback } from "react";
import {useAuth} from "@/context/auth-context";
import {logout} from "@/api/auth/auth-provider";
import {getToken} from "@/api/auth/token";

let apiUrl = process.env.REACT_APP_API_URL;

interface Config extends RequestInit {
  token?: string;
  data?: object;
}

export const http = async (
    endpoint: string,
    { data, headers, ...customConfig }: Config = {}
) => {
    const config = {
       method: "GET",
       headers: {
            Authorization: getToken()? `Bearer ${getToken()}`:"",
           'X-Requested-With':'XMLHttpRequest',
            "Content-Type": data ? "application/json" : ""
       },
       ...customConfig,
    } as Config;

    if (!config.method){
        config.method = "GET"
    }
    if (config.method.toUpperCase() === "GET") {
         endpoint += `?${qs.stringify(data)}`;
    } else {
         config.body = JSON.stringify(data || {});
    }

    if (endpoint && endpoint.startsWith("/")){
        endpoint = endpoint.substring(1);
    }
    if (apiUrl && apiUrl.endsWith("/")){
        apiUrl = apiUrl.substring(0,apiUrl.length-1)
    }


  // axios 和 fetch 的表现不一样，axios可以直接在返回状态不为2xx的时候抛出异常
  return window
      .fetch(`${apiUrl}/${endpoint}`, config)
      .then(async (res) => {
          if ( res.ok){
              return await res.json();
          }else {
              if (res.status === 401 || res.status === 401 || res.status === 403) {
                  if (res.statusText.indexOf('JWT认证失败') !== -1) {
                      await logout();
                      window.location.reload();
                      return Promise.reject({message: "请重新登录"});
                  } else {
                      return Promise.reject({message: "认证失败"});
                  }
              } else {
                    const data = await res.json();
                    return Promise.reject(data);
              }
          }
      })
};


export const useHttp = () => {
  const { user } = useAuth();
  return useCallback(
      (...[endpoint, config]: Parameters<typeof http>) =>
          http(endpoint, { ...config, token: user?.token }),
      [user?.token]
  );
};
