package services

import (
	"errors"
	"fmt"
	"strings"
	"time"
	
	"sf-paas/config"
	"sf-paas/middleware"
	"sf-paas/models"
	"sf-paas/pkg"
)

type Support struct{}

func (s *Support) FbObjs(param models.ParamFbObjs) ([]string, string, error) {
	var (
		//objs interface{}
		objs []string = make([]string, 0)
		err  error
		obj  *models.ObjectInfo = &models.ObjectInfo{}
	)
	//obj.ObjectLevel = "pro"
	switch param.FbType {
	case "员工点评":
		//objs, err = pkg.InServiceUsers()
		objs, err = pkg.GetOpUsers()
	case "产品点评":
		objs, err = obj.GetProNames()
	default:
		log.Error("fb type invalid")
		return nil, "反馈类型不合法", errors.New("反馈类型不合法")
	}
	if err != nil {
		log.Error("get fb objs by fb type failed:", err)
		return nil, "获取反馈对象失败", err
	}

	return objs, "", err
}
func (s *Support) FbTypes() ([]string, string, error) {
	if len(config.Conf.Base.FbTypes) <= 0 {
		log.Error("no fb type available")
		return nil, "无可用类型", errors.New("无可用类型")
	}
	return config.Conf.Base.FbTypes, "", nil
}

func (s *Support) MenuList(header models.ParamHeader) ([]*models.Base, string, error) {
	var (
		newmenus []*models.Base = make([]*models.Base,0 )
		err   error
		munus *models.SupportMenu = &models.SupportMenu{}
		h map[string]interface{} = make(map[string]interface{},0)
	)
	filepath := config.Conf.Base.SupportPath
	if filepath == "" {
		log.Error("create info config file is empty")
		filepath = "./conf/support_menu.json"
	}
	err = middleware.ReadFileToStruct(filepath, munus)
	if err != nil {
		log.Errorf("get create info config failed:", err)
		return nil, fmt.Sprintf("读取配置文件失败"), err
	}
	if len(munus.BaseInfo) > 0 {
		for _,v := range munus.BaseInfo {
			if !pkg.CheckOpRole(header.Emp) && !pkg.CheckManager(header.Emp) {
				if v.Key == "值班管理" {
					continue
				}
				newmenus = append(newmenus,v)
			}
			if strings.Contains(v.Key,"待办工单") {
				var m *models.Base = &models.Base{}
				h["UUAP-EMPLOYEE-NUMBER"] = header.Emp
				workorders,err := pkg.Fclient.RunningProcess(h)
				if err != nil {
					log.Errorf("get running workorder failed:", err)
					return nil, fmt.Sprintf("获取待办工单信息失败"), err
				}
				m.Value = v.Value
				if len(workorders) > 0 {
					m.Key = fmt.Sprintf("%s(%d)", "工单", len(workorders))
				}else {
					m.Key = fmt.Sprintf("%s", "工单")
				}
				m.ExternalJump = v.ExternalJump
				newmenus = append(newmenus,m)
			}else {
				newmenus = append(newmenus,v)
			}
		}
	}
	return newmenus, "", nil
}

func (s *Support) MyRecords(header models.ParamHeader) ([]*models.SupportInfo, string, error) {
	var (
		data []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err  error
	)

	err = db.Table("support_info").
		Where("submitter like ?", "%"+header.Emp+"%").
		Order("id desc").Find(&data).Error
	if err != nil {
		log.Error("get my record list failed:", err)
		return nil, fmt.Sprintf("获取我提交的记录失败"), err
	}

	return data, fmt.Sprintf("获取我提交的记录成功"), nil
}

func (s *Support) Submit(header models.ParamHeader, param models.SupportParams) (string, error) {
	var (
		support *models.SupportInfo = &models.SupportInfo{}
		err     error
	)

	support.SubmitTime = time.Now().Format("2006-01-02 15:04:05")
	support.FbObj = param.FbObj
	support.IsAnonymous = param.IsAnonymous
	support.Submitter = header.Cname + "(" + header.Emp + ")"
	support.FbType = param.FbType
	support.FbDesc = param.FbDesc

	err = support.Create()
	if err != nil {
		log.Error("create failed:", err)
		return fmt.Sprintf("提交失败"), err
	}
	return fmt.Sprintf("提交成功"), nil
}

