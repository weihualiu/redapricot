package socket

import (
	"errors"
	"net"
	log "github.com/Sirupsen/logrus"
)

// 处理接收的数据


type SocketResponse struct {
	Conn net.Conn
	Data chan []byte
	Timeout int32
}

type socketRegister struct {
	fun func([]byte, *SocketResponse)
	protocol uint8
}

// Scoket定义结构，存储协议值与对应的处理函数
type socketContext struct {
	reg map[uint8] func([]byte, *SocketResponse)
}

var (
	context *socketContext
)

func init() {
	sc := new(socketContext)
	context = sc
}

func (this *socketRegister)setProtocol (protocol uint8) {
	this.protocol = protocol
}

func (this *socketRegister)setFunc(funcHandler func([]byte, *SocketResponse)) {
	this.fun = funcHandler
}

func (this *socketContext)register(protocol uint8, funcHandler func([]byte, *SocketResponse)) {
	if this.reg == nil {
		this.reg = make(map[uint8]func([]byte, *SocketResponse))
	}
	sr := new(socketRegister)
	sr.setProtocol(protocol)
	sr.setFunc(funcHandler)
	this.reg[protocol] = funcHandler
}

func RegisterHandler(protocol uint8, funcHandler func([]byte, *SocketResponse)) {
	log.Debugln("protocol:", protocol)
	context.register(protocol, funcHandler)
}

func (this *socketContext)registerGetHandler(protocol uint8) (func([]byte, *SocketResponse), error) {
	if this.reg == nil || len(this.reg) == 0 {
		return nil,errors.New("register arr is empty")
	}
	funcHandler, ok := this.reg[protocol]
	if !ok {
		return nil, errors.New("register not found in array")
	}
	return funcHandler, nil
}

func RegisterGetHandler(protocol uint8) (func([]byte, *SocketResponse),error) {
	return context.registerGetHandler(protocol)
}

func NewSocketResponse(conn net.Conn, timeout int32) *SocketResponse {
	sr := new(SocketResponse)
	sr.Data = make(chan []byte)
	sr.Conn = conn
	sr.Timeout = timeout
	return sr
}

func (this *SocketResponse)WriteResponse(data []byte){
	this.Data <- data
}
