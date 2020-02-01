package service

import (
	"context"
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/encrypt"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/processvo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/processfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/domain"
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

		fmt.Printf("redis的配置.....%v", config.RedisConfig{})

		fmt.Printf("数据库的的配置.....%v", config.RedisConfig{})

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息：" + cacheUserInfoJson)

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

			reqVo := projectvo.CreateProjectReqVo{
				Input:  input,
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
			}

			project := projectfacade.CreateProject(reqVo)

			if project.Failure() {
				log.Info("创建Issue中..........创建项目失败")
				return
			}

			log.Infof("创建issueUnittest中的 project %v", project)

			convey.Convey("Test CreateIssue", func() {

				projectId := project.Project.ID

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

						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIssue input", func() {

							var remark = "这是我创建的第一个issue"
							var issueObjectId = int64(1)

							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}

							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue, err := CreateIssue(reqVo)
								convey.So(issue, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})
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

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
		cacheUserInfoJson, err := json.ToJson(cacheUserInfo)

		if err == nil {

		}

		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))

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

					v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
						OrgID: orgId,
						Type:  consts.PriorityTypeIssue})

					convey.Convey("mock createIssue input", func() {

						var remark = "这是我创建的第一个issue"
						var issueObjectId = int64(1)
						children := make([]*vo.IssueChildren, 0)

						firstChild := &vo.IssueChildren{
							Title: "alan子issue",
							// 负责人
							OwnerID: currentUserId,
							// 优先级
							PriorityID: v.ID}

						children = append(children, firstChild)

						cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
						if err != nil {
							log.Error(err)
							return
						}

						input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
							IssueObjectID: &issueObjectId, Children: children}

						reqVo := projectvo.CreateIssueReqVo{
							CreateIssue: input,
							OrgId:       cacheUserInfo.OrgId,
							UserId:      cacheUserInfo.UserId,
						}

						convey.Convey("createIssue done", func() {

							issue, err := CreateIssue(reqVo)
							convey.So(issue, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)
						})
					})
				})
			})
		})
	}))
}

//更新issue
func TestUpdateIssue(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))

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

						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})
						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createiteration input", func() {

							var iterationName = "alan创建的迭代" + string(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2020-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(cacheUserInfo.OrgId, cacheUserInfo.UserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							var remark = "这是我创建的第一个issue"
							var issueObjectId = int64(1)

							children := make([]*vo.IssueChildren, 0)

							firstChild := &vo.IssueChildren{
								Title: "alan子issue",
								// 负责人
								OwnerID: currentUserId,
								// 优先级
								PriorityID: v.ID}

							children = append(children, firstChild)
							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}
							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId, Children: children}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue, err := CreateIssue(reqVo)
								convey.So(issue, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)

								convey.Convey("mock updateIssue req", func() {

									var hour = 10
									randomStr := string(intn)
									var updateRemark = "这是我更新的第一个issue" + randomStr

									participantIds := make([]int64, 1)
									participantIds[0] = 1065
									followerIds := make([]int64, 1)
									followerIds[0] = 1070
									updateFields := make([]string, 0)

									updateFields = append(updateFields, "remark")
									updateFields = append(updateFields, "title")
									updateFields = append(updateFields, "planEndTime")
									updateFields = append(updateFields, "planStartTime")
									updateFields = append(updateFields, "planWorkHour")
									updateFields = append(updateFields, "priorityId")
									updateFields = append(updateFields, "iterationId")
									updateFields = append(updateFields, "sourceId")
									updateFields = append(updateFields, "issueObjectTypeId")
									updateFields = append(updateFields, "ownerId")
									updateFields = append(updateFields, "followerIds")
									updateFields = append(updateFields, "participantIds")
									updateFields = append(updateFields, "hour")

									updateInput := vo.UpdateIssueReq{
										ID:                issue.ID,
										Title:             &issue.Title,
										OwnerID:           &issue.Owner,
										PriorityID:        &issue.PriorityID,
										PlanStartTime:     &issue.PlanStartTime,
										PlanEndTime:       &issue.PlanEndTime,
										PlanWorkHour:      &hour,
										Remark:            &updateRemark,
										IterationID:       &modelsVoid.ID,
										SourceID:          &issue.SourceID,
										IssueObjectTypeID: &issue.IssueObjectTypeID,
										ParticipantIds:    participantIds,
										FollowerIds:       followerIds,
										UpdateFields:      updateFields}

									req := projectvo.UpdateIssueReqVo{
										Input:  updateInput,
										UserId: cacheUserInfo.UserId,
										OrgId:  cacheUserInfo.OrgId,
									}

									convey.Convey(" updateIssue done", func() {

										updateIssueResp, err := UpdateIssue(req)
										convey.So(updateIssueResp, convey.ShouldNotBeNil)
										convey.So(err, convey.ShouldBeNil)
									})
								})
							})
						})
					})
				})
			})
		})
	}))
}

