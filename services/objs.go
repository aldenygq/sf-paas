package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sf-paas/config"
	"sf-paas/middleware"
	"sf-paas/models"
	"sf-paas/pkg"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var o = &Objs{}

type Objs struct{}

var (
	db      = pkg.Sql
	log     = middleware.Logger
	rclient = pkg.RedisClient
)

func (o *Objs) Migrate(header models.ParamHeader, param models.ParamMigrateObj) (string, error) {
	var (
		obj *models.ObjectInfo = &models.ObjectInfo{}
		err error
	)
	obj.Id = param.Id
	err = obj.Get()
	if err != nil {
		log.Error("get obj info by id  failed:", err)
		return "迁移失败", err
	}
	if !pkg.CheckOpRole(header.Emp) && !strings.Contains(obj.Managers, header.Emp) && !strings.Contains(obj.Creator, header.Emp) {
		log.Error("curr user no perm")
		return "你没有编辑权限", errors.New("你没有编辑权限")
	}

	obj.ParId = param.DstParId
	err = obj.Update()
	if err != nil {
		log.Errorf("update failed:", err)
		return "迁移失败", err
	}

	return "迁移成功", nil
}
func (o *Objs) MigrateInfo(param models.ObjectParams) (*models.MigrateInfo, string, error) {
	var (
		mgrinfo *models.MigrateInfo = &models.MigrateInfo{}
		err     error
		app     *models.ObjectInfo = &models.ObjectInfo{}
		par     *models.ObjectInfo = &models.ObjectInfo{}
		obj     *models.ObjectInfo = &models.ObjectInfo{}
	)
	app.Id = param.Id
	err = app.Get()
	if err != nil {
		log.Errorf("get app failed:", err)
		return nil, "获取迁移信息失败", err
	}
	par.Id = app.ParId
	err = par.Get()
	if err != nil {
		log.Errorf("get par obj info  failed:", err)
		return nil, "获取迁移信息失败", err
	}
	obj.ObjectType = app.ObjectType
	obj.ObjectLevel = "pro"
	objs, err := obj.FrontObjs()
	if err != nil {
		log.Errorf("get pao list failed:", err)
		return nil, "获取迁移信息失败", err
	}
	mgrinfo.Pro = par.ObjectName
	mgrinfo.App = app.ObjectName
	mgrinfo.AppId = app.Id
	mgrinfo.DstProducts = objs
	return mgrinfo, "获取迁移信息成功", nil

}
func (o *Objs) CloudProSearch(param models.CloudProSearchParam) ([]*models.ProAndAppList, string, error) {
	var (
		obj       *models.ObjectInfo      = &models.ObjectInfo{}
		list      []*models.ProAndAppList = make([]*models.ProAndAppList, 0)
		pandas    []*models.ProAndApp     = make([]*models.ProAndApp, 0)
		err       error
		resultmap map[string]string = make(map[string]string, 0)
	)
	obj.ObjectType = "front"
	obj.ObjectLevel = "app"
	objs, err := obj.SearchFrontApp(param.KeyWord)
	if err != nil {
		log.Error("get app list failed:", err)
		return nil, fmt.Sprintf("搜索失败"), err
	}
	applist, _ := json.Marshal(objs)
	log.Info("apps:", string(applist))
	for _, app := range objs {
		var panda *models.ProAndApp = &models.ProAndApp{}
		var pro *models.ObjectInfo = &models.ObjectInfo{}
		pro.Id = app.ParId
		err = pro.Get()
		if err != nil {
			log.Error("get pro by id failed:", err)
			return nil, fmt.Sprintf("搜索失败"), err
		}
		panda.ObjectInfo = *pro
		panda.Child = app
		pandas = append(pandas, panda)
	}

	pas, _ := json.Marshal(pandas)
	log.Info("pas:", string(pas))
	for _, k := range pandas {
		if _, ok := resultmap[k.ObjectName]; !ok {
			var l *models.ProAndAppList = &models.ProAndAppList{}
			resultmap[k.ObjectName] = k.ObjectName
			l.Id = k.Id
			l.ParId = k.ParId
			l.Icon = k.Icon
			l.ObjectName = k.ObjectName
			l.ObjectType = k.ObjectType
			l.ObjectLevel = k.ObjectLevel
			l.Router = k.Router
			l.CreateTime = k.CreateTime
			l.Creator = k.Creator
			l.Managers = k.Managers
			l.ObjectDesc = k.ObjectDesc
			l.ObjectDoc = k.ObjectDoc
			l.IsJumpNewPage = k.IsJumpNewPage
			l.Children = make([]*models.ObjInfo, 0)
			list = append(list, l)
		}
		appkey := k.ObjectName + k.Child.ObjectName
		if _, ok := resultmap[appkey]; !ok {
			var a *models.ObjInfo = &models.ObjInfo{}
			resultmap[appkey] = appkey
			a.Id = k.Child.Id
			a.ParId = k.Child.ParId
			a.Icon = k.Child.Icon
			a.ObjectName = k.Child.ObjectName
			a.ObjectType = k.Child.ObjectType
			a.ObjectLevel = k.Child.ObjectLevel
			//a.Router = k.Child.Router
			a.CreateTime = k.Child.CreateTime
			a.Creator = k.Child.Creator
			a.Managers = strings.Split(k.Child.Managers, ",")
			a.ObjectDesc = k.Child.ObjectDesc
			a.ObjectDoc = k.Child.ObjectDoc
			a.IsJumpNewPage = k.Child.IsJumpNewPage
			a.ExternalJump = k.Child.IsJumpNewPage
			if k.Child.IsJumpNewPage == "N" {
				a.Router = "/paasapp/" + k.Child.ObjectPrefix + "#" + k.Child.Router
			} else {
				a.Router = k.Child.Router
			}
			for i, _ := range list {
				if list[i].Id == k.Id {
					list[i].Children = append(list[i].Children, a)
				}
			}
		}
	}

	return list, "搜索成功", nil
}

