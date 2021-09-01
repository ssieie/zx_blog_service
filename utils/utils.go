package utils

import (
	"encoding/json"
	"errors"
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
