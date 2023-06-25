package models

import (
	"errors"
	"fmt"
	
	//"sf-paas/models"
	"sf-paas/pkg"
	
	//"github.com/gpmgo/gopm/log"
	//"github.com/gpmgo/gopm/log"
	
	"github.com/google/martian/log"
	"github.com/jinzhu/gorm"
)

var db = pkg.Sql
type Oncall struct {
	Id            int64  `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id"`
	Users string `gorm:"column:users;varchar(1024)" json:"users,required"`
	Type string `gorm:"column:type;type:varchar(128)" json:"type,required"`
	Current string `gorm:"column:current;type:varchar(128)" json:"current,required"`
	Creator    string `gorm:"column:creator;type:varchar(256)" json:"creator,required"`
}
func (o *Oncall) Create() error {
	tx := db.Begin()
	
	err := tx.Table("oncall").Create(&o).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	
	return nil
}
func (o *Oncall) GetOncallList() ([]*Oncall,error) {
	var oncalls []*Oncall = make([]*Oncall,0)
	err := db.Table("oncall").Find(&oncalls).Error
	if err != nil {
		return nil,err
	}
	return oncalls,nil
}
func (o *Oncall) GetOncallByType() error {
	err := db.Table("oncall").Take(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("resord not exist")
	} else if err != nil {
		return err
	}
	return nil
}

type ObjectInfo struct {
	Id            int64  `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id"`
	ParId         int64  `gorm:"column:par_id;type:int(10)" json:"par_id,omitempty"`
	Icon          string `gorm:"column:icon;type:varchar(128)" json:"icon,omitempty"`
	ObjectName    string `gorm:"column:object_name;type:varchar(256)" json:"object_name,omitempty"`
	ObjectPrefix  string `gorm:"column:object_prefix;type:varchar(256)" json:"object_prefix,omitempty"`
	ObjectType    string `gorm:"column:object_type;type:varchar(256)" json:"object_type,omitempty"`
	ObjectLevel   string `gorm:"column:object_level;type:varchar(256)" json:"object_level,omitempty"`
	Router        string `gorm:"column:router;type:text" json:"router,omitempty"`
	CreateTime    string `gorm:"column:create_time;type:varchar(256)" json:"create_time,omitempty"`
	Creator       string `gorm:"column:creator;type:varchar(256)" json:"creator,omitempty"`
	Managers      string `gorm:"column:managers;type:text" json:"managers,omitempty"`
	ObjectDesc    string `gorm:"column:object_desc;type:text" json:"object_desc,omitempty"`
	ObjectDoc     string `gorm:"column:object_doc;type:text" json:"object_doc,omitempty"`
	IsJumpNewPage string `gorm:"column:is_jump_new_page;type:varchar(64)" json:"is_jump_new_page,omitempty"`
	IsChildren    string `gorm:"column:is_children;type:varchar(64)" json:"is_children,omitempty"`
	Children      string `gorm:"column:children;type:text" json:"children,omitempty"`
}
func (o *ObjectInfo) CheckObjExistByAppName() (bool,error) {
	err := db.Table("object_info").Where("object_level = ？ and object_name = ? and is_depand_perm = ?","app",o.ObjectName,"Y").Take(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false,nil
	} else if err != nil {
		return false,err
	}
	return true,nil
}
func (o *ObjectInfo) SearchFrontApp(keyword string) ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	err = db.Table("object_info").
		Where("object_level = ? and object_type= ? and object_name like ?", o.ObjectLevel, o.ObjectType, "%"+keyword+"%").Find(&objs).Error
	if err != nil {
		//log.Error("get object like object name failed:", err)
		return nil, err
	}

	return objs, nil
}
func (o *ObjectInfo) FrontObjs() ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	switch o.ObjectLevel {
	case "pro":
		err = db.Table("object_info").
			Where("object_type = ? and object_level = ?", o.ObjectType, o.ObjectLevel).
			Find(&objs).Error
	case "app":
		err = db.Table("object_info").
			Where("par_id = ? and object_type = ? and object_level = ?", o.ParId, o.ObjectType, o.ObjectLevel).
			Find(&objs).Error
	}
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (o *ObjectInfo) GetProNames() ([]string, error) {
	var objs []string = make([]string, 0)
	o.ObjectLevel = "pro"
	objlist, err := o.GetList()
	if err != nil {
		return nil, err
	}
	if len(objlist) <= 0 {
		return nil, nil
	}
	for _, obj := range objlist {
		objs = append(objs, obj.ObjectName)
	}
	return objs, nil

}
func (o *ObjectInfo) GetList() ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	log.Debugf("obj level:", o.ObjectLevel)
	switch o.ObjectLevel {
	case "pro":
		err = db.Table("object_info").
			Where("object_level = ?", o.ObjectLevel).
			Find(&objs).Error
	case "app":
		err = db.Table("object_info").
			//Where("par_id = ? and object_level = ?", o.ParId, o.ObjectLevel).
			Where("object_level = ?", o.ObjectLevel).
			Find(&objs).Error
	case "menu":
		err = db.Table("object_info").
			Where("par_id = ? and object_level = ?", o.ParId, o.ObjectLevel).
			Find(&objs).Error
	}
	if err != nil {
		return nil, err
	}

	return objs, nil
}
func (o *ObjectInfo) GetChildListById() ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	err = db.Table("object_info").
		Where("par_id = ?", o.Id).
		Find(&objs).Error
	if err != nil {
		return nil, err
	}

	return objs, nil
}
func (o *ObjectInfo) SearchPro(param ObjListParams) ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	switch param.KeyType {
	case "产品名称":
		err = db.Table("object_info").
			Where("object_level = ? and object_name like ?", "pro", "%"+param.KeyWord+"%").
			Find(&objs).Error
	case "创建人":
		err = db.Table("object_info").
			Where("object_level = ? and creator like ?", "pro", "%"+param.KeyWord+"%").
			Find(&objs).Error
	default:
		return nil, errors.New("obj type invalid")
	}
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (o *ObjectInfo) SearchApp(param ObjListParams) ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	switch param.KeyType {
	case "应用名称":
		err = db.Table("object_info").
			Where("object_level = ? and object_name like ?", "app", "%"+param.KeyWord+"%").
			Find(&objs).Error
	case "创建人":
		err = db.Table("object_info").
			Where("object_level = ? and creator like ?", "app", "%"+param.KeyWord+"%").
			Find(&objs).Error
	case "所属产品":
		objs, err = o.GetAppsByProName(param.KeyWord)
	case "转发路由":
		err = db.Table("object_info").
			Where("object_level = ? and router like ?", "app", "%"+param.KeyWord+"%").
			Find(&objs).Error
	default:
		//log.Error("object type invalid")
		return nil, errors.New("object type invalid")
	}
	if err != nil {
		//log.Error("search failed:", err)
		return nil, err
	}

	return objs, nil
}

