package test

import (
	"github.com/gw123/GMQ/bootstarp"
	"github.com/gw123/GMQ/core"
	"github.com/gw123/GMQ/services"
	"github.com/gw123/GMQ/test"
	"testing"
)

func TestGetGroupTags(t *testing.T) {
	bootstarp.SetConfigFile("./config.yml")
	config := bootstarp.GetConfig()
	App := core.NewApp(config)
	App.Start()

	s, err := services.NewResourceService(App)
	if err != nil {
		t.Error(err)
	}
	tags, err := s.GetGroupTags(180)
	if err != nil {
		t.Fail()
	}
	t.Log("Log ", tags)
}

func TestGroupCategory(t *testing.T) {
	app := test.GetAppFroTest()
	s, err := services.NewResourceService(app)
	if err != nil {
		t.Error(err)
	}
	tags, err := s.GetGroupCategory(180)
	if err != nil {
		t.Fail()
	}
	t.Log("Log ", tags)
}

func TestGetGroupChapter(t *testing.T) {
	app := test.GetAppFroTest()
	s, err := services.NewResourceService(app)
	if err != nil {
		t.Error(err)
	}
	tags, err := s.GetGroupChapter(180)
	if err != nil {
		t.Fail()
	}
	t.Log("Log ", tags)
}

func TestGetChapterResource(t *testing.T) {
	app := test.GetAppFroTest()
	s, err := services.NewResourceService(app)
	if err != nil {
		t.Error(err)
	}
	tags, err := s.GetChapterResource(475)
	if err != nil {
		t.Fail()
	}
	t.Log("Log ", tags)
}

func TestGetGroup(t *testing.T) {
	app := test.GetAppFroTest()
	s, err := services.NewResourceService(app)
	if err != nil {
		t.Error(err)
	}
	g, err := s.GetGroup(180)
	if err != nil {
		t.Error(err)
	}
	t.Log("Log ", g)
}

func TestGetCateogries(t *testing.T) {
	app := test.GetAppFroTest()
	s, err := services.NewResourceService(app)
	if err != nil {
		t.Error(err)
	}
	g, err := s.GetCategories()
	if err != nil {
		t.Error(err)
	}
	for _, i := range g {
		t.Log(*i)
	}
}

func TestGetIndexCtrl(t *testing.T) {
	app := test.GetAppFroTest()
	s, err := services.NewResourceService(app)
	if err != nil {
		t.Error(err)
	}
	g, err := s.GetIndexCtrl(0, 0)
	if err != nil {
		t.Error(err)
	}
	for _, i := range g {
		t.Log(*i)
	}
}
