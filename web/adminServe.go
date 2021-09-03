package web

import (
	"blog_service/model"
	X "blog_service/utils"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type AdminServe struct {
}

type AdminLogin struct {
	Id        int
	Name      string
	Password  string
	LoginTime string
	Address   string
	CityCode  string
	Token     string
}

func (a *AdminServe) Login(w http.ResponseWriter, r *http.Request) {

	data, err := X.ParsePostData(r.Body)

	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "请求参数错误",
		}))
		return
	}

	str := "select admin_id,admin_password from admin where admin_name =?"
	var u AdminLogin
	err = model.DB.QueryRow(str, data["username"]).Scan(&u.Id, &u.Password)
	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "用户名不存在",
		}))
		return
	}
	if fmt.Sprintf("%x", md5.Sum([]byte(data["password"].(string)))) != u.Password {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "密码错误",
		}))
		return
	}

	ipAddress := strings.Split(r.RemoteAddr, ":")[0]

	res, err := http.Get("https://api.map.baidu.com/location/ip?ak=5yHxHfaWylEVMVlY5cO1npKGeACFT7mn&ip=" + ipAddress + "&coor=bd09ll")
	if err != nil {
		log.Println("Get Address Err,", err.Error())
	}
	if res != nil {
		rBody, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println("Get Address Err,", err.Error())
		}

		var r interface{}
		err = json.Unmarshal(rBody, &r)
		if err != nil {
			log.Println("Err,", err.Error())
		}
		if res, ok := r.(map[string]interface{}); ok {
			if foo, ok := res["content"].(map[string]interface{}); ok {
				if bar, ok := foo["address_detail"].(map[string]interface{}); ok {
					fmt.Println(bar)
					if bar["province"] != nil {
						u.Address = bar["province"].(string)
					}
					if bar["adcode"] != nil {
						u.CityCode = bar["adcode"].(string)
					}
				}
			}
		}
	}

	now := time.Now()
	u.LoginTime = now.Format("2006-01-02 15:04:05")
	u.Token = fmt.Sprintf("%x", md5.Sum([]byte(u.LoginTime)))

	if u.Address == "" {
		str = "update admin set token=?,last_login_time=? where admin_id=?"
		_, err = model.DB.Exec(str, u.Token, u.LoginTime, u.Id)
		if err != nil {
			_, _ = w.Write(X.JSON(X.Z{
				"code":    1,
				"message": "登录失败,注意检查",
			}))
			return
		}
	} else {
		str = "update admin set token=?,last_login_time=?,last_login_address=?,last_login_address_code=? where admin_id=?"
		_, err = model.DB.Exec(str, u.Token, u.LoginTime, u.Address, u.CityCode, u.Id)
		if err != nil {
			_, _ = w.Write(X.JSON(X.Z{
				"code":    1,
				"message": "登录失败,注意检查",
			}))
			return
		}
	}

	_, _ = w.Write(X.JSON(X.Z{
		"code": 0,
		"data": X.Z{
			"token": u.Token,
		},
	}))

}
