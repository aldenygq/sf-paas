package pkg

import (
	"sf-paas/config"

	uuap "gitlab.sftcwl.com/sf-op-public/golang/uuap2"
)

func InServiceUsers() ([]string, error) {
	var users []string = make([]string, 0)
	userinfos, err := master.UuapClient.GetInServiceUserInfo()
	if err != nil {
		log.Error("获取在职员工信息失败:%v\n", err)
		return nil, err
	}

	for _, user := range userinfos {
		username := user.ChineseName + "(" + user.EmployeeNumber + ")"
		users = append(users, username)
	}

	return users, nil
}

func UserInfo(emp string) (uuap.UserInfo, error) {
	var (
		user uuap.UserInfo
		err  error
	)
	user, err = master.UuapClient.GetUserInfoByEmployeeNumber(emp)
	if err != nil {
		log.Error("get user info by emp failed:", err)
		return user, err
	}

	return user, nil
}

func GetOpUsers() ([]string, error) {
	var oplist []string = make([]string, 0)
	ops := config.Conf.Base.Opgroups
	for _, op := range ops {
		userinfo, err := UserInfo(op)
		if err != nil {
			log.Error("get user info by emp failed:", err)
			return nil, err
		}
		user := userinfo.ChineseName + "(" + userinfo.EmployeeNumber + ")"
		oplist = append(oplist, user)
	}

	return oplist, nil
}
