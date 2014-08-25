# pingo

Pings TCP services and stores results in [Amazon CloudWatch](http://aws.amazon.com/cloudwatch/).

## Installation

```
go get github.com/robinjmurphy/pingo/cmd/pingo
```

## Usage

Type `pingo -h` for full usage information.

## Examples

Ping `google.com` on port `80` every 10 seconds:

```
pingo google.com 80 --interval 10
```

Ping `twitter.com` on port `443` every 30 seconds, and publish the results to CloudWatch:

```
pingo twitter.com 443 --interval 30 --publish --region eu-west-1
```

## CloudWatch

To store results in CloudWatch, ensure the `AWS_ACCESS_KEY` and `AWS_SECRET_KEY` environment variables are set for an IAM user that has permission to `PutMetricData`. Metrics are stored in the `pingo` namespace with a value of `2` for successful pings and `1` for failed pings.