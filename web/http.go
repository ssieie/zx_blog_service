package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MyBlog struct {
}

func (m MyBlog) test(w http.ResponseWriter, r *http.Request) {
	Cross(w)
	fmt.Println(r.Method)

	test := make(map[string]interface{}, 10)
	test["name"] = "Zx"
	test["age"] = 18

	data, err := json.Marshal(test)
	processErr(err)
	_, _ = w.Write(data)
}

func InitHttp() (err error) {
	var m MyBlog
	http.HandleFunc("/test", m.test)
	err = http.ListenAndServe(":9999", nil)

	if err != nil {
		return err
	}
	return nil
}

func Cross(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
}

func processErr(e error) {
	if e != nil {
		log.Fatalf("%s \n", e.Error())
	}
}
