package bitrix

import (
	"encoding/json"
	"log"
)

func Task_add(task *Task) error {

	jsonData, _ := json.Marshal(task)
	response, err := Exec_api(api_url+"/tasks.task.add.json", jsonData)
	if err != nil {
		log.Printf("%s\n", string(response))
	}
	return nil
}
