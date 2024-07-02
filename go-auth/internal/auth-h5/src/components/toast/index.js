import React from 'react'
import Toast from './toast'
import './toast.css'
import {createRoot} from "react-dom/client";

function createNotification() {
    const container  = document.createElement("div")
    document.getElementById("root").appendChild(container)

    const root = createRoot(container)
    const notification = root.render(<Toast />)
    return {
        addNotice(notice) {
            notification.addNotice(notice)
        },
        destroy() {
            root.removeChild(container)
            document.getElementById("root").removeChild(container)
            root.unmount()
        }
    }
}

let notification
const notice = (type, content, duration = 2000, onClose) => {
    if (!notification) notification = createNotification()
    return notification.addNotice({ type, content, duration, onClose })
}

export default {
    info(content, duration, onClose) {
        return notice('info', content, duration, onClose)
    },
    success(content = '操作成功', duration, onClose) {
        return notice('success', content, duration, onClose)
    },
    error(content, duration , onClose) {

        return notice('error', content, duration, onClose)
    },
    loading(content = '加载中...', duration = 0, onClose) {
        return notice('loading', content, duration, onClose)
    }
}