func (o *Objs) ProAndAppList() ([]*models.ProAndAppList, string, error) {
	var (
		obj  *models.ObjectInfo      = &models.ObjectInfo{}
		list []*models.ProAndAppList = make([]*models.ProAndAppList, 0)
		err  error
	)
	obj.ObjectType = "front"
	obj.ObjectLevel = "pro"
	pros, err := obj.FrontObjs()
	if err != nil {
		log.Error("get front app list failed:", err)
		return nil, fmt.Sprintf("获取前台产品失败"), err
	}
	for _, p := range pros {
		var appobj *models.ObjectInfo = &models.ObjectInfo{}
		var applist []*models.ObjInfo = make([]*models.ObjInfo, 0)
		appobj.ParId = p.Id
		appobj.ObjectType = "front"
		appobj.ObjectLevel = "app"
		var pa *models.ProAndAppList = &models.ProAndAppList{}
		apps, err := appobj.FrontObjs()
		if err != nil {
			log.Error("get app list  by parent id failed:", err)
			return nil, fmt.Sprintf("获取子应用失败"), err
		}
		if len(apps) <= 0 {
			continue
		}
		for _, app := range apps {
			var a *models.ObjInfo = &models.ObjInfo{}
			a.Id = app.Id
			a.ParId = a.ParId
			a.ParName = p.ObjectName
			a.ObjectName = app.ObjectName
			if app.IsJumpNewPage == "N" {
				a.Router = "/paasapp/" + app.ObjectPrefix + "#" + app.Router
			} else {
				a.Router = app.Router
			}
			a.ObjectDesc = app.ObjectDesc
			a.ObjectDoc = app.ObjectDoc
			a.ExternalJump = app.IsJumpNewPage
			applist = append(applist, a)
		}
		pa.ObjectInfo = *p
		pa.Children = applist

		list = append(list, pa)
	}

	return list, fmt.Sprintf("获取产品/应用列表成功"), nil
}

