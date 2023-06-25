package services

import (
	"fmt"
	"strings"
	
	"sf-paas/config"
	"sf-paas/models"
)
type Oncall struct{}
func (o *Oncall) AddOncall(header models.ParamHeader,param models.ParamAddOncall) (string,error) {
	var (
		oncall *models.Oncall = &models.Oncall{}
		err error
	)
	oncall.Users = strings.Join(param.User,"|")
	oncall.Type = param.Type
	oncall.Creator = header.Cname + "(" + header.Emp + ")"
	err = oncall.Create()
	if err != nil {
		log.Errorf("添加运维值班成员:%v失败,失败原因:%v",param.User,err)
		return fmt.Sprintf("添加运维值班成员:%v失败,失败原因:%v",param.User,err),err
	}
	
	return fmt.Sprintf("添加运维值班成员:%v成功",param.User),nil
}

func (o *Oncall) OncallList() ([]*models.Oncall,string,error) {
	var (
		err error
		oncall *models.Oncall = &models.Oncall{}
	)
	oncalls,err := oncall.GetOncallList()
	if err != nil {
		log.Errorf("获取运维值班信息列表失败,失败原因:%v",err)
		return nil,fmt.Sprintf("获取运维值班信息列表失败,失败原因:%v",err),err
	}
	
	return oncalls,fmt.Sprintf("获取运维值班信息列表成功"),nil
}
func (o *Oncall) GetOncallByType(param models.ParamOncallByType) (*models.Oncall,string,error) {
	var (
		err error
		oncall *models.Oncall = &models.Oncall{}
	)
	oncall.Type = param.OncallType
	err = oncall.GetOncallByType()
	if err != nil {
		log.Errorf("获取%v值班信息列表失败,失败原因:%v",param.OncallType,err)
		return nil,fmt.Sprintf("获取运维值班信息列表失败,失败原因:%v",err),err
	}
	return oncall,"",nil
}
func (o *Oncall) GetOncall() (map[string]string,string,error) {
	var (
		err error
		oncall *models.Oncall = &models.Oncall{}
		oncalllist []string = make([]string,0)
		oncallinfo string
		oncallmap  map[string]string = make(map[string]string,0)
	)
	oncalls,err := oncall.GetOncallList()
	if err != nil {
		log.Errorf("获取运维值班信息列表失败,失败原因:%v",err)
		return nil,fmt.Sprintf("获取运维值班信息列表失败,失败原因:%v",err),err
	}
	if len(oncalls) >0 {
		for _,call := range oncalls {
			cur := fmt.Sprintf("%s(%s)",call.Current,call.Type)
			oncalllist = append(oncalllist,cur)
		}
		oncallinfo = strings.Join(oncalllist,"|")
	}
	oncallmap["oncall"] = oncallinfo
	oncallmap["info"] = config.Conf.Oncall.Info
	return oncallmap,fmt.Sprintf("获取运维值班信息列表成功"),nil
}
