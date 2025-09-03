package models

import (
	"HiChat/global"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
)

type Message struct {
	Model
	FormId   int64  `json:"userId"`   //信息发送者
	TargetId int64  `json:"targetId"` //信息接收者
	Type     int    //聊天类型：群聊 私聊 广播
	Media    int    //信息类型：文字 图片 音频
	Content  string //消息内容
	Pic      string `json:"url"` //图片相关
	Url      string //文件相关
	Desc     string //文件描述
	Amount   int    //其他数据大小
}

// MsgTableName 生成指定数据表名
func (m *Message) MsgTableName() string {
	return "message"
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// rock:Line Security
var rwLocker sync.RWMutex

// Chat    需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
func Chat(w http.ResponseWriter, r *http.Request) {
	//1.  获取参数信息发送者userId
	query := r.URL.Query()
	Id := query.Get("userId")
	userId, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		zap.S().Info("类型转换失败", err)
		return
	}

	//socket
	var isvalida = true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取socket连接,构造消息节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//将userId和Node绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//服务发送消息
	go sendProc(node)
	//服务接收消息
	go recProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

// sendProc 从node中获取信息并写入websocket中
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("写入消息失败", err)
				return
			}
			fmt.Println("数据发送socket成功")
		}
	}
}

// recProc 从websocket中将消息体拿出，然后进行解析，再进行信息类型判断， 最后将消息发送至目的用户的node中
func recProc(node *Node) {
	for {
		//get message
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("读取消息失败", err)
			return
		}
		//将消息体放入全局channel中
		broMsg(data)

		// //简单实现的一种方法
		// msg := Message{}
		// err = json.Unmarshal(data, &msg)

		// if err != nil {
		// 	zap.S().Info("json解析失败", err)
		// 	return
		// }
		// fmt.Println(msg)
		// if msg.Type==1{
		// 	zap.S().Info("这是一条私信:", msg.Content)
		// 	tarNode,ok:=clientMap[msg.TargetId]
		// 	if !ok{
		// 		zap.S().Info("不存在对应的node", msg.TargetId)
		//         return
		// 	}
		// 	tarNode.DataQueue<-data
		// 	fmt.Println("发送成功：", string(data))
		// }
	}
}

// global channel
var upSendChan chan []byte = make(chan []byte, 1024)

func broMsg(data []byte) {
	upSendChan <- data
}

// init方法，运行message包前调用
func init() {
	go UdpSendProc()
	go UpdRecProc()
}

// UdpSendProc 完成upd数据发送, 连接到udp服务端，将全局channel中的消息体，写入udp服务端
func UdpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6000,
		Zone: "",
	})
	if err != nil {
		zap.S().Info("拨号udp端口失败", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case data := <-upSendChan:
			_, err := conn.Write(data)
			if err != nil {
				zap.S().Info("写入udp消息失败", err)
				return
			}
		}
	}
}

// UpdRecProc 完成udp数据的接收，启动udp服务，获取udp客户端的写入的消息
func UpdRecProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6000,
	})
	if err != nil {
		zap.S().Info("监听udp端口失败", err)
		return
	}
	defer conn.Close()
	for {
		var buf [1024]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			zap.S().Info("读取udp数据失败", err)
			return
		}

		//处理发送逻辑
		dispatch(buf[0:n])
	}
}
func dispatch(data []byte) {
	//解析消息
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("消息解析失败", err)
		return
	}

	fmt.Println("解析数据:", msg, "msg.FormId", msg.FormId, "targetId:", msg.TargetId, "type:", msg.Type)

	switch msg.Type {
	case 1:
		//private
		sendMsgAndSave(msg.TargetId, data)
	case 2:
		//group
		sendGroupMsg(uint(msg.FormId), uint(msg.TargetId), data)
	}
}

// sendMsgTest 发送消息 并存储聊天记录到redis
func sendMsgAndSave(userId int64, msg []byte) {
	rwLocker.Lock()
	node, ok := clientMap[userId]
	rwLocker.Unlock()

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.FormId))

	//如果在线，需要即时推送
	if ok {
		node.DataQueue <- msg
	}

	//拼接记录名称
	var key string
	if userId > jsonMsg.FormId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	// 创建记录
	res, err := global.RedisDB.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	//将聊天记录写入redis缓存中
	score := float64(cap(res)) + 1
	ress, e := global.RedisDB.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if e != nil {
		fmt.Println(e)
		return
	}

	// 设置ZSET的过期时间为10天
	expirationTime := 24 * time.Hour * 7
	_, expireErr := global.RedisDB.Expire(ctx, key, expirationTime).Result()
	if expireErr != nil {
		fmt.Println(expireErr)
		return
	}

	fmt.Println("ZSET的过期时间已设置为10天")
	fmt.Println(ress)

	//将key放入全局并
	// addKeyToSaveKey(key)

}

func sendGroupMsg(formId, targetId uint, data []byte) (int, error) {
	//调用的 FindUsers(target)-community.go
	userIds, err := FindUsers(targetId)

	if err != nil {
		return -1, err
	}
	for _, v := range *userIds {
		if v != formId {
			sendMsgAndSave(int64(v), data)
		}
	}
	return 0, nil
}

// sendMs 向用户发送消息
func sendMsg(id int64, msg []byte) {
	rwLocker.Lock()
	node, ok := clientMap[id]
	rwLocker.Unlock()

	if !ok {
		zap.S().Info("userID没有对应的node")
		return
	}

	zap.S().Info("targetID:", id, "node:", node)
	if ok {
		node.DataQueue <- msg
	}
}

// MarshalBinary 需要重写此方法才能完整的msg转byte[]
func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// RedisMsg 获取缓存里面的聊天记录
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))

	//userIdStr和targetIdStr进行拼接唯一key
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}

	var rels []string
	var err error
	if isRev {
		rels, err = global.RedisDB.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = global.RedisDB.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	return rels
}
