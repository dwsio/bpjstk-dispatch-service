ARG BUILDER_IMAGE=golang:1.20-buster
ARG BASE_IMAGE=debian:buster-slim

# ############################ BUILDER IMAGE ############################
FROM ${BUILDER_IMAGE} AS builder

LABEL maintainer="dpti@bpjsketenagakerjaan.go.id"

# UPDATE BUILDER IMAGE
# RUN apt-get update && apt-get install -y xz-utils unzip

# CREATE A WORKING DIRECTORY
RUN mkdir /app
WORKDIR /app

#COPY SOURCE CODE
COPY . .

#BUILD BINARY FILE
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -o ./app ./tests/

########################### DISTRIBUTION IMAGE DEBIAN ############################
FROM ${BASE_IMAGE}

LABEL maintainer="dpti@bpjsketenagakerjaan.go.id"

# UPDATE DISTRIBUTION IMAGE
# RUN apt-get update && apt-get install -y xz-utils unzip
RUN apt-get update && apt-get install -y ca-certificates

# CLEAN UP
# RUN apt-get clean autoclean \
#     && apt-get autoremove --yes unzip wget \
#     && rm -rf /var/lib/{apt,dpkg,cache,log} \
#     && rm -rf /tmp/* /var/tmp/* \
#     && rm /var/log/lastlog /var/log/faillog

# SET TIMEZONE
ENV TZ="Asia/Jakarta"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ >/etc/timezone && dpkg-reconfigure -f noninteractive tzdata

#CREATE WORKDIR
RUN mkdir /usr/app
WORKDIR /usr/app

#COPY BINARY FILE FROM BUILDER
COPY --from=builder /app/app /usr/app/app

EXPOSE 5002 8002

ENTRYPOINT ["/usr/app/app"]
CMD [""]
