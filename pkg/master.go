package pkg

import (
	"sf-paas/config"

	uuap "gitlab.sftcwl.com/sf-op-public/golang/uuap2"
)

type Master struct {
	UuapClient *uuap.Uuap2Client
}

func (m *Master) GetUuapClient() *uuap.Uuap2Client {
	return m.UuapClient
}

var (
	master *Master
	Fclient *FerryClient
 	PClient *PermClient
	SClient *SvcClient
 	err error
)
func init() {
	master = &Master{}
	u, err := uuap.NewUuap2ClientForSSO(config.Conf.Util.UuapUrl, "", config.Conf.Util.Uuaptoken, 0)
	if err != nil {
		log.Error("new uuap client failed:", err)
		return
	}

	master.UuapClient = u
	Fclient,err = NewFerryClient()
	if err != nil {
		log.Error("new ferry client failed:", err)
		return
	}
	PClient,err = NewPermClient()
	if err != nil {
		log.Error("new ferry client failed:", err)
		return
	}
	SClient,err = NewSvcClient()
	if err != nil {
		log.Error("new svc client failed:", err)
		return
	}
	return
}
