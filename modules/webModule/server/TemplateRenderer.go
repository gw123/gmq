package server

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aymerick/raymond"
	"github.com/labstack/echo"
)

const handlebarsExtension = "handlebars"

// TemplateRenderer 自定义模板渲染器
type TemplateRenderer struct {
	staticFileUrl     string
	staticFileVersion string
	viewsRoot         string
	templates         *sync.Map
}

// NewTemplateRenderer 创建新的模板渲染器
// root是模板文件的根路径
func NewTemplateRenderer(viewsRoot, staticFileUrl, staticFileVersion string) echo.Renderer {
	t := &TemplateRenderer{
		viewsRoot:         viewsRoot,
		staticFileUrl:     staticFileUrl,
		templates:         new(sync.Map),
		staticFileVersion: staticFileVersion,
	}
	raymond.RegisterHelper("static",  t.StaticFileURL)
	if viewsRoot != "" {
		t.registerGlobalPartials()
	}
	return t
}

// Render 渲染模板
// 实现 echo.Renderer 接口
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	name = name + "." + handlebarsExtension
	tpl, err := t.Parse(name)
	if err != nil {
		return err
	}
	result, err := tpl.Exec(data)
	if err != nil {
		return err
	}
	w.Write([]byte(result))
	return nil
}

// Parse 根据name获取解析后的模板
// 模板解析成功后缓存在Map中，下次使用时直接从Map读取
func (t *TemplateRenderer) Parse(name string) (*raymond.Template, error) {
	cached, ok := t.templates.Load(name)
	if ok {
		return cached.(*raymond.Template), nil
	}

	filePath := path.Join(t.viewsRoot, name)
	tpl, err := raymond.ParseFile(filePath)
	if err != nil {
		return nil, err
	}

	t.templates.Store(name, tpl)
	return tpl, nil
}

// registerGlobalPartials 注册全局模板片段
func (t *TemplateRenderer) registerGlobalPartials() error {
	root := filepath.Join(t.viewsRoot, "_partials")
	prefix := root + "/"

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if root == path {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		parts := strings.Split(path, ".")
		if len(parts) != 2 || parts[1] != handlebarsExtension {
			return nil
		}

		name := strings.TrimPrefix(parts[0], prefix)
		tpl, err := raymond.ParseFile(path)
		if err != nil {
			return err
		}

		raymond.RegisterPartialTemplate(name, tpl)
		return nil
	})
}

// staticFileURL 拼接静态文件路径
func (this *TemplateRenderer) StaticFileURL(tp, uri string) string {
	if tp == "" {
		return this.staticFileUrl + "/" + uri + "?v=" + this.staticFileVersion
	} else {
		return this.staticFileUrl + "/" + tp + "/" + uri + "?v=" + this.staticFileVersion
	}
}
