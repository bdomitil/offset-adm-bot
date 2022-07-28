package bitrix

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
}

type Bitrix struct {
	Deal Deal
}

type Profile_body struct {
	Id        string `json:"ID"`
	Name      string `json:"NAME"`
	Last_name string `json:"LAST_NAME"`
	IsAdmin   bool   `json:"ADMIN"`
}

type Profile struct {
	Body Profile_body `json:"result"`
}
