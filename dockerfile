FROM golang:1.19-alpine AS build-dist
ENV GOPROXY='https://goproxy.cn,direct'

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -tags netcgo -installsuffix cgo -o /bin/app main.go


FROM nvidia/opencl:devel-ubuntu18.04 AS prod

ENV RUN_MODE='release'

COPY --from=build-dist /bin/app /bin/app

ADD source.list /etc/apt/sources.list

RUN apt-get update && apt-get install unzip wget libzip5  -y -q

WORKDIR /katago

RUN wget https://ghproxy.com/https://github.com/lightvector/KataGo/releases/download/v1.12.3/katago-v1.12.3-opencl-linux-x64.zip && \
    unzip katago-v1.12.3-opencl-linux-x64.zip && \
    cp katago /usr/bin/katago && \
    chmod +x /usr/bin/katago && \
    rm -rf katago-v1.12.3-opencl-linux-x64.zip

RUN rm -rf /katago

WORKDIR /project

#RUN wget https://media.katagotraining.org/uploaded/networks/models/kata1/kata1-b60c320-s7238447104-d3189358758.bin.gz

COPY kata1-b18c384nbt-s4975305984-d3174897359.bin.gz best.bin.gz
COPY analysis.cfg analysis.cfg

CMD ["/bin/app"]

EXPOSE 8080