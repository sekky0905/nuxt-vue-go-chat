package testutil

import (
	"fmt"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
)

// GenerateThreadHelper generates and returns Thread slice.
func GenerateThreadHelper(startNum, endNum int) []*model.Thread {
	num := endNum - startNum + 1

	threads := make([]*model.Thread, num, num)
	i := 0
	for j := startNum; j < endNum+1; j++ {
		thread := &model.Thread{
			ID:    uint32(j),
			Title: fmt.Sprintf("%s%d", model.TitleForTest, j),
			User: &model.User{
				ID:   model.UserValidIDForTest,
				Name: model.UserNameForTest,
			},
			CreatedAt: TimeNow(),
			UpdatedAt: TimeNow(),
		}
		threads[i] = thread
		i++
	}

	return threads
}
