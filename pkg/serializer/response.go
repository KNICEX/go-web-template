package serializer

type Response struct {
	Code  ResCode     `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

type PageResponse struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	List     any `json:"list"`
}

type InfiniteResponse struct {
	Count   int  `json:"count"`
	List    any  `json:"list"`
	HasMore bool `json:"has_more"`
}
