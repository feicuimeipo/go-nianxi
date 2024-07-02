package config

import "net/http"

func Cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                      // 可将将 * 替换为指定的域名
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,x-requested-with") //你想放行的header也可以在后面自行添加
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")                                                    //我自己只使用 get post 所以只放行它
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// 放行所有OPTIONS方法
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// 处理请求
		f(w, r)
	}
}
