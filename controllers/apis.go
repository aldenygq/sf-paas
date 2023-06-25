package controllers

import (
	"sf-paas/middleware"
	"sf-paas/models"
	"sf-paas/services"

	"github.com/gin-gonic/gin"
	//"github.com/gpmgo/gopm/log"
)

const (
	HTTP_SUCCESS_CODE = 0
	HTTP_FAIL_CODE    = 1001
)

type Users struct{}
type Objs struct{}
type Supports struct{}
type Oncalls struct{}

var (
	user    *services.User    = &services.User{}
	obj     *services.Objs    = &services.Objs{}
	support *services.Support = &services.Support{}
	oncall *services.Oncall = &services.Oncall{}
	log                       = middleware.Logger
)
func (o *Oncalls) GetOncallByType(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//参数校验
	var param models.ParamOncallByType
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}
	data,msg,err := oncall.GetOncallByType(param)
	if err != nil {
		log.Error("get oncall list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}
	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Oncalls) GetOncall(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	data,msg,err := oncall.GetOncall()
	if err != nil {
		log.Error("get oncall list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}
	
	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Oncalls) OncallList(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	data,msg,err := oncall.OncallList()
	if err != nil {
		log.Error("get oncall list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}
	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Oncalls) AddOncall(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}
	//参数校验
	var param models.ParamAddOncall
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}
	msg,err := oncall.AddOncall(header,param)
	if err != nil {
		log.Error("add oncall failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}
	ctx.Response(HTTP_SUCCESS_CODE, msg, "")
	return
}
/*
func (o *Oncalls) GetOncall(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	
	code,resp,err := user.CheckPerm(header,param)
	if err != nil {
		log.Error("get user perm failed:", err)
		ctx.Response(code, "", resp)
		return
	}
	ctx.Response(code, "", resp)
	return
}
 */
func (u *Users) CheckPerm(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}
	//参数校验
	var param models.ParamApp
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}
	
	code,resp,err := user.CheckPerm(header,param)
	if err != nil {
		log.Error("get user perm failed:", err)
		ctx.Response(code, "", resp)
		return
	}
	ctx.Response(code, "", resp)
	return
}
func (u *Users) CellPhone(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}
	data, msg, err := user.CellPhone(header)
	if err != nil {
		log.Error("get user cell phone  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (s *Supports) FbObjs(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//参数校验
	var param models.ParamFbObjs
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := support.FbObjs(param)
	if err != nil {
		log.Error("get fb objs failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (s *Supports) FbTypes(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//header校验
	data, msg, err := support.FbTypes()
	if err != nil {
		log.Error("get fnb types failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (s *Supports) Submit(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}
	//参数校验
	var param models.SupportParams
	err = ctx.ValidateJson(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	msg, err := support.Submit(header, param)
	if err != nil {
		log.Error("submit fb  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, "")
	return
}
func (s *Supports) MyRecords(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	data, msg, err := support.MyRecords(header)
	if err != nil {
		log.Error("get my support record list  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (s *Supports) AllRecords(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//参数校验
	var param models.SearchParam
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := support.AllRecords(param)
	if err != nil {
		log.Error("get support record list  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (s *Supports) MenuList(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}
	data, msg, err := support.MenuList(header)
	if err != nil {
		log.Error("get support menu list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

//func (u *Users) ClickOn(c *gin.Context) {
//	ctx := middleware.Context{Ctx: c}
//	//header校验
//	var header models.ParamHeader
//	err := ctx.ValidateHeader(&header)
//	if err != nil {
//		log.Error("request header invalid")
//		return
//	}
//
//	var param models.ClickOnParams
//	err = ctx.ValidateJson(&param)
//	if err != nil {
//		log.Error("request param invalid")
//		return
//	}
//
//	_ = obj.ClickOn(header, param)
//
//	ctx.Response(HTTP_SUCCESS_CODE, "", "")
//	return
//}
//
func (u *Users) UserPermTree(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	//参数校验
	var param models.SearchParam
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := user.GetPermTree(header.Emp, param.KeyWord)
	if err != nil {
		log.Error("get user perm tree failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (u *Users) Center(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	data, msg, err := user.Center(header)
	if err != nil {
		log.Error("get user info failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (u *Users) Info(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	data, msg, err := user.Info(header)
	if err != nil {
		log.Error("get user center info failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (u *Users) CommonlyApp(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	data, msg, err := user.CommonlyApp(header)
	if err != nil {
		log.Error("get user center info failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Objs) AppTree(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	//参数校验
	var param models.ParamAppTree
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := obj.AppTree(header, param)
	if err != nil {
		log.Error("get child menu list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Objs) Migrate(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	//参数校验
	var param models.ParamMigrateObj
	err = ctx.ValidateJson(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	msg, err := obj.Migrate(header, param)
	if err != nil {
		log.Error("migrate failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg,"")
	return
}
func (o *Objs) MigrateInfo(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//参数校验
	var param models.ObjectParams
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data,msg, err := obj.MigrateInfo(param)
	if err != nil {
		log.Error("get migrate info failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, data)
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Objs) Create(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Error("request header invalid")
		return
	}

	//参数校验
	var param models.CreateParams
	err = ctx.ValidateJson(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := obj.Create(header, param)
	if err != nil {
		log.Error("create product failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Objs) EditInfo(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//参数校验
	var param models.ObjectParams
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := obj.EditInfo(param)
	if err != nil {
		log.Error("get product edit info failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (o *Objs) Delete(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//参数校验
	var param models.ObjectParams
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	msg, err := obj.Delete(param)
	if err != nil {
		log.Error("delete product  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, "")
	return

}
func (o *Objs) CloudProSearch(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//参数校验
	var param models.CloudProSearchParam
	err := ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	data, msg, err := obj.CloudProSearch(param)
	if err != nil {
		log.Error("search failed failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
func (o *Objs) Edit(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Errorf("request header invalid")
		return
	}
	//参数校验
	var param models.EditParams
	err = ctx.ValidateJson(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}

	msg, err := obj.Edit(header, param)
	if err != nil {
		log.Error("edit product  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, "")
	return
}

func (o *Objs) ProAndAppList(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}

	data, msg, err := obj.ProAndAppList()
	if err != nil {
		log.Error("get product and app list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}
	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (o *Objs) CreateInfo(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	data, msg, err := obj.CreateInfo()
	if err != nil {
		log.Error("get create object info failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}
	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (o *Objs) ObjList(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Errorf("request header invalid")
		return
	}
	var param models.ObjListParams
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}
	data, msg, err := obj.ObjList(header, param)
	if err != nil {
		log.Error("get obj list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (o *Objs) SearchObj(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	//header校验
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Errorf("request header invalid")
		return
	}
	var param models.ObjListParams
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}
	data, msg, err := obj.SearchObj(header, param)
	if err != nil {
		log.Error("get obj list failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}

func (o *Objs) ObjInfo(c *gin.Context) {
	ctx := middleware.Context{Ctx: c}
	var header models.ParamHeader
	err := ctx.ValidateHeader(&header)
	if err != nil {
		log.Errorf("request header invalid")
		return
	}
	var param models.ObjListParams
	err = ctx.Validate(&param)
	if err != nil {
		log.Error("request param invalid")
		return
	}
	data, msg, err := obj.ObjInfo(header, param)
	if err != nil {
		log.Error("get obj info  failed:", err)
		ctx.Response(HTTP_FAIL_CODE, msg, "")
		return
	}

	ctx.Response(HTTP_SUCCESS_CODE, msg, data)
	return
}
