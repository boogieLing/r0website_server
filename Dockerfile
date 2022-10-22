FROM ubuntu:18.04 AS os_file_apply
LABEL maintainer="r0"
# 升级 apt-get
RUN sed -i "s@/archive.ubuntu.com/@/mirrors.tuna.tsinghua.edu.cn/@g" /etc/apt/sources.list \
    && sed -i "s@/security.ubuntu.com/@/mirrors.ustc.edu.cn/@g" /etc/apt/sources.list\
    && rm -Rf /var/lib/apt/lists/* \
    && apt-get update --fix-missing -o Acquire::http::No-Cache=True \
    && apt-get install libcurl4 openssl -y

# 以wget的方式准备文件到/tmp
${LOAD_FILE_WGET}

# 以COPY的方式准备文件到/tmp
${LOAD_FILE_COPY}


FROM os_file_apply AS go_apply
COPY --from=os_file_apply /tmp/go1.17.11.linux-amd64.tar.gz /tmp/
# 安装 golang1.17.11
RUN set -x; tar -zxvf /tmp/go1.17.11.linux-amd64.tar.gz -C /usr/local/\
    &&  echo "export GOROOT=/usr/local/go" >> ~/.bashrc\
    &&  echo "export GOBIN=\$GOROOT/bin" >> ~/.bashrc\
    &&  echo "export PATH=\$GOROOT/bin:\$PATH" >> ~/.bashrc\
    &&  /bin/bash -c "source ~/.bashrc" \
    && rm -rf /tmp/go1.17.11.linux-amd64.tar.gz

FROM go_apply AS mongo_apply
COPY --from=os_file_apply /tmp/mongodb-linux-x86_64-ubuntu1804-5.0.9.tgz /tmp/
# 安装 mongo
COPY ["mongod.conf", "/tmp/"]
RUN tar -zxvf /tmp/mongodb-linux-x86_64-ubuntu1804-5.0.9.tgz -C /usr/local/ \
    && mv /usr/local/mongodb-linux-x86_64-ubuntu1804-5.0.9 /usr/local/mongodb \
    && mv /tmp/mongod.conf /usr/local/mongodb/mongod.conf\
    && echo "export PATH=\$PATH:/usr/local/mongodb/bin" >> ~/.bashrc \
    && /bin/bash -c "source ~/.bashrc"\
    && mkdir -p /usr/local/mongodb/data/db \
    && mkdir -p /usr/local/mongodb/logs \
    && touch /usr/local/mongodb/logs/mongod.log \
    && rm -rf /tmp/mongodb-linux-x86_64-ubuntu1804-5.0.9.tgz

#mongodb的web端口
EXPOSE 28017
#连接端口
EXPOSE 27017

ENTRYPOINT ["/usr/local/mongodb/bin/mongod", "-f", "/usr/local/mongodb/mongod.conf"]

# 流程：构建->脱离模式运行容器->即时设置权限

# 构建镜像
# docker build -t ubuntu18-mongo5 .
# 以脱离模式运行容器
# docker run --name ubuntu18-mongo5 -d -p 27017:27017 ubuntu18-mongo5:latest
# 以交互终端运行容器
# docker run --name ubuntu18-mongo5 -it -p 27017:27017 ubuntu18-mongo5:latest /bin/bash

# docker attach ubuntu18-mongo5
# docker exec -it ubuntu18-mongo5 /bin/bash

# 停止、删除容器和镜像
# docker container stop ubuntu18-mongo5; docker container rm ubuntu18-mongo5 && docker image rm ubuntu18-mongo5