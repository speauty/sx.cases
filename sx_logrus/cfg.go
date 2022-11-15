package sx_logrus

// Cfg 日志配置结构体
// LogMaxAge和LogRotationCount只能采用一个, 实现中优先采用LogRotationCount
type Cfg struct {
	Key              string `json:"key,omitempty" yaml:"key" xml:"key"  mapstructure:"key" bson:"key"`                                                                           // 配置关键字
	Level            uint32 `json:"level,omitempty" yaml:"level" xml:"level" mapstructure:"level" bson:"level"`                                                                  // 日志等级 PanicLevel(0), FatalLevel(1), ErrorLevel(2), WarnLevel(3), InfoLevel(4), DebugLevel(5), TraceLevel(6)
	LogFile          string `json:"log_file,omitempty" yaml:"log_file" xml:"log_file" mapstructure:"log_file" bson:"log_file"`                                                   // 日志文件 ./runtime/gin.log
	LogRotationTime  int    `json:"log_rotation_time,omitempty" yaml:"log_rotation_time" xml:"log_rotation_time" mapstructure:"log_rotation_time" bson:"log_rotation_time"`      // 日志分割的时间，隔多久分割一次 单位:s
	LogMaxAge        int    `json:"log_max_age,omitempty" yaml:"log_max_age" xml:"log_max_age" mapstructure:"log_max_age" bson:"log_max_age"`                                    // 设置文件清理前的最长保存时间
	LogRotationCount uint   `json:"log_rotation_count,omitempty" yaml:"log_rotation_count" xml:"log_rotation_count" mapstructure:"log_rotation_count" bson:"log_rotation_count"` // 设置文件清理前最多保存的个数(优先采用)
}
