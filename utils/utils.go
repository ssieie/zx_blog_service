package utils

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func ParsePostData(data io.ReadCloser) (rData map[string]interface{}, err error) {

	res, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	var r interface{}

	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, err
	}

	rData, ok := r.(map[string]interface{})
	if !ok {
		return nil, errors.New("parse data err 断言失败")
	}

	return rData, nil
}

type Z map[string]interface{}

func JSON(res Z) (data []byte) {
	data, _ = json.Marshal(res)
	return
}

func VerifyPostParams(data map[string]interface{}) bool {

	var str = "key=ZhaoXin&"

	var arr []string

	for key := range data {
		arr = append(arr, key)
	}

	sort.Strings(arr)

	for _, key := range arr {
		if key == "sign" {
			continue
		}
		switch data[key].(type) {
		case string:
			str += key + "=" + data[key].(string) + "&"
		case int:
			str += key + "=" + strconv.Itoa(data[key].(int)) + "&"
		case bool:
			str += key + "=" + strconv.FormatBool(data[key].(bool)) + "&"
		case float64:
			str += key + "=" + strconv.FormatFloat(data[key].(float64), 'g', 10, 64) + "&"
		}
	}

	sign := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(strings.TrimSuffix(str, "&")))))

	if sign != data["sign"] {
		return false
	}
	return true
}
