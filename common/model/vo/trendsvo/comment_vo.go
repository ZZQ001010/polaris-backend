package trendsvo

import (
	"github.com/galaxy-book/polaris-backend/common/model/bo"
	"github.com/galaxy-book/polaris-backend/common/model/vo"
)

type CreateCommentReqVo struct {
	CommentBo bo.CommentBo `json:"commentBo"`
}

type CreateCommentRespVo struct {
	CommentId int64 `json:"data"`
	
	vo.Err
}