func (o *Objs) CreateInfo() (*models.CreateInfo, string, error) {
	var (
		err        error
		obj        *models.ObjectInfo = &models.ObjectInfo{}
		createinfo *models.CreateInfo = &models.CreateInfo{}
	)

	filepath := config.Conf.Base.TypePath
	if filepath == "" {
		log.Error("create info config file is empty")
		filepath = "./conf/create_info.json"
	}
	err = middleware.ReadFileToStruct(filepath, createinfo)
	if err != nil {
		log.Error("get create info config failed:", err)
		return nil, fmt.Sprintf("读取配置文件失败"), err
	}

	users, err := pkg.InServiceUsers()
	if err != nil {
		log.Error("get inservice user list  failed:", err)
		return nil, fmt.Sprintf("获取在职用户列表失败"), err
	}

	obj.ObjectLevel = "pro"
	pros, err := obj.GetList()
	if err != nil {
		log.Error("get all pros failed:", err)
		return nil, fmt.Sprintf("获取产品列表失败"), err
	}

	createinfo.Managers = users
	createinfo.Products = pros

	return createinfo, fmt.Sprintf("获取创建信息成功"), nil
}

func (o *Objs) ObjList(header models.ParamHeader, param models.ObjListParams) ([]*models.ObjList, string, error) {
	var (
		err     error
		objlist []*models.ObjList = make([]*models.ObjList, 0)
		errmsg  string
		msg     string
	)

	switch param.ObjType {
	case "pro":
		errmsg = fmt.Sprintf("获取产品列表信息失败")
		msg = fmt.Sprintf("获取产品列表信息成功")
	case "app":
		errmsg = fmt.Sprintf("获取应用列表信息失败")
		msg = fmt.Sprintf("获取应用列表信息成功")
	case "menu":
		errmsg = fmt.Sprintf("获取菜单列表信息失败")
		msg = fmt.Sprintf("获取菜单列表信息成功")
	default:
		return nil, fmt.Sprintf("类型不合法"), errors.New(fmt.Sprintf("类型不合法"))
	}

	list, err := o.ObjectInfos(param)
	if err != nil {
		log.Error("get obj list failed:", err)
		return nil, errmsg, err
	}
	for _, i := range list {
		var obj *models.ObjList = &models.ObjList{}
		var o *models.ObjInfo = &models.ObjInfo{}
		var parobj *models.ObjectInfo = &models.ObjectInfo{}
		var child []*models.SubMenu = make([]*models.SubMenu, 0)
		parobj.Id = i.ParId
		err = parobj.Get()
		if err != nil {
			log.Error("get pro obj info failed:", err)
			return nil, errmsg, err
		}
		o.Id = i.Id
		o.ParId = i.ParId
		o.ParName = parobj.ObjectName
		o.Icon = i.Icon
		o.ObjectName = i.ObjectName
		o.ObjectType = i.ObjectType
		o.ObjectLevel = i.ObjectLevel
		o.ObjectPrefix = i.ObjectPrefix
		o.Router = i.Router
		o.CreateTime = i.CreateTime
		o.Creator = i.Creator
		o.Managers = strings.Split(i.Managers, ",")
		o.ObjectDesc = i.ObjectDesc
		o.ObjectDoc = i.ObjectDoc
		o.IsJumpNewPage = i.IsJumpNewPage
		o.IsChildren = i.IsChildren
		if i.IsChildren == "Y" {
			err = json.Unmarshal([]byte(i.Children), &child)
			if err != nil {
				log.Error("string to json failed:", err)
				return nil, errmsg, err
			}
			o.Children = child
		}
		obj.ObjInfo = *o
		if pkg.CheckOpRole(header.Emp) || strings.Contains(i.Managers, header.Emp) {
			obj.IsOpreate = "Y"

		} else {
			obj.IsOpreate = "N"
		}
		objlist = append(objlist, obj)
	}

	return objlist, msg, nil
}

