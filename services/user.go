package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	
	"sf-paas/config"
	"sf-paas/middleware"
	"sf-paas/models"
	"sf-paas/pkg"
	
	"github.com/go-redis/redis"
	//"google.golang.org/appengine/log"
)

const (
	SYSTEM_NAME = "paas"
	SYSTEM_TYPE = "perm"
)

type User struct{}

func (u *User) GetPermTree(emp, keyword string) ([]*models.Tree, string, error) {
	var (
		trees []*models.Tree = make([]*models.Tree, 0)
		err   error
	)
	if keyword == "" {
		trees, err = u.getalltrees(emp)
	} else {
		trees, err = u.searchtrees(emp, keyword)
	}
	if err != nil {
		log.Error("get user perm tree failed:", err)
		return nil, fmt.Sprintf("获取用户权限树列表失败"), err
	}
	if len(trees) <= 0 {
		return nil, fmt.Sprintf("当前用户权限树为空"), errors.New(fmt.Sprintf("当前用户权限树为空"))
	}
	return trees, fmt.Sprintf("获取用户权限树列表成功"), nil
}

func (u *User) getalltrees(emp string) ([]*models.Tree, error) {
	var trees []*models.Tree = make([]*models.Tree, 0)
	key := SYSTEM_TYPE + "-" + SYSTEM_NAME + "-" + emp

	val, err := rclient.Get(key).Result()
	if val == "" || err == redis.Nil {
		log.Error("user(%s) perm tree is empty", emp)
		return nil, nil
	} else if err != nil {
		log.Error("get user(%s) perm tree failed:", emp, err)
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &trees)
	if err != nil {
		log.Error("data decoding failed", err)
		return nil, err
	}

	return trees, nil
}

