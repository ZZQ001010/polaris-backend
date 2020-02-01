package domain

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/tests"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//func TestCreateIssue(t *testing.T) {
//
//	config.LoadConfig("/Users/tree/work/08_all_star/01_src/go/polaris-backend/polaris-server/configs", "application")
//
//	pid, _ := idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue")
//
//	orgId := int64(1000)
//	projectId := int64(1)
//
//	code, _ := idfacade.ApplyCode(orgId, "PC", "")
//
//	issueBo := &bo.IssueBo{
//		Id:                  pid,
//		OrgId:               orgId,
//		Code:                code,
//		ProjectId:           projectId,
//		ProjectObjectTypeId: 1234,
//		Title:               "1AAA",
//		Owner:               345,
//		PriorityId:          234,
//		SourceId:            345,
//		IssueObjectTypeId:   456,
//		PlanStartTime:       types.NowTime(),
//		PlanEndTime:         types.NowTime(),
//		StartTime:           types.NowTime(),
//		EndTime:             types.NowTime(),
//		PlanWorkHour:        4,
//		IterationId:         0,
//		VersionId:           0,
//		ModuleId:            0,
//		ParentId:            0,
//		Status:              567,
//		Creator:             345,
//		CreateTime:          types.NowTime(),
//		Updator:             345,
//		UpdateTime:          types.NowTime(),
//		Version:             1,
//	}
//
//	pdid, _ := idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_detail")
//
//	issueDetailBo := &bo.IssueDetailBo{
//		Id:         pdid,
//		OrgId:      orgId,
//		IssueId:    pid,
//		ProjectId:  projectId,
//		StoryPoint: 0,
//		Tags:       "tags",
//		Remark:     consts.TcRemark,
//		Status:     123,
//		Creator:    345,
//		CreateTime: types.NowTime(),
//		Updator:    345,
//		UpdateTime: types.NowTime(),
//		Version:    1,
//	}
//
//	issueBo.IssueDetailBo = *issueDetailBo
//
//	prid, _ := idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_relation")
//	issueBo.OwnerInfo = &bo.IssueUserBo{
//		IssueRelationBo: bo.IssueRelationBo{
//			Id:           prid,
//			OrgId:        orgId,
//			IssueId:      pid,
//			RelationId:   345,
//			RelationType: 2,
//			Creator:      345,
//			CreateTime:   types.NowTime(),
//			Updator:      345,
//			UpdateTime:   types.NowTime(),
//			Version:      1,
//		},
//	}
//
//	relationBos := make([]bo.IssueUserBo, 10)
//
//	for i := int64(0); i < 10; i++ {
//		prid, _ = idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_relation")
//		relationBos[i] = bo.IssueUserBo{
//			IssueRelationBo: bo.IssueRelationBo{
//				Id:           prid,
//				OrgId:        orgId,
//				IssueId:      pid,
//				RelationId:   i + 10,
//				RelationType: 2,
//				Creator:      345,
//				CreateTime:   types.NowTime(),
//				Updator:      345,
//				UpdateTime:   types.NowTime(),
//				Version:      1,
//			},
//		}
//	}
//
//	relationBos2 := make([]bo.IssueUserBo, 10)
//
//	for i := int64(0); i < 10; i++ {
//		prid, _ = idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_relation")
//		relationBos2[i] = bo.IssueUserBo{
//			IssueRelationBo: bo.IssueRelationBo{
//				Id:           prid,
//				OrgId:        orgId,
//				IssueId:      pid,
//				RelationId:   i + 20,
//				RelationType: 2,
//				Creator:      345,
//				CreateTime:   types.NowTime(),
//				Updator:      345,
//				UpdateTime:   types.NowTime(),
//				Version:      1,
//			},
//		}
//	}
//
//	issueBo.FollowerInfos = &relationBos2
//	issueBo.ParticipantInfos = &relationBos
//
//	err3 := CreateIssue(issueBo)
//
//	fmt.Println(err3)
//	fmt.Println(issueBo.Id)
//}

func TestGetIssueBo(t *testing.T) {
	convey.Convey("获取任务bo", t, tests.StartUp(func() {
		convey.Convey("获取任务bo", func() {
			issueBo, _ := GetIssueBo(1000, 1795)
			fmt.Println(issueBo.PlanStartTime)
			fmt.Println(issueBo.PlanEndTime)

			newIssueBo := &bo.IssueBo{}
			_ = copyer.Copy(issueBo, newIssueBo)

			fmt.Println(newIssueBo.PlanStartTime)
			fmt.Println(newIssueBo.PlanEndTime)
		})
	}))

}

func TestConvertIssueTagBosToMapGroupByIssueId(t *testing.T) {

	issueTagBos := []bo.IssueTagBo{
		{
			IssueId:1,
			TagId:1,
			TagName:"hello",
		},
		{
			IssueId:1,
			TagId:2,
			TagName:"world",
		},
		{
			IssueId:2,
			TagId:1,
			TagName:"hello",
		},
		{
			IssueId:2,
			TagId:2,
			TagName:"world",
		},
	}
	lm := ConvertIssueTagBosToMapGroupByIssueId(issueTagBos)
	t.Log(json.ToJsonIgnoreError(lm))
}

func TestGetCalendarInfo(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetCalendarInfo(1082, 1744))
	}))

}