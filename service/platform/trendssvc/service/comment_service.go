package service

import (
	"github.com/galaxy-book/common/core/util/copyer"
	"github.com/galaxy-book/polaris-backend/common/core/consts"
	"github.com/galaxy-book/polaris-backend/common/core/errs"
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/facade/idfacade"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/dao"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/po"
)

func CreateComment(commentBo bo.CommentBo) (int64, errs.SystemErrorInfo) {
	commentPo := &po.PpmTreComment{}
	_ = copyer.Copy(commentBo, commentPo)

	commentId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableComment)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}

	commentPo.Id = commentId

	err := dao.InsertComment(*commentPo)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return commentId, nil
}
