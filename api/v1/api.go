package v1

type ListTableResponse struct {
	Tables []string `json:"tables"`
}

type GetItemResponse struct {
	Item string `json:"item"`
}
