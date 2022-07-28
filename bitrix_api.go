package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var api_url string = os.Getenv("bitrix_api_url")

func exec_api(url string, jsonData []byte) ([]byte, error) {

	if api_url == "" {
		log.Fatalln("no bitrix api url token exported")
	}

	fmt.Println(string(jsonData))

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

func api_test() {
	var testReport reportForm
	testReport.comments = "здесь тестовое описание задачи"
	deals := get_deals()
	id, _ := get_deal_id_by_name(deals, "TEST")
	add_task_to_deal(id, &testReport)
}
