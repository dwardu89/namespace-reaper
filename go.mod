module github.com/dwardu89/k8s-namespace-reaper

go 1.14

replace k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1

replace k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v11.0.0+incompatible
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
