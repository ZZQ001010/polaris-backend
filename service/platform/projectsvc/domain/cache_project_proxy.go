package domain

import (
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/slice"
	"github.com/galaxy-book/common/library/cache"
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	sconsts "github.com/galaxy-book/polaris-backend/service/platform/projectsvc/consts"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProjectAuthBoBatch(orgId int64, projectIds []int64) ([]bo.ProjectAuthBo, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(projectIds))
	for i, projectId := range projectIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseProjectInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
		})
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0{
		list, err := cache.MGet(keys...)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		resultList = list
	}
	projectAuthBos := make([]bo.ProjectAuthBo, 0)
	validProjectIds := make([]int64, 0)
	for _, projectInfoJson := range resultList {
		projectCacheInfo := &bo.ProjectAuthBo{}

		err := json.FromJson(projectInfoJson, projectCacheInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		projectAuthBos = append(projectAuthBos, *projectCacheInfo)
		validProjectIds = append(validProjectIds, projectCacheInfo.Id)
	}
	//找不存在的
	if len(projectIds) != len(validProjectIds) {
		for _, projectId := range projectIds {
			exist, _ := slice.Contain(validProjectIds, projectId)
			if !exist {
				projectAuthBo, err := LoadProjectAuthBo(orgId, projectId)
				if err != nil {
					log.Error(err)
					continue
				}
				projectAuthBos = append(projectAuthBos, *projectAuthBo)
			}
		}
	}
	return projectAuthBos, nil
}

func LoadProjectAuthBo(orgId int64, projectId int64) (*bo.ProjectAuthBo, errs.SystemErrorInfo) {
	key, projectInfoJson, err5 := getCacheProjectInfo(orgId, projectId)

	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	if projectInfoJson != "" {
		projectCacheInfo := &bo.ProjectAuthBo{}
		err := json.FromJson(projectInfoJson, projectCacheInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return projectCacheInfo, nil
	} else {
		project := &po.PpmProProject{}
		err := mysql.SelectOneByCond(project.TableName(), db.Cond{
			consts.TcId:       projectId,
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, project)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
		}

		projectMemberInfos := &[]po.PpmProProjectRelation{}
		err = mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
			consts.TcProjectId:    projectId,
			consts.TcOrgId:        orgId,
			consts.TcStatus:       consts.AppStatusEnable,
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcRelationType: db.In([]int{consts.IssueRelationTypeFollower, consts.IssueRelationTypeParticipant}),
		}, projectMemberInfos)

		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		participants := &[]int64{}
		followers := &[]int64{}

		//for _, projectMemberInfo := range *projectMemberInfos{
		//	if projectMemberInfo.RelationType == consts.IssueRelationTypeFollower{
		//		*followers = append(*followers, projectMemberInfo.RelationId)
		//	}else if projectMemberInfo.RelationType == consts.IssueRelationTypeParticipant{
		//		*participants = append(*participants, projectMemberInfo.RelationId)
		//	}
		//}
		dealProjectMemberInfo(projectMemberInfos, followers, participants)

		projectCacheInfo := &bo.ProjectAuthBo{
			Id:           project.Id,
			Name:         project.Name,
			Creator:      project.Creator,
			Owner:        project.Owner,
			Participants: *participants,
			Followers:    *followers,
			Status:       project.Status,
			IsFilling:    project.IsFiling,
			PublicStatus: project.PublicStatus,
			ProjectType:  project.ProjectTypeId,
		}

		projectInfoJson, err = json.ToJson(projectCacheInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, projectInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return projectCacheInfo, nil
	}
}

func dealProjectMemberInfo(projectMemberInfos *[]po.PpmProProjectRelation, followers, participants *[]int64) {
	for _, projectMemberInfo := range *projectMemberInfos {
		if projectMemberInfo.RelationType == consts.IssueRelationTypeFollower {
			*followers = append(*followers, projectMemberInfo.RelationId)
		} else if projectMemberInfo.RelationType == consts.IssueRelationTypeParticipant {
			*participants = append(*participants, projectMemberInfo.RelationId)
		}
	}
}

func getCacheProjectInfo(orgId, projectId int64) (string, string, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseProjectInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return "", "", err5
	}

	projectInfoJson, err := cache.Get(key)
	if err != nil {
		return key, "", errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return key, projectInfoJson, nil
}

func RefreshProjectAuthBo(orgId int64, projectId int64) error {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseProjectInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	mod, err := cache.Del(key)
	if err != nil {
		log.Errorf("RefreshProjectInfo 发生异常 %v", err)
		return err
	}
	if mod == 0 {
		log.Error("RefreshProjectInfo 刷新失败，key不存在")
	}

	//_, err = LoadProjectAuthBo(orgId, projectId)
	//if err != nil {
	//	log.Errorf("LoadProjectInfo 失败，err %v", err)
	//}
	return nil
}

func GetProjectCalendarInfo(orgId, projectId int64) (*bo.CacheProjectCalendarInfoBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectCalendarInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	infoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	infoBo := &bo.CacheProjectCalendarInfoBo{}
	if infoJson != "" {
		err = json.FromJson(infoJson, infoBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return infoBo, nil
	} else {
		//查看项目是否支持导出日历
		detail, err := GetProjectDetailByProjectIdBo(projectId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		//查看项目日历对应的id
		relation, err := GetProjectRelationByType(projectId, []int64{consts.IssueRelationTypeCalendar})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var calendarId string
		for _, v := range *relation {
			calendarId = v.RelationCode
		}
		infoBo.IsSyncOutCalendar = detail.IsSyncOutCalendar
		infoBo.CalendarId = calendarId

		setErr := cache.SetEx(key, json.ToJsonIgnoreError(infoBo), consts.GetCacheBaseExpire())
		if setErr != nil {
			log.Error(setErr)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		return infoBo, nil
	}
}

func DeleteProjectCalendarInfo(orgId, projectId int64) errs.SystemErrorInfo {
	key, err := util.ParseCacheKey(sconsts.CacheProjectCalendarInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	_, err1 := cache.Del(key)

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	return nil
}
