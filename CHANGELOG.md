## vNext

- **Potential breaking change**: Updated commands to output logging messages to STDERR, while JSON
  and other output goes to STDOUT
- Added `version` command to support viewing the current release version of the running binary
- Add `dce leases login` command (without Lease ID)
- Allow DCE version used by `dce system deploy` to be configured as CLI flag, in YAML config, or as env var
- Accept YAML configuration and env vars for `dce system deploy` params for namespace and budget notification emails.
- Fix bug where modified configuration values used by `dce system deploy` would not be reflected in the terraform deployment
- Fix pre-run credentials check: was accepting expired credentials.

## v0.4.0
- Added `--tf-init-options` and `--tf-apply-options` for greater control over underlying provisioner
- Added `--save-options` flag to persist the `--tf-init-options` and `--tf-apply-options` to the
  configuration file when specified. Default behavior is false.
- Persist API host name and path after running `system deploy`
- **Potential breaking change**: explicitly configured logging output to go to STDOUT; command output should go to STDOUT

## v0.3.1
- Moved quick start to readthedocs
- Download deployment assets from the public url rather than using Github's GraphQL API
- Use local terraform backend by default; located in `~/.dce/.cache/module/main.tf`
- Use terraform binary directly, downloaded to `~/.dce/.cache/terraform/${terraform_version}/` folder.
- Output from terraform now redirected to `~/.dce/deploy.log`
- Added `--noprompt` flag for easier scripting

## v0.3.0
- Modified dce auth command to prompt for input and accept base64 encoded credentials string
- Added cognito auth documentation to quickstart

## v0.2.0

**BREAKING CHANGES**
- Fix `--principle-id` flag to be `principal-id`
- Remove MasterAccountCreds from dce.yml. Use default AWS CLI creds instead
- Indent JSON output from CLI commands

## v0.1.0

Initial release
