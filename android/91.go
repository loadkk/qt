package android

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesDecrypt(text string) string {
	skey, _ := hex.DecodeString("304d44515251675949613965335a6e73")

	block, _ := aes.NewCipher([]byte(skey))

	src, _ := base64.StdEncoding.DecodeString(text)
	/**AES 解密**/
	buffer := bytes.NewBufferString("")
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Decrypt(tmpData, src[index:index+block.BlockSize()])
		buffer.Write(tmpData)
	}

	// 去掉末尾非打印控制字符
	var deByte []byte
	for i := len(buffer.Bytes()); i > 0; i-- {
		if buffer.Bytes()[i-1] >= 32 {
			deByte = buffer.Bytes()[:i]
			break
		}
	}
	return strings.TrimSpace(string(deByte))
}

func AesEncrypt(origData []byte) (string, error) {
	skey, _ := hex.DecodeString("304d44515251675949613965335a6e73")

	// key只能是 16 24 32长度
	block, err := aes.NewCipher([]byte(skey))
	if err != nil {
		return "", err
	}
	//padding
	origData = padding(origData, block.BlockSize())
	//存储每次加密的数据

	//分组分块加密
	buffer := bytes.NewBufferString("")
	tmpData := make([]byte, block.BlockSize()) //存储每次加密的数据
	for index := 0; index < len(origData); index += block.BlockSize() {
		block.Encrypt(tmpData, origData[index:index+block.BlockSize()])
		buffer.Write(tmpData)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func getData(code string) string {
	retval := make(map[string]interface{})
	//m := map[string]interface{}{"code": "123456789"}
	m := map[string]interface{}{"code": code}
	jsonstr, _ := json.Marshal(m)
	data := fmt.Sprintf("%s", jsonstr)
	encrypt, _ := AesEncrypt([]byte(data))
	retval["data"] = encrypt
	retval["handshake"] = "v20200610" // 固定值 估计是版本之类的
	retvalStr, _ := json.Marshal(retval)
	fmt.Printf("getData: %s\n", retvalStr)
	return fmt.Sprintf("%s", retvalStr)
}

func xToken(deviceId, useToken string) string {
	m := make(map[string]string)
	//m["device_no"] = "fec2d283-bd58-37d0-9b7c-e2f3b0e008e7" // device_id
	m["device_no"] = deviceId // device_id
	m["device_type"] = "A"    // 固定值
	//m["token"] = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczpcL1wvZ2Z3Lmdvb2dsZS5jb20iLCJhdWQiOiJhcHAtdXNlcnMiLCJleHAiOjE2NTYzNjQ2MjEsImp0aSI6ImU0YWZhOGRiMjY0ODA2ZGZkMTJmODJlYjE2N2MzYzQ4IiwiaWF0IjoxNjU2MjU2NjIxLCJ1c2VyX2lkIjoiVDNFOVIiLCJ2ZXJzaW9uIjoiMS4wLjQiLCJkZXZpY2Vfbm8iOiJmZWMyZDI4My1iZDU4LTM3ZDAtOWI3Yy1lMmYzYjBlMDA4ZTciLCJkZXZpY2VfdHlwZSI6IkEiLCJpcCI6IjE3NS4xNTMuMTY5LjEwNyJ9.HLEc9UdZGTj-sHM7tkW-ZyJhqw-K45XzAf15iNOJF40"
	m["token"] = useToken
	m["version"] = "1.0.4" // 固定值
	mj, _ := json.Marshal(m)
	data := fmt.Sprintf("%s", mj)
	encrypt, _ := AesEncrypt([]byte(data))
	fmt.Println("xToken:", encrypt)
	return encrypt
}

func sendRequests(params string, token string) {
	client := &http.Client{}
	data := strings.NewReader(params)
	url := "http://api.ccav69api.com/app/api/user/bindcode"
	req, _ := http.NewRequest("POST", url, data)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; U; Android 2.1; en-us; Nexus One Build/ERD62) AppleDart/530.17 (KHTML, like Gecko) Version/4.0 Mobile Safari/530.17")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-TOKEN", token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bodyText, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s\n", bodyText)
}

func random_andrvice_device_id() string {
	// fec2d283-bd58-37d0-9b7c-e2f3b0e008e7
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := md5.New()
	_, _ = io.WriteString(h, t)
	hashID := hex.EncodeToString(h.Sum(nil))
	sb := strings.Builder{}
	sb.WriteString(hashID[:8])
	sb.WriteString("-")
	sb.WriteString(hashID[8:12])
	sb.WriteString("-")
	sb.WriteString(hashID[12:16])
	sb.WriteString("-")
	sb.WriteString(hashID[16:20])
	sb.WriteString("-")
	sb.WriteString(hashID[20:])
	//fmt.Println(sb.String())
	return sb.String()
}

func getJWT(id string) string {
	jsonData := make(map[string]string)
	jsonData["channel"] = ""
	jsonData["code"] = ""
	//jsonData["device_no"] = "fec2d283-bd58-37d0-9b7c-e2f3b0e008e7"
	jsonData["device_no"] = id
	jsonData["device_type"] = "A"
	jsonData["version"] = "1.0.4"
	jsonDataM, _ := json.Marshal(jsonData)
	encryptData := fmt.Sprintf("%s", jsonDataM)
	//fmt.Println(encryptData_b)
	encryptData, _ = AesEncrypt([]byte(encryptData))

	jsonData = map[string]string{}
	jsonData["data"] = encryptData
	jsonData["handshake"] = "v20200610"
	jsonDataM, _ = json.Marshal(jsonData)

	url := "http://api.ccav69api.com/app/api/auth/login/device"
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", url, bytes.NewReader(jsonDataM))
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Linux; U; Android 2.1; en-us; Nexus One Build/ERD62) AppleDart/530.17 (KHTML, like Gecko) Version/4.0 Mobile Safari/530.17")
	reqest.Header.Add("Content-Type", "application/json")
	response, _ := client.Do(reqest)
	defer response.Body.Close()
	bodyText, _ := io.ReadAll(response.Body)
	//fmt.Printf("%s\n", bodyText)
	return fmt.Sprintf("%s", bodyText)
}

type JWT struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      string `json:"data"`
	Handshake string `json:"handshake"`
}

type Auth struct {
	Auth struct {
		Token     string `json:"token"`
		ExpiredAt string `json:"expired_at"`
	} `json:"auth"`
}

func main() {
	deviceId := random_andrvice_device_id()
	JWTStr := getJWT(deviceId)
	var token JWT
	var _token Auth
	_ = json.Unmarshal([]byte(JWTStr), &token)
	_ = json.Unmarshal([]byte(AesDecrypt(token.Data)), &_token)

	useToken := _token.Auth.Token
	fmt.Println("JWT.Token:", useToken)
	cdoe := "T3E9r" // 邀请码
	data := getData(cdoe)
	XToken := xToken(deviceId, useToken)
	sendRequests(data, XToken)
}
