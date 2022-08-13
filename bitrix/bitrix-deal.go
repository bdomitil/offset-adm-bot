package bitrix

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

func Get_deal_id_by_name(d Deal, deal_name string) (int, error) {
	for _, deal := range d.Result {
		if strings.Contains(deal.Title, deal_name) {
			return strconv.Atoi(deal.Id)
		}
	}
	return 0, errors.New("No deal found by " + deal_name)
}

func Get_deals() (d Deal) {
	values := map[string][]string{"select": {"ID", "TITLE"}}
	jsonData, _ := json.Marshal(values)
	response, err := Exec_api(api_url+"/crm.deal.list/", jsonData) //TODO: handle error
	if err != nil {
		log.Println(string(response))
	}
	err = json.Unmarshal(response, &d)
	if err != nil {
		log.Println(err.Error())
	}
	return
}
