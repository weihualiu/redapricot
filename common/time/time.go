package time

import "time"

func TimeToBytes(datetime time.Time) []byte {
	buf := make([]byte, 7)
	format := datetime.Format("20060102150405")
	for i,j := 0,0; i< len(format); i, j = i+2, j+1 {
		buf[j] = format[i]*10 + format[i+1]
	}

	return buf
}
