module github.com/Optum/dce-cli

go 1.13

// Two bugs arise during the normal build. The first is the ambiguous import error shown here...
//
//	cannot load github.com/ugorji/go/codec: ambiguous import: found github.com/ugorji/go/codec in multiple modules:
//		github.com/ugorji/go v1.1.4 (/go/pkg/mod/github.com/ugorji/go@v1.1.4/codec)
//		github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 (/go/pkg/mod/github.com/ugorji/go/codec@v0.0.0-20181204163529-d75b2dcb6bc8)
//
// ...which can be fixed by replacing the first module with the second, as suggested in https://github.com/gin-gonic/gin/issues/1673#issuecomment-502203637
// The following error manifests after fixing the first bug...
//
//	panic: codecgen version mismatch: current: 8, need 10. Re-generate file: /go/pkg/mod/github.com/coreos/etcd@v3.3.10+incompatible/client/keys.generated.go
//
// ...which can be overcome by deleting the indicated file and rebuilding, as suggested in the panic itself and here https://github.com/spf13/viper/issues/644#issuecomment-466287597

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go v0.0.0-20181204163529-d75b2dcb6bc8

require (
	github.com/aws/aws-sdk-go v1.25.16
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.7.0 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.7
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.3
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/hashicorp/terraform v0.12.10 // indirect
	github.com/magiconair/properties v1.8.0
	github.com/manifoldco/promptui v0.3.2
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/rjeczalik/interfaces v0.1.0 // indirect
	github.com/shurcooL/githubv4 v0.0.0-20191006152017-6d1ea27df521 // indirect
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/thriftrw v1.20.2
	go.uber.org/zap v1.10.0 // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20191105091915-95d230a53780 // indirect
	gopkg.in/yaml.v2 v2.2.4
)
