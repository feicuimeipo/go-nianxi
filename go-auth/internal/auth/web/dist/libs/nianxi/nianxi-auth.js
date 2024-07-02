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


function showIframe(url,w,h){
    //添加iframe
    var if_w = w;
    var if_h = h;
    //allowTransparency='true' 设置背景透明
    $("<iframe width='" + if_w + "' height='" + if_h + "' id='YuFrame1' name='YuFrame1' style='position:absolute;z-index:4;'  frameborder='no' marginheight='0' marginwidth='0' allowTransparency='true'></iframe>").prependTo('body');
    var st=document.documentElement.scrollTop|| document.body.scrollTop;//滚动条距顶部的距离
    var sl=document.documentElement.scrollLeft|| document.body.scrollLeft;//滚动条距左边的距离
    var ch=document.documentElement.clientHeight;//屏幕的高度
    var cw=document.documentElement.clientWidth;//屏幕的宽度
    var objH=$("#YuFrame1").height();//浮动对象的高度
    var objW=$("#YuFrame1").width();//浮动对象的宽度
    var objT=Number(st)+(Number(ch)-Number(objH))/2;
    var objL=Number(sl)+(Number(cw)-Number(objW))/2;
    $("#YuFrame1").css('left',objL);
    $("#YuFrame1").css('top',objT);

    $("#YuFrame1").attr("src", url)

    //添加背景遮罩
    $("<div id='YuFrame1Bg' style='background-color: Gray;display:block;z-index:3;position:absolute;left:0px;top:0px;filter:Alpha(Opacity=30);/* IE */-moz-opacity:0.4;/* Moz + FF */opacity: 0.4; '/>").prependTo('body');
    var bgWidth = Math.max($("body").width(),cw);
    var bgHeight = Math.max($("body").height(),ch);
    $("#YuFrame1Bg").css({width:bgWidth,height:bgHeight});

    //点击背景遮罩移除iframe和背景
    $("#YuFrame1Bg").click(function() {
        $("#YuFrame1").remove();
        $("#YuFrame1Bg").remove();
    });

   // $("<div id='nianxi-auth-iframe' style='display:block;z-index:0;position:absolute;left:0px;top:0px;filter:Alpha(Opacity=30);/* IE */-moz-opacity:0.4;/* Moz + FF */opacity: 0.4; '/>").prependTo('body');
   // $("#nianxi-auth-iframe").attr("src", url)
}

function removeIframe(url){
    //点击背景遮罩移除iframe和背景
    $("#YuFrame1").remove();
    $("#YuFrame1Bg").remove();
}