func (o *Objs) ObjectInfos(params models.ObjListParams) ([]*models.ObjectInfo, error) {
	var (
		err  error
		obj  *models.ObjectInfo   = &models.ObjectInfo{}
		objs []*models.ObjectInfo = make([]*models.ObjectInfo, 0)
	)
	if params.ObjType == "pro" {
		obj.ObjectLevel = params.ObjType
	} else {
		obj.ObjectLevel = params.ObjType
		obj.ParId = params.Id
	}
	objs, err = obj.GetList()
	if err != nil {
		log.Error("get object list failed:", err)
		return nil, err
	}

	return objs, nil
}
func (o *Objs) SearchObj(header models.ParamHeader, param models.ObjListParams) ([]*models.ObjList, string, error) {
	var (
		err     error
		objlist []*models.ObjList    = make([]*models.ObjList, 0)
		list    []*models.ObjectInfo = make([]*models.ObjectInfo, 0)
		msg     string
		obj     *models.ObjectInfo = &models.ObjectInfo{}
	)

	switch param.ObjType {
	case "pro":
		list, err = obj.SearchPro(param)
	case "app":
		list, err = obj.SearchApp(param)
	case "menu":
		list, err = obj.SearchMenu(param)
	default:
		log.Error("obj type invalid")
		return nil, "对象类型不匹配", errors.New("对象类型不匹配")
	}
	if err != nil {
		log.Error("search obj failed:", err)
		return nil, "搜索失败", err
	}

	for _, i := range list {
		var obj *models.ObjList = &models.ObjList{}
		var o *models.ObjInfo = &models.ObjInfo{}
		var parobj *models.ObjectInfo = &models.ObjectInfo{}
		var child []*models.SubMenu = make([]*models.SubMenu, 0)
		parobj.Id = i.ParId
		err = parobj.Get()
		if err != nil {
			log.Error("get pro obj info failed:", err)
			return nil, "搜索失败", err
		}
		o.ParName = parobj.ObjectName
		o.Id = i.Id
		o.ParId = i.ParId
		o.Icon = i.Icon
		o.ObjectName = i.ObjectName
		o.ObjectType = i.ObjectType
		o.ObjectLevel = i.ObjectLevel
		o.Router = i.Router
		o.CreateTime = i.CreateTime
		o.Creator = i.Creator
		o.Managers = strings.Split(i.Managers, ",")
		o.ObjectDesc = i.ObjectDesc
		o.ObjectDoc = i.ObjectDoc
		o.IsJumpNewPage = i.IsJumpNewPage
		o.IsChildren = i.IsChildren
		if i.IsChildren == "Y" {
			err = json.Unmarshal([]byte(i.Children), &child)
			if err != nil {
				log.Error("string to map failed:", err)
				return nil, "搜索失败", err
			}
			o.Children = child
		}
		obj.ObjInfo = *o
		if pkg.CheckOpRole(header.Emp) || strings.Contains(i.Managers, header.Emp) {
			obj.IsOpreate = "Y"
		} else {
			obj.IsOpreate = "N"
		}
		objlist = append(objlist, obj)
	}

	return objlist, msg, nil
}

func (o *Objs) ObjInfo(header models.ParamHeader, param models.ObjListParams) (interface{}, string, error) {
	var (
		err  error
		data interface{}
		msg  string
	)
	switch param.ObjType {
	case "app":
		data, msg, err = o.AppInfo(header, param)
	case "menu":
		data, msg, err = o.MenuInfo(param)
	default:
		return nil, fmt.Sprintf("类型不合法"), errors.New(fmt.Sprintf("类型不合法"))
	}

	if err != nil {
		log.Error("get obj info failed:", err)
		return nil, msg, err
	}

	return data, msg, nil
}

