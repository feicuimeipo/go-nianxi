import Cookies from "@/utils/cookies";
import {nanoid} from "nanoid";

const CanoIdKey = '__auth_provider_canoId__'
const TokenKey = '__auth_provider_token__'
export function setToken(token:string) {
    //return window.navigator.cookieEnabled?  Cookies.set(TokenKey, token): localStorage.setItem(TokenKey,token)
    return localStorage.setItem(TokenKey,token)
}

export function removeToken() {
    //return window.navigator.cookieEnabled ?  Cookies.remove(TokenKey) : localStorage.removeItem(TokenKey)
    return localStorage.removeItem(TokenKey)
}

export function getToken():string|null {
    // let token = window.navigator.cookieEnabled? Cookies.get(TokenKey) :  localStorage.get(TokenKey)
    // return token === undefined? null : token
    return localStorage.getItem(TokenKey)
}

export function getCanoIdKey():string {
    //let canoId = window.navigator.cookieEnabled? Cookies.get(CanoIdKey) :  localStorage.get(CanoIdKey)
    let canoId =localStorage.getItem(CanoIdKey)
    if (!canoId || canoId == undefined || canoId === "" || canoId == null){
        canoId = nanoid();
        Cookies.set(canoId, canoId)
    }
    return canoId
}
