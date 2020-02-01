package api

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
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestCreateIssue(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息：" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
			cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		}

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

			reqVo := projectvo.CreateProjectReqVo{
				Input:  input,
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
			}

			resp := projectfacade.CreateProject(reqVo)

			if resp.Failure() {
				log.Info("创建Issue中..........创建项目失败")
				return
			}

			log.Infof("创建issueUnittest中的 project %v", resp.Project)

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

						req := projectvo.CreatePriorityReqVo{
							vo.CreatePriorityReq{
								OrgID: orgId,
								Type:  consts.PriorityTypeIssue},
							cacheUserInfo.UserId,
						}

						priority := postGreeter.CreatePriority(req)

						if priority.Failure() {
							log.Info("创建Issue中..........创建优先级失败")
							return
						}

						convey.Convey("mock createIssue input", func() {

							var remark = "这是我创建的第一个issue"
							var issueObjectId = int64(1)

							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}

							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: priority.Void.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue := postGreeter.CreateIssue(reqVo)
								convey.So(issue.Failure(), convey.ShouldBeFalse)
							})
						})
					})
				})
			})
		})
	}))
}

func TestCreateIssueForTest(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息：" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1004), CorpId: "1", OrgId: 1002}
			//cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		}

		convey.Convey("Test CreateIssue", func() {

			for i := 0; i < 3000; i++ {

				var projectId = int64(1003)

				rand.Seed(time.Now().Unix())
				_ = rand.Intn(10000)
				//projectObjectId
				typeId := int64(1006)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
					Id:                  processId,
					OrgId:               orgId,
					ProjectId:           projectId,
					ProjectObjectTypeId: typeId,
					ProcessId:           5, //这个是一开始默认的不能写死
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

				req := projectvo.CreatePriorityReqVo{
					vo.CreatePriorityReq{
						OrgID: orgId,
						Type:  consts.PriorityTypeIssue},
					cacheUserInfo.UserId,
				}

				priority := postGreeter.CreatePriority(req)

				if priority.Failure() {
					log.Info("创建Issue中..........创建优先级失败")
					return
				}

				var remark = "这是我创建的第一个issue"
				var issueObjectId = int64(1001)

				input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: priority.Void.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
					IssueObjectID: &issueObjectId}

				reqVo := projectvo.CreateIssueReqVo{
					CreateIssue: input,
					OrgId:       cacheUserInfo.OrgId,
					UserId:      cacheUserInfo.UserId,
				}


					issue := postGreeter.CreateIssue(reqVo)
					fmt.Println("这个项目内容", json.ToJsonIgnoreError(issue))

			}
		})
	}))
}

func TestCreateChildIssue(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
			cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
		}

		cacheUserInfoJson, err := json.ToJson(cacheUserInfo)

		if err == nil {

		}

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

			reqVo := projectvo.CreateProjectReqVo{
				Input:  input,
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
			}

			resp := projectfacade.CreateProject(reqVo)

			if resp.Failure() {
				log.Info("创建Issue中..........创建项目失败")
				return
			}

			log.Infof("创建issueUnittest中的 project %v", resp.Project)

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

					req := projectvo.CreatePriorityReqVo{
						vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue},
						cacheUserInfo.UserId,
					}

					priority := postGreeter.CreatePriority(req)

					if priority.Failure() {
						log.Info("创建Issue中..........创建优先级失败")
						return
					}

					convey.Convey("mock createIssue input", func() {

						var remark = "这是我创建的第一个issue"
						var issueObjectId = int64(1)

						children := make([]*vo.IssueChildren, 0)

						firstChild := &vo.IssueChildren{
							Title: "alan子issue",
							// 负责人
							OwnerID: currentUserId,
							// 优先级
							PriorityID: priority.Void.ID}

						children = append(children, firstChild)

						cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
						if err != nil {
							log.Error(err)
							return
						}

						input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: priority.Void.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
							IssueObjectID: &issueObjectId, Children: children}

						reqVo := projectvo.CreateIssueReqVo{
							CreateIssue: input,
							OrgId:       cacheUserInfo.OrgId,
							UserId:      cacheUserInfo.UserId,
						}

						convey.Convey("createIssue done", func() {

							issue := postGreeter.CreateIssue(reqVo)
							convey.So(issue.Failure(), convey.ShouldBeFalse)
						})
					})
				})
			})
		})
	}))
}
