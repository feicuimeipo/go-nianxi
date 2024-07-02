import {fileType} from "@/api/home/fileApi";


const ModuleRouter = {
    goClub: () => {
        window.location.href = "../club/index.html"
    },

    goMember: () => {
        window.location.href = "../member/index.html"
    },

    goEditor: (fileId: string,fileType:fileType) => {
        if (fileId == null || fileId === "") {
            window.location.href = "../editor/index.html?fileId="+ fileId +"&fileType="+fileType;
        }
    },

    goHome: () => {
        window.location.href = "../home/index.html";
    },
}

export default ModuleRouter
