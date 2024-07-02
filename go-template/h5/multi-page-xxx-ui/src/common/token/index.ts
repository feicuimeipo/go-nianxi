import Cookies from "js-cookie"

const TokenKey:string = 'Authorization'

export function getToken():string|null {
  const token  = Cookies.get(TokenKey)
  if (token == undefined){
      return  null;
  }
  return token;
}

export function setToken(token:string) {
  return Cookies.set(TokenKey, token)
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}
