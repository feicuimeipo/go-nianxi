import {ReducerStateWithoutAction, useCallback, useReducer, useState} from "react";
import { useMountedRef } from "@/utils/index";

interface State<D> {
    error: Error | null;
    data: D | null;
    stat: "idle" | "loading" | "error" | "success";
}

const defaultInitialState = {
    stat: "idle",
    data: null,
    error: null,
} as State<null>

const defaultConfig = {
    throwOnError: false,
};

const useSafeDispatch = <D>(dispatch: (...args: State<D>[]) => void) => {
    const mountedRef = useMountedRef();
    return useCallback(
        (...args: State<D>[]) => (mountedRef.current ? dispatch(...args) : void 0),
        [dispatch, mountedRef]
    );
};


export const useAsync = <D>(
      initialState?: State<D>,
      initialConfig?: typeof defaultConfig
) => {

    const config = { ...defaultConfig, ...initialConfig };

    const initState = {
        ...defaultInitialState,
        ...initialState,
    } as never

    const [state, dispatch] = useReducer (
        (state:State<D>, action:Partial<State<D>>) => ({...state,...action})
        ,initState);


    const safeDispatch = useSafeDispatch(dispatch);
    // useState直接传入函数的含义是：惰性初始化；所以，要用useState保存函数，不能直接传入函数
    const [retry, setRetry] = useState(() => () => {});

    const setData = useCallback(
        (data: D) =>
            safeDispatch({
                data,
                stat: "success",
                error: null,
            } ),
        [safeDispatch]
    );

    const setError = useCallback(
        (error: Error) =>
            safeDispatch({
                error,
                stat: "error",
                data: null,
            } ),
        [safeDispatch]
    );

      // run 用来触发异步请求
      const run = useCallback(
        (fn: Promise<D>, runConfig?: { retry: () => Promise<D> }) => {
          if (!fn || !fn.then) {
                throw new Error("请传入 Promise 类型数据");
          }
          setRetry(() => () => {
            if (runConfig?.retry) {
                run(runConfig.retry(), runConfig);
            }
          });
          safeDispatch({data:null,error:null, stat: "loading" });
          return fn
            .then((data) => {
                console.log("useAsync=",JSON.stringify(data))
                setData(data);
                return data;
            })
            .catch((error) => {
                    //console.log("error=",JSON.stringify(error))
                  setError(error);
                  if (config.throwOnError) return Promise.reject(error);
                  return error;
            });
        },
        [config.throwOnError, setData, setError, safeDispatch]
      );



      return {
            isIdle: state.stat === "idle",
            isLoading: state.stat === "loading",
            isError: state.stat === "error",
            isSuccess: state.stat === "success",
            run,
            setData,
            setError,
            // retry 被调用时重新跑一遍run，让state刷新一遍
            retry,
            ...state,
      };
};
