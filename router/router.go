package router

import (
	"fmt"
	
	"sf-paas/config"
	"sf-paas/controllers"
	"sf-paas/middleware"
	
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var (
	usergroup     *gin.RouterGroup
	objectgroup   *gin.RouterGroup
	supportgroup  *gin.RouterGroup
	menutreegroup *gin.RouterGroup
	oncallgroup *gin.RouterGroup
	log                                 = middleware.Logger
	user          *controllers.Users    = &controllers.Users{}
	obj           *controllers.Objs     = &controllers.Objs{}
	support       *controllers.Supports = &controllers.Supports{}
	oncall       *controllers.Oncalls = &controllers.Oncalls{}
)

const (
	HTTP_NOT_FOUND_CODE = 404
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.Conf.Server.Mode)
	r := gin.Default()
	r.Use(middleware.Log())
	pprof.Register(r)

	if err := middleware.TransInit("zh"); err != nil {
		log.Error("init trans failed:", err)
		return nil
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		ctx := middleware.Context{Ctx: c}
		path := c.Request.URL.Path
		method := c.Request.Method

		ctx.Response(HTTP_NOT_FOUND_CODE, fmt.Sprintf("%s %s not found", method, path), nil)
	})

	//用户相关
	userRouter(r)
	//产品/应用/菜单相关
	objectRouter(r)
	//支持反馈相关
	supportRouter(r)

	//左侧菜单列表
	menutreeRouter(r)
	
	//值班管理
	oncallRouter(r)
	
	return r
}
func oncallRouter(r *gin.Engine) {
	oncallgroup = r.Group("/paas/v1/oncall")
	{
		oncallgroup.GET("oncallbytype", oncall.GetOncallByType)
		oncallgroup.GET("oncall", oncall.GetOncall)
		oncallgroup.POST("addoncall", oncall.AddOncall)
		oncallgroup.GET("oncalllist", oncall.OncallList)
	}
}
func menutreeRouter(r *gin.Engine) {
	menutreegroup = r.Group("/paas/v1/applocal")
	{
		menutreegroup.GET("appTree", obj.AppTree)
	}
}
func userRouter(r *gin.Engine) {
	usergroup = r.Group("/paas/v1/user")
	{
		usergroup.GET("permtree", user.UserPermTree)
		usergroup.GET("info", user.Info)
		usergroup.GET("center", user.Center)
		usergroup.GET("commonlyapp", user.CommonlyApp)
		usergroup.GET("cellphone", user.CellPhone)
		usergroup.GET("checkperm", user.CheckPerm)
	}
}
func objectRouter(r *gin.Engine) {
	objectgroup = r.Group("/paas/v1/object")
	{
		objectgroup.GET("pandalist", obj.ProAndAppList)
		objectgroup.GET("cloudprosearch", obj.CloudProSearch)
		objectgroup.GET("createinfo", obj.CreateInfo)
		objectgroup.POST("create", obj.Create)
		objectgroup.GET("editinfo", obj.EditInfo)
		objectgroup.DELETE("delete", obj.Delete)
		objectgroup.POST("edit", obj.Edit)
		objectgroup.GET("objlist", obj.ObjList)
		objectgroup.GET("searchobj", obj.SearchObj)
		objectgroup.GET("objinfo", obj.ObjInfo)
		objectgroup.POST("migrate", obj.Migrate)
		objectgroup.GET("migrateinfo", obj.MigrateInfo)
	}
}

func supportRouter(r *gin.Engine) {
	supportgroup = r.Group("/paas/v1/support")
	{
		supportgroup.GET("fbtypes", support.FbTypes)
		supportgroup.GET("fbobjs", support.FbObjs)
		supportgroup.GET("menulist", support.MenuList)
		supportgroup.GET("allrecords", support.AllRecords)
		supportgroup.GET("myrecords", support.MyRecords)
		supportgroup.POST("submit", support.Submit)
	}

}
