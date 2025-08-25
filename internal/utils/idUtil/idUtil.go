package idUtil

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go-chat/internal/utils/logUtil"
)

// GenerateId 生成一个长度为 10 的纯数字id 不保证唯一
func GenerateId() string {
	id, err := gonanoid.Generate("1234567890", 10)
	if err != nil {
		logUtil.Errorf("生成id失败:%s", err.Error())
	}
	return id
}
