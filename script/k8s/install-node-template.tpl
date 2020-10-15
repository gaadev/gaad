#!/usr/bin/env bash



set +e
set -o noglob

bold=; underline=; reset=; red=; green=; white=; tan=; blue=;
#echo "the value of TERM is ${TERM}"
if test ! dumb = "${TERM}" ; then
        bold=$(tput bold)
        underline=$(tput sgr 0 1)
        reset=$(tput sgr0)
        red=$(tput setaf 1)
        green=$(tput setaf 76)
        white=$(tput setaf 7)
fi

function underline() { printf "${underline}${bold}%s${reset}\n" "$@"
}
function h1() { printf "\n${bold}${blue}%s${reset}\n" "$@"
}
function h2() { printf "\n${bold}${white}%s${reset}\n" "$@"
}
function debug() { printf "${white}%s${reset}\n" "$@"
}
function info() { printf "${white}➜ %s${reset}\n" "$@"
}
function success() { printf "${green}✔ %s${reset}\n" "$@"
}
function error() { printf "${red}✖ %s${reset}\n" "$@"
}
function warn() { printf "${tan}➜ %s${reset}\n" "$@"
}
function bold() { printf "${bold}%s${reset}\n" "$@"
}
function note() { printf "\n${underline}${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$@"
}
function hlog() { echo -e "\033[32;32m $1 \033[0m \n"
}
set -e
set +o noglob

function uninstall_k8s() {
  #hash判断命令是否存在
  hash docker >/dev/null 2>&1 && systemctl restart docker
  hash kubectl >/dev/null 2>&1 && kubectl delete cm kubeadm-config -n kube-system
  hash kubeadm >/dev/null 2>&1 && kubeadm reset -f
  modprobe -r ipip
  lsmoda >/dev/null 2>&1
  rm -rf ~/.kube/
  rm -rf /etc/kubernetes/
  rm -rf /etc/systemd/system/kubelet.service.d
  rm -rf /etc/systemd/system/kubelet.service
  rm -rf /usr/bin/kube*
  rm -rf /var/lib/cni/
  rm -rf /var/lib/kubelet/*
  rm -rf /etc/cni
}

function check_linux_system() {
  hlog "检测linux系统类型"
  linux_version=$(cat /etc/redhat-release)
  if [[ ${linux_version} =~ "CentOS" ]]; then
    info "系统为 ${linux_version}"
  else
    info "系统不是CentOS,该脚本只支持CentOS环境"
    exit 1
  fi
}

function set_hostname() {
  hlog "设置主机名: $1"
  if [ -n "$1" ]; then
    grep $1 /etc/hostname &&
      info "主机名已设置，退出设置主机名步骤" && return

    hostnamectl set-hostname $1
    echo "$1" >/etc/hostname
    echo "127.0.0.1 $1" >>/etc/hosts
  else
    info "主机名不存在，退出" && exit 1
  fi
}

function install_docker() {
  hlog "开始安装docker"
  #先清理镜像
  sudo rm -rf /etc/yum.repos.d/*docker* /etc/yum.repos.d/*kubernetes*
  #step1: 卸载docker相关
  sudo yum remove -y docker \
    docker-client \
    docker-client-latest \
    docker-common \
    docker-latest \
    docker-latest-logrotate \
    docker-logrotate \
    docker-engine
  #step2: 按照yum-utils
  sudo yum install -y yum-utils

  #step3: 配置镜像源为阿里
  sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

  #step4: 更新yum
  sudo yum makecache fast
  #step5: 按照docker-ce
  sudo yum install -y docker-ce docker-ce-cli containerd.io
  #step6:启动docker
  sudo systemctl start docker
  #step7: 测试是否成功
  sudo docker run hello-world

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
  systemctl stop firewalld
  systemctl disable firewalld
  #关闭交换内存
  exist=$(cat /etc/sysctl.conf | grep "vm.swappiness = 0")

  if test -z "$exist"; then
    echo "vm.swappiness = 0" >>/etc/sysctl.conf
  fi
  swapoff -a && swapon -a && swapoff -a
  sysctl -p
}

function install_kubeadm_kubectl_kubelet() {
  hlog 安装 kubeadm kubectl kubelet
  version=$1
  cat <<EOF >/etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
EOF

  yum makecache fast
  yum -y remove kubelet kubeadm kubectl
  if test -z $version; then
    #为空使用默认版本
    version=$(yum list kubelet kubeadm kubectl | grep kubeadm | awk '{print $2}')
  fi
  yum install -y kubeadm-$version kubectl-$version kubelet-$version --disableexcludes=kubernetes && systemctl enable --now kubelet

}
uninstall_k8s
init_env
check_linux_system
set_hostname ?HOSTNAME
install_docker
install_kubeadm_kubectl_kubelet

?JOIN_COMMAND