func (o *Objs) MenuInfo(param models.ObjListParams) (*models.MenuInfo, string, error) {
	var (
		menuinfo  *models.MenuInfo = &models.MenuInfo{}
		err       error
		msg       string
		menu      *models.ObjectInfo = &models.ObjectInfo{}
		app       *models.ObjectInfo = &models.ObjectInfo{}
		pro       *models.ObjectInfo = &models.ObjectInfo{}
		bases     []*models.Base     = make([]*models.Base, 0)
		childs    []*models.Base     = make([]*models.Base, 0)
		basemap   map[string]string  = make(map[string]string, 0)
		childmap  map[string]string  = make(map[string]string, 0)
		childmenu []*models.SubMenu  = make([]*models.SubMenu, 0)
	)
	menu.Id = param.Id
	err = menu.Get()
	if err != nil {
		log.Error("get menu info failed:", err)
		return nil, fmt.Sprintf("获取菜单详情失败"), err
	}
	app.Id = menu.ParId
	err = app.Get()
	if err != nil {
		log.Error("get obj info by id failed:", err)
		return nil, fmt.Sprintf("获取应用详情失败"), err
	}

	pro.Id = app.ParId
	err = pro.Get()
	if err != nil {
		log.Error("get obj info by id failed:", err)
		return nil, fmt.Sprintf("获取产品详情失败"), err
	}

	var n int64 = 1
	if menu.IsChildren == "Y" {
		err = json.Unmarshal([]byte(menu.Children), &childmenu)
		if err != nil {
			log.Error("string to json failed:", err)
			return nil, msg, err
		}

		for _, m := range childmenu {
			nkey := strconv.FormatInt(n, 10)
			childmap["菜单"+nkey] = m.ObjectName
			childmap["转发路由"+nkey] = m.Router
			n = n + 1
		}
		for k, v := range childmap {
			if v == "" {
				continue
			}
			var child *models.Base = &models.Base{}
			child.Key = k
			child.Value = v
			childs = append(childs, child)
		}
		menuinfo.Children = childs
	}

	basemap["菜单"] = menu.ObjectName
	basemap["所属应用"] = app.ObjectName
	basemap["转发路由"] = menu.Router
	basemap["管理员"] = menu.Managers
	basemap["描述"] = menu.ObjectDesc
	for k, v := range basemap {
		if v == "" {
			continue
		}
		var base *models.Base = &models.Base{}
		base.Key = k
		base.Value = v
		bases = append(bases, base)
	}

	menuinfo.BaseInfo = bases
	//menuinfo.Children = childs
	menuinfo.Path = pro.ObjectName + "/" + app.ObjectName + "/" + menu.ObjectName
	return menuinfo, fmt.Sprintf("获取菜单详情成功"), nil
}

func (o *Objs) AppInfo(header models.ParamHeader, param models.ObjListParams) (*models.AppInfo, string, error) {
	var (
		appinfo *models.AppInfo = &models.AppInfo{}
		err     error
		msg     string
		menus   []*models.ObjList  = make([]*models.ObjList, 0)
		app     *models.ObjectInfo = &models.ObjectInfo{}
		info    *models.ObjectInfo = &models.ObjectInfo{}
		basemap map[string]string  = make(map[string]string, 0)
		bases   []*models.Base     = make([]*models.Base, 0)
	)
	//app信息
	app.Id = param.Id
	err = app.Get()
	if err != nil {
		log.Error("get app info failed:", err)
		return nil, msg, err
	}

	//父节点信息
	info.Id = app.ParId
	err = info.Get()
	if err != nil {
		log.Error("get par info failed", err)
		return nil, msg, err
	}

	param.ObjType = "menu"
	//菜单信息	//父节点信息
	menulist, err := o.ObjectInfos(param)
	if err != nil {
		log.Errorf("get menu list failed:", err)
		return nil, fmt.Sprintf("获取菜单列表信息失败"), err
	}
	basemap["应用"] = app.ObjectName
	basemap["所属产品"] = info.ObjectName
	basemap["转发路由"] = app.Router
	basemap["管理员"] = app.Managers
	basemap["使用手册"] = app.ObjectDoc
	basemap["描述"] = app.ObjectDesc
	for k, v := range basemap {
		if v == "" {
			continue
		}
		var base *models.Base = &models.Base{}
		base.Key = k
		base.Value = v
		bases = append(bases, base)
	}
	for _, m := range menulist {
		var menu *models.ObjInfo = &models.ObjInfo{}
		var me *models.ObjList = &models.ObjList{}
		var par *models.ObjectInfo = &models.ObjectInfo{}
		var child []*models.SubMenu = make([]*models.SubMenu, 0)
		par.Id = m.ParId
		err = par.Get()
		if err != nil {
			log.Error("get par node info by par id  failed:", err)
			return nil, fmt.Sprintf("获取菜单列表信息失败"), err
		}
		menu.Id = m.Id
		menu.ParId = m.ParId
		menu.ParName = par.ObjectName
		menu.ObjectName = m.ObjectName
		menu.ObjectType = m.ObjectType
		menu.ObjectLevel = m.ObjectLevel
		menu.Router = m.Router
		menu.CreateTime = m.CreateTime
		menu.Creator = m.Creator
		menu.Managers = strings.Split(m.Managers, ",")
		menu.ObjectDesc = m.ObjectDesc
		menu.ObjectDoc = m.ObjectDoc
		menu.IsJumpNewPage = m.IsJumpNewPage
		menu.IsChildren = m.IsChildren
		if m.IsChildren == "Y" {
			err = json.Unmarshal([]byte(m.Children), &child)
			if err != nil {
				log.Error("string to map failed:", err)
				return nil, fmt.Sprintf("获取菜单列表信息失败"), err
			}
			menu.Children = child
		}
		if pkg.CheckOpRole(header.Emp) || strings.Contains(m.Managers, header.Emp) {
			me.IsOpreate = "Y"
		} else {
			me.IsOpreate = "N"
		}
		me.ObjInfo = *menu
		menus = append(menus, me)
	}
	appinfo.BaseInfo = bases
	appinfo.Menus = menus
	appinfo.Path = info.ObjectName + "/" + app.ObjectName

	return appinfo, fmt.Sprintf("获取app信息成功"), err
}

