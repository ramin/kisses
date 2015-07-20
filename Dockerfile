from ubuntu:latest

COPY ./kisses-amd64 /home/kisses

RUN apt-get update
RUN apt-get install -y ca-certificates

ENV AWS_ACCESS_KEY <add-your-own>
ENV AWS_SECRET_KEY <add-your-own>

EXPOSE 443
# CMD /home/kisses -stream nord_customer_activities_stream -region us-west-2