package main

type Deal_body struct {
	Id    string `json:"ID"`
	Title string `json:"TITLE"`
}

type Deal struct {
	Body []Deal_body `json:"result"`
}

type Task_body struct {
	Description   string `json:"DESCRIPTION"`
	Title         string `json:"TITLE"`
	Id            string `json:"ID"`
	Deal_id       string `json:"UF_CRM_TASK"`
	Responible_id string `json:"RESPONSIBLE_ID"`
}
type Task struct {
	Body Task_body `json:"fields"`
	// Description string    `json:"DESCRIPTION"`
}

type Bitrix struct {
	Deal Deal
}
