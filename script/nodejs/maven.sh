#!/bin/bash

h1 "install maven begin"

install_path=$(dirname $(readlink -f "$0"))
source $install_path/log.sh
source $install_path/env.sh

if [ $# -ne 2 ]; then
  operation_type='install' #安装类型
  MAVEN_VERSION=3.6.3      #默认版本
else
  operation_type=$1 #安装类型
  MAVEN_VERSION=$2
fi

function check_env_by_cmd_v() {
  command -v $1 >/dev/null 2>&1 && (error "Installed ##$1## command,Are you want to update?" && exit 1)
}
#https://mirror.bit.edu.cn/apache/maven/maven-3/3.6.3/binaries/apache-maven-3.6.3-bin.tar.gz
info "check maven operation: $operation_type,version: $MAVEN_VERSION"
#如果是选择安装MAVEN,则检查MAVEN是否存在，已存在，则直接退出
if [[ $operation_type == "install" ]]; then
  check_env_by_cmd_v mvn
fi

if [[ ${MAVEN_VERSION:0:1} -ne '3' ]]; then
  error "maven version must be more than 3." && exit 1
fi
info "download maven file"
MAVEN_DOWNLOAD_FILE_NAME="apache-maven-$MAVEN_VERSION-bin.tar.gz"
MAVEN_FOLDER_NAME="apache-maven-$MAVEN_VERSION"
#删除与maven相关的文件
rm -rf $install_path/apach-maven.*
#下载文件
wget https://mirror.bit.edu.cn/apache/maven/maven-3/$MAVEN_VERSION/binaries/$MAVEN_DOWNLOAD_FILE_NAME
tar -zxvf $MAVEN_DOWNLOAD_FILE_NAME

info "Add environment variables for Maven"

echo "" >$install_path/env.sh
echo '#!/bin/bash' >>$install_path/env.sh
echo 'export MAVEN_HOME='$install_path/$MAVEN_FOLDER_NAME >>$install_path/env.sh
echo 'export PATH=$PATH:$MAVEN_HOME/bin' >>$install_path/env.sh

success "The installation was successful for Maven"
