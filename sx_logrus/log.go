// Package sx_logrus 日志封装
package sx_logrus

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	api  *Factory
	once sync.Once
)

// NewApi 获取当前服务工厂接口
func NewApi() *Factory {
	once.Do(func() { api = new(Factory) })
	return api
}

type Factory struct {
	sets sync.Map
}

// New 初始化一个日志服务, 返回日志指针
func (factory *Factory) New(cfg *Cfg) *Log {
	if _, exist := factory.sets.Load(cfg.Key); !exist {
		if "" == cfg.LogFile {
			cfg.LogFile = defaultLogFile
		}
		tmpLog := new(Log)
		tmpLog.init(cfg)
		factory.sets.Store(cfg.Key, tmpLog)
		fmt.Println(fmt.Sprintf("[%s]日志服务[key: %s]注册成功", now(), cfg.Key))
	}
	currentDb, _ := factory.sets.Load(cfg.Key)
	return currentDb.(*Log)
}

// Get 根据关键字获取对应日志指针
func (factory *Factory) Get(key string) (*Log, error) {
	if "" == key {
		key = "default"
	}
	if currentLog, exist := factory.sets.Load(key); exist {
		return currentLog.(*Log), nil
	}
	return nil, fmt.Errorf("[%s]当前日志服务[key: %s]暂未注册", now(), key)
}

type Log struct {
	client *logrus.Logger
	cfg    *Cfg
}

// GetClient 获取客户端
func (log *Log) GetClient() *logrus.Logger {
	return log.client
}

// AddHook 添加钩子
func (log *Log) AddHook(hook logrus.Hook) {
	log.client.Hooks.Add(hook)
}

func (log *Log) init(cfg *Cfg) {
	if nil == log.client {
		log.client = logrus.New()
		log.cfg = cfg
		log.set()
	}
}

func (log *Log) set() {
	if log.cfg.Level < 7 {
		log.client.SetLevel(logrus.Level(log.cfg.Level))
	}
	log.setDefaultFormatter()
	if log.cfg.LogFile != "" { // 如果日志文件非空, 将日志打到对应文件
		var fd *rotatelogs.RotateLogs
		var err error
		optLogFileFmt := log.cfg.LogFile + ".%Y%m%d"
		optWithLinkName := rotatelogs.WithLinkName(log.cfg.LogFile)
		optWithMaxAge := rotatelogs.WithMaxAge(time.Duration(log.cfg.LogMaxAge) * time.Second)
		optWithRotationCount := rotatelogs.WithRotationCount(log.cfg.LogRotationCount)
		if log.cfg.LogRotationTime == 0 {
			log.cfg.LogRotationTime = 60 * 60 * 24
		}
		optWithRotationTime := rotatelogs.WithRotationTime(time.Duration(log.cfg.LogRotationTime) * time.Second)
		var opts []rotatelogs.Option
		if runtime.GOOS != osWindows { // windows系统添加软链有异常, 暂时未处理
			opts = append(opts, optWithLinkName)
		}
		opts = append(opts, optWithRotationTime) // 多久分割一次
		if log.cfg.LogRotationCount > 0 {        // 优先采用保留文件数
			opts = append(opts, optWithRotationCount)
		} else {
			opts = append(opts, optWithMaxAge)
		}

		fd, err = rotatelogs.New(optLogFileFmt, opts...)
		if nil != err {
			fmt.Println(fmt.Sprintf("[%s]日志服务[key: %s]新建文件分割实例异常[已将输出流切到标准输出流], 错误: %s", now(), log.cfg.Key, err.Error()))
			log.client.SetOutput(os.Stdout)
			return
		}

		log.client.SetFormatter(&logrus.JSONFormatter{TimestampFormat: defaultTimestampFormat})
		log.client.SetOutput(fd)
	} else {
		fmt.Println(fmt.Sprintf("[%s]日志服务[key: %s]暂无配置日志文件, 已将输出流切到标准输出流", now(), log.cfg.Key))
		log.client.SetOutput(os.Stdout)
		return
	}
}

// 设置默认格式化
func (log *Log) setDefaultFormatter() {
	log.client.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: defaultTimestampFormat, DisableColors: false,
		ForceColors: true, FullTimestamp: true})
}
