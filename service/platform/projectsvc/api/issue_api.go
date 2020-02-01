package api

import (
	"github.com/galaxy-book/common/library/db/mysql"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/core/util/format"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/projectvo"
	"github.com/galaxy-book/polaris-backend/service/platform/projectsvc/service"
	"strings"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

func (PostGreeter) CreateIssue(reqVo projectvo.CreateIssueReqVo) projectvo.IssueRespVo {
	input := reqVo.CreateIssue

	//校验标题
	title := strings.TrimSpace(input.Title)
	//checkTitleLenErr := util.CheckIssueTitleLen(title)
	//if checkTitleLenErr != nil{
	//	log.Error(checkTitleLenErr)
	//	return projectvo.IssueRespVo{Err: vo.NewErr(checkTitleLenErr), Issue: nil}
	//}
	isTitleRight := format.VerifyIssueNameFormat(title)
	if !isTitleRight {
		log.Error(errs.IssueTitleError)
		return projectvo.IssueRespVo{Err: vo.NewErr(errs.IssueTitleError), Issue: nil}
	}
	input.Title = title

	//校验描述
	if input.Remark != nil {
		remark := *input.Remark
		//checkRemarkLenErr := util.CheckIssueRemarkLen(remark)
		//if checkRemarkLenErr != nil{
		//	log.Error(checkRemarkLenErr)
		//	return projectvo.IssueRespVo{Err: vo.NewErr(checkRemarkLenErr), Issue: nil}
		//}
		isRemarkRight := format.VerifyIssueRemarkFormat(remark)
		if !isRemarkRight {
			log.Error(errs.IssueRemarkLenError)
			return projectvo.IssueRespVo{Err: vo.NewErr(errs.IssueRemarkLenError), Issue: nil}
		}
		input.Remark = &remark
	}

	//截止时间
	if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() && input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {
		if time.Time(*input.PlanEndTime).Before(time.Time(*input.PlanStartTime)) {
			return projectvo.IssueRespVo{Err: vo.NewErr(errs.PlanEndTimeInvalidError), Issue: nil}
		}
	}

	reqVo.CreateIssue = input
	res, err := service.CreateIssue(reqVo)
	return projectvo.IssueRespVo{Err: vo.NewErr(err), Issue: res}
}

func (PostGreeter) UpdateIssue(reqVo projectvo.UpdateIssueReqVo) projectvo.UpdateIssueRespVo {
	input := reqVo.Input

	//校验标题
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		//checkTitleLenErr := util.CheckIssueTitleLen(title)
		//if checkTitleLenErr != nil{
		//	log.Error(checkTitleLenErr)
		//	return projectvo.UpdateIssueRespVo{Err: vo.NewErr(checkTitleLenErr), UpdateIssue: nil}
		//}
		isTitleRight := format.VerifyIssueNameFormat(title)
		if !isTitleRight {
			log.Error(errs.IssueTitleError)
			return projectvo.UpdateIssueRespVo{Err: vo.NewErr(errs.IssueTitleError), UpdateIssue: nil}
		}
		input.Title = &title
	}

	//校验描述
	if input.Remark != nil {
		remark := *input.Remark
		//checkRemarkLenErr := util.CheckIssueRemarkLen(remark)
		//if checkRemarkLenErr != nil{
		//	log.Error(checkRemarkLenErr)
		//	return projectvo.UpdateIssueRespVo{Err: vo.NewErr(checkRemarkLenErr), UpdateIssue: nil}
		//}
		isRemarkRight := format.VerifyIssueRemarkFormat(remark)
		if !isRemarkRight {
			log.Error(errs.IssueRemarkLenError)
			return projectvo.UpdateIssueRespVo{Err: vo.NewErr(errs.IssueRemarkLenError), UpdateIssue: nil}
		}
		input.Remark = &remark
	}

	//截止时间
	if input.PlanStartTime != nil && input.PlanStartTime.IsNotNull() && input.PlanEndTime != nil && input.PlanEndTime.IsNotNull() {
		if time.Time(*input.PlanEndTime).Before(time.Time(*input.PlanStartTime)) {
			return projectvo.UpdateIssueRespVo{Err: vo.NewErr(errs.PlanEndTimeInvalidError), UpdateIssue: nil}
		}
	}

	reqVo.Input = input
	res, err := service.UpdateIssue(reqVo)
	return projectvo.UpdateIssueRespVo{Err: vo.NewErr(err), UpdateIssue: res}
}

func (PostGreeter) DeleteIssue(reqVo projectvo.DeleteIssueReqVo) projectvo.IssueRespVo {
	res, err := service.DeleteIssue(reqVo)
	return projectvo.IssueRespVo{Err: vo.NewErr(err), Issue: res}
}

func (GetGreeter) IssueInfo(reqVo projectvo.IssueInfoReqVo) projectvo.IssueInfoRespVo {
	res, err := service.IssueInfo(reqVo.OrgId, reqVo.UserId, reqVo.IssueID, reqVo.SourceChannel)
	return projectvo.IssueInfoRespVo{Err: vo.NewErr(err), IssueInfo: res}
}

func (PostGreeter) GetIssueRestInfos(reqVo projectvo.GetIssueRestInfosReqVo) projectvo.GetIssueRestInfosRespVo {
	res, err := service.GetIssueRestInfos(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.GetIssueRestInfosRespVo{Err: vo.NewErr(err), GetIssueRestInfos: res}
}

func (PostGreeter) UpdateIssueStatus(reqVo projectvo.UpdateIssueStatusReqVo) projectvo.IssueRespVo {
	res, err := service.UpdateIssueStatus(reqVo)
	return projectvo.IssueRespVo{Err: vo.NewErr(err), Issue: res}
}

func (PostGreeter) UpdateIssueProjectObjectType(reqVo projectvo.UpdateIssueProjectObjectTypeReqVo) vo.CommonRespVo {
	res, err := service.UpdateIssueProjectObjectType(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) IssueLarkInit(reqVo projectvo.LarkIssueInitReqVo) vo.VoidErr {
	respVo := vo.VoidErr{}
	_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := service.LarkIssueInit(reqVo.OrgId, reqVo.ZhangsanId, reqVo.LisiId, reqVo.ProjectId, reqVo.OperatorId, tx)
		respVo.Err = vo.NewErr(err)
		return err
	})

	return respVo
}

func (PostGreeter) GetIssueInfoList(reqVo projectvo.IssueInfoListReqVo) projectvo.IssueInfoListRespVo {
	res, err := service.GetIssueInfoList(reqVo.IssueIds)
	return projectvo.IssueInfoListRespVo{
		Err:        vo.NewErr(err),
		IssueInfos: res,
	}
}

func (PostGreeter) UpdateIssueSort(reqVo projectvo.UpdateIssueSortReqVo) vo.CommonRespVo {
	res, err := service.UpdateIssueSort(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

