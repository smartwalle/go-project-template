# internal/config

配置文件的加载与校验。

| 文件 | 职责 |
| --- | --- |
| `config.go` | 配置文件的结构，以及应用自身的字段校验 |
| `load.go` | 读取文件、严格解码、校验 |

```go
// Config 是配置文件的结构。
type Config struct {
	// App 是应用元信息。
	App AppConfig `yaml:"app"`

	// HTTP 是 HTTP 服务配置，由传输层维护。
	HTTP http.Config `yaml:"http"`

	// Logging 是日志配置，由日志模块维护。
	Logging logging.Config `yaml:"logging"`

	// Database 是数据库配置，由平台层维护。
	Database database.Config `yaml:"database"`
}
```

- 配置结构由各模块自己维护：传输层维护 http 段，平台层维护 database、logging 段，这里只负责组合
- 解码是严格的（`KnownFields(true)`）：配置里出现未知字段会报错，避免 key 写错被静默忽略
- 校验错误带具体配置路径，便于定位问题
