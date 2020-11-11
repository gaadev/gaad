#!/bin/bash
#定义全局安装路径 后续按照所有的脚本规划，放置到别的地方
install_path=$(dirname $(readlink -f "$0"))
source $install_path/log.sh
source $install_path/env.sh
#定义目录
NODEJS_PATH=$install_path/nodejs

function check_env_by_cmd_v() {
  command -v $1 >/dev/null 2>&1 && (error "Installed ##$1## command,Are you want to update?" && exit 1)
}

function check_cpu_type() {
  kernel_type=$(arch)
  if [[ ${kernel_type} =~ 'x86_64' || ${kernel_type} =~ 'amd64' ]]; then
    NODE_FILE_NAME="node-v$NODE_VERSION-linux-x64"
  elif [[ ${kernel_type} =~ 'arm' ]]; then
    NODE_FILE_NAME="node-v$NODE_VERSION-linux-arm64"
  else
    error "Not found file;cpu type: $kernel_type" && exit 1
  fi
}

function download_and_install() {
  NODE_FILE_SUFFIX=".tar.gz"
  wget https://nodejs.org/dist/$NODE_VERSION/$NODE_FILE_NAME$NODE_FILE_SUFFIX
  tar -zxvf $NODE_FILE_NAME$NODE_FILE_SUFFIX
  mv $NODE_FILE_NAME $NODEJS_PATH
  #检查是否已添加到环境变量
  echo "" >$install_path/env.sh
  echo '#!/bin/bash' >>$install_path/env.sh
  echo 'export NODE_PATH='$NODEJS_PATH >>$install_path/env.sh
  echo 'export PATH=$PATH:$NODE_PATH/bin' >>$install_path/env.sh
  succss "The installation was successful for Nodejs"
}

if [ $# -ne 2 ]; then
  operation_type='install' #安装类型
  NODE_VERSION=14.15.0     #默认版本
else
  operation_type=$1 #安装类型
  NODE_VERSION=$2
fi

cd $install_path

info "nodejs: operation: $operation_type,version: $NODE_VERSION"
#如果是选择安装node,则检查node是否存在，已存在，则直接退出
if [[ $operation_type == "install" ]]; then
  check_env_by_cmd_v node
fi

#创建nodejs目录
rm -rf $NODEJS_PATH
mkdir $NODEJS_PATH
#检查cpu内核
check_cpu_type
#下载文件
download_and_install
