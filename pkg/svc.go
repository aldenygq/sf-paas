package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	
	"sf-paas/config"
	"sf-paas/tools"
)

type NodeInfoResponse struct {
	Errno int64     `json:"err_no"`
	ErrMsg string    `json:"err_msg"`
	Data  interface{} `json:"data"`
}
type NodeInfor struct {
	Id                    int64                  `json:"id,omitempty" description:"节点id"`
	NodeEnglishName       string                 `json:"node_english_name,omitempty" description:"节点英文名称"`
	NodeChineseName       string                 `json:"node_chinese_name,omitempty" description:"节点中文名称"`
	AppName               string                 `json:"app_name,omitempty" description:"app名称"`
	NodeType              string                 `json:"node_type,omitempty" description:"节点类型"`
	ParentNodeId          int64                  `json:"parent_node_id,omitempty" description:"父节点id"`
	ParentNodeChineseName string                 `json:"parent_node_chineses_name,omitempty" description:"父节点中文名称"`
	ParentNodeType        string                 `json:"parent_node_type,omitempty" description:"父节点类型"`
	ParentNodeEnglishName string                 `json:"parent_node_english_name,omitempty" description:"父节点英文名称"`
	NodeManager           string                 `json:"node_manager,omitempty" description:"节点管理员"`
	NodePm                string                 `json:"node_pm,omitempty" description:"节点pm"`
	NodeRdManager         string                 `json:"node_rd_pm,omitempty" description:"节点研发管理员"`
	NodeOpManager         string                 `json:"node_op_manager,omitempty" description:"节点运维管理员"`
	NodeQaManager         string                 `json:"node_qa_manager,omitempty" description:"节点测试管理员"`
	NodeDesc              string                 `json:"node_desc,omitempty" description:"节点描述"`
	NodeConfig            map[string]interface{} `json:"node_config,omitempty" description:"节点配置"`
	ServiceType           string                 `json:"service_type,omitempty" description:"服务类型"`
	Tag                   string                 `json:"tag,omitempty" description:"标签"`
	SystemCode            string                 `json:"system_code,omitempty" description:"系统编码"`
	SfSystemCode          string                 `json:"sf_system_code,omitempty" description:"顺丰总部系统编码"`
	ResourceType          string                 `json:"resource_type,omitempty" description:"app类型,common或空表示常规服务,container表示容器服务"`
}
type ParentNodeInfo struct {
	Com *NodeInfor `json:"com,omitempty" description:"公司节点"`
	Dep *NodeInfor `json:"dep,omitempty" description:"部门节点"`
	Org *NodeInfor `json:"org,omitempty" description:"组织节点"`
	Pro *NodeInfor `json:"pro,omitempty" description:"产品线节点"`
	Ser *NodeInfor `json:"ser,omitempty" description:"服务节点"`
}

type SvcClient struct {
	httpClient *tools.Client
}
func NewSvcClient() (*SvcClient,error) {
	c := tools.NewClient(60)
	if c == nil {
		return nil,errors.New("new perm client invalid")
	}
	return &SvcClient{c},nil
}
const (
	HTTP_PARENT_NODE_INFO_API = "/service-magage/v1/parentnodebyid"
)
func (s *SvcClient) getParentNodeInfoUrl(nodeid string) string {
	return fmt.Sprintf("%s%s?nodeid=%s", config.Conf.Util.SvcUrl, HTTP_PARENT_NODE_INFO_API,nodeid)
}
func (s *SvcClient) ParentNode(nodeid string) (*ParentNodeInfo,error) {
	url := s.getParentNodeInfoUrl(nodeid)
	log.Infof("url:%v",url)
	var resp NodeInfoResponse
	if dat, err := s.httpClient.Get(url, nil, nil); err != nil {
		log.Errorf("get parent node info failed:%v", err)
		return nil, err
	} else {
		if err := json.Unmarshal(dat, &resp); err != nil {
			log.Errorf("decode parent node info failed:%v", err)
			return nil, err
		}
		if resp.Errno != 0 {
			log.Errorf("check user perm failed")
			return nil, errors.New(fmt.Sprintf("获取父节点信息失败,失败原因:%v",resp.ErrMsg))
		}
	}
	parent,err := TransInterToModInfo(resp.Data)
	if err != nil {
		return nil, err
	}
	return parent, nil
	
}
func TransInterToModInfo(data interface{}) (*ParentNodeInfo,error) {
	var parent *ParentNodeInfo = &ParentNodeInfo{}
	d,err :=json.Marshal(data)
	if err != nil {
		log.Errorf("decode interface{} to json failed:%v",err)
		return nil,err
	}
	err =json.Unmarshal(d,&parent)
	if err != nil {
		log.Errorf("decode json to struct failed:%v",err)
		return nil,err
	}
	return parent,nil
}