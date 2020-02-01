package service

import (
	"context"
	"github.com/galaxy-book/polaris-backend/common/model/vo/msgvo"
	"github.com/galaxy-book/polaris-backend/service/basic/msgsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"gotest.tools/assert"
	"testing"
)

func TestSendEmail(t *testing.T) {
	convey.Convey("TestSendEmail", t, test.StartUp(func(ctx context.Context) {
		err := SendMail(msgvo.SendMailReqVo{
			Input: msgvo.SendMailReqData{
				Emails: []string{"ainililia@163.com"},
				Subject: "text",
				Content: "body",
			},
		})
		assert.Equal(t, err, nil)
	}))

}