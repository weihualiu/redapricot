package p

import (
	"errors"
	"time"
	"bytes"
	"encoding/binary"

	m_time "github.com/weihualiu/redapricot/common/time"
	m_bytes "github.com/weihualiu/redapricot/common/bytes"
)

// 组装数据包
// package len|package cmd type | pacakge content
// package content = fname | fcontent

// 组包接口
type Packager interface {
	Build() ([]byte,error)
	Parse([]byte) error
	// 设置业务数据报文
	SetBody(byte,[]byte) error
	// 获取业务数据报文
	GetBody() []byte
	GetCmd() byte
}

type PackageCommon struct {
	header  byte
	len     uint32 //4bytes
	cmdType byte   //数据类型 0 heartbeat 1 api 2 log
	date    []byte //7bytes
	body    []byte
	tail    byte
}

func NewPackageCommon() *PackageCommon {
	this := new(PackageCommon)
	this.header = byte(0xF0)
	this.tail = byte(0xFE)
	// 7bytes 20171030151201
	this.date = m_time.TimeToBytes(time.Now())

	return this
}

func (this PackageCommon)Build() ([]byte, error) {
	this.len = uint32(1 + 4 + 1 + 7) + uint32(len(this.body)) + uint32(1)

	//buf := make([]byte, this.len)
	buffer := new(bytes.Buffer)
	buffer.WriteByte(this.header)
	binary.Write(buffer, binary.BigEndian, this.len)
	buffer.WriteByte(this.cmdType)
	buffer.Write(this.date)
	buffer.Write(this.body)
	buffer.WriteByte(this.tail)

	return buffer.Bytes(), nil
}

func (this *PackageCommon)Parse(data []byte) error {
	//packComm := new(PackCommon)
	this.header = byte(0xF0)
	this.tail = byte(0xFE)

	if data[0] != this.header || data[len(data)-1] != this.tail {
		return errors.New("data struct parse failed, err data header!")
	}

	this.len = m_bytes.BytesToUInt32(data[1:5])
	if this.len != uint32(len(data)) {
		return errors.New("data struct parse failed, package length failed!")
	}
	this.cmdType = byte(data[5])
	beforeLen := 1 + 4 + 1
	this.date = data[beforeLen : beforeLen+7]
	this.body = data[beforeLen+7 : len(data)-1]
	//输出解析的内容
	//log.Println("this.cmdType: ", this.cmdType)
	//log.Println("this.date: ", this.date)
	//log.Println("this.body: ", this.body)
	return nil

}

func (this *PackageCommon)SetBody(cmd byte, data []byte) error {
	this.cmdType = cmd
	this.body = data
	return nil
}

func (this PackageCommon)GetBody() []byte {
	return this.body
}

func (this PackageCommon)GetCmd() byte {
	return this.cmdType
}
