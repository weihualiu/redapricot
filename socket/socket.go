package socket

import (
	"net"
	log "github.com/Sirupsen/logrus"
	
	g "github.com/weihualiu/redapricot/cfg"
	//"io"
	"syscall"
	"github.com/weihualiu/redapricot/p"
	"bytes"
	_ "time"
	m_bytes "github.com/weihualiu/redapricot/common/bytes"
	"time"
)

func Start() {
	if !g.Config().Socket.Enabled {
		return
	}

	addr := g.Config().Socket.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("socket listening", addr)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		go socketHandler(conn)
	}
}

func socketHandler(conn net.Conn) {
	//解析报文
	log.Debugln("new conenction come in.")
	// one package buffer
	packageData := make([]byte, 0)
	// loop receive
	for {
		readBuf := make([]byte, 1024)
		readLen, err := conn.Read(readBuf)
		log.Debugln("data:", readBuf)
		switch err {
		case nil:
			packageData = append(packageData, readBuf[:readLen]...)
			stick(packageData, conn)
		//case io.EOF:
		//	goto DISCONNECT
		case syscall.EAGAIN:
			continue
		default:
			log.Println(err)
			goto DISCONNECT
		}
	}

DISCONNECT:
	packageData = nil
	err := conn.Close()
	if err != nil {
		//log.Fatal(err)
		log.Debugln("connnection is closed! ", err)
	}
	log.Debugln("Close connection: ", conn.RemoteAddr().String())
}

// 粘包处理
func stick(packageData []byte, conn net.Conn) {
	//flag := true
	//for flag {
		if packageData == nil {
			//flag = false
			return
		}
		pdLen := len(packageData)
		if pdLen == 0 || pdLen < 14 {
			//flag = false
			return
		} else if packageData[0] == byte(0xF0) {
			packageLen := m_bytes.BytesToUInt32(packageData[1:5])
	//		log.Debugln("packageLen:",packageLen)
	//		log.Debugln("data packageLen:", uint32(len(packageData)))
	//		if packageLen < 0 {
	//			packageData = nil
	//			flag = false
	//		} else
	 		if uint32(len(packageData)) >= packageLen {
	//			//如果数据满足一个完整包则进入下一步处理
	//			parseData := make([]byte, packageLen)
	//			copy(parseData, packageData[0:packageLen])
				process(packageData, conn)
	//			//减去完整包
	//			log.Debugln("packageLen:", packageLen, ",pacakgeData len:", len(packageData))
	//			if packageLen == uint32(len(packageData)) {
	//				packageData = nil
	//				flag = false
	//				//goto DISCONNECT
	//			} else {
	//				log.Println("package next read from network buffer")
	//				newSize := len(packageData) - int(packageLen)
	//				tmp := make([]byte, newSize, newSize+1)
	//				tmp = packageData[packageLen:len(packageData)]
	//				packageData = make([]byte, newSize, newSize+1)
	//				copy(packageData, tmp)
	//			}
	//		} else {
	//			flag = false
			}
	//	} else {
	//		log.Errorln("data error! ", packageData)
	//		//错误数据，抛弃
	//		packageData = nil
	//		//readBuf = make([]byte, 0)
	//		flag = false
		}
	//}
}

func process(data []byte, conn net.Conn) {
	log.Debugln("process starting......")
	// defer 处理err异常的情况
	var err error
	defer func() {
		if err != nil {
			log.Println("process is error", err.Error())

			pc := p.NewPackageCommon()
			pc.SetBody(0x00,bytes.NewBuffer([]byte{p.FAILED}).Bytes())
			d, _ := pc.Build()
			conn.Write(d)
		}
	}()

	pc := new(p.PackageCommon)
	err = pc.Parse(data)
	if err != nil {
		return
	}

	handler, err := RegisterGetHandler(uint8(pc.GetCmd()))
	if err != nil {
		return
	}

	response := NewSocketResponse(conn, int32(g.Config().Socket.Timeout))
	// 启动一个goroutine 用于send message
	go receive(response)
	handler(pc.GetBody(), response)
}

func receive(response *SocketResponse) {
	log.Debugln("receive starting......")
	//超时处理 select for
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(response.Timeout) * time.Second)
		timeout <- true
		// 关闭channel，防止memory leak
		close(timeout)
	}()

	for {
		time.Sleep(time.Duration(10) * time.Millisecond)

		select {
		case data, _ := <-response.Data:
			log.Debug("data:", data)
			p := p.NewPackageCommon()
			p.SetBody(0x00, data)
			d, _ := p.Build()
			response.Conn.Write(d)
			//response.Conn.Close()
			// 使用break不能跳出for select循环
			goto End
		case _, ok := <-timeout:
			//超时处理
			if ok {
				//p := p.NewPackageCommon()
				//p.SetBody(0x00,bytes.NewBuffer([]byte{0xFE}).Bytes())
				//d, _ := p.Build()
				//response.Conn.Write(d)
				//response.Conn.Close()
				goto End
			}
		default:
		}
	}
	End:
		response.Conn.Close()

}