//删除迭代 有子任务不可以删除的
func TestDeleteIssueExistChildren(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))

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

						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})
						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIssue input", func() {

							var remark = "来删除的issue含字"
							var issueObjectId = int64(1)

							children := make([]*vo.IssueChildren, 0)

							firstChild := &vo.IssueChildren{
								Title: "alan子issue",
								// 负责人
								OwnerID: currentUserId,
								// 优先级
								PriorityID: v.ID}

							children = append(children, firstChild)

							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}

							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId, Children: children}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue, err := CreateIssue(reqVo)
								convey.So(issue, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)

								convey.Convey("mock deleteIssue req", func() {

									deleteInput := vo.DeleteIssueReq{
										ID: issue.ID}

									req := projectvo.DeleteIssueReqVo{
										Input:  deleteInput,
										OrgId:  cacheUserInfo.OrgId,
										UserId: cacheUserInfo.UserId,
									}

									convey.Convey(" DeleteIssue done", func() {

										deleteIssueResp, err := DeleteIssue(req)
										convey.So(deleteIssueResp, convey.ShouldBeNil)
										convey.So(err, convey.ShouldNotBeNil)

									})
								})

							})
						})
					})
				})
			})
		})
	}))
}

//删除迭代 有子任务不可以删除的
func TestDeleteIssueNotExistChildren(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))

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

						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})

						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIssue input", func() {

							var remark = "来删除的issue"
							var issueObjectId = int64(1)

							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}

							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue, err := CreateIssue(reqVo)
								convey.So(issue, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)

								convey.Convey("mock deleteIssue req", func() {

									deleteInput := vo.DeleteIssueReq{
										ID: issue.ID}

									req := projectvo.DeleteIssueReqVo{
										Input:  deleteInput,
										OrgId:  cacheUserInfo.OrgId,
										UserId: cacheUserInfo.UserId,
									}

									convey.Convey(" DeleteIssue done", func() {

										deleteIssueResp, err := DeleteIssue(req)
										convey.So(deleteIssueResp, convey.ShouldNotBeNil)
										convey.So(err, convey.ShouldBeNil)
									})
								})

							})
						})
					})
				})
			})
		})
	}))
}

//问题想详情
func TestIssueInfo(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))

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

						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})
						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIssue input", func() {

							var remark = "来删除的issue"
							var issueObjectId = int64(1)

							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}

							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue, err := CreateIssue(reqVo)
								convey.So(issue, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)

								convey.Convey("Test IssueInfo", t, func() {
									convey.Convey("IssueInfo mock req", func() {

										resp, err := IssueInfo(cacheUserInfo.OrgId, cacheUserInfo.UserId, issue.ID, "")
										convey.So(resp, convey.ShouldNotBeNil)

										convey.So(err, convey.ShouldBeNil)
									})
								})
							})
						})
					})
				})
			})
		})
	}))
}

func TestGetIssueRestInfos(t *testing.T) {

	convey.Convey("Test GetIssueRestInfos", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("GetIssueRestInfos mock req", func() {
			input := &vo.IssueRestInfoReq{}
			page := 1
			size := 20
			resp, err := GetIssueRestInfos(cacheUserInfo.OrgId, page, size, input)

			convey.So(resp, convey.ShouldNotBeNil)

			convey.So(err, convey.ShouldBeNil)
		})
	}))
}

