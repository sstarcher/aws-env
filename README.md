AWS Environment
================

[![CircleCI](https://circleci.com/gh/sstarcher/job-runner.svg?style=svg)](https://circleci.com/gh/sstarcher/job-runner)
[![](https://imagelayers.io/badge/sstarcher/aws-env:latest.svg)](https://imagelayers.io/?images=sstarcher/aws-env:latest 'Get your own badge on imagelayers.io')
[![Docker Registry](https://img.shields.io/docker/pulls/sstarcher/aws-env.svg)](https://registry.hub.docker.com/u/sstarcher/aws-env)&nbsp;

This repo outputs AWS Metadata for a server to a file that can be sourced.

Project: [https://github.com/sstarcher/aws-env]
(https://github.com/sstarcher/aws-env)

Docker image: [https://registry.hub.docker.com/u/sstarcher/aws-env/]
(https://registry.hub.docker.com/u/sstarcher/aws-env/)



```docker run -v /etc/aws:/etc/aws sstarcher:aws-env:2.0```

* /etc/aws
```
AWS_INSTANCE_ID=i-xxxxxxxx
AWS_AVAILABLITY_ZONE=us-east-xx
AWS_REGION=us-east-1
AWS_ACCOUNT_ALIAS=account_alias
AWS_TAG_MYTAG=myvalue
AWS_TAG_NAME=ServerTag
```