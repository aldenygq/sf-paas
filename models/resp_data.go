package models

type Tree struct {
	NodeValue string  `json:"value,omitempty"`
	NodeLabel string  `json:"label,omitempty"`
	Type      string  `json:"type, omitempty"`
	NodeId    int64   `json:"node_id,omitempty"`
	NextNodes []*Tree `json:"children, omitempy"`
}

type UserCenter struct {
	BaseInfo []*Base `json:"menus,omitempty"`
}

type AppInfor struct {
	Id           int64  `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id,omitempty"`
	ObjectName   string `gorm:"column:object_name;type:varchar(64)" json:"app,omitempty"`
	Router       string `gorm:"column:router;type:varchar(1024)" json:"router,omitempty"`
	ObjectDoc    string `gorm:"object_doc;type:text"  json:"doc,omitempty"`
	ExternalJump string ` json:"external_jump,omitempty"`
}

type ClickInfo struct {
	Id       string `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id"`
	Emp      string `gorm:"column:emp;type:varchar(64)" json:"emp"`
	ClickNum int64  `gorm:"column:click_num;type:int(10)" json:"click_num"`
	ObjectId int64  `gorm:"column:object_id;type:int(10)" json:"object_id"`
}

type ProAndAppList struct {
	ObjectInfo
	Children []*ObjInfo `json:"children"`
}
type ProAndApp struct {
	ObjectInfo
	Child *ObjectInfo `json:"child"`
}
type CreateInfo struct {
	ObjectTypes []*ProductType `json:"object_types,omitempty"`
	Managers    []string       `json:"managers,omitempty"`
	Products    []*ObjectInfo  `json:"products,omitempty"`
}
type ProductType struct {
	CnName string `json:"cn_name,omitempty"`
	EnName string `json:"en_name,omitempty"`
}
type MigrateInfo struct {
	Pro         string        `json:"product,omitempty"`
	App         string        `json:"app,omitempty"`
	AppId       int64         `json:"app_id",omitempty`
	DstProducts []*ObjectInfo `json:"products,omitempty"`
}
type EditInfoResp struct {
	ObjectInfo
	NewManagers []string `json:"new_managers,omitempty"`
}

type ObjList struct {
	ObjInfo
	IsOpreate string `json:"is_opreate"`
}

type AppInfo struct {
	BaseInfo []*Base    `json:"base_info"`
	Menus    []*ObjList `json:"menus"`
	Path     string     `json:"path"`
}

type Base struct {
	Key          string `json:"key,omitempty"`
	Value        string `json:"value,omitempty"`
	ExternalJump string `json:"external_jump,omitempty"`
}

type MenuInfo struct {
	BaseInfo []*Base `json:"base_info"`
	Children []*Base `json:"children"`
	Path     string  `json:"path"`
}

type SupportMenu struct {
	BaseInfo []*Base `json:"menus,omitempty"`
}
type UserInfo struct {
	CnName    string `json:"cn_name,omitempty"`
	Emp       string `json:"emp,omitempty"`
	Email     string `json:"email,omitempty"`
	CellPhone string `json:"cell_phone,omitempty"`
	Position  string `json:"position,omitempty"`
	Dep       string `json:"dep,omitempty"`
	Leader    string `json:"leader,omitempty"`
	EntryTime string `json:"entry_time,omitempty"`
}

type ObjInfo struct {
	Id            int64      `json:"id"`
	ParId         int64      `json:"par_id,omitempty"`
	ParName       string     `json:"par_name,omitempty"`
	Icon          string     `json:"icon,omitempty"`
	ObjectPrefix  string     `json:"object_prefix,omitempty"`
	ObjectName    string     `json:"object_name,omitempty"`
	ObjectType    string     `json:"object_type,omitempty"`
	ObjectLevel   string     `json:"object_level,omitempty"`
	Router        string     `json:"router,omitempty"`
	CreateTime    string     `json:"create_time,omitempty"`
	Creator       string     `json:"creator,omitempty"`
	Managers      []string   `json:"managers,omitempty"`
	ObjectDesc    string     `json:"object_desc,omitempty"`
	ObjectDoc     string     `json:"object_doc,omitempty"`
	IsJumpNewPage string     `json:"is_jump_new_page,omitempty"`
	IsChildren    string     `json:"is_children,omitempty"`
	ExternalJump  string     `json:"external_jump,omitempty"`
	Children      []*SubMenu `json:"children,omitempty"`
}

type AppTree struct {
	CateName string      `json:"cate_name"`
	CateId   int64       `json:"cateid"`
	Children []*MenuTree `json:"children"`
}

type MenuTree struct {
	AppName  string      `json:"appName"`
	AppId    int64       `json:"appid"`
	CateId   int64       `json:"cateid"`
	Path     string      `json:"path"`
	Children []*MenuTree `json:"children"`
}


type RespCheckPermUser struct {
	Result bool `json:"result"`
	Title string `json:"title"`
	Desc string `json:"desc"`
	Domain string `json:"domain"`
	SkipDesc string `json:"skip_desc"`
}