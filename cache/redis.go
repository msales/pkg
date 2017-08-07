package cache

import "strconv"

var redisDecoder = stringDecoder{}

type stringDecoder struct{}

func (d stringDecoder) Int64(v []byte) (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

func (d stringDecoder) Float64(v []byte) (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}
