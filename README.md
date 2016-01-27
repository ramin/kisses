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
./kisses -stream customer_activities -region us-west-2
```

# Docker

If you are so inclined

```
docker build -t kisses .
docker run -i -t kisses
```

and watch it fly...