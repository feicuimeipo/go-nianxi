let AuthApiBaseUrl  =  "http://localhost:8000"
let LocalApiBaseUrl = "http://localhost:8000"
let LogoutReturnUrl = "/login.html"
const HeaderTokenName = "Authorization"


function getLocationUrl(){
    const currentWwwPath=window.document.location.href;
    const pathName=window.document.location.pathname;
    const position=currentWwwPath.indexOf(pathName);
    const localhostPath=currentWwwPath.substring(0,position);
    return localhostPath;
}




/**
 * 得到查询
 * @param name
 * @returns {string|number}
 */
function getQuery(name){
    if (!window.location.search || window.location.search==undefined || window.location.search==""){
        return ""
    }
    let paramsString = window.location.search;
    let searchParams = new URLSearchParams(paramsString);
    console.log(searchParams.get(name))
    return searchParams.get(name)
}


function AddTokenToHeader(token){
    var head = new Headers();
    head.append(HeaderTokenName,token)
}




//请求图片get事件
function AjaxPost(data,path,resolve,reject){
    if (!path.startsWith("/") ){
        path = "/" + path
    }
    let apiUrl = ""
    if (path.startsWith("https://") || path.startsWith("https://")) {
        apiUrl = path;
    }else{
        apiUrl = AuthApiBaseUrl + path
    }

    $.ajax({
        type : "post",
        contentType: "application/json;charset=UTF-8",
        url : apiUrl,
        data :JSON.stringify(data),
        cache: false,
        crossDomain: true === !(document.all),
        async: false,
        success:function(res,status){
            resolve(res)
            // try{
            //     const r = eval(res)
            //     resolve(r)
            // }catch (e){
            //     console.log("数据成功返回，但是解释返回值异常,返回值:\n"+res+",\n 错误信息：\n"+e.message)
            // }
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            console.log("error-XMLHttpRequest.status="+XMLHttpRequest.status)
            if (XMLHttpRequest.status == 401 || XMLHttpRequest.status==402){
                //TOOD 授权不正确
            }
        },
        fail: function(err) {
            console.log("fail=\n"+JSON.stringify(err))
            reject(err)
        }
    })
}

//请求图片get事件
function AjaxGet(data,path,resolve,reject){
    if (!path.startsWith("/") ){
        path = "/" + path
    }
    let apiUrl = ""
    if (path.startsWith("https://") || path.startsWith("https://")) {
        apiUrl = path;
    }else{
        apiUrl = AuthApiBaseUrl + path
    }


    $.ajax({
        type : "get",
        contentType: "application/json;charset=UTF-8",
        url : apiUrl,
        data :JSON.stringify(data),
        cache: false,
        crossDomain: true === !(document.all),
        async: false,
        success:function(res,status){
            resolve(res)
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            console.log("error-XMLHttpRequest.status="+XMLHttpRequest.status)
            if (XMLHttpRequest.status == 401 || XMLHttpRequest.status==402){
                //TOOD 授权不正确
            }
        },
        fail: function(err) {
            console.log("fail=\n"+JSON.stringify(err))
            reject(err)
        }
    })
}


function initAuthPanel(parentId,authApiBaseUrl,localApiBaseUrl,logoutReturnUrl){
    AuthApiBaseUrl = authApiBaseUrl
    LocalApiBaseUrl = localApiBaseUrl
    LogoutReturnUrl = logoutReturnUrl
    if (LocalApiBaseUrl.endsWith("/")){
        LocalApiBaseUrl = LocalApiBaseUrl.substring(0,LocalApiBaseUrl.length-1)
    }
    if (AuthApiBaseUrl.endsWith("/")){
        AuthApiBaseUrl = AuthApiBaseUrl.substring(0,AuthApiBaseUrl.length-1)
    }

    const html = `
            <!-- 头 像 -->
            <a class="layui-icon layui-icon-username" href="javascript:;" >&nbsp;nianxi</a>
            <!-- 功 能 菜 单 -->
            <dl class="layui-nav-child">
                <dd><a href="javascript:void(0);" onclick='onAuthLogout()'>退出</a></dd>
            </dl>
                `
    $("#"+parentId).append(html)
}


function onAuthLogout(){
    const logoutReturnUrl =  LocalApiBaseUrl + LogoutReturnUrl

    showIframe(AuthApiBaseUrl + "/auth/logout",0,0)
    AjaxGet({},  LocalApiBaseUrl + "/auth/client/logout",function(res){
        location.href = logoutReturnUrl;
    },function (err){

    })

}


function showIframe(url){
    $("<div id='nianxi-auth-iframe' style='display:block;z-index:0;position:absolute;left:0px;top:0px;filter:Alpha(Opacity=30);/* IE */-moz-opacity:0.4;/* Moz + FF */opacity: 0.4; '/>").prependTo('body');
    $("#nianxi-auth-iframe").attr("src", url)
}
