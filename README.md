# Push SNS notification to Slack

This script will push a SNS notification to slack when a new message is pushed to the SNS topic it is subscribed to.

## How to build a Go deployment package for AWS Lambda

Run the following commands:

```
# set build os to linux and name executable "main"
GOOS=linux GOARCH=amd64 go build -o main main.go

# then zip the executable file
zip main.zip main
```

Finally, upload directly to a lambda function or save on s3 and point lambda function to s3 url.