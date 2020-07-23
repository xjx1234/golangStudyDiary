/**
* @Author: XJX
* @Description: websocket服务端示例
* @File: websocket_server.go
* @Date: 2020/7/14 14:24
 */

package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	maxConnId int64
	allWsConn map[int64]*WsConn
)

//连接结构体
type WsConn struct {
	Id         int64
	WsSocket   *websocket.Conn
	mutex      sync.Mutex
	InChan     chan *WsMessage
	OutChan    chan *WsMessage
	ClosedChan chan byte
	isClosed   bool
}

//消息结构体
type WsMessage struct {
	Type int
	Data []byte
}

//升级websocket参数配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 10 * time.Second,
	// 控制跨域，可以根据需要写入自己需要的代码
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsConn *WsConn) procLoop() {

	//心跳机制，为了处理无效连接，这机制大部分场景不需要实现，只是为了演示去实现
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := wsConn.wsWrite(websocket.TextMessage, []byte("heartbeat")); err != nil {
				wsConn.wsClose() //关闭连接
				break
			}
		}
	}()
}

func (wsConn *WsConn) wsReadLoop() {
	for {
		msgType, data, err := wsConn.WsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &WsMessage{
			msgType,
			data,
		}
		select {
		case wsConn.InChan <- req:
		case <-wsConn.ClosedChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

func (wsConn *WsConn) wsWriteLoop() {
	for {
		select {
		case msg := <-wsConn.OutChan:
			if err := wsConn.WsSocket.WriteMessage(msg.Type, msg.Data); err != nil {
				goto error
			}
		case <-wsConn.ClosedChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

//关闭连接
func (wsConn *WsConn) wsClose() {
	wsConn.WsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.ClosedChan)
	}
}

func (wsConn *WsConn) wsRead() (*WsMessage, error) {
	select {
	case msg := <-wsConn.InChan:
		return msg, nil
	case <-wsConn.ClosedChan:
	}
	return nil, errors.New("websocket closed")
}

//写数据
func (wsConn *WsConn) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.OutChan <- &WsMessage{messageType, data}:
	case <-wsConn.ClosedChan:
		return errors.New("webscoket closed")
	}
	return nil
}

//处理请求
func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket升级失败" + err.Error())
		return
	}
	maxConnId++
	conn := &WsConn{
		Id:       maxConnId,
		WsSocket: wsSocket,
		InChan:   make(chan *WsMessage, 1000),
		OutChan:  make(chan *WsMessage, 1000),
		isClosed: false,
	}
	allWsConn[maxConnId] = conn
	log.Println("total online number:", len(allWsConn))
	go conn.procLoop()
	go conn.wsReadLoop()
	go conn.wsWriteLoop()

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:8888", mux)

}
