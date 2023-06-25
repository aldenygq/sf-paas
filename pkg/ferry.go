package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strconv"
	//"strings"
	"time"
	
	"sf-paas/config"
	"sf-paas/tools"
)
type JSONTime struct {
	time.Time
}
type WorkOrderInfo struct {
	//Model
	Sfns 		  string  `gorm:"column:sfns; type:varchar(128)" json:"sfns" form:"sfns"`                                  // 关联sfns
	Title         string          `gorm:"column:title; type:varchar(128)" json:"title" form:"title"`                                  // 工单标题
	Priority      int             `gorm:"column:priority; type:int(11)" json:"priority" form:"priority"`                              // 工单优先级 1，正常 2，紧急 3，非常紧急
	Process       int             `gorm:"column:process; type:int(11)" json:"process" form:"process"`                                 // 流程ID
	Classify      int             `gorm:"column:classify; type:int(11)" json:"classify" form:"classify"`                              // 分类ID
	IsEnd         int             `gorm:"column:is_end; type:int(11); default:0" json:"is_end" form:"is_end"`                         // 是否结束， 0 未结束，1 已结束
	IsDenied      int             `gorm:"column:is_denied; type:int(11); default:0" json:"is_denied" form:"is_denied"`                // 是否被拒绝， 0 没有，1 有
	State         json.RawMessage `gorm:"column:state; type:json" json:"state" form:"state"`                                          // 状态信息
	RelatedPerson json.RawMessage `gorm:"column:related_person; type:json" json:"related_person" form:"related_person"`               // 工单所有处理人
	Creator       int             `gorm:"column:creator; type:int(11)" json:"creator" form:"creator"`                                 // 创建人
	UrgeCount     int             `gorm:"column:urge_count; type:int(11); default:0" json:"urge_count" form:"urge_count"`             // 催办次数
	UrgeLastTime  int             `gorm:"column:urge_last_time; type:int(11); default:0" json:"urge_last_time" form:"urge_last_time"` // 上一次催促时间
}
type Model struct {
	Id        int                `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id" form:"id"`
	CreatedAt JSONTime  `gorm:"column:create_time" json:"create_time" form:"create_time"`
	UpdatedAt JSONTime  `gorm:"column:update_time" json:"update_time" form:"update_time"`
	DeletedAt *JSONTime `gorm:"column:delete_time" sql:"index" json:"-"`
}
type FerryRunningProcessResponse struct {
	Code  int64     `json:"code"`
	Msg string    `json:"msg"`
	Data  []*WorkOrderInfo `json:"data"`
}
type FerryClient struct {
	httpClient *tools.Client
}
func NewFerryClient() (*FerryClient,error) {
	c := tools.NewClient(60)
	if c == nil {
		return nil,errors.New("ferry client invalid")
	}
	return &FerryClient{c},nil
}
const (
	HTTP_RUNNING_PROCESS_API = "/ferry/api/v1/workorder/mytodo"
)
func (f *FerryClient) getRunningProcess() string {
	return fmt.Sprintf("%s%s", config.Conf.Util.FerryUrl, HTTP_RUNNING_PROCESS_API)
}

func (f *FerryClient) RunningProcess(header map[string]interface{}) ([]*WorkOrderInfo,error) {
	//var wokorders []*WorkOrderInfo = make([]*WorkOrderInfo,0)
	url := f.getRunningProcess()
	log.Infof("url:%v",url)
	var resp FerryRunningProcessResponse
	if dat, err := f.httpClient.Get(url, nil, header); err != nil {
		log.Errorf("get running process failed:%v", err)
		return nil, err
	} else {
		if err := json.Unmarshal(dat, &resp); err != nil {
			log.Errorf("decode  workorder info failed:%v", err)
			return nil, err
		}
		if resp.Code != 200 {
			log.Errorf("get  running process by ferry failed:%v", err)
			return nil, errors.New("获取待办工单信息失败")
		}
	}
	return resp.Data, nil
}

