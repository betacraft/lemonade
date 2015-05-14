# Using wheezy from the official golang docker repo 
FROM golang:1.4.2-wheezy
MAINTAINER Akshay Deo <akshay@rainingclouds.com>

# Setting up working directory
WORKDIR /go/src/github.com/rainingclouds/lemonades
Add . /go/src/github.com/rainingclouds/lemonades/

# Get godeps from main repo
RUN go get github.com/tools/godep

# Restore godep dependencies
RUN godep restore

# Install
RUN go install github.com/rainingclouds/lemonades

# Setting up environment variables
ENV ENV prod
ENV RESSY_PG_URL host=aasprs0iuoe3eo.cqoummlog10q.ap-southeast-1.rds.amazonaws.com user=akshay password=amdeo3116 dbname=ebdb sslmode=require
ENV LEMN_MG_URI mongodb://web-app:web-app123456@ds031192-a0.mongolab.com:31192,ds031192-a1.mongolab.com:31192/lmnd?replicaSet=rs-ds031192
ENV LEMN_MG_DB_NAME lmnd
ENV AWS_ACCESS AKIAINECWOX2MEE4UOSA
ENV AWS_SECRET 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

EXPOSE 80

ENTRYPOINT ["/go/bin/lemonades"]