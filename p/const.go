package p

// const define

var (
	PROTOCOL_LOG uint8 = 0x00
	PROTOCOL_CPU uint8 = 0x01
)

// 错误码定义
var (
	FAILED uint8 = 0xA0
	NETWORK_FAILED uint8 = 0xA1
)

// 正常值
var NORMAL uint8 = 0x00