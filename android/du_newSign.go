package android

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"io"
	"strings"
)

func encryptPassword(data string) (result string) {
	salt := "du"
	h := md5.New()
	_, _ = io.WriteString(h, data)
	_, _ = io.WriteString(h, salt)
	result = hex.EncodeToString(h.Sum(nil))
	return
}
func padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func AesEncrypt(origData []byte) string {
	skey := "d245a0ba8d678a61"

	block, _ := aes.NewCipher([]byte(skey))
	origData = padding(origData, block.BlockSize())

	buffer := bytes.NewBufferString("")
	tmpData := make([]byte, block.BlockSize()) //存储每次加密的数据
	for index := 0; index < len(origData); index += block.BlockSize() {
		block.Encrypt(tmpData, origData[index:index+block.BlockSize()])
		buffer.Write(tmpData)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}
func newSign(uuid, use, timestamp, pwd string) string {
	sb := strings.Builder{}
	//sb.WriteString("countryCode86loginTokenpassword9ec4df545311da06c312bbecf8d82b51platformandroidtimestamp1655993108827typepwduserName16688889999uuidca69583c3bbbbd72v4.94.5")
	sb.WriteString("countryCode86loginTokenpassword")
	sb.WriteString(pwd)
	sb.WriteString("platformandroidtimestamp")
	sb.WriteString(timestamp)
	sb.WriteString("typepwduserName")
	sb.WriteString(use)
	sb.WriteString("uuid")
	sb.WriteString(uuid)
	sb.WriteString("v4.94.5")

	m := md5.New()
	_, _ = io.WriteString(m, AesEncrypt([]byte(sb.String())))
	Sign := hex.EncodeToString(m.Sum(nil))
	return Sign
}

func main() {
	use := "16688889999"
	pwd := "a12345678"
	uuidWithHyphen, _ := uuid.NewRandom()
	u := strings.Replace(uuidWithHyphen.String(), "-", "", -1)[:16]
	//timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	pwd = encryptPassword(pwd)
	timeStamp := "1655993108827"
	u = "ca69583c3bbbbd72"
	Sign := newSign(u, use, timeStamp, pwd)
	fmt.Printf("%s\n%s", pwd, Sign)
	// 9ec4df545311da06c312bbecf8d82b51
	// 6cac57368950fcf7fcf26f02648963de
}
