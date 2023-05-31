package connection

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"webchat/user"
)

// 用户中心，维护多个用户的connection
var DefaultH = &Hub{
	c: make(map[*Connection]bool),
	u: make(chan *Connection),
	b: make(chan []byte),
	r: make(chan *Connection),
}

// 用户连接结构体
type Connection struct {
	ws   *websocket.Conn
	h    *Hub
	sc   chan []byte
	data *Data
}

// 用户在线名单列表
var userList []string

type Data struct {
	Ip       string   `json:"ip"`
	User     string   `json:"user"`
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
	Token    string   `json:"token"`
}

// websocket服务
func Upgrade(ctx *gin.Context) {
	// 协议升级
	upgrader := &websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "websocket upgrade error")
		return
	}
	// 创建连接
	c := &Connection{sc: make(chan []byte, 256), ws: conn, h: DefaultH, data: &Data{}}
	// connection加入hub管理
	DefaultH.r <- c
	go c.writer()
	c.reader()
	// 退出登录
	defer logout(c)
}

// 数据写入器
func (c *Connection) writer() {
	// 取出发送信息并写入
	for message := range c.sc {
		// fmt.Println(message, "\n")
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	c.ws.Close()
}

// 数据读取器
func (c *Connection) reader() {
	for {
		// 接收ws信息
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.h.r <- c
			break
		}
		json.Unmarshal(message, &c.data)
		// 解析信息类型
		switch c.data.Type {
		// 用户登录
		case "login":
			token := c.data.Token
			userInfo, err := user.Oauth(token)
			if err != nil {
				// 推出到登录页面
				c.data.Type = "relogin"
				data_b, _ := json.Marshal(c.data)
				// 推出到登录页面
				c.sc <- data_b
				return
			}
			c.data.User = userInfo.Name
			c.data.From = userInfo.Name
			// 在线人数增加
			userList = append(userList, c.data.User)
			c.data.UserList = userList
			data_b, _ := json.Marshal(c.data)
			// 发送信息
			c.h.b <- data_b
		case "user":
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)
			c.h.b <- data_b
		// 用户登出
		case "logout":
			c.data.Type = "logout"
			// 在线人数减少
			userList = del(userList, c.data.User)
			data_b, _ := json.Marshal(c.data)
			// 删除连接
			c.h.b <- data_b
			// 发送用户离线信息
			c.h.r <- c
		default:
			fmt.Print("========default================")
		}
	}
}

// 删除登出的用户，维护在线用户名单
func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return n_slice
}

// 退出
func logout(c *Connection) {
	c.data.Type = "logout"
	userList = del(userList, c.data.User)
	c.data.UserList = userList
	c.data.Content = c.data.User
	data_b, _ := json.Marshal(c.data)
	c.h.b <- data_b
	c.h.r <- c
}

type Hub struct {
	// 当前在线connection信息
	c map[*Connection]bool
	// 删除connection
	u chan *Connection
	// 传递数据
	b chan []byte
	// 加入connection
	r chan *Connection
}

func (h *Hub) Run() {
	for {
		select {
		// 用户连接，添加connection信息
		case c := <-h.r:
			h.c[c] = true
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			c.data.UserList = userList
			data_b, _ := json.Marshal(c.data)
			// 发送给写入器
			c.sc <- data_b
		// 删除指定用户连接
		case c := <-h.u:
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.sc)
			}
		// 向聊天室在线人员发送信息
		case data := <-h.b:
			for c := range h.c {
				select {
				// 发送数据
				case c.sc <- data:
				// 发送不成功则删除connection信息
				default:
					delete(h.c, c)
					close(c.sc)
				}
			}
		}
	}
}
