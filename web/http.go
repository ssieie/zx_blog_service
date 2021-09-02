package web

import (
	"blog_service/middleware"
	"fmt"
	"net/http"
)

func InitHttp() (err error) {
	var admin AdminServe

	// admin
	http.Handle("/adminLogin", middleware.Cross(middleware.Verify(http.HandlerFunc(admin.Login))))

	fmt.Println("go server start running...")
	err = http.ListenAndServe(":9999", nil)

	if err != nil {
		return err
	}
	return nil
}