func (o *Objs) Create(header models.ParamHeader, createinfo models.CreateParams) (interface{}, string, error) {
	var (
		objectinfo *models.ObjectInfo = &models.ObjectInfo{}
		err        error
		errmsg     string
		msg        string
		existmsg   string
		menumap    map[string]string = make(map[string]string, 0)
	)
	objectinfo.ObjectName = createinfo.ObjectName
	objectinfo.ObjectType = createinfo.ObjectType
	objectinfo.ObjectLevel = createinfo.ObjectLevel
	objectinfo.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	objectinfo.Creator = header.Cname + "(" + header.Emp + ")"
	objectinfo.Managers = strings.Join(createinfo.Managers, ",")
	objectinfo.ObjectDesc = createinfo.ObjectDesc
	switch createinfo.ObjectLevel {
	case "pro":
		msg = fmt.Sprintf("创建产品成功")
		errmsg = fmt.Sprintf("创建产品失败")
		existmsg = fmt.Sprintf("产品已存在")
		objectinfo.Icon = createinfo.Icon
	case "app":
		msg = fmt.Sprintf("创建app成功")
		errmsg = fmt.Sprintf("创建app失败")
		existmsg = fmt.Sprintf("app已存在")
		objectinfo.ParId = createinfo.ParId
		objectinfo.IsJumpNewPage = createinfo.IsJumpNewPage
		objectinfo.ObjectDoc = createinfo.ObjectDoc
		if objectinfo.IsJumpNewPage == "N" {
			objectinfo.ObjectPrefix = createinfo.ObjectPrefix
		}
		//objectinfo.ObjectPrefix = createinfo.ObjectPrefix
		objectinfo.Router = createinfo.Router
	case "menu":
		msg = fmt.Sprintf("创建菜单成功")
		errmsg = fmt.Sprintf("创建菜单失败")
		existmsg = fmt.Sprintf("菜单已存在")
		var app *models.ObjectInfo = &models.ObjectInfo{}
		objectinfo.ParId = createinfo.ParId
		app.Id = createinfo.ParId
		err = app.Get()
		if err != nil {
			return "", errmsg, err
		}
		objectinfo.IsJumpNewPage = createinfo.IsJumpNewPage
		objectinfo.IsChildren = createinfo.IsChildren
		//if objectinfo.IsJumpNewPage == "N" {
		objectinfo.ObjectPrefix = app.ObjectPrefix
		//}
		objectinfo.Router = createinfo.Router
		if len(createinfo.Children) > 0 {
			for _, v := range createinfo.Children {
				if _, ok := menumap[v.ObjectName]; !ok {
					menumap[v.ObjectName] = v.ObjectName
				} else {
					log.Error("child menu %v exist", v.ObjectName)
					return "", fmt.Sprintf("子菜单%v重复", v.ObjectName), errors.New(fmt.Sprintf("子菜单%v重复", v.ObjectName))
				}
			}
			data, err := middleware.StructToString(createinfo.Children)
			if err != nil {
				log.Error("struct to json  failed:", err)
				return "", errmsg, err
			}
			objectinfo.Children = string(data)
		}
	}
	s, _ := json.Marshal(objectinfo)
	log.Info("obj:", string(s))
	result, err := objectinfo.Exist()
	if err != nil {
		log.Error("check obj exist failed:", err)
		return "", errmsg, err
	}
	if result {
		log.Error("obj %v exist", createinfo.ObjectName)
		return "", existmsg, errors.New(existmsg)
	}

	err = objectinfo.Create()
	if err != nil {
		log.Error("create object failed:", err)
		return "", errmsg, err
	}

	return "", msg, err
}

