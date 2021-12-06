package Protocol

import "errors"

/**
 convert 4 bytes to int32 with big endian
 */
func bytesToInt32BE(bytes []byte) (i int32,err error) {
	if len(bytes) != 4 {
		return 0 , errors.New("bytes length must be 4 ")
	}
	return int32(bytes[0])<<24 |
		int32(bytes[1])<<16 |
		int32(bytes[2])<<8 |
		int32(bytes[3]) , nil
}

/**
convert 8 bytes to int32 with big endian
*/
func bytesToInt64BE(bytes []byte) (result int64,err error) {
	if len(bytes) != 8 {
		return 0 , errors.New("bytes length must be 8 ")
	}
	for i, b := range bytes {
		result |= int64(b)<< (8 * (7 - i))
	}
	return result , nil
}