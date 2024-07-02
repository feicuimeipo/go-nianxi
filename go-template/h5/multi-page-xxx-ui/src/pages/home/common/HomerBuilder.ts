import {NavMenu} from "@/common/model/navMenu";
import {DownMenuTree} from "@/common/model/downMenuTree";
import {fileType} from "@/api/homer/fileApi";

const HomerBuilder = {
    Init: {
        navMenuList : <NavMenu[]> [
            {id:1,name: "最近",router: 'recently', icon:'la la-crosshairs',rightIcon:null},
            {id:2,name: "我的文件",router: 'recently', icon:'la la-folder-open',rightIcon:null},
            {id:3,name: "正在跟进",router: 'recently', icon:'la la-hand-point-right',rightIcon:null},
            {id:4,name: "资源社区",router: 'recently', icon:'la la-globe',rightIcon:'la la-arrow-right'}
        ],
        addFileMenuList: <DownMenuTree[]>[
            {id: 1,name:"创建文件夹",parentId:0,divider:false,icon: 'la la-folder-plus',  code: "-"                  ,child:null},
            {id: -1,name:"创建文件"  ,parentId:0,divider:true,icon: '',                   code: '-'                  ,child:null},
            {id: 1,name: "在线设计",parentId:0,divider:false,icon: 'la la-pencil-ruler',  code: fileType.Design     ,child:null},
            {id: 2,name: "视频合成",parentId:0,divider:false,icon: 'la la-film',          code: fileType.Video      ,child:null},
            {id: 3,name: "个性制造",parentId:0,divider:false,icon: 'la la-crop',          code: fileType.PersonalizedManufacturing  ,child:null},
            {id: 4,name: "企业印刷",parentId:0,divider:false,icon: 'la la-stamp',         code: fileType.Enterprise_Print           ,child: null},
            {id: 5,name: "#"      ,parentId:0,divider:true,icon: '',                     code: "-"                     ,child:null},
            {id: 6,name: "网页制作",parentId:0,divider:false,icon: 'la la-globe',         code: fileType.Web            ,child:null},
            {id: 7,name: "#",parentId:0,divider:true,icon: '',                           code: "-"                     ,child:null},
            {id: 8,name: "智能写作",parentId:0,divider:false,icon: 'la la-edit',          code: fileType.Witting        ,child:null},
        ]
    }
}

export default HomerBuilder
