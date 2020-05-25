package restapi

type listResponse struct {
	Size int         `json:"size"`
	Data interface{} `json:"data"`
}
