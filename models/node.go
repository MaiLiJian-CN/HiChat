package models

import (
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

//Node 构造连接
type Node struct {
    Conn      *websocket.Conn //socket连接
    Addr      string          //客户端地址
    DataQueue chan []byte     //消息内容
    GroupSets set.Interface   //好友 / 群
}