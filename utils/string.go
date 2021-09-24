package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxMask = 1<<6 - 1 // All 1-bits, as many as 6
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for 10 characters!
	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		b[i] = letterBytes[int(cache&letterIdxMask)%len(letterBytes)]
		i--
		cache >>= 6
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

//Md5Encode Md5加密
func Md5Encode(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

/*
	testData := "Mm4443332221"
	testKey := "test! key"
	encode := utils.EncryptData(testData, testKey)
	fmt.Println("加密后的代码：",encode)
	decode := utils.DecryptData(encode, testKey)
	fmt.Println("解秘后的代码：",decode)
*/

//EncryptData 加密
func EncryptData(codeData string, callbackKey string) string {
	dataArr := []rune(codeData)
	keyArr := []byte(callbackKey)
	keyLen := len(keyArr)

	var tmpList []int

	for index, value := range dataArr {
		base := int(value)
		dataString := base + int(0xFF&keyArr[index%keyLen])
		tmpList = append(tmpList, dataString)
	}

	var str string

	for _, value := range tmpList {
		str += "@" + fmt.Sprintf("%d", value)
	}
	return str
}

//DecryptData 解密
func DecryptData(ntData string, callbackKey string) string {
	strLen := len(ntData)
	newData := []rune(ntData)
	resultData := string(newData[1:strLen])
	dataArr := strings.Split(resultData, "@")
	keyArr := []byte(callbackKey)
	keyLen := len(keyArr)

	var tmpList []int

	for index, value := range dataArr {
		base, _ := strconv.Atoi(value)
		dataString := base - int(0xFF&keyArr[index%keyLen])
		tmpList = append(tmpList, dataString)
	}

	var str string

	for _, val := range tmpList {
		str += string(rune(val))
	}
	return str
}
