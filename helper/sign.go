package helper

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"sort"
)

// Sign 微信公众号 url 签名.
func Sign(token, timestamp, nonce string) (signature string) {
	strArr := sort.StringSlice{token, timestamp, nonce}
	strArr.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce))
	buf = append(buf, strArr[0]...)
	buf = append(buf, strArr[1]...)
	buf = append(buf, strArr[2]...)

	hashSum := sha1.Sum(buf)
	return hex.EncodeToString(hashSum[:])
}

// MsgSign 微信公众号/企业号 消息体签名.
func MsgSign(token, timestamp, nonce, encryptedMsg string) (signature string) {
	strArr := sort.StringSlice{token, timestamp, nonce, encryptedMsg}
	strArr.Sort()

	h := sha1.New()

	bufWriter := bufio.NewWriterSize(h, 128) // sha1.BlockSize 的整数倍
	_, _ = bufWriter.WriteString(strArr[0])
	_, _ = bufWriter.WriteString(strArr[1])
	_, _ = bufWriter.WriteString(strArr[2])
	_, _ = bufWriter.WriteString(strArr[3])
	_ = bufWriter.Flush()

	hashSum := h.Sum(nil)
	return hex.EncodeToString(hashSum)
}
