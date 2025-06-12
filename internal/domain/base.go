package domain

type BaseResponse struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status_code"`
	Message    any  `json:"message"`
}
