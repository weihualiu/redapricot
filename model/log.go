package model

// 分析上传日志消耗

//  UseTime 4bytes
//  Tag  100bytes
//  MainTag 从Tag中分析得出
//  AppName 20bytes
//  Level 从Tag中分析得出
//  ExecTime 7bytes
//  Content Nbytes

import (
	"github.com/weihualiu/redapricot/common/bytes"
	"strings"
	"github.com/weihualiu/redapricot/socket"
	"github.com/weihualiu/redapricot/common/db"
	"fmt"
	"github.com/weihualiu/redapricot/p"
	log "github.com/Sirupsen/logrus"
)

type LogData struct {
	UseTime uint32 //消耗时间
	Tag string // 标识
	MainTag string // 主标识
	ClientTag string // 客户端请求业务标识，可涵盖多个服务端请求接口
	Content []byte // 日志内容
	Appname string // 应用名称
	Level string // 调用链层级
	ExecTime string // 执行时间点
	FuncName string // 业务功能标识
}

func NewLogData() *LogData {
	return new(LogData)
}

func (this *LogData)Parse(data []byte) error {
	this.UseTime = bytes.BytesToUInt32(data[0:4])
	this.Tag = bytes.BytesToString(data[4:104])
	this.ClientTag = bytes.BytesToString(data[104:204])
	this.Appname = bytes.BytesToString(data[204:224])
	this.ExecTime = bytes.TimeBytes2String(data[224:231])
	this.Content = data[231:]

	arr := strings.Split(this.Tag, "_")
	// XXXX_XXX_1
	this.MainTag = arr[0]
	this.FuncName = arr[1]
	this.Level = arr[2]

	log.Debugln("data size:", len(data), ",data:", data)
	log.Debugln("tag:", this.Tag)
	log.Debugln("client tag:", this.ClientTag)
	log.Debugln("appname:", this.Appname)
	log.Debugln("exectime:", this.ExecTime)
	log.Debugln("content:", this.Content)

	return  nil
}

func Handler(data []byte, response *socket.SocketResponse) {
	log.Debugln("log handler process......")
	this := new(LogData)
	this.Parse(data)

	Sql := fmt.Sprintf("insert into data_log (tag, client_tag, funcname, level, exectime, appname, maintag, usetime, content) values ('%s', '%s', '%s', %s, '%s', '%s', '%s', %d, ?)",
		this.Tag, this.ClientTag, this.FuncName, this.Level, this.ExecTime, this.Appname, this.MainTag, this.UseTime)
	_, err := db.DB.Exec(Sql, this.Content)
	if err != nil {
		log.Println("exec", Sql, "fail", err)
		response.Data <- []byte{p.FAILED}
	}
	response.Data <- []byte{p.NORMAL}
}

func init() {
	socket.RegisterHandler(p.PROTOCOL_LOG, Handler)
}