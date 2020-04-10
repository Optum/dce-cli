## vNext

## v0.5.0

- **Potential breaking change**: Updated commands to output logging messages to STDERR, while JSON
  and other output goes to STDOUT
- Added `version` command to support viewing the current release version of the running binary
- Add `dce leases login` command (without Lease ID)
- Allow DCE version used by `dce system deploy` to be configured as CLI flag, in YAML config, or as env var
- Accept YAML configuration and env vars for `dce system deploy` params for namespace and budget notification emails.
- Fix bug where modified configuration values used by `dce system deploy` would not be reflected in the terraform deployment
- Fix pre-run credentials check: was accepting expired credentials.
- Fix parsing of `--expires-on` flag for `dce leases create` command
- Upgraded the default backend DCE version from 0.23.0 to 0.29.0

## v0.4.0
- Added `--tf-init-options` and `--tf-apply-options` for greater control over underlying provisioner
- Added `--save-options` flag to persist the `--tf-init-options` and `--tf-apply-options` to the
  configuration file when specified. Default behavior is false.
- Persist API host name and path after running `system deploy`
- **Potential breaking change**: explicitly configured logging output to go to STDOUT; command output should go to STDOUT

## v0.3.1

- **BREAKING CHANGE:** Move `~/.dce.yaml` file location to `~/.dce/config.yaml`
- Use local terraform backend by default; located in `~/.dce/.cache/module/main.tf`
- Use terraform binary directly, downloaded to `~/.dce/.cache/terraform/${terraform_version}/` folder.
- Output from terraform now redirected to `~/.dce/deploy.log`
- Moved quick start to readthedocs
- Download deployment assets from the public url rather than using Github's GraphQL API
- Added `--noprompt` flag for easier scripting

**Migration Notes**

dce-cli v0.3.1 introduces a breaking change, in that it expects your yaml configuration to be located at `~/.dce/config.yaml`, instead of `~/.dce.yaml`. If you are migrating from v0.3.0 or lower, you will need to manually move your existing config file to the new location before running any dce-cli commands.

If you do not have an existing config file at `~/.dce.yaml`, you should be able to upgrade to v0.3.1 without problem.

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
