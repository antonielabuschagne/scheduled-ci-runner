# Jenkins

### Running Jenkins

You can run Jenkins locally using docker by following [these instructions](https://www.jenkins.io/doc/book/installing/docker/)

### Config

To trigger a build remotely, you would need to configure access to Jenkins following [these instructions](https://narenchejara.medium.com/trigger-jenkins-job-remotely-using-jenkins-api-20973618a493)

Once you have a token, you need to add it to AWS Secrets Manager with the following key: `JENKINS_API_TOKEN`

You can do it via the aws-cli as follows:
```
aws-vault exec personal -- aws secretsmanager create-secret \
   --name JENKINS_API_TOKEN
   --description "Jenkins API token for env: dev"
   --secret-string "{\"JENKINS_API_TOKEN\":\"Shhhhh!\"}"
```

### Remote Access

Once you have Jenkins running locally, you will need to be able to access your instance remotely. However, your instance won't be publicly  accessible. So, I used `ngrok` to resolve this, but you would need to check with your IT department before you do this in case there are security concerns.

Running ngok:
```
npm install -g ngrok
```

```
ngrok http 8080
```

Update your `.env` file to set the `JENKINS_JOB_ENDPOINT` to the ngrok version so your lambda can access.