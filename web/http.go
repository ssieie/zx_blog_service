package web

import (
	"fmt"
	"net/http"
)

func middlewareCross(next http.Handler) http.Handler {
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

func InitHttp() (err error) {
	var admin AdminServe

	// admin
	http.Handle("/adminLogin", middlewareCross(http.HandlerFunc(admin.Login)))

	fmt.Println("go server start running...")
	err = http.ListenAndServe(":9999", nil)

	if err != nil {
		return err
	}
	return nil
}
