package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

//问题首页
func TestHomeIssues(t *testing.T) {

	convey.Convey("Test HomeIssues", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("HomeIssues mock req", func() {
			orderBy := 2
			status := 2
			input := &vo.HomeIssueInfoReq{}
			input.FollowerIds = []int64{1080, 1058}
			input.OrderType = &orderBy
			input.Status = &status
			page := 1
			size := 20

			resp, err := HomeIssues(cacheUserInfo.OrgId, cacheUserInfo.UserId, page, size, input)

			t.Log(json.ToJsonIgnoreError(resp.List))
			convey.So(resp, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
	}))
}

func TestIssueStatusTypeStat(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
				return
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIssue", func() {

				var projectId = int64(1204)

				rand.Seed(time.Now().Unix())
				intn := rand.Intn(10000)
				//projectObjectId
				typeId := int64(intn)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss", func() {

					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
						Id:                  processId,
						OrgId:               orgId,
						ProjectId:           projectId,
						ProjectObjectTypeId: typeId,
						ProcessId:           1120, //这个是一开始默认的不能写死
						Creator:             currentUserId,
						CreateTime:          time.Now(),
						Updator:             currentUserId,
						UpdateTime:          time.Now(),
						Version:             1,
					}

					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
						if err4 != nil {
							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						}
						return nil
					})

					convey.Convey("insertPriority.... ", func() {

						_, _ = CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createiteration input", func() {

							var iterationName = "alan迭代" + string(intn)
							planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, _ := time.Parse(consts.AppTimeFormat, "2020-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(cacheUserInfo.OrgId, cacheUserInfo.UserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("IssueStatusTypeStat", func() {

								id := modelsVoid.ID

								input := vo.IssueStatusTypeStatReq{
									ProjectID:   &projectId,
									IterationID: &id}

								resp, err := IssueStatusTypeStat(cacheUserInfo.OrgId, cacheUserInfo.UserId, &input)

								convey.So(resp, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})

					})
				})
			})
		})
	}))
}

func TestIssueStatusTypeStat2(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		input := vo.IssueStatusTypeStatReq{}

		resp, err := IssueStatusTypeStat(1003, 1007, &input)

		fmt.Println("返回的数据........", json.ToJsonIgnoreError(resp))

		convey.So(resp, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)

	}))
}

func TestIssueStatusTypeStatDetail(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIssue", func() {

				var projectId = int64(1204)

				rand.Seed(time.Now().Unix())
				intn := rand.Intn(10000)
				//projectObjectId
				typeId := int64(intn)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss", func() {

					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
						Id:                  processId,
						OrgId:               orgId,
						ProjectId:           projectId,
						ProjectObjectTypeId: typeId,
						ProcessId:           1120, //这个是一开始默认的不能写死
						Creator:             currentUserId,
						CreateTime:          time.Now(),
						Updator:             currentUserId,
						UpdateTime:          time.Now(),
						Version:             1,
					}

					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
						if err4 != nil {
							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						}
						return nil
					})

					convey.Convey("insertPriority.... ", func() {

						_, _ = CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createiteration input", func() {

							var iterationName = "alan迭代" + string(intn)
							planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, _ := time.Parse(consts.AppTimeFormat, "2020-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(cacheUserInfo.OrgId, cacheUserInfo.UserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("IssueStatusTypeDetailStat", func() {

								id := modelsVoid.ID

								input := vo.IssueStatusTypeStatReq{
									ProjectID:   &projectId,
									IterationID: &id}

								resp, err := IssueStatusTypeStatDetail(cacheUserInfo.OrgId, cacheUserInfo.UserId, &input)

								convey.So(resp, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})
		})
	}))
}

func TestGetIssueIds(t *testing.T) {

	convey.Convey("TestGetIssueIds", t, test.StartUp(func(ctx context.Context) {
		beforePlanEndTime := "2015-01-01 11:11:11"
		afterPlanEndTIme := "2020-01-01 11:11:11"

		resp, err := GetIssueRemindInfoList(projectvo.GetIssueRemindInfoListReqVo{
			Page: 1,
			Size: 10,
			Input: projectvo.GetIssueRemindInfoListReqData{
				BeforePlanEndTime: &beforePlanEndTime,
				AfterPlanEndTime:  &afterPlanEndTIme,
			},
		})
		t.Log(json.ToJsonIgnoreError(resp))
		assert.Equal(t, err, nil)
	}))

}
