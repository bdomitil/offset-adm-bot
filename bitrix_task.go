package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func add_task_to_deal(deal_id int, data *reportForm) error {

	var task Task
	task.Body.Description = data.comments
	task.Body.Title = "Задача созданная телеграм ботом"
	task.Body.Deal_id = fmt.Sprintf("D_%d", deal_id)
	task.Body.Responible_id = "14"
	jsonData, _ := json.Marshal(task)
	response, err := exec_api(api_url+"/tasks.task.add.json", jsonData)
	if err != nil {
		log.Printf("%s\n", string(response))
	}
	return nil
}