func (o *ObjectInfo) SearchMenu(param ObjListParams) ([]*ObjectInfo, error) {
	var (
		objs []*ObjectInfo = make([]*ObjectInfo, 0)
		err  error
	)
	switch param.KeyType {
	case "应用菜单":
		err = db.Table("object_info").
			Where("object_level = ? and object_name like ?", "menu", "%"+param.KeyWord+"%").
			Find(&objs).Error
	case "创建人":
		err = db.Table("object_info").
			Where("object_level = ? and creator like ?", "menu", "%"+param.KeyWord+"%").
			Find(&objs).Error
	case "转发路由":
		err = db.Table("object_info").
			Where("object_level = ? and router like ?", "menu", "%"+param.KeyWord+"%").
			Find(&objs).Error
	default:
		return nil, errors.New("object type invalid")
	}
	if err != nil {
		return nil, err
	}

	return objs, nil
}
func (o *ObjectInfo) Get() error {
	err := db.Table("object_info").
		Take(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("查询对象不存在"))
	} else if err != nil {
		return err
	}
	return nil
}

func (o *ObjectInfo) GetAppsByProName(name string) ([]*ObjectInfo, error) {
	var applist []*ObjectInfo = make([]*ObjectInfo, 0)
	pros, err := o.GetProByName(name)
	if err != nil {
		//log.Error("get pros by name  failed:", err)
		return nil, err
	}
	if len(pros) <= 0 {
		//log.Error("pro not exist")
		return nil, errors.New("查找对象不存在")
	}
	for _, pro := range pros {
		var apps []*ObjectInfo = make([]*ObjectInfo, 0)
		err = db.Table("object_info").
			Where("par_id = ?", pro.Id).
			Find(&apps).Error
		if err != nil {
			//log.Errorf("get app list by par id failed:", err)
			return nil, err
		}
		applist = append(applist, apps...)
	}

	return applist, nil
}
func (o *ObjectInfo) GetProByName(name string) ([]*ObjectInfo, error) {
	var objs []*ObjectInfo = make([]*ObjectInfo, 0)
	err := db.Table("object_info").
		Where("object_level = ？and object_name like ?", "pro", "%"+name+"%").
		Find(&objs).Error
	if err != nil {
		//log.Error("get pro by name failed:", err)
		return nil, err
	}

	return objs, nil
}
func (o *ObjectInfo) Exist() (bool, error) {
	var err error
	switch o.ObjectLevel {
	case "pro":
		err = db.Table("object_info").
			Where("object_level = ? and object_name = ? and id != ?", o.ObjectLevel, o.ObjectName, o.Id).
			Take(&o).Error
	case "app":
		err = db.Table("object_info").
			Where("par_id = ? and object_level = ? and object_name = ? and id != ?", o.ParId, o.ObjectLevel, o.ObjectName, o.Id).
			Take(&o).Error
	case "menu":
		err = db.Table("object_info").
			Where("par_id = ? and object_level = ? and object_name = ? and id != ?", o.ParId, o.ObjectLevel, o.ObjectName, o.Id).
			Take(&o).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return true, err
	}
	return true, nil
}