func (o *Objs) EditInfo(param models.ObjectParams) (*models.EditInfoResp, string, error) {
	var (
		editinfo *models.EditInfoResp = &models.EditInfoResp{}
		err      error
		proinfo  *models.ObjectInfo = &models.ObjectInfo{}
	)

	proinfo.Id = param.Id
	err = proinfo.Get()
	if err != nil {
		log.Error("get object info failed:", err)
		return nil, fmt.Sprintf("获取编辑信息失败"), err
	}
	users, err := pkg.InServiceUsers()
	if err != nil {
		log.Error("get inservice user list  failed:", err)
		return nil, fmt.Sprintf("获取在职用户列表失败"), err
	}
	switch proinfo.ObjectType {
	case "back":
		proinfo.ObjectType = "后台管理"
	case "front":
		proinfo.ObjectType = "前台应用"
	}
	editinfo.ObjectInfo = *proinfo
	editinfo.NewManagers = users
	return editinfo, fmt.Sprintf("获取编辑信息失败"), nil
}

func (o *Objs) Delete(param models.ObjectParams) (string, error) {
	var (
		obj *models.ObjectInfo = &models.ObjectInfo{}
		err error
	)

	obj.Id = param.Id
	err = obj.Get()
	if err != nil {
		log.Error("get object info failed:", err)
		return fmt.Sprintf("删除失败"), err
	}
	err = obj.DeleteObjById()
	if err != nil {
		log.Errorf("delete object failed:", err)
		return fmt.Sprintf("删除失败"), err
	}

	return fmt.Sprintf("删除成功"), err

}

