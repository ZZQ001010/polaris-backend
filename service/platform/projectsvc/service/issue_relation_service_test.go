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
	"github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestCreateIssueComment(t *testing.T) {

	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		currentUserId := cacheUserInfo.UserId
		orgId := cacheUserInfo.OrgId

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
				OrgId: cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input: input,
			})

			projectId := project.ID
			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("createIssue insert Project proecss", func() {

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				typeId := int64(intn)

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

						convey.Convey("createIssue insert Project proecss", func() {

							processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

							projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
								Id:                  processId,
								OrgId:               orgId,
								ProjectId:           projectId,
								ProjectObjectTypeId: 1117,
								ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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

							convey.Convey("mock createIteration input", func() {

								var iterationName = "alan创建的迭代" + string(intn)
								planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
								planEndTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

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
								//先删除缓存里面的
								cache.Del(fmt.Sprintf("polaris:org_%d:priority_list", orgId))
								cache.Del(fmt.Sprintf("polaris:org_%d:process_list", orgId))

								children := make([]*vo.IssueChildren, 0)

								firstChild := &vo.IssueChildren{
									Title: "alan子issue",
									// 负责人
									OwnerID: currentUserId,
									// 优先级
									PriorityID: v.ID}

								children = append(children, firstChild)

								input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
									IssueObjectID: &issueObjectId, Children: children}

								convey.Convey("createIssue done", func() {

									reqVo := projectvo.CreateIssueReqVo{
										CreateIssue:input,
										UserId:cacheUserInfo.UserId,
										OrgId:cacheUserInfo.OrgId,
										SourceChannel:consts.AppSourceChannelDingTalk,
									}



									issue, err := CreateIssue(reqVo)
									convey.So(issue, convey.ShouldNotBeNil)
									convey.So(err, convey.ShouldBeNil)

									convey.Convey("CreateIssueComment", func() {

										input := vo.CreateIssueCommentReq{
											IssueID: issue.ID,
											// 评论信息
											Comment: "我喜欢乱评论"}

										resp, err := CreateIssueComment(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)

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

func TestCreateIssueResource(t *testing.T) {

	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		currentUserId := cacheUserInfo.UserId
		orgId := cacheUserInfo.OrgId

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
				OrgId: cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input: input,
			})

			projectId := project.ID
			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("createIssue insert Project proecss", func() {

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				typeId := int64(intn)

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

						convey.Convey("createIssue insert Project proecss", func() {

							processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

							projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
								Id:                  processId,
								OrgId:               orgId,
								ProjectId:           projectId,
								ProjectObjectTypeId: 1117,
								ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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

							convey.Convey("mock createIteration input", func() {

								var iterationName = "alan创建的迭代" + string(intn)
								planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
								planEndTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

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
								//先删除缓存里面的
								cache.Del(fmt.Sprintf("polaris:org_%d:priority_list", orgId))
								cache.Del(fmt.Sprintf("polaris:org_%d:process_list", orgId))

								children := make([]*vo.IssueChildren, 0)

								firstChild := &vo.IssueChildren{
									Title: "alan子issue",
									// 负责人
									OwnerID: currentUserId,
									// 优先级
									PriorityID: v.ID}

								children = append(children, firstChild)

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

									convey.Convey("CreateIssueResource", func() {

										md := "qweqweqwe"
										bucketName := "bucket"

										input := vo.CreateIssueResourceReq{
											IssueID:      issue.ID,
											ResourcePath: "suibiandeluj",
											ResourceSize: 1024,
											FileName:     "suibian",
											FileSuffix:   ".txt",
											Md5:          &md,
											BucketName:   &bucketName}

										resp, err := CreateIssueResource(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)

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

func TestCreateIssueRelationIssue(t *testing.T) {

	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		currentUserId := cacheUserInfo.UserId
		orgId := cacheUserInfo.OrgId

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
				OrgId: cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input: input,
			})

			projectId := project.ID
			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("createIssue insert Project proecss", func() {

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				typeId := int64(intn)

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

						convey.Convey("createIssue insert Project proecss", func() {

							processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

							projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
								Id:                  processId,
								OrgId:               orgId,
								ProjectId:           projectId,
								ProjectObjectTypeId: 1117,
								ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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

							convey.Convey("mock createIteration input", func() {

								var iterationName = "alan创建的迭代" + string(intn)
								planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
								planEndTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

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
								//先删除缓存里面的
								cache.Del(fmt.Sprintf("polaris:org_%d:priority_list", orgId))
								cache.Del(fmt.Sprintf("polaris:org_%d:process_list", orgId))

								children := make([]*vo.IssueChildren, 0)

								firstChild := &vo.IssueChildren{
									Title: "alan子issue",
									// 负责人
									OwnerID: currentUserId,
									// 优先级
									PriorityID: v.ID}

								children = append(children, firstChild)

								input := vo.CreateIssueReq{ProjectID: projectId, Title: fmt.Sprintf("%s%d", "alan创建的issue", intn), PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
									IssueObjectID: &issueObjectId, Children: children}

								input2 := vo.CreateIssueReq{ProjectID: projectId, Title: fmt.Sprintf("%s%d", "alan的issue2", intn), PriorityID: v.ID, TypeID: &typeId, OwnerID: currentUserId, Remark: &remark,
									IssueObjectID: &issueObjectId, Children: children}

								reqVo := projectvo.CreateIssueReqVo{
									CreateIssue: input,
									OrgId:       cacheUserInfo.OrgId,
									UserId:      cacheUserInfo.UserId,
								}

								reqVo2 := projectvo.CreateIssueReqVo{
									CreateIssue: input2,
									OrgId:       cacheUserInfo.OrgId,
									UserId:      cacheUserInfo.UserId,
								}

								convey.Convey("createIssue done", func() {

									issue, err := CreateIssue(reqVo)
									convey.So(issue, convey.ShouldNotBeNil)
									convey.So(err, convey.ShouldBeNil)

									issue2, err := CreateIssue(reqVo2)

									convey.Convey("CreateIssueComment", func() {

										addRelateIssueIds := make([]int64, 0)

										addRelateIssueIds = append(addRelateIssueIds, issue2.ID)

										input := vo.UpdateIssueAndIssueRelateReq{
											IssueID:           issue.ID,
											AddRelateIssueIds: addRelateIssueIds}

										resp, err := CreateIssueRelationIssue(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)
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

func TestDeleteIssueResource(t *testing.T) {

	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		currentUserId := cacheUserInfo.UserId
		orgId := cacheUserInfo.OrgId

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
				OrgId: cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input: input,
			})

			projectId := project.ID
			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("createIssue insert Project proecss", func() {

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				typeId := int64(intn)

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

						convey.Convey("createIssue insert Project proecss", func() {

							processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

							projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
								Id:                  processId,
								OrgId:               orgId,
								ProjectId:           projectId,
								ProjectObjectTypeId: 1117,
								ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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

							convey.Convey("mock createIteration input", func() {

								var iterationName = "alan创建的迭代" + string(intn)
								planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
								planEndTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

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
								//先删除缓存里面的
								cache.Del(fmt.Sprintf("polaris:org_%d:priority_list", orgId))
								cache.Del(fmt.Sprintf("polaris:org_%d:process_list", orgId))

								children := make([]*vo.IssueChildren, 0)

								firstChild := &vo.IssueChildren{
									Title: "alan子issue",
									// 负责人
									OwnerID: currentUserId,
									// 优先级
									PriorityID: v.ID}

								children = append(children, firstChild)

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

									convey.Convey("CreateIssueResource", func() {

										md := "qweqweqwe"
										bucketName := "bucket"

										input := vo.CreateIssueResourceReq{
											IssueID:      issue.ID,
											ResourcePath: "suibiandeluj",
											ResourceSize: 1024,
											FileName:     "suibian",
											FileSuffix:   ".txt",
											Md5:          &md,
											BucketName:   &bucketName}

										resp, err := CreateIssueResource(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)

										convey.So(resp, convey.ShouldNotBeNil)
										convey.So(err, convey.ShouldBeNil)

										convey.Convey("DeleteIssueResource", func() {

											deleteResourceIds := make([]int64, 0)

											deleteResourceIds = append(deleteResourceIds, resp.ID)

											input := vo.DeleteIssueResourceReq{
												IssueID:            issue.ID,
												DeletedResourceIds: deleteResourceIds}

											deleteResp, err := DeleteIssueResource(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)

											convey.So(deleteResp, convey.ShouldNotBeNil)
											convey.So(err, convey.ShouldBeNil)
										})

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

func TestRelatedIssueList(t *testing.T) {

	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		currentUserId := cacheUserInfo.UserId
		orgId := cacheUserInfo.OrgId

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
				OrgId: cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input: input,
			})

			projectId := project.ID
			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("createIssue insert Project proecss", func() {

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				typeId := int64(intn)

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

						convey.Convey("createIssue insert Project proecss", func() {

							processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

							projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
								Id:                  processId,
								OrgId:               orgId,
								ProjectId:           projectId,
								ProjectObjectTypeId: 1117,
								ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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

							convey.Convey("mock createIteration input", func() {

								var iterationName = "alan创建的迭代" + string(intn)
								planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
								planEndTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

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
								//先删除缓存里面的
								cache.Del(fmt.Sprintf("polaris:org_%d:priority_list", orgId))
								cache.Del(fmt.Sprintf("polaris:org_%d:process_list", orgId))

								children := make([]*vo.IssueChildren, 0)

								firstChild := &vo.IssueChildren{
									Title: "alan子issue",
									// 负责人
									OwnerID: currentUserId,
									// 优先级
									PriorityID: v.ID}

								children = append(children, firstChild)

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

									convey.Convey("CreateIssueResource", func() {

										md := "qweqweqwe"
										bucketName := "bucket"

										input := vo.CreateIssueResourceReq{
											IssueID:      issue.ID,
											ResourcePath: "suibiandeluj",
											ResourceSize: 1024,
											FileName:     "suibian",
											FileSuffix:   ".txt",
											Md5:          &md,
											BucketName:   &bucketName}

										resp, err := CreateIssueResource(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)

										convey.So(resp, convey.ShouldNotBeNil)
										convey.So(err, convey.ShouldBeNil)

										convey.Convey("RelatedIssueList", func() {

											input := vo.RelatedIssueListReq{IssueID: issue.ID}
											resp2, err := RelatedIssueList(cacheUserInfo.OrgId, input)

											convey.So(resp2, convey.ShouldNotBeNil)
											convey.So(err, convey.ShouldBeNil)
										})

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

func TestIssueResources(t *testing.T) {

	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		currentUserId := cacheUserInfo.UserId
		orgId := cacheUserInfo.OrgId

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
				OrgId: cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input: input,
			})

			projectId := project.ID
			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("createIssue insert Project proecss", func() {

				processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

				typeId := int64(intn)

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

						convey.Convey("createIssue insert Project proecss", func() {

							processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())

							projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
								Id:                  processId,
								OrgId:               orgId,
								ProjectId:           projectId,
								ProjectObjectTypeId: 1117,
								ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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

							convey.Convey("mock createIteration input", func() {

								var iterationName = "alan创建的迭代" + string(intn)
								planStartTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
								planEndTime, _ := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

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
								//先删除缓存里面的
								cache.Del(fmt.Sprintf("polaris:org_%d:priority_list", orgId))
								cache.Del(fmt.Sprintf("polaris:org_%d:process_list", orgId))

								children := make([]*vo.IssueChildren, 0)

								firstChild := &vo.IssueChildren{
									Title: "alan子issue",
									// 负责人
									OwnerID: currentUserId,
									// 优先级
									PriorityID: v.ID}

								children = append(children, firstChild)

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

									convey.Convey("CreateIssueResource", func() {

										md := "qweqweqwe"
										bucketName := "bucket"

										input := vo.CreateIssueResourceReq{
											IssueID:      issue.ID,
											ResourcePath: "suibiandeluj",
											ResourceSize: 1024,
											FileName:     "suibian",
											FileSuffix:   ".txt",
											Md5:          &md,
											BucketName:   &bucketName}

										resp, err := CreateIssueResource(cacheUserInfo.OrgId, cacheUserInfo.UserId, input)

										convey.So(resp, convey.ShouldNotBeNil)
										convey.So(err, convey.ShouldBeNil)

										convey.Convey("IssueResources", func() {

											input := vo.GetIssueResourcesReq{IssueID: issue.ID}

											resp2, err := IssueResources(cacheUserInfo.OrgId, 1, 1, &input)

											convey.So(resp2, convey.ShouldNotBeNil)
											convey.So(err, convey.ShouldBeNil)
										})

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

func TestCreateIssueComment2(t *testing.T) {
	convey.Convey("Test UpdateIssue", t, test.StartUp(func(ctx context.Context) {
		t.Log(CreateIssueComment(1070, 1131, vo.CreateIssueCommentReq{IssueID:5647, Comment:"aaa", MentionedUserIds:[]int64{}}))
	}))

}
