package logic

import (
	"fmt"
	"testing"
)

func TestCheckSentence(t *testing.T) {
	sensitiveList := []string{"中国", "中国人"}
	input := "我来自中1国cd中2国  中3国"

	util := NewDFAUtil(sensitiveList)
	newInput := util.CheckSentence(input)
	fmt.Println(newInput)
}
