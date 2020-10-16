#!/bin/bash

source log.sh

function check_linux_system() {
  hlog "检测linux系统类型"
  #由于Ubuntu和CentOS系统文件不同，故可通过文件进行区分：CentOS:/etc/redhat-release  Ubuntu:/etc/lsb-release
  filePath="/etc/redhat-release"
  if [ -f "$filePath" ]; then
    linux_version='CentOS'
  else
    filePath="/etc/lsb-release"
    if [ -f "$filePath" ]; then
      linux_version='Ubuntu'
    else
      info '无法识别linux系统类型'
      exit 1
    fi
  fi
  info "系统为 ${linux_version}"
}

function config_yum_source() {
  if [[ ${linux_version} =~ "CentOS" ]]; then
    config_yum_source_centos
  else
    config_yum_source_ubuntu
  fi
}

#配置yum源
function config_yum_source_centos() {
  sleep 5
  sudo rm -rf /etc/yum.repos.d/*docker* \
  /etc/yum.repos.d/*epel* /etc/yum.repos.d/*elrepo*
  wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
  wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo
  #配置docker 镜像源
   #配置docker 镜像源
  wget -O /etc/yum.repos.d/docker-ce.repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
#  yum update -y
  yum makecache fast
}

#配置ubuntu源为阿里源
function config_yum_source_ubuntu() {
  sleep 5
  #将源指定到阿里源
  sudo sed -i "s/cn.archive.ubuntu.com/mirrors.aliyun.com/g" /etc/apt/sources.list
  sudo sed -i "s/security.ubuntu.com/mirrors.aliyun.com/g" /etc/apt/sources.list
  sudo apt update -y
  sudo apt -y upgrade
}

function install_docker() {
  if [[ ${linux_version} =~ "CentOS" ]]; then
    install_docker_centos
  else
    install_docker_ubuntu
  fi
}

function install_docker_ubuntu() {
  info "sleep 25 防止频繁刷aliyun镜像而造成下载速度缓慢"
  sleep 25
  hlog "Ubuntu开始安装docker"
  #step1: 卸载docker相关
  sudo apt-get remove  docker* containerd  runc || :
  #step2: 安装必要的一些系统工具
#  sudo apt-get update -y
  sudo apt-get -y install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
  # step 3: 安装GPG证书
  curl -fsSL http://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
  #step4: 写入软件源信息
  sudo  add-apt-repository "deb [arch=amd64] http://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
  #step5: 更新并安装 Docker-CE
  sudo apt-get -y update
  version=$1
  if test -n "$version"; then
#    sudo  apt-get install -y docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io
  else
    sudo apt-get -y install docker-ce
  fi

  #step6:启动docker
  sudo systemctl start docker

  #step7:可选，镜像加速
  sudo mkdir -p /etc/docker
  sudo tee /etc/docker/daemon.json <<-'EOF'
	{
		 "registry-mirrors": ["https://tsuz7ym7.mirror.aliyuncs.com"]
	}
EOF
  sudo systemctl daemon-reload
  sudo systemctl restart docker
  systemctl enable docker.service
}

function install_docker_centos() {
  info "sleep 25 防止频繁刷aliyun镜像而造成下载速度缓慢"
  sleep 25
  hlog "CentOS开始安装docker"
  #step1: 卸载docker相关
  sudo yum remove -y containerd.io docker*  || :
  #step5: 按照docker-ce
  #yum list docker-ce --showduplicates|sort -r #可以查看镜像
  version=$1
  if test -n "$version"; then
    sudo yum install -y docker-ce-cli-$version
    sudo yum install -y docker-ce-$version
  else
    sudo yum install -y docker-ce
  fi
  #step6:启动docker
  sudo systemctl start docker
  #step8:可选，镜像加速
  sudo mkdir -p /etc/docker
  sudo tee /etc/docker/daemon.json <<-'EOF'
	{
		 "registry-mirrors": ["https://tsuz7ym7.mirror.aliyuncs.com"]
	}
EOF
  sudo systemctl daemon-reload
  sudo systemctl restart docker
  systemctl enable docker.service
}

function init_env() {
  hlog "关闭防火墙，关闭swap"
  #初始化顺序
  #关闭防火墙
  if [[ ${linux_version} =~ "CentOS" ]]; then
    sudo systemctl stop firewalld
    sudo systemctl disable firewalld
  else
    sudo ufw disable
  fi
}

check_linux_system
init_env
config_yum_source
install_docker $DOCKER_VERSION
