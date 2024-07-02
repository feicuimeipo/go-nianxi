

# 首先引入
```angular2html
import Toast from './components/toast'
```

# JSX中事件调用：
```
<button onClick={() => { Toast.info('普通提示') }}>普通提示</button>
```

# JS中方法调用：
```angular2html
Toast.info('普通提示')
```

# 回调方法：
```angular2html
const hideLoading = Toast.loading('加载中...', 0, () => {
  Toast.success('加载完成')
})
setTimeout(hideLoading, 2000)
```

# 调用规则：

3个参数
- content 提示内容 string（loading方法为可选）
- duration 提示持续时间 number，单位ms（可选）
- onClose 提示关闭时的回调函数（可选）

```angular2html
Toast.info("普通",2000)
Toast.success("成功",1000,() => {
  console.log('回调方法')
}))
Toast.error("错误")
Toast.loading()
```
