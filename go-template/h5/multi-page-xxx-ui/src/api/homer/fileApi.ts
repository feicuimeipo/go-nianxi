import request from "@/common/utils/request";

export enum fileType {
    Design=1,
    Video=2,
    PersonalizedManufacturing,
    Witting,
    Enterprise_Print,
    Web,
}

export interface createFileDTO{
    userId: number
    fileType: fileType
    isPersonal: boolean
}

export interface fileInfo{
    userId: number
    fileType: fileType
    isPersonal: boolean
}

export const FileApi = {
    //返回文件id
    createFile: (data: createFileDTO) =>{
        return request({
            url: '/file/add',
            method: 'post',
            data
        })
    },
    getFile: (fileId:number) =>{
        return request({
            url: '/file/'+fileId,
            method: 'get',
        })
    },
}