func (s *Support) AllRecords(param models.SearchParam) ([]*models.SupportInfo, string, error) {
	var (
		data []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err  error
		msg  string
	)
	if param.KeyWord == "" {
		data, msg, err = s.allrecords()
	} else {
		switch param.KeyType {
		case "反馈对象":
			data, msg, err = s.recordsbyfbobj(param.KeyWord)
		case "提交人":
			data, msg, err = s.recordsbysubmitter(param.KeyWord)
		case "反馈类型":
			data, msg, err = s.recordsbyfbtype(param.KeyWord)
		case "反馈详情":
			data, msg, err = s.recordsbyfbdesc(param.KeyWord)
		default:
			log.Error("search type invalid")
			return nil, fmt.Sprintf("搜索类型不合法"), errors.New(fmt.Sprintf("搜索类型不合法"))
		}
	}

	if err != nil {
		log.Error("get record failed:", err)
		return nil, msg, err
	}

	return data, msg, nil
}

func (s *Support) allrecords() ([]*models.SupportInfo, string, error) {
	var (
		data    []*models.SupportInfo = make([]*models.SupportInfo, 0)
		records []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err     error
	)

	err = db.Table("support_info").
		Order("id desc").Find(&data).Error
	if err != nil {
		log.Error("get all records failed:", err)
		return nil, fmt.Sprintf("获取反馈记录列表失败"), err
	}
	for _, v := range data {
		var record *models.SupportInfo = &models.SupportInfo{}
		record.Id = v.Id
		record.SubmitTime = v.SubmitTime
		record.FbObj = v.FbObj
		if v.IsAnonymous == "Y" {
			record.Submitter = "匿名反馈"
		} else {
			record.Submitter = v.Submitter
		}
		record.FbType = v.FbType
		record.FbDesc = v.FbDesc
		records = append(records, record)
	}
	return records, fmt.Sprintf("获取反馈记录列表成功"), nil
}

func (s *Support) recordsbyfbobj(keyword string) ([]*models.SupportInfo, string, error) {
	var (
		data []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err  error
	)

	err = db.Table("support_info").
		Where("fb_obj like ?", "%"+keyword+"%").
		Order("id desc").Find(&data).Error
	if err != nil {
		log.Error("get records by fb obj failed:", err)
		return nil, fmt.Sprintf("获取反馈记录列表失败"), err
	}

	return data, fmt.Sprintf("获取反馈记录列表成功"), nil
}

func (s *Support) recordsbysubmitter(keyword string) ([]*models.SupportInfo, string, error) {
	var (
		data []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err  error
	)

	err = db.Table("support_info").
		Where("submitter like ?", "%"+keyword+"%").
		Order("id desc").Find(&data).Error
	if err != nil {
		log.Error("get records by submitter failed:", err)
		return nil, fmt.Sprintf("获取反馈记录列表失败"), err
	}

	return data, fmt.Sprintf("获取反馈记录列表成功"), nil
}

func (s *Support) recordsbyfbtype(keyword string) ([]*models.SupportInfo, string, error) {
	var (
		data []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err  error
	)

	err = db.Table("support_info").
		Where("fb_type like ?", "%"+keyword+"%").
		Order("id desc").Find(&data).Error
	if err != nil {
		log.Error("get records by fb type failed:", err)
		return nil, fmt.Sprintf("获取反馈记录列表失败"), err
	}

	return data, fmt.Sprintf("获取反馈记录列表成功"), nil
}

func (s *Support) recordsbyfbdesc(keyword string) ([]*models.SupportInfo, string, error) {
	var (
		data []*models.SupportInfo = make([]*models.SupportInfo, 0)
		err  error
	)

	err = db.Table("support_info").
		Where("fb_desc like ?", "%"+keyword+"%").
		Order("id desc").Find(&data).Error
	if err != nil {
		log.Error("get records by fb desc failed:", err)
		return nil, fmt.Sprintf("获取反馈记录列表失败"), err
	}

	return data, fmt.Sprintf("获取反馈记录列表成功"), nil
}
