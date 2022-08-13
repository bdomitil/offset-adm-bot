package bitrix

type Deal struct {
	Result []struct {
		Id    string `json:"ID"`
		Title string `json:"TITLE"`
	} `json:"result"`
}

type Task struct {
	Fields struct {
		Description   string `json:"DESCRIPTION"`
		Title         string `json:"TITLE"`
		Id            string `json:"ID"`
		Deal_id       string `json:"UF_CRM_TASK"`
		Responible_id string `json:"RESPONSIBLE_ID"`
	} `json:"fields"`
}

type Profile struct {
	Result struct {
		Id        string `json:"ID"`
		Name      string `json:"NAME"`
		Last_name string `json:"LAST_NAME"`
		IsAdmin   bool   `json:"ADMIN"`
	} `json:"result"`
}
