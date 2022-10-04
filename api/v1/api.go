package v1

type ListTableResponse struct {
	Tables []string `json:"tables"`
}

type GetItemResponse struct {
	Item string `json:"item"`
}

type RollResponse struct {
	Result      int    `json:"result"`
	Description string `json:"description"`
}
