module p00q.cn/video_cdn/server

go 1.17

require (
	github.com/aceld/zinx v1.0.0
	github.com/gogf/gf v1.16.6
	github.com/grafov/m3u8 v0.11.1
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
	p00q.cn/video_cdn/comm v0.0.0
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/clbanning/mxj v1.8.5-0.20200714211355-ff02cfb8ea28 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-ping/ping v0.0.0-20210911151512-381826476871 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gomodule/redigo v1.8.5 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grokify/html-strip-tags-go v0.0.0-20190921062105-daaa06bf1aaf // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	go.opentelemetry.io/otel v1.0.0-RC2 // indirect
	go.opentelemetry.io/otel/trace v1.0.0-RC2 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace p00q.cn/video_cdn/comm => ../comm
