package node

import (
	"p00q.cn/video_cdn/server/global"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	go func() {
		time.Sleep(time.Second)
		whetherThereIsAWaitingRecipient(1, Msg{SessionCode: 1, Data: "test"})
	}()
	response := sendAMessageAndWaitForAResponse(nil, 1, Msg{SessionCode: 1}, time.Minute)
	global.Logs.Info(response.Data)
}
