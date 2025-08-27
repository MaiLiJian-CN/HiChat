package common

import "net/http"


type H struct{
	Code int
	Msg	string
	Data interface{}
	Rows interface{}
	Total interface{}
}

func Resp(w http.ResponseWriter,code int,data interface{},msg string){
	w.Header().Set("Content-Type","applicaion/json")
	w.WriteHeader(http.StatusOK)
	// h:=H{
	// 	Code: code,
	// 	Data: data,
	// 	Msg: msg,
	// }

}