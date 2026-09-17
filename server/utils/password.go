package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// 注意：演示用简单哈希（SHA256 + 固定盐）。生产环境请替换为 bcrypt/scrypt 等强哈希。
var passwordPepper = "drone-server@2024"

// HashPassword 计算密码哈希
func HashPassword(pw string) string {
	h := sha256.Sum256([]byte(pw + passwordPepper))
	return hex.EncodeToString(h[:])
}

// CheckPassword 校验密码
func CheckPassword(pw, hash string) bool {
	return HashPassword(pw) == hash
}