func (o *Objs) Edit(header models.ParamHeader, params models.EditParams) (string, error) {
	var (
		obj      *models.ObjectInfo = &models.ObjectInfo{}
		errmsg   string
		msg      string
		existmsg string
		err      error
		menumap  map[string]string = make(map[string]string, 0)
	)
	obj.Id = params.Id
	err = obj.Get()
	if err != nil {
		log.Error("get obj info fialed:", err)
		return "对象不存在", err
	}
	if !pkg.CheckOpRole(header.Emp) && !strings.Contains(obj.Managers, header.Emp) && !strings.Contains(obj.Creator, header.Emp) {
		log.Error("curr user no perm")
		return "你没有编辑权限", errors.New("你没有编辑权限")
	}
	switch params.ObjectLevel {
	case "pro":
		msg = fmt.Sprintf("编辑产品信息成功")
		errmsg = fmt.Sprintf("编辑产品信息失败")
		existmsg = fmt.Sprintf("产品已存在")
		obj.Icon = params.Icon
	case "app":
		msg = fmt.Sprintf("编辑app信息成功")
		errmsg = fmt.Sprintf("编辑app信息失败")
		existmsg = fmt.Sprintf("app已存在")
		obj.ParId = params.ParId
		obj.IsJumpNewPage = params.IsJumpNewPage
		obj.ObjectDoc = params.ObjectDoc
		if obj.IsJumpNewPage == "N" {
			obj.ObjectPrefix = params.ObjectPrefix
		}
		//obj.ObjectPrefix = params.ObjectPrefix
		obj.Router = params.Router
	case "menu":
		msg = fmt.Sprintf("编辑菜单信息成功")
		errmsg = fmt.Sprintf("编辑菜单信息失败")
		existmsg = fmt.Sprintf("菜单已存在")
		obj.ParId = params.ParId
		obj.IsJumpNewPage = params.IsJumpNewPage
		obj.IsChildren = params.IsChildren
		obj.Router = params.Router
		if len(params.Children) > 0 {
			for _, v := range params.Children {
				if _, ok := menumap[v.ObjectName]; !ok {
					menumap[v.ObjectName] = v.ObjectName
				} else {
					log.Error("child menu %v exist", v.ObjectName)
					return fmt.Sprintf("子菜单%v重复", v.ObjectName), errors.New(fmt.Sprintf("子菜单%v重复", v.ObjectName))
				}
			}
			data, err := middleware.StructToString(params.Children)
			if err != nil {
				log.Error("struct to json  failed:", err)
				return errmsg, err
			}
			obj.Children = data
		}
	}
	//obj.Id = params.Id
	obj.ObjectName = params.ObjectName
	obj.ObjectType = params.ObjectType
	obj.ObjectLevel = params.ObjectLevel
	obj.Managers = strings.Join(params.Managers, ",")
	obj.ObjectDesc = params.ObjectDesc
	result, err := obj.Exist()
	if err != nil {
		log.Error("check obj exist failed:", err)
		return errmsg, err
	}
	if result {
		log.Error("obj %v exist", params.ObjectName)
		return existmsg, errors.New(existmsg)
	}

	err = obj.Update()
	if err != nil {
		log.Error("update object failed:", err)
		return errmsg, err
	}

	return msg, nil
}
func (o *Objs) ClickOn(emp string, appid int64) error {

	var (
		clickrecord *models.ClickVolume = &models.ClickVolume{}
		clickon     *models.ClickVolume = &models.ClickVolume{}
		err         error
	)
	err = db.Table("click_volume").
		Where("emp = ? and obj_id = ?", emp, appid).
		Take(clickrecord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var clicknum int64 = 0
		clickon.Emp = emp
		clickon.ClickNum = clicknum + 1
		clickon.ObjId = appid
		_ = clickon.Create()
		return nil
	} else if err != nil {
		log.Error("click opn failed")
		return err
	}
	clickrecord.ClickNum = clickrecord.ClickNum + 1
	err = clickrecord.Update()
	if err != nil {
		log.Error("update click info failed:%v\n", err)
		return err
	}

	return nil
}
func (o *Objs) AppTree(header models.ParamHeader, param models.ParamAppTree) ([]*models.AppTree, string, error) {
	var (
		err   error
		app   *models.ObjectInfo = &models.ObjectInfo{}
		menus []*models.MenuTree = make([]*models.MenuTree, 0)
		//menu    *models.ObjectInfo = &models.ObjectInfo{}
		apptree  *models.AppTree   = &models.AppTree{}
		apptrees []*models.AppTree = make([]*models.AppTree, 0)
	)

	_ = o.ClickOn(header.Emp, param.AppId)

	app.Id = param.AppId
	err = app.Get()
	if err != nil {
		log.Error("get app by id failed:", err)

		return nil, "获取菜单列表失败", err
	}

	list, err := app.GetChildListById()
	if err != nil {
		log.Error("get child menu list failed:", err)
		return nil, "获取菜单列表失败", err
	}

	for _, l := range list {
		var menu *models.MenuTree = &models.MenuTree{}
		var submenus []*models.SubMenu = make([]*models.SubMenu, 0)
		menu.AppName = l.ObjectName
		menu.AppId = l.Id
		menu.CateId = app.Id
		menu.Path = l.Router
		if l.IsChildren == "Y" {
			var ts []*models.MenuTree = make([]*models.MenuTree, 0)
			err := json.Unmarshal([]byte(l.Children), &submenus)
			if err != nil {
				log.Error("get child menu failed:", err)
				return nil, "获取菜单列表失败", err
			}

			for _, m := range submenus {
				var t *models.MenuTree = &models.MenuTree{}
				t.AppName = m.ObjectName
				t.Path = m.Router
				ts = append(ts, t)
			}
			menu.Children = ts
		} else {
			menu.Children = nil
		}
		menus = append(menus, menu)
	}
	apptree.CateName = app.ObjectName
	apptree.CateId = app.Id
	apptree.Children = menus
	apptrees = append(apptrees, apptree)
	return apptrees, "", nil
}
