package bitrix

import (
	"encoding/json"
)

func Task_add(task *Task) (err error) {

	jsonData, _ := json.Marshal(task)
	_, err = Exec_api(api_url+"/tasks.task.add.json", jsonData)
	if err != nil {
		return
	}
	return nil
}
