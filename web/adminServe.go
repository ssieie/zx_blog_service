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
	fmt.Println(data)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	res, err := json.Marshal(data)
	processErr(err)
	_, _ = w.Write(res)

	//fmt.Printf("%x", md5.Sum([]byte("haha")))
}

func processErr(e error) {
	if e != nil {
		log.Printf("%s \n", e.Error())
	}
}
