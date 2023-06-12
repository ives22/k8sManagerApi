package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// ToMd5 计算输入值的 MD5 哈希并将哈希作为字符串返回。
func ToMd5(value string) (md5ValueStr string) {
	data := []byte(value)
	md5Hash := md5.New()
	md5Hash.Write(data)
	md5Value := md5Hash.Sum(nil)
	md5ValueStr = hex.EncodeToString(md5Value)
	return md5ValueStr
}

// ToSha1 计算输入值的 SHA-1 哈希并将哈希作为字符串返回。
func ToSha1(value string) (sha1ValueStr string) {
	data := []byte(value)
	sha1Hash := sha1.New()
	sha1Hash.Write(data)
	sha1Value := sha1Hash.Sum(nil)
	sha1ValueStr = hex.EncodeToString(sha1Value)
	return sha1ValueStr
}

// ToSha256 计算输入值的 SHA-256 哈希并将哈希作为字符串返回。
func ToSha256(value string) (sha256ValueStr string) {
	data := []byte(value)
	sha256Hash := sha256.New()
	sha256Hash.Write(data)
	sha256Value := sha256Hash.Sum(nil)
	sha256ValueStr = hex.EncodeToString(sha256Value)
	return sha256ValueStr
}