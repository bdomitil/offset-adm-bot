package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

func get_deal_id_by_name(deals Deal, deal_name string) (int, error) {
	for _, deal := range deals.Body {
		if strings.Contains(deal.Title, deal_name) {
			return strconv.Atoi(deal.Id)
		}
	}
	return 0, errors.New("No deal found by " + deal_name)
}

func get_deals() (deals Deal) {
	values := map[string][]string{"select": {"ID", "TITLE"}}
	jsonData, _ := json.Marshal(values)
	response, err := exec_api(api_url+"/crm.deal.list/", jsonData) //TODO: handle error
	if err != nil {
		log.Println(string(response))
	}
	err = json.Unmarshal(response, &deals)
	if err != nil {
		log.Println(err.Error())
	}
	return
}
