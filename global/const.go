package global

const (
	DefaultDatabaseUserTable  = "users"
	DefaultDatabaseGroupTable = "groups"
	DefaultDatabaseMediaTable = "media"
	//DefaultDatabaseLabelTable     = "label"
	DefaultDatabaseLabelTable = "label"

	DefaultAdminGroup    = "administrators"
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin@BUAA"

	EUserNotExist       = "用户不存在"
	EUserAlreadyExisted = "用户已存在"
	EUserForbidden      = "禁止的用户操作"

	EGroupNotExisted              = "用户组不存在"
	EGroupAlreadyExisted          = "用户组已存在"
	EGroupMemberAlreadyExisted    = "用户已存在于组内"
	EGroupMediaAlreadyInThisGroup = "媒体已存在于组内"
	EGroupMediaNotExistInGroup    = "组内无对应媒体"
	EGroupUserAlreadyExisted      = "用户已属于组"
	EGroupUserNotExist            = "用户不属于组"
	EGroupForbidden               = "禁止的组操作"

	EMediaNotExist       = "媒体不存在"
	EMediaAlreadyExisted = "媒体已存在"
	EMediaRawHashNull    = "无原始哈希"
	EMediaUnknownType    = "未知数据类型"
	EMediaForbidden      = "无权访问"

	ELabelDBLabedNotExist   = "标注数据不存在"
	ELabelDBLabelExisted    = "标注数据已存在"
	ELabelDBReviewDifferent = "数据正在被其他审阅人编辑"
	ELabelDBNoData          = "未提供数据"
	ELabelInvalidType       = "无效标注类别"

	EFFMPEGConvert = "media convert"

	MediaTypeUltrasonicImage = "us_image"
	MediaTypeUltrasonicVideo = "us_video"

	LabelTypeAuthor = "label_author"
	LabelTypeReview = "label_review"
	LabelTypeFinal  = "label_final"

	LabelProgressAuthoring = "author_ing"
	LabelProgressAuthored  = "author_finish"
	LabelProgressReviewing = "review_ing"
	LabelProgressReviewed  = "review_finish"

	DefGroupUngrouped     = "Ungrouped"
	DefGroupUngroupedName = "未分组"
)

const TimeFormat = "2006-01-02 15:04:05"
const SqlDBFile = "application/database/db.sqlite"
