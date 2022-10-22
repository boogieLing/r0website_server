#!/bin/bash

files=("go1.17.11.linux-amd64.tar.gz" "mongodb-linux-x86_64-ubuntu1804-5.0.9.tgz")
files_trans=("\"go1.17.11.linux\-amd64.tar.gz\"" "\"mongodb\-linux-x86_64\-ubuntu1804\-5.0.9.tgz\"")
files_from_wget=("https:\/\/golang.google.cn\/dl\/go1.17.11.linux\-amd64.tar.gz" "https:\/\/fastdl.mongodb.org\/linux\/mongodb\-linux\-x86_64\-ubuntu1804\-5.0.9.tgz")
wget_files=()
copy_files=()
copy_flag="false"
copy_apply_atr=""
wget_flag="false"
wget_apply_str=""

# 读取文件列表，不存在的文件就加入wget队列，存在就加入copy队列
for (( i=0;i<${#files[@]};i++ ))
do
  if [ ! -f ${files[$i]} ];then
    echo -e "\033[41;30m XXX ${files[$i]} not exist ?? XXX \033[0m"
    wget_flag="true"
    wget_files[${#wget_files[*]}]=${files_from_wget[$i]}
  else
    echo -e "\033[46;30m √√√ ${files[$i]} exist !! √√√ \033[0m"
    copy_flag="true"
    copy_files[${#copy_files[*]}]=${files_trans[$i]}
  fi
done

#如果存在需要wget的文件，需要安装wget
if [ "${wget_flag}"x == "true"x ];then
  wget_apply_str="RUN apt\-get install \-\-assume-yes apt\-utils \&\& apt\-get install wget \-y"
  for wget in ${wget_files[@]}
  do
    connect_str_prefix=" \&\& wget "
    connect_str_suffix=" \-p \/tmp "
    wget_apply_str=$wget_apply_str$connect_str_prefix$wget$connect_str_suffix
  done
else
  wget_apply_str=""
fi
sed -i "s/\${LOAD_FILE_WGET}/$wget_apply_str/g" ./Dockerfile

#如果存在需要COPY的文件，提供docker的COPY指令
if [ "${copy_flag}"x == "true"x ];then
  connect_str_prefix="COPY ["
  for copy in ${copy_files[@]}
  do
    connect_str=", "
    copy_apply_atr=$copy_apply_atr$copy$connect_str
  done
  connect_str_suffix="\"\/tmp\/\"]"
  copy_apply_atr=$connect_str_prefix$copy_apply_atr$connect_str_suffix
else
  copy_apply_atr=""
fi
sed -i "s/\${LOAD_FILE_COPY}/$copy_apply_atr/g" ./Dockerfile

# cat ./Dockerfile
# 构建镜像
docker build -t ubuntu18-mongo5 . &&
# 以脱离模式运行容器
docker run --name ubuntu18-mongo5 -d -p 27017:27017 ubuntu18-mongo5:latest