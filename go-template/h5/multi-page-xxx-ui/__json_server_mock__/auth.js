const localStoreKey = "__auth_provider_token__";
module.exports = {
  authAllow: (path)=>{return authAllow(path)},
  aclAllow: (path,action,uid)=>{return aclAllow(path,action,uid);},
  saveToken: (res,token,remember) => {saveToken(res,token,remember);},
  getToken: (req)=>{
      return req.header("Authorization")
      //return "r487394793";
      //const token = getCookieValue(req,localStoreKey);
      //return token;
  }
};


function getReturnAuthInfo(loginName,recruit){
    const data = {
        loginName: loginName,
        username: "carmen",
        name: "卡门",
        token: "r487394793",
        openId: "",
        lastTimeLoginDate: "111",
        isRecruit: recruit
    };
    return data;
}
const tokenList = () => {
    return [
        {
        token: "111111",
        expire: 119191991
        },
        {
            token: "222222",
            expire: 119191991
        },
        {
            token: "r487394793",
            expire: -1
        },
    ]
}

const authAllow = (path) =>{
   return true;
}

//Access Control List
const aclAllow = (path,action,uid) =>{
    return true;
}

const saveToken = (res,token,remember) => {
    if (remember) {
       /*
        const expires = BigInt((new Date()).valueOf() + 24 * 60 * 60 * 1000 * 7);
        res.cookie(localStoreKey,token,{maxAge:expires,httpOnly:true,path:"/"})
        */
        const cookies = "'"+localStoreKey+"="+token+"'"
        //res.setHeader("Set-Cookie",[cookies]);
    }else {
        const cookies = "'"+localStoreKey+"="+token+"'"
        //res.setHeader("Set-Cookie",[cookies]);// .cookie(localStoreKey,token,{httpOnly:true,path:"/"});
        //res.cookie(localStoreKey,token,{httpOnly:true,path:"/"})
    }
    //req.setHeader()

}


const getCookieValue = (req,name) =>{
    let Cookies ={};
    if (req.headers.cookie != null) {
        req.headers.cookie.split(';').forEach(l => {
            const parts = l.split['='];
            const key = parts[0].trim();
            const value = (parts[1] || '').trim();;
            if (key === name){
                return value;
            }
        });
    }
    return null;
}


const fail = (code,msg) =>{
    return {
        code: code,
        msg: msg,
        data: null
    };
}

const success = (val) =>{
    return {
        code: 200,
        data: val?val:{},
        msg: "",
    };
}
