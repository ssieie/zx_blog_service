package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	data["timestamp"] = int(data["timestamp"].(float64))

	var sign string

	for key, value := range data {
		if key == "sign" {
			sign = data[key].(string)
			delete(data, key)
		}
		fmt.Println(key, value)
	}

	fmt.Println(data, sign)
	return true
}
