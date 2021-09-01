package web

import (
	"blog_service/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AdminServe struct {
}

func (a *AdminServe) Login(w http.ResponseWriter, r *http.Request) {

	data, err := utils.ParsePostData(r.Body)
	processErr(err)
	fmt.Println(data["username"])

	res, err := json.Marshal(data)
	processErr(err)
	_, _ = w.Write(res)
}

func processErr(e error) {
	if e != nil {
		log.Printf("%s \n", e.Error())
	}
}
func processErrExit(e error) {
	if e != nil {
		log.Fatalf("%s \n", e.Error())
	}
}
