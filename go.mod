module github.com/cloud-native-taiwan/controller101

go 1.12

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/thoas/go-funk v0.4.0
	k8s.io/api v0.0.0-20191005115622-2e41325d9e4b
	k8s.io/apimachinery v0.0.0-20191005115455-e71eb83a557c
	k8s.io/client-go v0.0.0-20191005115821-b1fd78950135
	k8s.io/code-generator v0.0.0-20191003035328-700b1226c0bd
	k8s.io/klog v1.0.0
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20181025213731-e84da0312774
	golang.org/x/lint => golang.org/x/lint v0.0.0-20181217174547-8f45f776aaf1
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	golang.org/x/sync => golang.org/x/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190209173611-3b5209105503
	golang.org/x/text => golang.org/x/text v0.3.1-0.20181227161524-e6919f6577db
	golang.org/x/time => golang.org/x/time v0.0.0-20161028155119-f51c12702a4d
	k8s.io/api => k8s.io/api v0.0.0-20191005115622-2e41325d9e4b
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191005115455-e71eb83a557c
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191005115821-b1fd78950135
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191003035328-700b1226c0bd
)
