package util

type CommonError struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
}
