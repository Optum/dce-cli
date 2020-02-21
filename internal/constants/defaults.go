package constants

const CommandShortName = "dce"

const DefaultConfigFileName string = "config.yaml"
const GlobalTFTagDefaults string = `"Terraform":"True","AppName":"DCE","Source":"github.com/Optum/dce//modules"`
const TerraformBinName = "terraform"
const TerraformBinVersion = "0.12.18"
const TerraformBinDownloadURLFormat = "https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip"

// Default version of DCE to deploy, using `dce system deploy`
const DefaultDCEVersion = "0.23.0"
// Default DCE Location
const DefaultDCELocation = "github.com/Optum/dce"
