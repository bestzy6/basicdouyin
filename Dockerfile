FROM golang:1.18

RUN echo "deb [check-valid-until=no] http://archive.debian.org/debian jessie-backports main" > /etc/apt/sources.list.d/jessie-backports.list
# As suggested by a user, for some people this line works instead of the first one. Use whichever works for your case
# RUN echo "deb [check-valid-until=no] http://archive.debian.org/debian jessie main" > /etc/apt/sources.list.d/jessie.list
RUN sed -i '/deb http:\/\/deb.debian.org\/debian jessie-updates main/d' /etc/apt/sources.list

#更新源
RUN apt-get -o Acquire::Check-Valid-Until=false update

#安装ffmpeg
RUN apt-get -y --force-yes install yasm ffmpeg


# 环境变量，设置goproxy
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"
# 开放端口
EXPOSE 3000
# 设置工作区
WORKDIR /bastictiktok
# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# 将代码复制到容器中
COPY . .
# 编译程序为二进制可执行文件app
RUN go build -o app .
# 运行程序
ENTRYPOINT ["./app"]