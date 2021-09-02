package middleware

import (
	"blog_service/utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Cross(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("content-type", "application/json")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		var bodyBytes []byte

		if request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(request.Body)
		}

		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var r interface{}

		err := json.Unmarshal(bodyBytes, &r)
		if err != nil {
			_, _ = writer.Write(utils.JSON(utils.Z{
				"code":    1,
				"message": "请求参数错误",
			}))
			return
		}

		params, ok := r.(map[string]interface{})

		if !ok {
			_, _ = writer.Write(utils.JSON(utils.Z{
				"code":    1,
				"message": "请求参数错误",
			}))
			return
		}

		ok = utils.VerifyPostParams(params)
		if !ok {
			_, _ = writer.Write(utils.JSON(utils.Z{
				"code":    1,
				"message": "验证不通过",
			}))
			return
		}

		next.ServeHTTP(writer, request)
	})
}
