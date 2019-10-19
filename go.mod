module github.com/Optum/dce-cli

go 1.13

// Two bugs arise during the normal build. The first is the ambiguous import error shown here...
//
//	cannot load github.com/ugorji/go/codec: ambiguous import: found github.com/ugorji/go/codec in multiple modules:
//		github.com/ugorji/go v1.1.1 (/go/pkg/mod/github.com/ugorji/go@v1.1.1/codec)
//		github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 (/go/pkg/mod/github.com/ugorji/go/codec@v0.0.0-20181204163529-d75b2dcb6bc8)
//
// ...which can be fixed by replacing the first module with the second, as suggested in https://github.com/gin-gonic/gin/issues/1673#issuecomment-502203637
// Fixing the first bug results in the following error...
//
//	panic: codecgen version mismatch: current: 8, need 10. Re-generate file: /go/pkg/mod/github.com/coreos/etcd@v3.3.10+incompatible/client/keys.generated.go
//
// ...which can be overcome by deleting the indicated file and rebuilding, as suggested in the panic itself and here https://github.com/spf13/viper/issues/644#issuecomment-466287597

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go v0.0.0-20181204163529-d75b2dcb6bc8

require (
	github.com/aws/aws-sdk-go v1.25.16
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/hashicorp/terraform v0.12.10
	github.com/manifoldco/promptui v0.3.2
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/mitchellh/cli v1.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/shurcooL/githubv4 v0.0.0-20191006152017-6d1ea27df521
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/yaml.v2 v2.2.4
)
