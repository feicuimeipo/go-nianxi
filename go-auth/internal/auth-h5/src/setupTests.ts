// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import "@testing-library/jest-dom";
import { server } from './mocks/server'

//我们将在所以测试用例开始前将服务端msw启动监听
beforeAll(() => {
    // 开启api监听
    server.listen()
})

//每个测试之后重置处理程序
afterEach(() => {
    // 重置处理程序
    server.resetHandlers()
})

//测试结束后结束msw
afterAll(() => {
    // 关闭监听
    server.close()
})
