
import React, { useState, useEffect, useRef,  useImperativeHandle } from 'react'
import { Button } from 'antd'


type selfProps = {
    codeStyle: object;
};

export const CountDown = (props, ref) => {
    const intervalRef = useRef<any>(null)
    const [count, changeCount] = useState(0)
    const { codeStyle } = props // 样式

    // 组件卸载时清除计时器
    useEffect(() => {
        return () => {
            clearInterval(intervalRef.current)
        }
    }, [])

    useEffect(() => {
        if (count === 59) {
            intervalRef.current = setInterval(() => {
                changeCount((preCount) => preCount - 1)
            }, 1000)
        } else if (count === 0) {
            clearInterval(intervalRef.current)
        }
    }, [count])

    // 暴露的子组件方法，给父组件调用
    useImperativeHandle(ref, () => {
        return {
            _childFn() {
                changeCount(59)
            }
        }
    })

    return (
        <Button disabled={!!count} style={ codeStyle }>
            {count ? `${count} s后重新发送` : '发送验证码'}
        </Button>
    )
}


