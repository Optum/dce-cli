# AWS Disposable Cloud Environments (DCE) CLI

This is the CLI for [DCE](https://github.com/Optum/Redbox) by Optum. For usage information, view the complete [command reference](./docs/dce.md).

# Quick Start

1. Download the appropriate executabl3 for your OS from the [latest release](https://github.com/Optum/dce-cli/releases/latest). e.g. for mac, you should download dce_darwin_amd64.zip

2. Unzip the artifact and move the executable to a directory on your PATH, e.g.

    ```
    # Download the zip file
    wget https://github.com/aws/aws-cdk/releases/download/v1.14.0/aws-cdk-1.14.0.zip

    # Unzip to a directory on your path
    unzip aws-cdk-1.14.0.zip -d /usr/local/bin
    ```

3. Test the dce command by typing `dce`
    ```
    ➜  ~ dce
    Disposable Cloud Environment (DCE)

      The DCE cli allows:

      - Admins to provision DCE to a master account and administer said account
      - Users to lease accounts and execute commands against them

    Usage:
      dce [command]

    Available Commands:
      accounts    Manage dce accounts
      auth        Login to dce
      help        Help about any command
      init        First time DCE cli setup. Creates config file at ~/.dce.yaml
      leases      Manage dce leases
      system      Deploy and configure the DCE system

    Flags:
          --config string   config file (default is $HOME/.dce.yaml)
      -h, --help            help for dce

    Use "dce [command] --help" for more information about a command.
    ```

4. Type `dce init` to configure dce via an interactive prompt. This will generate a config file at ~/.dce.yaml

5. Type `dce system deploy` to deploy dce to the AWS account specied in the previous step. This will be your new "DCE Master Account"

6. Edit your dce config file with the api gateway url that was just deployed to your master account. This can be found in the master account under `API Gateway > (The API with "dce" in the title) > Stages`. It is listed as the "Invoke URL".

7. Prepare a second AWS account to be your first "DCE Child Account"
    - Create an IAM role with `AdministratorAccess` and a trust relationship to your DCE Master Accounts
    - Create an account alias by clicking the 'customize' link in the IAM dashboard of the child account

8. Use the `dce accounts add` command to add your child account to the "DCE Accounts Pool"

```
dce accounts add --account-id XXXXXXXXXXXX --admin-role-arn arn:aws:iam::XXXXXXXXXXXX:role/DCEMasterAccess
```

9. Now that your accounts pool isn't emtpy, you can create your first lease using the `dce leases create` command

```
dce leases create --budget-amount 100.0 --budget-currency USD --email jane.doe@email.com --principle-id jdoe99
```g
