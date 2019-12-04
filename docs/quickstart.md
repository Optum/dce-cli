# DCE CLI Quickstart


## Installing the DCE CLI

1. Download the appropriate executable for your OS from the [latest release](https://github.com/Optum/dce-cli/releases/latest). e.g. for mac, you should download dce_darwin_amd64.zip

1. Unzip the artifact and move the executable to a directory on your PATH, e.g.

    ```
    # Download the zip file
    wget https://github.com/Optum/dce-cli/releases/download/<VERSION>/dce_darwin_amd64.zip

    # Unzip to a directory on your path
    unzip dce_darwin_amd64.zip -d /usr/local/bin
    ```

1. Test the dce command by typing `dce`
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

1. Type `dce init` to generate a new config file at ~/.dce.yaml. Leave everything blank for now.

## Deploying DCE

1. [Download and install the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)

1. Choose an AWS account to be your new "DCE Master Account" and configure the AWS CLI with user credentials that have AdministratorAccess in that account.

    ```
    ➜  ~ aws configure set aws_access_key_id default_access_key
    ➜  ~ aws configure set aws_secret_access_key default_secret_key
    ```

1. Type `dce system deploy` to deploy dce to the AWS account specified in the previous step.

1. Edit your dce config file with the host and base url from the api gateway that was just deployed to your master account. This can be found in the master account under `API Gateway > (The API with "dce" in the title) > Stages > "Invoke URL: https://<host>/<baseurl>"`. Your config file should look something like this:

    ```
    api:
      host: "abcdefghij.execute-api.us-east-1.amazonaws.com"
      basepath: "/api"
    region: us-east-1
    ```
   
## Authenticating with the DCE System

DCE uses AWS Cognito to manager authentication and authorization.

1. Open the AWS Console in your DCE Master Account and Navigate to AWS Cognito by typing `Cognito` in the search bar

    ![A test image](./cognito.png)

1. Select `Manage User Pools` and click on the dce user pool.

    ![A test image](./manageuserpools.png)

1. Select `Users and groups`

    ![A test image](./usersandgroups.png)
    
1. Create a user

    ![A test image](./createuser.png)

1. Name the user and provide a temporary password. You may uncheck all of the boxes and leave the other fields blank. This user will not have admin priviliges.

    ![A test image](./quickstartuser.png)
    
1. Create a second user to serve as a system admin. Follow the same steps as you did for creating the first user, but name this one something appropriate for their role as an administrator.
   
    ![A test image](./quickstartadmin.png)

1. Create a group

    ![A test image](./creategroup.png)

1. Users in this group will be granted admin access to DCE. The group name must contain the term `Admin`. Choose a name and click on the `Create group` button.

    ![A test image](./groupname.png)
    
1. Add your admin user to the admin group to grant them admin privileges. 
    ![A test image](./quickstartadmindetail.png)
    ![A test image](./addtogroup.png)

1. Type `dce auth` in your command terminal. This will open a browser with a login screen. Enter the username and password for the non-admin user that you created. Reset the password as prompted.

    ![A test image](./quickstartuserlogin.png)

1. Enter the username and password for the non-admin user that you created. Reset the password as prompted.
   
    ![A test image](./quickstartuserlogin.png)

1. Upon successfully logging in, you will be redirected to a credentials page containing a temporary authentication code. Click the button to copy the auth code to your clipboard.
   
    ![A test image](./credspage.png)
    
1. Return to your command terminal and paste the auth code into the prompt, then press enter.
    
    ```
    dce auth
    ✔ Enter API Token: : █ 
    ```

1. You are now authorized as a DCE User. Test that you have proper authorization by typing `dce leases list`.
This will return an empty list indicating that there are currently no leases which you can view. 
If you are not properly authorized as a user, you will see a permissions error.

    ```
    ➜  ~ dce leases list
    []
    ```
   
1. Users are not authorized to list child accounts in the accounts pool. Type `dce accounts list` to verify that you get a permissions error when trying to
view information you do not have access to.

    ```
    ➜  ~ dce accounts list
    err:  [GET /accounts][403] getAccountsForbidden
    ```

1. You will need to be authenticated as an admin before continuing to the next section. Type `dce auth` to login as a different user. Enter the username 
and password for the admin that you created. As before, copy the auth code and paste it in the prompt in your command terminal.

    ![A test image](./quickstartadminlogin.png)
    
1. Test that you have admin authorization by typing `dce accounts list`. You should see an empty list now instead of a permissions error.
   
    ```
    ➜  ~ dce accounts list
    []
    ```

## Adding a child account

1. Prepare a second AWS account to be your first "DCE Child Account"
    - Create an IAM role with `AdministratorAccess` and a trust relationship to your DCE Master Accounts
    - Create an account alias by clicking the 'customize' link in the IAM dashboard of the child account. This must not include the terms "prod" or "production".

1. Use the `dce accounts add` command to add your child account to the "DCE Accounts Pool"

    ```
    ➜  ~ dce accounts add --account-id 555555555555 --admin-role-arn arn:aws:iam::555555555555:role/DCEMasterAccess
    ```

1. Type `dce accounts list` to verify that your account has been added
➜  ~ dce accounts list
[{"accountStatus":"NotReady","adminRoleArn":"arn:aws:iam::555555555555:role/MasterAcctAcces","createdOn":1572562180,"id":"555555555555","lastModifiedOn":1572637591,"principalPolicyHash":"\"852ee9abbf1220a111c435a8c0e65490\"","principalRoleArn":"arn:aws:iam::555555555555:role/DCEPrincipal-dcelogin"}]

    The account status will initially say `NotReady`. It may take up to 5 minutes for the new account to be processed. Once the account status is `Ready`, you may proceed with creating a lease.

## Leasing a DCE Account

1. Now that your accounts pool isn't emtpy, you can create your first lease using the `dce leases create` command

    ```
    ➜  ~ dce leases create --budget-amount 100.0 --budget-currency USD --email jane.doe@email.com --principal-id jdoe99
    ```

1. Type `dce leases list` to verify that a lease has been created

    ```
    ➜  ~ dce leases list
    [{"accountId":"555555555555","budgetAmount":100,"budgetCurrency":"USD","budgetNotificationEmails":["jane.doe@email.com "],"createdOn":1572562298,"id":"d326bddf-af36-44d9-bbd0-70b9f7e55356","lastModifiedOn":1572637591,"leaseStatus":"Active","leaseStatusModifiedOn":1572637591,"principalId":"jdoe99"}]
    ```

## Logging into a leased account

1. Access your leased account programmatically via the `dce leases login` command. E.g.

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
1. Access your leased account in a web browser via the `dce leases login` command with the `--open-browser` flag

    ```
    ➜  ~ dce leases login d326bddf-af36-44d9-bbd0-70b9f7e55356 --open-browser
    Opening AWS Console in Web Browser
    ```

1. You can end a lease using the `dce leases end` command

    ```
    ➜  ~ dce leases end --account-id 555555555555 --principal-id jdoe99
    ```

## Removing a Child Account

1. You can remove an account from the accounts pool using the `dce accounts remove` command

    ```
    ➜  ~ dce accounts remove 555555555555
    ```