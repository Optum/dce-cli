# Disposable Cloud Environment (DCE) CLI
> **DCE is your playground in the cloud**

This is the CLI for [DCE](https://github.com/Optum/dce). For usage information, view the complete [command reference](./docs/dce.md).

Disposable Cloud Environment (DCE) manages ephemeral AWS accounts for easy and secure access to the cloud.

DCE users can "lease" an AWS account for a defined period of time and with a limited budget.

At the end of the lease, or if the lease's budget is reached, the account is wiped clean and returned to the account pool so it may be leased again.

## Getting Started & Documentation

Deploy your own Disposable Cloud Environment by following our [quick start guide](./docs/quickstart.md), available on our documentation website:

[dce.readthedocs.io]()

# Feature Availability

| Feature                                 | *nix        | Windows                                                   | Notes                                                             |
| -----------                             | ----------- | -----------                                               | -----------                                                       |
| Deployment (`dce system deploy`)        | Available   | [Unavailable](https://github.com/Optum/dce-cli/issues/21) |                                                                   |
| CLI Initialization (`dce init`)         | Available   | Available                                                 |                                                                   |
| Authentication (`dce auth`)             | Available   | Available                                                 | [Unavailable on Firefox](https://github.com/Optum/dce/issues/166) |
| Account Management (`dce accounts ...`) | Available   | Available                                                 |                                                                   |
| Lease Management (`dce leases ...`)     | Available   | Available                                                 |                                                                   |
| Usage (`dce usage ...`)                 | Available   | Available                                                 |                                                                   |

# Logging

To change the logging level of the DCE CLI, set the DCE_LOG_LEVEL environment variable to `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`, or `PANIC`. When the log level is `INFO`, only un-prefixed log message will be output. This is the default behavior.

## Contributing to DCE

DCE was born at Optum, but belongs to the community. Improve your cloud experience and [open a PR](https://github.com/Optum/dce-cli/pulls).

[Contributor Guidelines](./CONTRIBUTING.md)