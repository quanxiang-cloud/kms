module kms

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/quanxiang-cloud/cabin v0.0.6
	github.com/quanxiang-cloud/qtcc v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.22.4
	sigs.k8s.io/controller-runtime v0.11.0
)

replace github.com/quanxiang-cloud/qtcc => ../qtcc
