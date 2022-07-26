package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// PostJSON post json 数据请求
func PostJSON(uri string, obj []byte) ([]byte, error) {
	var body = strings.NewReader(string(obj))
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}
