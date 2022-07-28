package bitrix

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var api_url string
var executor Profile

func Exec_api(url string, jsonData []byte) ([]byte, error) {

	if api_url == "" {
		log.Fatalln("no bitrix api url token exported")
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil

}

func Init(apiUrl string) (profile Profile, x error) {
	api_url = strings.Clone(apiUrl)
	response, err := Exec_api(apiUrl+"/profile", nil)
	if err == nil {
		err = json.Unmarshal(response, &executor)
	}
	return executor, err
}