func (o *ObjectInfo) Create() error {
	tx := db.Begin()

	err := tx.Table("object_info").Create(&o).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (o *ObjectInfo) DeleteObjById() error {
	var (
		c *ClickVolume = &ClickVolume{}
		m *ObjectInfo  = &ObjectInfo{}
		a *ObjectInfo  = &ObjectInfo{}
	)
	tx := db.Begin()
	switch o.ObjectLevel {
	case "app":
		err := tx.Table("click_volume").Where("obj_id = ?", o.Id).Delete(c).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Table("object_info").Where("par_id = ?", o.Id).Delete(m).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	case "pro":
		apps, err := o.GetChildListById()
		if err != nil {
			tx.Rollback()
			//log.Error("get child list by id failed:",err)
			return err
		}
		if len(apps) > 0 {
			for _, app := range apps {
				err = tx.Table("object_info").Where("par_id = ?", app.ParId).Delete(m).Error
				if err != nil {
					//log.Error("delete memu failed:",err)
					tx.Rollback()
					return err
				}
				err = tx.Table("object_info").Where("id = ?", app.Id).Delete(a).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	err := tx.Table("object_info").Delete(&o).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (o *ObjectInfo) Update() error {
	tx := db.Begin()
	err := tx.Table("object_info").
		Model(o).
		Updates(o).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

type SupportInfo struct {
	Id          int64  `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id"`
	SubmitTime  string `gorm:"column:submit_time;type:varchar(64)" json:"submit_time"`
	FbObj       string `gorm:"column:fb_obj;type:varchar(128)" json:"fb_obj,omitempty"`
	Submitter   string `gorm:"column:submitter;type:varchar(128)" json:"submitter,omitempty"`
	FbType      string `gorm:"column:fb_type;type:varchar(128)" json:"fb_type,omitempty"`
	FbDesc      string `gorm:"column:fb_desc;type:text" json:"fb_desc,omitempty"`
	IsAnonymous string `gorm:"column:is_anonymous;type:varchar(64)" json:"is_anonymous,omitempty"`
}

func (s *SupportInfo) Create() error {
	tx := db.Begin()
	err := tx.Table("support_info").Create(s).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

type ClickVolume struct {
	Id       int64  `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id"`
	Emp      string `gorm:"column:emp;type:varchar(128)" json:"emp"`
	ClickNum int64  `gorm:"column:click_num;type:int(10)" json:"click_num"`
	ObjId    int64  `gorm:"column:obj_id;type:int(10)" json:"obj_id"`
}

func (c *ClickVolume) Create() error {
	tx := db.Begin()
	err := tx.Table("click_volume").
		Create(&c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (c *ClickVolume) Update() error {
	tx := db.Begin()
	err := tx.Table("click_volume").
		Model(c).
		Updates(c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
func (c *ClickVolume) Delete() error {
	tx := db.Begin()
	err := tx.Table("click_volume").Delete(c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

type CrtInfo struct {
	Id         int64  `gorm:"column:id;PRIMARY_KEY;type:int(10)" json:"id"`
	Emp        string `gorm:"column:emp;type:varchar(64)" json:"emp"`
	Record     string `gorm:"column:record;type:varchar(64)" json:"record"`
	LastRecord string `gorm:"column:last_record;type:varchar(64)" json:"last_record"`
	UpdateTime string `gorm:"column:update_time;type:varchar(64)" json:"update_time"`
}

func (c *CrtInfo) Exist() (bool, error) {
	err := db.Table("crt_info").
		Where("emp =  ?", c.Emp).
		Take(&c).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
