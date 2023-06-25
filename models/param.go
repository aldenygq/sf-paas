package models

//header
type ParamHeader struct {
	Emp   string `header:"UUAP-EMPLOYEE-NUMBER" binding:"required,min=1" label:"请求header(工号)"`
	Email string `header:"UUAP-EMAIL" binding:"omitempty,min=1" label:"请求header(邮箱)"`
	Cname string `header:"UUAP-USERCNAME" binding:"omitempty,min=1" label:"请求header(用户中文名)"`
	Uname string `header:"UUAP-USERNAME" binding:"omitempty,min=1" label:"请求header(英文名)"`
}
type CloudProSearchParam struct {
	KeyWord string `form:"key_word" json:"key_word" binding:"omitempty,min=1" label:"搜索关键字"`
}

//search param
type SearchParam struct {
	KeyWord string `form:"key_word" json:"key_word" binding:"omitempty,min=1" label:"搜索关键字"`
	KeyType string `form:"key_type" json:"key_type" binding:"omitempty,min=1" label:"搜索类型"`
}

//create object
type CreateParams struct {
	ParId         int64      `form:"par_id" json:"par_id" binding:"omitempty,gt=0" label:"父节点id"`
	Icon          string     `form:"icon" json:"icon" binding:"omitempty,min=1" label:"图标"`
	ObjectName    string     `form:"object_name" json:"object_name" binding:"required,min=1,max=10" label:"产品名称"`
	ObjectPrefix  string     `form:"object_prefix" json:"object_prefix" binding:"omitempty,min=1,max=10" label:"应用前缀"`
	ObjectType    string     `form:"object_type" json:"object_type" binding:"omitempty,min=1" label:"产品类型"`
	Router        string     `form:"router" json:"router" binding:"omitempty,min=1" label:"转发路由"`
	ObjectLevel   string     `form:"object_level" json:"object_level" binding:"required,min=1" label:"级别"`
	Managers      []string   `form:"managers" json:"managers" binding:"required,min=1" label:"管理员"`
	ObjectDesc    string     `form:"object_desc" json:"object_desc" binding:"required,min=1,max=20" label:"描述"`
	ObjectDoc     string     `form:"object_doc" json:"object_doc" binding:"omitempty,min=1" label:"说明"`
	IsJumpNewPage string     `form:"is_jump_new_page" json:"is_jump_new_page" binding:"omitempty,min=1" label:"是否跳转新页面"`
	IsChildren    string     `form:"is_children" json:"is_children" binding:"omitempty,min=1" label:"是否有子菜单"`
	Children      []*SubMenu `form:"children" json:"children" binding:"omitempty,dive" label:"子菜单"`
}

type SubMenu struct {
	ObjectName string `form:"object_name" json:"object_name" binding:"omitempty,min=1,max=10" label:"子菜单名称"`
	Router     string `form:"router" json:"router" binding:"omitempty,min=1" label:"子菜单转发路由"`
}

//create product
type EditParams struct {
	ParId         int64      `form:"par_id" json:"par_id" binding:"omitempty,gt=0" label:"父节点id"`
	Id            int64      `form:"id" json:"id" binding:"required,gt=0" label:"对象id"`
	Icon          string     `form:"icon" json:"icon" binding:"omitempty,min=1" label:"图标"`
	ObjectPrefix  string     `form:"object_prefix" json:"object_prefix" binding:"omitempty,min=1,max=10" label:"应用前缀"`
	ObjectLevel   string     `form:"object_level" json:"object_level" binding:"required,min=1" label:"级别"`
	Router        string     `form:"router" json:"router" binding:"omitempty,min=1" label:"转发路由"`
	ObjectName    string     `form:"object_name" json:"object_name" binding:"required,min=1" label:"对象名称"`
	ObjectType    string     `form:"object_type" json:"object_type" binding:"omitempty,min=1" label:"对象类型"`
	Managers      []string   `form:"managers" json:"managers" binding:"required,min=1" label:"管理员"`
	ObjectDesc    string     `form:"object_desc" json:"object_desc" binding:"required,min=1" label:"产品描述"`
	ObjectDoc     string     `form:"object_doc" json:"object_doc" binding:"omitempty,min=1" label:"说明"`
	IsJumpNewPage string     `form:"is_jump_new_page" json:"is_jump_new_page" binding:"omitempty,min=1" label:"是否跳转新页面"`
	IsChildren    string     `form:"is_children" json:"is_children" binding:"omitempty,min=1" label:"是否有子菜单"`
	Children      []*SubMenu `form:"children" json:"children" binding:"omitempty,dive" label:"子菜单"`
}

// pro get method
type ObjectParams struct {
	Id int64 `form:"id" json:"id" binding:"required,gt=0" label:"id"`
}

type ObjListParams struct {
	Id      int64  `form:"id" json:"id" binding:"omitempty,gt=0" label:"id"`
	ObjType string `form:"obj_type" json:"obj_type" binding:"required,min=1" label:"对象类型"`
	KeyWord string `form:"key_word" json:"key_word" binding:"omitempty,min=1" label:"搜索关键字"`
	KeyType string `form:"key_type" json:"key_type" binding:"omitempty,min=1" label:"搜索类型"`
}

type SupportParams struct {
	FbType      string `form:"fb_type" json:"fb_type" binding:"required,min=1" label:"反馈类型"`
	FbObj       string `form:"fb_obj" json:"fb_obj" binding:"required,min=1" label:"反馈对象"`
	FbDesc      string `form:"fb_desc" json:"fb_desc" binding:"required,min=1" label:"反馈详情"`
	IsAnonymous string `form:"is_anonymous" json:"is_anonymous" binding:"required,min=1" label:"是否匿名"`
}

type ClickOnParams struct {
	AppId int64 `form:"id" json:"id" binding:"required,gt=0" label:"应用id"`
}
type ParamAppTree struct {
	AppId int64 `form:"appid" json:"appid" binding:"required,gt=0" label:"应用id"`
}
type ParamFbObjs struct {
	FbType string `form:"fb_type" json:"fb_type" binding:"required,min=1" label:"反馈类型"`
}
type ParamMigrateObj struct {
	Id       int64 `form:"id" json:"id" binding:"required,gt=0" label:"应用id"`
	DstParId int64 `form:"dst_par_id" json:"dst_par_id" binding:"required,gt=0" label:"目标产品id"`
}

type ParamApp struct {
	AppName string `form:"app_name" json:"app_name" binding:"required,min=1" label:"应用名称"`
	NodeId int64 `form:"node_id" json:"node_id" binding:"required,gt=1" label:"节点id"`
}

type ParamOncallByType struct {
	OncallType string `form:"type" json:"type" binding:"required,min=1" label:"值班类型"`
}
type ParamAddOncall struct {
	User []string `form:"user" json:"user" binding:"required,min=1" label:"值班人员"`
	Type string `form:"type" json:"type" binding:"required,min=1" label:"类型" description:"DBA,SRE,PE"`
	Creator string `form:"cname" json:"cname" binding:"required,min=1" label:"姓名"`
}
