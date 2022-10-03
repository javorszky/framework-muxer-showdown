FROM golang:1.19

ENV GOCACHE /tmp/.cache

RUN apt-get update -y && apt-get upgrade -y && apt-get install jq -y

## Install gci
RUN go install github.com/daixiang0/gci@v0.8.0

# Install golang-ci lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

# Install mockery
WORKDIR /temp

RUN wget -q -O - https://api.github.com/repos/vektra/mockery/releases/latest \
    | jq '.assets | .[] | .browser_download_url | select(test(".*Linux_arm.*"))' \
    | tr -d '"' \
    | xargs -n1 wget

RUN tar -xzvf *.tar.gz && chmod +x mockery && mv mockery /usr/local/bin/

WORKDIR /app
