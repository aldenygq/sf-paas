package middleware

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
	//"gopkg.in/go-playground/validator.v10"
)

// ShouldBind
func (c *Context) Validate(p interface{}) error {
	// 参数绑定
	if err := c.Ctx.ShouldBindWith(p, binding.Query); err != nil {
		Logger.Error("param bind error:", err)
		errs, _ := err.(validator.ValidationErrors)
		c.Response(1001, removeTopStruct(errs.Translate(trans)), nil)
		return err
	}

	return nil
}

//ShouldBindJson
func (c *Context) ValidateJson(p interface{}) error {
	// 参数绑定
	if err := c.Ctx.ShouldBindWith(p, binding.JSON); err != nil {
		Logger.Error("param bind error:", err)
		errs, _ := err.(validator.ValidationErrors)
		c.Response(1001, removeTopStruct(errs.Translate(trans)), nil)
		return err
	}

	return nil
}

//ShouldBindHeader
func (c *Context) ValidateHeader(p interface{}) error {
	// 参数绑定
	if err := c.Ctx.ShouldBindWith(p, binding.Header); err != nil {
		Logger.Error("param bind error:", err)
		errs, _ := err.(validator.ValidationErrors)
		c.Response(1001, removeTopStruct(errs.Translate(trans)), nil)
		return err
	}

	return nil
}

// ShouldBindUri
func (c *Context) ValidateRouter(p interface{}) error {
	if err := c.Ctx.ShouldBindUri(p); err != nil {
		Logger.Error("param ShouldBindUri err:", err)
		c.Response(1001, err.Error(), nil)
		return err
	}
	return nil
}

var trans ut.Translator

// loca 通常取决于 http 请求头的 'Accept-Language'
func TransInit(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		//注册一个函数，获取struct tag里自定义的label作为字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})
		//添加额外翻译
		_ = v.RegisterTranslation("required_with", trans, func(ut ut.Translator) error {
			return ut.Add("required_with", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_with", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
			return ut.Add("required_without", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("required_without_all", trans, func(ut ut.Translator) error {
			return ut.Add("required_without_all", "{0} 为必填字段!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_without_all", fe.Field())
			return t
		})
		//register translate
		// 注册翻译器
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return nil
	}
	return nil
}

//func removeTopStruct(fields map[string]string) map[string]string {
func removeTopStruct(fields map[string]string) string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	//var resp map[string]interface{} = make(map[string]interface{},0)
	var resp []string = make([]string, 0)
	for _, value := range res {
		resp = append(resp, value)
	}
	data, _ := json.Marshal(res)
	Logger.Info("data:", string(data))
	response := strings.Join(resp, ",")
	return response
}