func TestUpdateIssueStatus(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))

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
						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
							OrgID: orgId,
							Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIssue input", func() {

							var remark = "这是我创建的第一个issue"
							var issueObjectId = int64(1)

							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
							if err != nil {
								log.Error(err)
								return
							}

							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
								IssueObjectID: &issueObjectId}

							reqVo := projectvo.CreateIssueReqVo{
								CreateIssue: input,
								OrgId:       cacheUserInfo.OrgId,
								UserId:      cacheUserInfo.UserId,
							}

							convey.Convey("createIssue done", func() {

								issue, err := CreateIssue(reqVo)
								convey.So(issue, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)

								req := processvo.GetProcessStatusListByCategoryReqVo{
									OrgId:    orgId,
									Category: consts.ProcessStatusCategoryIssue}

								category := processfacade.GetProcessStatusListByCategory(req)

								bos := category.CacheProcessStatusBoList

								//processStatusList := &[]po.PpmPrsProcessStatus{}
								//err = mysql.SelectAllByCond(processStatus.TableName(), db.Cond{
								//	consts.TcOrgId:    orgId,
								//	consts.TcIsDelete: consts.AppIsNoDelete,
								//	consts.TcStatus:   consts.AppStatusEnable,
								//	consts.TcCategory: consts.ProcessStatusCategoryIssue}, processStatusList)
								//if err != nil {
								//	return
								//}

								var nextStatusId int64

								issueBo, err := domain.GetIssueBo(orgId, issue.ID)

								for _, status := range bos {
									if issueBo.Status != status.StatusId {
										nextStatusId = status.StatusId
										break
									}
								}

								convey.Convey("UpdateIssueStatus mock req", func() {

									input := vo.UpdateIssueStatusReq{
										ID:           issue.ID,
										NextStatusID: &nextStatusId}

									resp, err := UpdateIssueStatus(projectvo.UpdateIssueStatusReqVo{
										OrgId:  cacheUserInfo.OrgId,
										UserId: cacheUserInfo.UserId,
										Input:  input,
									})

									convey.So(resp, convey.ShouldNotBeNil)

									convey.So(err, convey.ShouldBeNil)
								})
							})
						})
					})
				})

			})
		})
	}))
}

func TestIssueReport(t *testing.T) {

	convey.Convey("Test IssueReport", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("IssueReport mock req", func() {

			resp, err := IssueReport(cacheUserInfo.OrgId, cacheUserInfo.UserId, consts.DailyReport)

			convey.So(resp, convey.ShouldNotBeNil)

			convey.So(err, convey.ShouldBeNil)

		})
	}))
}

func TestIssueReportDetail(t *testing.T) {
	convey.Convey("Test IssueReportDetail", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("IssueReportDetail mock req", func() {

			resp, _ := IssueReport(cacheUserInfo.OrgId, cacheUserInfo.UserId, consts.DailyReport)
			//str := resp.ShareID  因为每次状态要变
			str := "201"
			aesStr, _ := encrypt.AesEncrypt(str)

			resp, err2 := IssueReportDetail(aesStr)

			convey.So(resp, convey.ShouldNotBeNil)

			convey.So(err2, convey.ShouldBeNil)
		})
	}))
}

func TestGetIssueInfoList(t *testing.T) {
	convey.Convey("Test GetIssueInfoList", t, test.StartUp(func(ctx context.Context) {
		//fmt.Printf("数据库的的配置.....%s", json.ToJsonIgnoreError(config.GetMysqlConfig()))
		fmt.Println(GetIssueInfoList([]int64{1001, 1002}))
	}))

}

func TestUpdateIssueProjectObjectType(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		t.Log(UpdateIssueProjectObjectType(1003, 1007, vo.UpdateIssueProjectObjectTypeReq{ID:11603, ProjectObjectTypeID:1013}))
	}))
}