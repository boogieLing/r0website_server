FROM ubuntu:18.04 AS os_file_apply
LABEL maintainer="r0"

# 升级 apt-get
RUN sed -i "s@/archive.ubuntu.com/@/mirrors.tuna.tsinghua.edu.cn/@g" /etc/apt/sources.list \
    && sed -i "s@/security.ubuntu.com/@/mirrors.ustc.edu.cn/@g" /etc/apt/sources.list\
    && rm -Rf /var/lib/apt/lists/* \
    && apt-get update --fix-missing -o Acquire::http::No-Cache=True \
    && apt-get install libcurl4 openssl -y

# 复制当前项目所有文件
COPY ["./", "/home/r0website_server"]
RUN true
COPY ["go1.17.11.linux-amd64.tar.gz", "/tmp/"]
RUN true

FROM os_file_apply AS server_apply
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

WORKDIR /home/r0website_server

# 安装 golang1.17.11，编译
RUN set -x; tar -zxvf /tmp/go1.17.11.linux-amd64.tar.gz -C /usr/local/\
    && echo "export GOROOT=/usr/local/go" >> ~/.bashrc\
    && echo "export GOBIN=\$GOROOT/bin" >> ~/.bashrc\
    && echo "export PATH=\$GOROOT/bin:\$PATH" >> ~/.bashrc\
    && /bin/bash -c "source ~/.bashrc" \
    && rm -rf /tmp/go1.17.11.linux-amd64.tar.gz \
    && /usr/local/go/bin/go build -o ./r0Website-server ./main/main.go

# 暴露server端口
EXPOSE 8202
ENTRYPOINT ["./r0Website-server"]
# docker exec -it r0website_server /bin/bash