import { useState, useEffect, useRef } from 'react'

/**
 * 定时的次数：比如60秒
 * 定时器想做的事情
 * @param Num
 * @param callBack
 * @return
 *  返回的第一参数为可变值
 *  返回第二参数为触发函数
 */
export function useTimer (Num, callBack = () => {}) {
    const [num, setNum] = useState(Num)
    const ref = useRef(null)
    // 调用这个方法,开始倒数
    const start = () => {
        setNum(Num) // 重新赋值
        // 定时器
        // @ts-ignore
        ref.current = setInterval(() => {
            setNum((num) => num - 1)
        }, 1000)
    }

    // 倒数为0 关闭定时器
    useEffect(
        () => {
            if (num === 0) {
                return () => {
                    // @ts-ignore
                    clearInterval(ref.current) // 关闭定时器
                    callBack() // 自己想做的事
                }
            }
        },
        [num]
    )

    // 解决当正在计数的组件开始倒数 删除组件导致无法取消定时器
    useEffect(() => {
        return () => {
            // @ts-ignore
            clearInterval(ref.current)
        }
    }, [])
    // 返回可变值 与 触发函数
    return {
        num, // 可变值
        start // 触发函数
    }
}
