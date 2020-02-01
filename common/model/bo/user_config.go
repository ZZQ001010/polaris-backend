package bo

type UserConfigBo struct {
	ID                              int64 `json:"id"`
	DailyReportMessageStatus        int   `json:"dailyReportMessageStatus"`
	OwnerRangeStatus                int   `json:"ownerRangeStatus"`
	ParticipantRangeStatus          int   `json:"participantRangeStatus"`
	AttentionRangeStatus            int   `json:"attentionRangeStatus"`
	CreateRangeStatus               int   `json:"createRangeStatus"`
	RemindMessageStatus             int   `json:"remindMessageStatus"`
	CommentAtMessageStatus          int   `json:"commentAtMessageStatus"`
	ModifyMessageStatus             int   `json:"modifyMessageStatus"`
	RelationMessageStatus           int   `json:"relationMessageStatus"`
	DailyProjectReportMessageStatus int   `json:"dailyProjectReportMessageStatus"`
	DefaultProjectId				int64 `json:"defaultProjectId"`
	DefaultProjectObjectTypeId					int64		`db:"default_project_object_type_id,omitempty" json:"defaultProjectObjectTypeId"`
	PcNoticeOpenStatus					int		`db:"pc_notice_open_status,omitempty" json:"pcNoticeOpenStatus"`
	PcIssueRemindMessageStatus					int		`db:"pc_issue_remind_message_status,omitempty" json:"pcIssueRemindMessageStatus"`
	PcOrgMessageStatus					int		`db:"pc_org_message_status,omitempty" json:"pcOrgMessageStatus"`
	PcProjectMessageStatus					int		`db:"pc_project_message_status,omitempty" json:"pcProjectMessageStatus"`
	PcCommentAtMessageStatus					int		`db:"pc_comment_at_message_status,omitempty" json:"pcCommentAtMessageStatus"`
}
