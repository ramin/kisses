cli utility to watch a kinesis stream and pipe into other unix utlities

# Requirements

AWS_ACCESS_KEY and AWS_SECRET_KEY environment variables set, test these with

```
echo $AWS_ACCESS_KEY
echo $AWS_SECRET_KEY
```

# Flags

stream and region

```
./kisses -stream nord_customer_activities_stream -region us-east-1
```