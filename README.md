# AWS Disposable Cloud Environments (DCE) CLI

This is the CLI for [DCE](https://github.com/Optum/dce) by Optum. For usage information, view the complete [command reference](./docs/dce.md).

# Quick Start

1. Download the appropriate executable for your OS from the [latest release](https://github.com/Optum/dce-cli/releases/latest). e.g. for mac, you should download dce_darwin_amd64.zip

2. Unzip the artifact and move the executable to a directory on your PATH, e.g.

    ```
    # Download the zip file
    wget https://github.com/Optum/dce-cli/releases/download/<VERSION>/dce_darwin_amd64.zip

    # Unzip to a directory on your path
    unzip dce_darwin_amd64.zip -d /usr/local/bin
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

6. Edit your dce config file with the host and base url from the api gateway that was just deployed to your master account. This can be found in the master account under `API Gateway > (The API with "dce" in the title) > Stages > "Invoke URL: https://<host>/<baseurl>"`.

7. Prepare a second AWS account to be your first "DCE Child Account"
    - Create an IAM role with `AdministratorAccess` and a trust relationship to your DCE Master Accounts
    - Create an account alias by clicking the 'customize' link in the IAM dashboard of the child account

8. Use the `dce accounts add` command to add your child account to the "DCE Accounts Pool"

```
➜  ~ dce accounts add --account-id 555555555555 --admin-role-arn arn:aws:iam::555555555555:role/DCEMasterAccess
```

9. Type `dce accounts list` to verify that your account has been added
➜  ~ dce accounts list
[{"accountStatus":"NotReady","adminRoleArn":"arn:aws:iam::555555555555:role/MasterAcctAcces","createdOn":1572562180,"id":"555555555555","lastModifiedOn":1572637591,"principalPolicyHash":"\"852ee9abbf1220a111c435a8c0e65490\"","principalRoleArn":"arn:aws:iam::555555555555:role/DCEPrincipal-dcelogin"}]

    The account status will initially say `NotReady`. It may take up to 5 minutes for the new account to be processed. Once the account status is `Ready`, you may proceed with creating a lease.

10. Now that your accounts pool isn't emtpy, you can create your first lease using the `dce leases create` command

```
➜  ~ dce leases create --budget-amount 100.0 --budget-currency USD --email jane.doe@email.com --principle-id jdoe99
```

11. Type `dce leases list` to verify that a lease has been created

```
➜  ~ dce leases list
[{"accountId":"555555555555","budgetAmount":100,"budgetCurrency":"USD","budgetNotificationEmails":["jane.doe@email.com "],"createdOn":1572562298,"id":"d326bddf-af36-44d9-bbd0-70b9f7e55356","lastModifiedOn":1572637591,"leaseStatus":"Active","leaseStatusModifiedOn":1572637591,"principalId":"jdoe99"}]
```

12. Access your leased account programmatically via the `dce leases login` command. E.g.

```
# View your leased account credentials

➜  ~ dce leases login d326bddf-af36-44d9-bbd0-70b9f7e55356
aws configure set aws_access_key_id ABCDEFGHIJABCDEFGHIJ;aws configure set aws_secret_access_key ABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJ;aws configure set aws_session_token FLoGEqQvYXdzEP7//////////wEaDBkto0JF2d31gUZAMSLuAWwTAD9/ABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHIJ+DV8jsiOjuhvQYQ9oFUotYe+G+6snwOljzs6ZjHaMQMMmEpSsBNkpwVDBCWnaKHlXgZiv4oiudQ5rxQv0MpXyGxYMmuMTujM/zEyEAjfXoD/fmRFgVK64lkhrW1vgqpKGTMPfX0Ii9ubtL9wYTsDQEG2CpQK28arBl4yTM3DLjPFm3oujslchJrHMGRgSxkrunfdnKvizt5ik7B6B1Hrvt28g5rvV5xG4gIHBK6G1UGfM5IuhcF9hCHQHl5AzayKJiso7rTy7QU==

# Set aws cli credentials

➜  ~ aws ec2 describe-tags
Unable to locate credentials. You can configure credentials by running "aws configure".
➜  ~ eval $(dce leases login d326bddf-af36-44d9-bbd0-70b9f7e55356)
➜  ~ aws ec2 describe-tags
{
    "Tags": []
}
```
13. Access your leased account in a web browser via the `dce leases login` command with the `--open-browser` flag

```
➜  ~ dce leases login d326bddf-af36-44d9-bbd0-70b9f7e55356 --open-browser
Opening AWS Console in Web Browser
```

14. You can end a lease using the `dce leases end` command
```
➜  ~ dce leases end --account-id 555555555555 --principle-id jdoe99
```

15. You can remove an account from the accounts pool using the `dce accounts remove` command
```
➜  ~ dce accounts remove 555555555555
```
# Logging

To change the logging level of the DCE CLI, set the DCE_LOG_LEVEL environment variable to `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`, or `PANIC`. When the log level is `INFO`, only un-prefixed log message will be output. This is the default behavior.