func (u *User) searchtrees(emp, keyword string) ([]*models.Tree, error) {
	//var trees []*models.Tree = make([]*models.Tree, 0)
	var searchTree []*models.Tree = make([]*models.Tree, 0)
	trees, err := u.getalltrees(emp)
	if err != nil {
		log.Error("get user(%v) all perm tree failed:", emp, err)
		return nil, err
	}

	var result map[int64]int64 = make(map[int64]int64, 0)
	for _, tree := range trees {
		for _, org := range tree.NextNodes {
			for _, pro := range org.NextNodes {
				if strings.Contains(pro.NodeLabel, keyword) {
					var searchorgtree *models.Tree = &models.Tree{}
					var searchdeptree *models.Tree = &models.Tree{}
					searchorgtree.NodeLabel = org.NodeLabel
					searchorgtree.NodeValue = org.NodeValue
					searchorgtree.Type = org.Type
					searchorgtree.NodeId = org.NodeId
					searchorgtree.NextNodes = append(searchorgtree.NextNodes, pro)
					searchdeptree.NodeLabel = tree.NodeLabel
					searchdeptree.NodeValue = tree.NodeValue
					searchdeptree.Type = tree.Type
					searchdeptree.NodeId = tree.NodeId
					searchdeptree.NextNodes = append(searchdeptree.NextNodes, searchorgtree)
					searchTree = append(searchTree, searchdeptree)
				}
			}
		}
	}

	var searchTrees []*models.Tree = make([]*models.Tree, 0)
	for _, tree := range searchTree {
		if _, ok := result[tree.NodeId]; !ok {
			result[tree.NodeId] = tree.NodeId
			var deptree *models.Tree = &models.Tree{}
			deptree.NodeId = tree.NodeId
			deptree.NodeLabel = tree.NodeLabel
			deptree.NodeValue = tree.NodeValue
			deptree.Type = tree.Type
			deptree.NextNodes = make([]*models.Tree, 0)
			searchTrees = append(searchTrees, deptree)
		}
		for _, org := range tree.NextNodes {
			if _, ok := result[org.NodeId]; !ok {
				result[org.NodeId] = org.NodeId
				var orgtree *models.Tree = &models.Tree{}
				orgtree.NodeId = org.NodeId
				orgtree.NodeLabel = org.NodeLabel
				orgtree.NodeValue = org.NodeValue
				orgtree.Type = org.Type
				orgtree.NextNodes = make([]*models.Tree, 0)
				for i, _ := range searchTrees {
					if searchTrees[i].NodeId == tree.NodeId {
						searchTrees[i].NextNodes = append(searchTrees[i].NextNodes, orgtree)
					}
				}
			}
			for _, pro := range org.NextNodes {
				if _, ok := result[pro.NodeId]; !ok {
					for i, _ := range searchTrees {
						if searchTrees[i].NodeId == tree.NodeId {
							for j, _ := range searchTrees[i].NextNodes {
								if searchTrees[i].NextNodes[j].NodeId == org.NodeId {
									searchTrees[i].NextNodes[j].NextNodes = append(searchTrees[i].NextNodes[j].NextNodes, pro)
								}
							}
						}
					}
				}
			}
		}
	}
	return searchTrees, nil
}
func (u *User) CheckPerm(header models.ParamHeader,param models.ParamApp) (int,*models.RespCheckPermUser,error) {
	var (
		result bool
		err error
		obj *models.ObjectInfo = &models.ObjectInfo{}
		resp *models.RespCheckPermUser = &models.RespCheckPermUser{}
	)
	
	parent,err := pkg.SClient.ParentNode(strconv.FormatInt(param.NodeId,10))
	if err != nil {
		log.Infof("获取节点数据失败，节点id:%v,失败原因:%v",param.NodeId,err)
		resp.Title = fmt.Sprintf("抱歉,校验权限失败，请联系开发人员处理。")
		resp.Result = false
		return 1002,resp,err
	}
	path := fmt.Sprintf("%s>%s>%s>%s",parent.Com.NodeChineseName,parent.Dep.NodeChineseName,parent.Org.NodeChineseName,parent.Pro.NodeChineseName)
	resp.Title = fmt.Sprintf("抱歉,您没有权限访问当前页面。")
	resp.Desc = fmt.Sprintf("请您前往权限系统申请对应节点的权限。权限树为:%v",path)
	resp.SkipDesc = "前往权限系统"
	resp.Domain = config.Conf.Util.PermDomain
	
	obj.ObjectName = param.AppName
	result,err = obj.CheckObjExistByAppName()
	if err != nil {
		log.Infof("校验系统:%v权限失败,失败原因:%v",err)
		resp.Result = false
		return 1001,resp,err
	}
	if !result {
		log.Infof("系统:%v不依赖权限")
		resp.Result = true
		return 0,resp,nil
	}
	
	presult,err := pkg.PClient.CheckUserperm(header.Emp,strconv.FormatInt(param.NodeId,10))
	if err != nil {
		log.Infof("校验用户:%v权限失败,失败原因:%v",err)
		resp.Result = false
		return 1001,resp,err
	}
	if !presult {
		log.Infof("当前用户:%v无权限",header.Emp)
		resp.Result = false
		return 0,resp,nil
	}
	resp.Result = true
	return 0,resp,nil
}
func (u *User) CellPhone(header models.ParamHeader) (string, string, error) {
	var err error

	userInfo, err := pkg.UserInfo(header.Emp)
	if err != nil {
		log.Error("get user info by emp failed:", err)
		return "", fmt.Sprintf("获取用户手机号失败"), err
	}

	return userInfo.CellphoneNumber, "", nil
}
func (u *User) Center(header models.ParamHeader) (map[string]interface{}, string, error) {
	var (
		centerinfo map[string]interface{} = make(map[string]interface{}, 0)
		user       *models.UserInfo       = &models.UserInfo{}
		err        error
		crt        *models.CrtInfo = &models.CrtInfo{}
	)

	userInfo, err := pkg.UserInfo(header.Emp)
	if err != nil {
		log.Error("get user info by emp failed:", err)
		return nil, fmt.Sprintf("获取用户中心信息失败"), err
	}

	user.CnName = userInfo.ChineseName
	user.Emp = userInfo.EmployeeNumber
	user.Email = userInfo.EmailAddr
	user.CellPhone = userInfo.CellphoneNumber
	user.Position = userInfo.PositionName
	user.Dep = userInfo.Department
	user.Leader = userInfo.DirectLeaderChineseName + "(" + userInfo.DirectLeaderEmployeeNumber + ")"
	user.EntryTime = userInfo.EntryTime[0:4] + "-" + userInfo.EntryTime[4:6] + "-" + userInfo.EntryTime[6:8]
	
	crt.Emp = header.Emp
	result, err := crt.Exist()
	if err != nil {
		log.Error("get user crt  failed:", err)
		return nil, fmt.Sprintf("获取用户信息失败"), nil
	}
	if result {
		centerinfo["has_crt"] = "Y"
	} else {
		centerinfo["has_crt"] = "N"
	}
	centerinfo["user_info"] = user
	centerinfo["msg"] = fmt.Sprintf("当前用户未申请门神权限，请前往权限系统申请门神权限！")
	return centerinfo, fmt.Sprintf("获取用户信息成功"), nil
}
func (u *User) Info(header models.ParamHeader) (map[string]interface{}, string, error) {
	var (
		userinfo map[string]interface{} = make(map[string]interface{}, 0)
		menus    *models.UserCenter     = &models.UserCenter{}
		err      error
		bases    []*models.Base = make([]*models.Base, 0)
	)

	filepath := config.Conf.Base.UCenterPath
	if filepath == "" {
		log.Error("user center config file is empty")
		filepath = "./conf/user_center.json"
	}

	err = middleware.ReadFileToStruct(filepath, menus)
	if err != nil {
		log.Errorf("get user center config failed:", err)
		return nil, fmt.Sprintf("读取配置文件失败"), err
	}
	if !pkg.CheckOpRole(header.Emp) && !pkg.CheckManager(header.Emp) {
		for _, menu := range menus.BaseInfo {
			if menu.Key == "应用中心" {
				continue
			}
			bases = append(bases, menu)
		}
		menus.BaseInfo = bases
	}
	if pkg.CheckOpRole(header.Emp) {
		userinfo["role"] = "OP"
	} else {
		userinfo["role"] = "RD"
	}

	userinfo["user"] = header.Cname + "(" + header.Emp + ")"
	userinfo["children"] = menus.BaseInfo
	return userinfo, fmt.Sprintf("获取个人信息成功"), nil
}

