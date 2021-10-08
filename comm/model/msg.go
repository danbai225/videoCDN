package model

type Msg struct {
	SessionCode uint64      `json:"session_code"`
	Err         string      `json:"err"`
	Data        interface{} `json:"data"`
}
