package service

import (
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/orgsvc/domain"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

func OrgInit(corpId string, permanentCode string, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.OrgInit(corpId, permanentCode, tx)
}

func OrgOwnerInit(orgId int64, owner int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	return domain.OrgOwnerInit(orgId, owner, tx)
}

func OrgSysConfigInit(tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
	return domain.OrgSysConfigInit(tx, orgId)
}

func TeamInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.TeamInit(orgId, tx)
}

func TeamOwnerInit(teamId int64, owner int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	return domain.TeamOwnerInit(teamId, owner, tx)
}

func TeamUserInit(orgId int64, teamId int64, userId int64, isRoot bool, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	return domain.TeamUserInit(orgId, teamId, userId, isRoot, tx)
}

func UserInitByOrg(userId string, corpId string, orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.UserInitByOrg(userId, corpId, orgId, tx)
}

func InitDepartment(orgId int64, corpId string, sourceChannel string, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	return domain.InitDepartment(orgId, corpId, sourceChannel, 0, 0, tx)
}

//初始化lark通用数据
func LarkInit(orgId, userId int64, sourceChannel string, sourcePlatform string) errs.SystemErrorInfo {
	departmentInfo, err := domain.GetTopDepartmentInfoList(orgId)
	if err != nil {
		log.Error("获取部门信息错误 " + strs.ObjectToString(err))
		return err
	}
	var departmentId int64
	for _, v := range departmentInfo {
		departmentId = v.Id
		break
	}

	//用户初始化
	zhangsanId, lisiId, err := domain.LarkUserInit(orgId, sourceChannel, sourcePlatform, departmentId)
	if err != nil {
		return err
	}
	log.Info("用户初始化成功")

	//项目初始化
	preCode := "SLXM"
	remark := "你可参考示例项目，快速熟悉北极星协作工具"
	start := types.NowTime()
	endTime, _ := time.Parse(consts.AppTimeFormat, "2099-12-12 12:00:00")
	end := types.Time(endTime)
	//获取普通项目
	var projectTypeId int64
	projectTypesResp := projectfacade.ProjectTypes(projectvo.ProjectTypesReqVo{
		OrgId: orgId,
	})
	if projectTypesResp.Failure() {
		return projectTypesResp.Error()
	}
	for _, v := range projectTypesResp.ProjectTypes {
		if v.LangCode == consts.ProjectTypeLangCodeNormalTask {
			projectTypeId = v.ID
		}
	}

	projectInfo := projectfacade.CreateProject(projectvo.CreateProjectReqVo{Input: vo.CreateProjectReq{
		Name:          "示例项目",
		PreCode:       &preCode,
		PublicStatus:  consts.PublicProject,
		Remark:        &remark,
		Owner:         zhangsanId,
		MemberIds:     []int64{zhangsanId, lisiId},
		ResourcePath:  "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/project/undraw_Projectpicture_update_jjgk.png",
		ResourceType:  consts.OssResource,
		PlanStartTime: &start,
		PlanEndTime:   &end,
		ProjectTypeID: &projectTypeId,
	},
		UserId: userId,
		OrgId:  orgId,
	})
	if projectInfo.Failure() {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, projectInfo.Error())
	}
	log.Info("项目初始化成功")

	objectType := 2
	projectObjectTypeErr1 := initProjectObjectType(orgId, userId, projectInfo.Project.ID, objectType, "需求")
	if projectObjectTypeErr1 != nil {
		log.Error(projectObjectTypeErr1)
		return projectObjectTypeErr1
	}
	projectObjectTypeErr2 := initProjectObjectType(orgId, userId, projectInfo.Project.ID, objectType, "设计")
	if projectObjectTypeErr2 != nil {
		log.Error(projectObjectTypeErr2)
		return projectObjectTypeErr2
	}

	//任务初始化
	issueInitErr := projectfacade.IssueLarkInit(projectvo.LarkIssueInitReqVo{
		ZhangsanId: zhangsanId,
		LisiId:     lisiId,
		OrgId:      orgId,
		ProjectId:  projectInfo.Project.ID,
		OperatorId: userId,
	})
	if issueInitErr.Failure() {
		return issueInitErr.Error()
	}
	log.Info("任务初始化成功")

	return nil
}

func initProjectObjectType(orgId, userId int64, projectId int64, objectType int, name string) errs.SystemErrorInfo {
	projectType1 := projectfacade.CreateProjectObjectType(projectvo.CreateProjectObjectTypeReqVo{
		Input: vo.CreateProjectObjectTypeReq{
			ProjectID:  projectId,
			Name:       name,
			ObjectType: objectType,
			BeforeID:   0,
		},
		OrgId:  orgId,
		UserId: userId,
	})
	if projectType1.Failure() {
		return errs.BuildSystemErrorInfo(errs.BaseDomainError, projectType1.Error())
	}
	log.Infof("项目对象类型-%s初始化成功", name)
	return nil
}
