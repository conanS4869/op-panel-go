package service

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"op-panel-go/define"
	"op-panel-go/helper"
	"op-panel-go/models"
	"os"
	"strconv"
)

func WebList(c echo.Context) error {
	var (
		index, _ = strconv.Atoi(c.QueryParam("index"))
		size, _  = strconv.Atoi(c.QueryParam("size"))
		wb       = make([]*models.WebBasic, 0)
		cnt      int64
	)
	size = helper.If(size == 0, define.PageSize, size).(int)
	index = helper.If(index == 0, 1, index).(int)
	err := models.DB.Model(new(models.WebBasic)).Count(&cnt).Offset((index - 1) * size).Limit(size).Find(&wb).Error
	if err != nil {
		log.Println("[DB ERROR]" + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "加载成功",
		"data": wb,
	})
}

func WebAdd(c echo.Context) error {
	var (
		name   = c.FormValue("name")
		domain = c.FormValue("domain")
		cnt    int64
	)

	if name == "" || domain == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	// 判断域名是否已存在
	err := models.DB.Model(new(models.WebBasic)).Where("domain = ?", domain).Count(&cnt).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	if cnt > 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "域名已存在",
		})
	}
	identity := helper.GetUUID()
	wb := &models.WebBasic{
		Identity: identity,
		Name:     name,
		Domain:   domain,
		Dir:      define.DefaultWebDir + domain,
		ConfPath: define.NginxConfigDir + identity + ".conf",
	}
	// 创建网站目录
	err = os.MkdirAll(wb.Dir, 0666)
	if err != nil {
		log.Println("[CREATE DIR ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	// 创建 Nginx 配置目录
	err = os.MkdirAll(define.NginxConfigDir, 0666)
	if err != nil {
		log.Println("[CREATE DIR ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	// 创建新的网站记录
	err = models.DB.Create(wb).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}

	// TODO: 创建nginx配置文件 & 重启nginx加载配置
	return c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "新增成功",
	})
}

func WebEdit(c echo.Context) error {
	var (
		identity = c.FormValue("identity")
		name     = c.FormValue("name")
		domain   = c.FormValue("domain")
		cnt      int64
	)

	if identity == "" || name == "" || domain == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	// 判断域名是否已存在
	err := models.DB.Model(new(models.WebBasic)).Where("domain = ? AND identity != ?", domain, identity).Count(&cnt).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	if cnt > 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "域名已存在",
		})
	}
	wb := &models.WebBasic{
		Name:   name,
		Domain: domain,
		Dir:    define.DefaultWebDir + domain,
	}
	// 创建网站目录
	err = os.MkdirAll(wb.Dir, 0666)
	if err != nil {
		log.Println("[CREATE DIR ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	// 更新网站记录
	err = models.DB.Where("identity = ?", identity).Updates(wb).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}

	// TODO: 创建nginx配置文件 & 重启nginx加载配置
	return c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "编辑成功",
	})
}

func WebDelete(c echo.Context) error {
	identity := c.FormValue("identity")
	if identity == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	err := models.DB.Where("identity = ?", identity).Delete(new(models.WebBasic)).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	// TODO: 删除nginx配置文件 & 重启nginx加载配置
	return c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "删除成功",
	})
}
