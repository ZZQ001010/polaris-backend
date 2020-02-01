package schedule

import (
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
	"github.com/galaxy-book/polaris-backend/facade/projectfacade"
	"github.com/Jeffail/tunny"
	"time"
)

func StatisticProjectIssueBurnDownChart(pool tunny.Pool){
	statDate := times.GetYesterdayDate()

	//获取所有的组织
	orgBoListRespVo := orgfacade.GetOrgBoList()
	if orgBoListRespVo.Failure(){
		log.Error(orgBoListRespVo.Message)
		return
	}
	orgBos := orgBoListRespVo.OrganizationBoList

	//策略: 先查组织，后查项目，再根据组织和项目定位迭代
	//原因：防止一次性查太多导致内存崩溃，所以采用这种局部处理方式
	for _, orgBo := range orgBos{
		orgId := orgBo.Id
		log.Infof("项目任务燃尽图统计-开始统计迭代的组织信息：orgId: %v, orgName: %s", orgId, orgBo.Name)

		//获取所有项目
		projectBos, err := projectfacade.GetProjectBoListByProjectTypeLangCodeRelaxed(orgId, nil)
		if err != nil{
			log.Errorf("获取组织 %v 下的项目列表时出现问题，跳过该组织", orgId)
		}

		for _, projectBo := range projectBos{
			projectId := projectBo.Id
			log.Infof("项目任务燃尽图统计-开始统计迭代的项目信息：orgId: %v, orgName: %s， 项目信息: projectId: %v, projectName：%s", orgId, orgBo.Name, projectId, projectBo.Name)
			if projectBo.IsFiling == consts.AppIsFilling{
				log.Infof("项目任务燃尽图统计-项目已归档，不需要继续统计：orgId: %v, orgName: %s， 项目信息: projectId: %v, projectName：%s", orgId, orgBo.Name, projectId, projectBo.Name)
				continue
			}
			pd := projectBo
			go pool.ProcessTimed(func () error{
				statRespVo := projectfacade.AppendProjectDayStat(projectvo.AppendProjectDayStatReqVo{
					ProjectBo: pd,
					Date: statDate,
				})
				if statRespVo.Failure(){
					log.Errorf("项目 %v 燃尽图统计失败", projectBo.Id)
				}
				return nil
			}, time.Duration(5) * time.Minute)
		}
	}


}

