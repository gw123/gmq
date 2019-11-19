package debugModule

import (
	"fmt"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
	"github.com/pkg/profile"
	"os"
	"time"
)

type DebugModule struct {
	base.BaseModule
	MemProfile            interface{ Stop() }
	CPUProfile            interface{ Stop() }
	MutexProfile          interface{ Stop() }
	ThreadcreationProfile interface{ Stop() }
	TraceProfile          interface{ Stop() }
	count                 int
	filename              string
	mode                  string
	period                int
}

func NewDebugModule() *DebugModule {
	this := new(DebugModule)
	return this
}

func (debug *DebugModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	if err := debug.BaseModule.Init(app, debug, config); err != nil {
		return err
	}

	debug.period = debug.Config.GetIntItem("period")
	if debug.period == 0 {
		debug.period = 60
	}
	debug.mode = debug.Config.GetStringItem("mode")
	debug.Debug("mode : %s ,period :%d", debug.mode, debug.period)
	return nil
}

func (debug *DebugModule) Handle(event interfaces.Msg) error {
	return nil
}

func (debug *DebugModule) Watch(index int) {
	if index%debug.period == 0 {
		if debug.MemProfile != nil {
			debug.Debug("StopProfile: %s", debug.filename)
			debug.MemProfile.Stop()
			debug.MemProfile = nil
		}
		if debug.CPUProfile != nil {
			debug.Debug("StopProfile: %s", debug.filename)
			debug.CPUProfile.Stop()
			debug.CPUProfile = nil
		}
		if debug.TraceProfile != nil {
			debug.Debug("StopProfile: %s", debug.filename)
			debug.TraceProfile.Stop()
			debug.TraceProfile = nil
		}

		if debug.ThreadcreationProfile != nil {
			debug.Debug("StopProfile: %s", debug.filename)
			debug.ThreadcreationProfile.Stop()
			debug.ThreadcreationProfile = nil
		}

		if debug.MutexProfile != nil {
			debug.Debug("StopProfile: %s", debug.filename)
			debug.MutexProfile.Stop()
			debug.MutexProfile = nil
		}
		//pprof 会把有:的当做url地址 所以文件路径不能有:
		//path := fmt.Sprintf("./logs/pprof/%d_%d:%d", time.Now().Day(), time.Now().Hour(), time.Now().Minute())
		path := fmt.Sprintf("./logs/pprof/%d_%d_%d", time.Now().Day(), time.Now().Hour(), time.Now().Minute())

		err := os.MkdirAll(path, 0755)
		if err == nil {
			switch debug.mode {
			case "cpu":
				debug.CPUProfile = profile.Start(profile.CPUProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook,
				)
				debug.Debug("recode %s Profile: %s",debug.mode ,debug.filename)
			case "mem":
				debug.MemProfile = profile.Start(
					profile.MemProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook,
				)
				debug.Debug("recode %s Profile: %s",debug.mode ,debug.filename)
			case "mutex":
				debug.MutexProfile = profile.Start(profile.MutexProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook,
				)
				debug.Debug("recode %s Profile: %s",debug.mode ,debug.filename)
			case "thread":
				debug.ThreadcreationProfile = profile.Start(profile.ThreadcreationProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook, )
				debug.Debug("recode %s Profile: %s",debug.mode ,debug.filename)
			case "trace":
				debug.TraceProfile = profile.Start(profile.TraceProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook, )
				debug.Debug("recode %s Profile: %s",debug.mode ,debug.filename)
			}
		} else {
			debug.Error("创建文件失败:%s", err)
		}

	}

	return
}
