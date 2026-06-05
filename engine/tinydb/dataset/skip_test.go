package dataset

import (
	"fmt"
	"math/rand"
	"testing"
	_ "testing"
)

func TestAdd(t *testing.T) {
	skipList := NewSkipList[*PassStuct]()
	for i := 0; i < 1000000; i++ {
		// 生成1-100之间的随机数
		randomNum := rand.Intn(30000) + 1
		skipList.Insert(randomNum, &PassStuct{
			Name: fmt.Sprintf("张-%d", randomNum),
		})
	}
	skipList.Get(100)
	skipList.GetWith(">=", 100, 2)
	skipList.GetWith("like", "1%", 1)
}
