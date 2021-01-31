# serverless-cloudfront-page kit
Simple kit for serverless CloudFront page using AWS Lambda.


## Dependence
- aws-lambda-go
- aws-sdk-go-v2


## Requirements
- AWS (Lambda, API Gateway, Amazon CloudFront)
- aws-sam-cli
- golang environment


## Usage

### Edit View
##### HTML
- Edit templates/index.html

##### CSS
- Edit static/css/main.css

##### Javascript
- Edit static/js/main.js

##### Image
- Add image file into static/img/
- Edit templates/header.html like as 'favicon.ico'.

### Deploy
```bash
make clean build
AWS_PROFILE={profile} AWS_DEFAULT_REGION=us-east-1 make bucket={bucket} stack={stack name} deploy
```
