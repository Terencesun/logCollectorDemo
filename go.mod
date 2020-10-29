module logCollector

go 1.15

require (
	github.com/Shopify/sarama v1.27.2
	github.com/astaxie/beego v1.12.2
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hpcloud/tail v1.0.0
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/genproto v0.0.0-20201015140912-32ed001d685c // indirect
	google.golang.org/grpc v1.32.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace collectore/config => "./src/config"
