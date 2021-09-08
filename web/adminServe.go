package web

import (
	"blog_service/config"
	"blog_service/model"
	X "blog_service/utils"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

// 管理员登录

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

// 首页信息

func (a *AdminServe) HomeInfo(w http.ResponseWriter, r *http.Request) {
	data, err := X.ParsePostData(r.Body)
	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "请求参数错误",
		}))
		return
	}

	var (
		loginTime   string
		address     string
		addressCode string
	)

	str := "select last_login_time,last_login_address,last_login_address_code from admin where token=?"
	err = model.DB.QueryRow(str, data["token"]).Scan(&loginTime, &address, &addressCode)
	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    9,
			"message": "登录过期",
		}))
		return
	}

	_, _ = w.Write(X.JSON(X.Z{
		"code": 0,
		"data": map[string]string{
			"last_time":         loginTime,
			"last_address":      address,
			"last_address_code": addressCode,
		},
	}))
}

// 获取当地天气

func (a *AdminServe) getWeather(w http.ResponseWriter, r *http.Request) {
	data, err := X.ParsePostData(r.Body)
	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "参数错误",
		}))
		return
	}

	url := fmt.Sprintf("https://api.map.baidu.com/weather/v1/?district_id=%s&data_type=all&ak=5yHxHfaWylEVMVlY5cO1npKGeACFT7mn", data["code"].(string))

	get, err := http.Get(url)
	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "天气获取错误",
		}))
		return
	}
	defer get.Body.Close()

	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "天气获取错误",
		}))
		return
	}

	_, _ = w.Write(all)
}

// 富文本图片上传处理

func (a *AdminServe) ImageLoad(w http.ResponseWriter, r *http.Request) {

	var param = make(map[string]interface{})
	param["token"] = r.FormValue("token")
	param["timestamp"] = r.FormValue("timestamp")
	param["sign"] = r.FormValue("sign")
	isOk := X.VerifyPostParams(param)

	if !isOk {
		_, _ = w.Write(X.JSON(X.Z{
			"code":    1,
			"message": "非法调用",
		}))
		return
	}

	_ = r.ParseMultipartForm(32 << 20)

	if r.MultipartForm.File != nil {
		// formName 文件的名字
		var files []map[string]string

		for formName := range r.MultipartForm.File {
			file, header, err := r.FormFile(formName)
			if err != nil {
				log.Println(err)
				break
			}
			defer func(file multipart.File) {
				_ = file.Close()
			}(file)
			saveName := fmt.Sprintf("%d%s", time.Now().UnixNano(), header.Filename)

			var filePath string
			if config.Dev == "local" {
				filePath = config.LocalFilePath
			} else {
				filePath = config.RemoteFilePath
			}
			destFile, err := os.Create(filePath + saveName)

			if err != nil {
				log.Println(err)
				break
			}
			defer func(destFile *os.File) {
				_ = destFile.Close()
			}(destFile)

			_, err = io.Copy(destFile, file)
			if err != nil {
				log.Println(err)
				break
			}

			fileDetail := make(map[string]string)
			if config.Dev == "local" {
				fileDetail["url"] = config.LocalAddress + saveName
			} else {
				fileDetail["url"] = config.RemoteAddress + saveName
			}
			fileDetail["alt"] = saveName
			fileDetail["href"] = ""
			files = append(files, fileDetail)
		}
		_, _ = w.Write(X.JSON(X.Z{
			"errno": 0,
			"data":  files,
		}))
	}
}

// 发布文章

type Article struct {
	ArTitle      string
	ArCategoryId int
	ArTag        string
	ArAuthor     string
	ArCreTime    string
	ArImages     string
}

func (a *AdminServe) ReleaseArticle() {

}
