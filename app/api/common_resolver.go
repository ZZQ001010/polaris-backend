package api

import (
	"context"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/polaris-backend/common/core/buildinfo"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/commonvo"
	"github.com/galaxy-book/polaris-backend/facade/commonfacade"
	"github.com/galaxy-book/polaris-backend/facade/orgfacade"
)

func (r *queryResolver) GetBaseConfig(ctx context.Context) (*vo.BasicConfigResp, error) {
	BuildInfo := vo.BuildInfoDefine{
		GitCommitLog:   buildinfo.GitCommitLog,
		GitStatus:      buildinfo.GitStatus,
		BuildGoVersion: buildinfo.BuildGoVersion,
		BuildTime:      buildinfo.BuildTime,
	}
	result := vo.BasicConfigResp{
		RunMode:   config.GetApplication().RunMode,
		BuildInfo: &BuildInfo,
	}

	return &result, nil
}

func (r *queryResolver) AreaLinkageList(ctx context.Context, input vo.AreaLinkageListReq) (*vo.AreaLinkageListResp, error) {
	_, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := commonfacade.AreaLinkageList(commonvo.AreaLinkageListReqVo{
		Input: input,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.AreaLinkageListResp, nil

}

func (r *queryResolver) IndustryList(ctx context.Context) (*vo.IndustryListResp, error) {

	_, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := commonfacade.IndustryList()

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.IndustryList, nil
}
