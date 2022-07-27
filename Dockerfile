#FROM golang:1.17 as builder
#RUN apt install make
#
#WORKDIR /code/
#
#ADD go.mod .
#COPY Makefile .
#COPY . .
#
#RUN go env -w GO111MODULE=on
#RUN go env -w GOPROXY=https://goproxy.cn,direct
#
#RUN make dependencies
#RUN go mod download
#RUN make build


FROM ubuntu:20.04
WORKDIR /app
#COPY --from=builder /code/bin/dlocator ./dlocator
#COPY --from=builder /code/configs .

RUN mkdir "configs"
COPY ./configs/config.yml ./configs/
COPY bin/dlocator ./dlocator
ENTRYPOINT ["/app/dlocator"]

EXPOSE 6379 9000

