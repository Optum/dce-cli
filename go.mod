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
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.7.0 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.7
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.3
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.8.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.uber.org/thriftrw v1.20.2
	golang.org/x/crypto v0.1.0 // indirect
	gopkg.in/yaml.v2 v2.2.4
)
