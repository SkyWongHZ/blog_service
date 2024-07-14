package errcode

var (
	ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
	ErrorCreateTagFail  = NewError(20010002, "创建标签失败")
	ErrorUpdateTagFail  = NewError(20010003, "更新标签失败")
	ErrorDeleteTagFail  = NewError(20010004, "删除标签失败")
	ErrorCountTagFail   = NewError(20010005, "统计标签失败")

	ErrorGetArticleFail    = NewError(20020001, "获取单个文章失败")
	ErrorGetArticlesFail   = NewError(20020002, "获取多个文章失败")
	ErrorCreateArticleFail = NewError(20020003, "创建文章失败")
	ErrorUpdateArticleFail = NewError(20020004, "更新文章失败")
	ErrorDeleteArticleFail = NewError(20020005, "删除文章失败")

	ErrorListUserFail     = NewError(20030001, "获取用户列表失败")
	ErrorRegisterUserFail = NewError(20030002, "创建用户列表失败")
	ErrorCountUserFail    = NewError(20030003, "统计用户总数失败")
	ErrorUpdateUserFail   = NewError(20030004, "更新用户失败")
	ErrorDeleteUserFail=NewError(20030005,"删除用户失败")
)

