package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetProjectAttachment(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		var filetype int = 1
		var keyWord string = "好"
		input := vo.ProjectAttachmentReq{
			ProjectID: 1006,
			FileType:  &filetype,
			KeyWord:   &keyWord,
		}
		res, err := GetProjectAttachment(1003, 1003, 1, 10, input)
		fmt.Println(err)
		fmt.Println(res.Total)
		for _, value := range res.List {
			fmt.Println("resource:")
			fmt.Println(value)
			fmt.Println("issue:")
			for _, dv := range value.IssueList {
				fmt.Println(dv.ID)
			}
		}
	}))
}

func TestDeleteProjectAttachment(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		input := vo.DeleteProjectAttachmentReq{
			ProjectID:   1006,
			ResourceIds: []int64{10504},
		}
		res, err := DeleteProjectAttachment(1003, 1003, input)
		fmt.Println(err)
		fmt.Println(res.ResourceIds)
	}))
}
