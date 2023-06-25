package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	
	"sf-paas/config"
	"sf-paas/tools"
)
const (
	PERM_SYSTEM_PAAS = "Paas"
)
type PermInfoResponse struct {
	Errno int64     `json:"err_no"`
	ErrMsg string    `json:"err_msg"`
	Data  interface{} `json:"data"`
}
type PermClient struct {
	httpClient *tools.Client
}
func NewPermClient() (*PermClient,error) {
	c := tools.NewClient(60)
	if c == nil {
		return nil,errors.New("new perm client invalid")
	}
	return &PermClient{c},nil
}
const (
	HTTP_CHECK_USER_PERM_API = "/auth/v2/checkuserperm"
)
func (p *PermClient) getCheckUserpermUrl(emp,nodeid string) string {
	return fmt.Sprintf("%s%s?emp=%s&nodeid=%s&system_name=%s", config.Conf.Util.PermUrl, HTTP_CHECK_USER_PERM_API,emp,nodeid,PERM_SYSTEM_PAAS)
}
func (p *PermClient) CheckUserperm(emp,nodeid string) (bool,error) {
	url := p.getCheckUserpermUrl(emp,nodeid)
	log.Infof("url:%v",url)
	var resp PermInfoResponse
	if dat, err := p.httpClient.Get(url, nil, nil); err != nil {
		log.Errorf("check user perm failed:%v", err)
		return false, err
	} else {
		if err := json.Unmarshal(dat, &resp); err != nil {
			log.Errorf("decode user perm info failed:%v", err)
			return false, err
		}
		if resp.Errno != 0 {
			log.Errorf("check user perm failed")
			return false, errors.New(fmt.Sprintf("用户%v权限校验失败,失败原因:%v",emp,resp.ErrMsg))
		}
	}
	return resp.Data.(bool), nil
}

func CheckOpRole(emp string) bool {
	opgroups := config.Conf.Base.Opgroups
	for _, user := range opgroups {
		if emp == user {
			return true
		} else {
			continue
		}
	}

	return false
}
func CheckManager(emp string) bool {
	managers := config.Conf.Base.ObjManager
	for _, user := range managers {
		if emp == user {
			return true
		} else {
			continue
		}
	}

	return false
}
