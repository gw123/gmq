package admin

import (
	"errors"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/gw123/GMQ/modules/webModule/db_models"
	"github.com/labstack/echo"
)

type TaskController struct {
	controllers.BaseController
}

func NewTaskController(module interfaces.Module) *TaskController {
	temp := new(TaskController)
	temp.BaseController.Module = module
	return temp
}

/***
 * 登录并且检测是否有新的版本
 */
func (t *TaskController) Login(ctx echo.Context) error {
	version := ctx.QueryParam("version")
	formData := make(map[string]string, 0)
	formData["version"] = version

	return t.Success(ctx, formData)
}

/***
 * 上传新的版本
 */
func (t *TaskController) AddTask(ctx echo.Context) error {
	reqServer := new(db_models.TaskDetail)

	if err := ctx.Bind(reqServer); err != nil {
		return t.Fail(ctx, 0, "参数解析失败", err)
	}

	if err := ctx.Validate(reqServer); err != nil {
		return t.Fail(ctx, 0, "参数验证失败", err)
	}

	task := new(db_models.Task)

	db, err := t.Module.GetApp().GetDefaultDb()
	if err != nil {
		return t.Fail(ctx, 0, "保存失败 001", err)
	}

	result := db.First(task, map[string]interface{}{"name": reqServer.Name, "version": reqServer.Version})

	if result.Error != nil && !result.RecordNotFound() {
		return t.Fail(ctx, 0, "保存失败 003", result.Error)
	}

	if !result.RecordNotFound() {
		return t.Fail(ctx, 0, "版本已经存在", errors.New("版本已经存在"))
	}

	if db.Save(reqServer).Error != nil {
		return t.Fail(ctx, 0, "保存失败 002", db.Save(reqServer).Error)
	}

	return t.Success(ctx, reqServer)
}

func (t *TaskController) AddClientTask(ctx echo.Context) error {
	reqServer := new(db_models.ClientTask)

	if err := ctx.Bind(reqServer); err != nil {
		return t.Fail(ctx, 0, "参数解析失败", err)
	}

	if err := ctx.Validate(reqServer); err != nil {
		return t.Fail(ctx, 0, "参数验证失败", err)
	}

	task := new(db_models.ClientTask)

	db, err := t.Module.GetApp().GetDefaultDb()
	if err != nil {
		return t.Fail(ctx, 0, "保存失败 001", err)
	}

	result := db.First(task, map[string]interface{}{"client_id": reqServer.ClientId, "task_id": reqServer.TaskId})
	if result.Error != nil && !result.RecordNotFound() {
		return t.Fail(ctx, 0, "保存失败 002", result.Error)
	}

	if !result.RecordNotFound() {
		return t.Fail(ctx, 0, "版本已经存在", errors.New("版本已经存在"))
	}

	if db.Save(reqServer).Error != nil {
		return t.Fail(ctx, 0, "保存失败 001", db.Save(reqServer).Error)
	}

	return t.Success(ctx, reqServer)
}

func (t *TaskController) QueryTasksByName(ctx echo.Context) error {
	key := ctx.QueryParam("key")
	db, err := t.Module.GetApp().GetDefaultDb()
	if err != nil {
		return t.Fail(ctx, 0, "查找失败 001", err)
	}
	var tasks []db_models.Task

	//db.LogMode(true)
	res := db.Where("name like ?", "%"+key+"%").Find(&tasks)
	if res.Error != nil && !res.RecordNotFound() {
		return t.Fail(ctx, controllers.ErrorDb, "查找失败 002", err)
	}
	return t.Success(ctx, tasks)
}

/***
 * 下载新的版本
 */
func (t *TaskController) Download(ctx echo.Context) error {

	return t.Success(ctx, nil)
}

