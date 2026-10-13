package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

// 链接模块
type Connection struct {
	//当前链接的socket tcp套接字
	Conn *net.TCPConn

	//链接的id
	ConnID uint32

	//当前链接状态
	isClosed bool

	//当前链接所绑定的处理业务方法api
	handleAPI ziface.HandleFunc

	//告知当前链接已经推出的channel
	ExitChan chan bool
}

// 初始化链接方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// 链接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader goRoutine is running...")
	defer fmt.Println("connID =", c.ConnID, "Reader is exit, remoteAddr =", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端数据到buf
		buf := make([]byte, 1024)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err =", err)
			continue
		}
		//调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID =", c.ConnID, "Handle is err =", err)
			break
		}

	}
}

func (c *Connection) Start() {
	fmt.Println("Connection Start...ConnID = ", c.ConnID)
	//启动业务
	go c.StartReader()
	//TODO 启动从当前链接写数据

}
func (c *Connection) Stop() {
	fmt.Println("Connection Stop...ConnID = ", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭链接
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前模块的链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的tcp状态
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) Send(date []byte) error {
	return nil
}