func (u *User) CommonlyApp(header models.ParamHeader) ([]*models.AppInfor, string, error) {
	var (
		clickinfos []*models.ClickVolume = make([]*models.ClickVolume, 0)
		err        error
		appinfos   []*models.AppInfor = make([]*models.AppInfor, 0)
	)
	err = db.Table("click_volume").
		Where("emp = ?", header.Emp).
		Order("click_num desc").
		Limit(5).Find(&clickinfos).Error
	if err != nil {
		log.Error("get user commonly app failed:", err)
		return nil, fmt.Sprintf("获取用户常用应用失败"), err
	}
	if len(clickinfos) <= 0 {
		log.Error("curr user no have common app")
		return nil, fmt.Sprintf("当前用户无常用应用"), err
	}
	comapp, _ := json.Marshal(clickinfos)
	log.Info("common app:", string(comapp))
	for _, clickinfo := range clickinfos {
		var appinfo *models.ObjectInfo = &models.ObjectInfo{}
		log.Info("obj id:", clickinfo.ObjId)
		appinfo.Id = clickinfo.ObjId
		err = appinfo.Get()
		//appinfo, _, err := o.ObjInfoById(clickinfo.ObjId)
		if err != nil {
			log.Error("get app info by id failed:", err)
			continue
			//return nil, fmt.Sprintf("获取用户常用应用失败"), err
		}
		//ainfo, _ := json.Marshal(appinfo)
		//log.Info("app info:", string(ainfo))
		var app *models.AppInfor = &models.AppInfor{}
		app.Id = appinfo.Id
		app.ObjectName = appinfo.ObjectName
		app.ExternalJump = appinfo.IsJumpNewPage
		if appinfo.IsJumpNewPage == "N" {
			app.Router = "/paasapp/" + appinfo.ObjectPrefix + "#" + appinfo.Router
		} else {
			app.Router = appinfo.Router
		}
		//app.Router = appinfo.Router
		appinfos = append(appinfos, app)
	}

	return appinfos, fmt.Sprintf("获取用户常用应用成功"), nil
}

//func (u *User) ClickOn(emp string, appid int64) error {
//	var (
//		clickrecord *models.ClickVolume = &models.ClickVolume{}
//		clickon     *models.ClickVolume = &models.ClickVolume{}
//		err         error
//	)
//	err = db.Table("click_volume").
//		Where("emp = ? and obj_id = ?", emp, appid).
//		Take(clickrecord).Error
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		var clicknum int64 = 0
//		clickon.Emp = emp
//		clickon.ClickNum = clicknum + 1
//		clickon.ObjId = appid
//		_ = clickon.Create()
//		return nil
//	} else if err != nil {
//		log.Error("click opn failed")
//		return err
//	}
//	clickrecord.ClickNum = clickrecord.ClickNum + 1
//	err = clickrecord.Update()
//	if err != nil {
//		log.Error("update click info failed:%v\n", err)
//		return err
//	}
//
//	return nil
//}
