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

func TaskList(c echo.Context) error {
	var (
		index, _ = strconv.Atoi(c.QueryParam("index"))
		size, _  = strconv.Atoi(c.QueryParam("size"))
		tb       = make([]*models.TaskBasic, 0)
		cnt      int64
	)

	size = helper.If(size == 0, define.PageSize, size).(int)
	index = helper.If(index == 0, 1, index).(int)

	err := models.DB.Model(new(models.TaskBasic)).Count(&cnt).Offset((index - 1) * size).Limit(size).Find(&tb).Error
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
		"data": echo.Map{
			"list":  tb,
			"count": cnt,
		},
	})
}

func TaskDetail(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	data := new(TaskDetailResponse)
	err := models.DB.Model(new(models.TaskBasic)).Select("id, name, spec, shell_path data").Where("id = ?", id).Find(&data).Error
	if err != nil {
		log.Println("[DB ERROR]" + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
	}
	b, err := os.ReadFile(data.Data)
	if err != nil {
		log.Println("[READ_FILE ERROR]" + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
		})
	}
	data.Data = string(b)
	return c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "加载成功",
		"data": data,
	})
}

func TaskAdd(c echo.Context) error {
	name := c.FormValue("name")
	spec := c.FormValue("spec")
	data := c.FormValue("data")
	if name == "" || spec == "" || data == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	shellName := helper.GetUUID()
	tb := &models.TaskBasic{
		Name:      name,
		Spec:      spec,
		ShellPath: define.ShellDir + "/" + shellName + ".sh",
		LogPath:   define.LogDir + "/" + shellName + ".log",
	}
	err := models.DB.Create(tb).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}

	err = helper.TouchFile(tb.ShellPath, data)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "保存脚本异常" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "新增成功",
	})
	// 重启
	helper.SendSIGINT()
	return nil
}

func TaskEdit(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")
	spec := c.FormValue("spec")
	data := c.FormValue("data")
	if id == "" || name == "" || spec == "" || data == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	tb := new(models.TaskBasic)
	err := models.DB.Where("id = ?", id).First(tb).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	tb.Name = name
	tb.Spec = spec
	err = models.DB.Where("id = ?", id).Updates(tb).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	err = helper.TouchFile(tb.ShellPath, data)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "保存脚本异常" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "编辑成功",
	})

	helper.SendSIGINT()
	return nil
}

func TaskDelete(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "必填参不能为空",
		})
	}
	err := models.DB.Where("id = ?", id).Delete(new(models.TaskBasic)).Error
	if err != nil {
		log.Println("[DB ERROR] : " + err.Error())
		return c.JSON(http.StatusOK, echo.Map{
			"code": -1,
			"msg":  "系统异常 : " + err.Error(),
		})
	}
	c.JSON(http.StatusOK, echo.Map{
		"code": 200,
		"msg":  "删除成功",
	})

	helper.SendSIGINT()
	return nil
}
