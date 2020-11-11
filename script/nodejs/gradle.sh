#!/bin/bash

install_path=$(dirname $(readlink -f "$0"))

source $install_path/log.sh
source $install_path/env.sh

h1 "install gradle begin"




if [ $# -ne 2 ]; then
  operation_type='install' #安装类型
  GRADLE_VERSION=6.7      #默认版本
else
  operation_type=$1 #安装类型
  GRADLE_VERSION=$2
fi

function check_env_by_cmd_v() {
  command -v $1 >/dev/null 2>&1 && (error "Installed ##$1## command,Are you want to update?" && exit 1)
}
#https://services.gradle.org/distributions/gradle-6.7-all.zip
info "check gradle operation: $operation_type,version: $GRADLE_VERSION"
#如果是选择安装MAVEN,则检查MAVEN是否存在，已存在，则直接退出
if [[ $operation_type == "install" ]]; then
  check_env_by_cmd_v gradle
fi

info "download gradle file"
GRADLE_DOWNLOAD_FILE_NAME="gradle-$GRADLE_VERSION-all.zip"
GRADLE_FOLDER_NAME="gradle-$GRADLE_VERSION"
#删除与maven相关的文件
rm -rf $install_path/gradle*.zip
#下载文件
wget https://services.gradle.org/distributions/$GRADLE_DOWNLOAD_FILE_NAME
unzip $GRADLE_DOWNLOAD_FILE_NAME

info "Add environment variables for Gradle"

echo "" >$install_path/env.sh
echo '#!/bin/bash' > $install_path/env.sh
echo 'export GRADLE_HOME='$install_path/$GRADLE_FOLDER_NAME >>$install_path/env.sh
echo 'export PATH=$PATH:$GRADLE_HOME/bin' >>$install_path/env.sh

success "The installation was successful for Gradle"



