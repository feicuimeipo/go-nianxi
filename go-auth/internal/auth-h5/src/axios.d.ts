// axios.d.ts
import {AxiosRequestConfig} from "axios";

declare module 'axios' {
    interface AxiosInstance {
        (config: AxiosRequestConfig): Promise<any>
    }
}
