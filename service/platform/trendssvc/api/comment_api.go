package api

import (
	"github.com/galaxy-book/polaris-backend/common/model/vo"
	"github.com/galaxy-book/polaris-backend/common/model/vo/trendsvo"
	"github.com/galaxy-book/polaris-backend/service/platform/trendssvc/service"
)

func (PostGreeter) CreateComment(req trendsvo.CreateCommentReqVo) trendsvo.CreateCommentRespVo{
	id, err := service.CreateComment(req.CommentBo)
	return trendsvo.CreateCommentRespVo{CommentId: id, Err: vo.NewErr(err)}
}