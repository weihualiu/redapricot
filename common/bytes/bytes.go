package bytes

import (
	"encoding/binary"
	"strconv"
	log "github.com/Sirupsen/logrus"
)

func BytesToUInt32(buf []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(buf))
}

func BytesToString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func BytesTrim(c []byte) []byte {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return c[:n+1]
}

func TimeBytes2String(c []byte) string {
	if len(c) != 7 {
		return ""
	}
	var s string
	for i := 0; i < len(c); i++ {
		log.Debugln("%d", int(c[i]))
		 s += strconv.FormatInt(int64(c[i]), 10)
	}
	log.Debugln("s:",s)
	return  s
}

