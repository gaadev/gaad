#!/bin/bash
source log.sh
source config

function uninstall_k8s() {
  #hash判断命令是否存在
  #hash docker >/dev/null 2>&1 && systemctl restart docker
  hash kubectl >/dev/null 2>&1 && kubectl delete cm kubeadm-config -n kube-system || :
  hash kubeadm >/dev/null 2>&1 && kubeadm reset -f
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

function config_yum_source() {
  sleep 3
  sudo rm -rf /etc/yum.repos.d/*docker* /etc/yum.repos.d/*kubernetes* \
    /etc/yum.repos.d/*epel* /etc/yum.repos.d/*elrepo*
  wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
  wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo
  #配置docker 镜像源
  sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
  cat <<EOF >/etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
EOF

  yum makecache fast
}

function install_docker() {
  info "sleep 25 防止频繁刷aliyun镜像而造成下载速度缓慢"
  sleep 25
  hlog "开始安装docker"
  #step1: 卸载docker相关
  sudo yum remove -y docker* containerd.io
  #step2: 按照yum-utils
  sudo yum install -y yum-utils

  #step4: 更新yum
  sudo yum makecache fast
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
  info "sleep 20 防止频繁刷aliyun镜像而造成下载速度缓慢"
  sleep 20
  hlog 安装 kubeadm kubectl kubelet
  version=$1
  yum -y remove kubelet kubeadm kubectl
  if test -n "$version"; then
    yum install -y kubeadm-$version kubectl-$version kubelet-$version --disableexcludes=kubernetes && systemctl enable --now kubelet
  else
    #为空使用默认版本
    yum install -y kubeadm kubectl kubelet --disableexcludes=kubernetes && systemctl enable --now kubelet
  fi
}

function pull_k8s_images() {
  hlog "接取k8s所需镜像"
  set -e
  k8s_version=$1
  if test -z "$k8s_version"; then
    versions=$(kubeadm config images list)
  else
    versions=$(kubeadm config images list --kubernetes-version $k8s_version)
  fi
  #versions=`kubeadm config images list`
  apiserver_version_tmp=$(echo $versions | sed 's/ /\n/g' | grep k8s.gcr.io/kube-apiserver)
  pause_version_tmp=$(echo $versions | sed 's/ /\n/g' | grep k8s.gcr.io/pause)
  etcd_version_tmp=$(echo $versions | sed 's/ /\n/g' | grep k8s.gcr.io/etcd)
  coredns_version_tmp=$(echo $versions | sed 's/ /\n/g' | grep k8s.gcr.io/coredns)
  apiserver_version=${apiserver_version_tmp#*:}
  pause_version=${pause_version_tmp#*:}
  etcd_version=${etcd_version_tmp#*:}
  coredns_version=${coredns_version_tmp#*:}

  #KUBE_VERSION=v1.16.3
  #KUBE_PAUSE_VERSION=3.1
  #ETCD_VERSION=3.3.15-0
  #CORE_DNS_VERSION=1.6.2

  KUBE_VERSION=$apiserver_version
  KUBE_PAUSE_VERSION=$pause_version
  ETCD_VERSION=$etcd_version
  CORE_DNS_VERSION=$coredns_version

  GCR_URL=k8s.gcr.io
  ALIYUN_URL=registry.cn-hangzhou.aliyuncs.com/google_containers

  images=(kube-proxy:${KUBE_VERSION}
    kube-scheduler:${KUBE_VERSION}
    kube-controller-manager:${KUBE_VERSION}
    kube-apiserver:${KUBE_VERSION}
    pause:${KUBE_PAUSE_VERSION}
    etcd:${ETCD_VERSION}
    coredns:${CORE_DNS_VERSION})

  # shellcheck disable=SC2068
  for imageName in ${images[@]}; do
    docker pull $ALIYUN_URL/$imageName
    docker tag $ALIYUN_URL/$imageName $GCR_URL/$imageName
    docker rmi $ALIYUN_URL/$imageName
  done

}

function init_k8s_master() {
  kubeadm init \
    --apiserver-advertise-address $1 \
    --kubernetes-version=$2 \
    --pod-network-cidr=10.244.0.0/16

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config
}

function config_network() {
  kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
}

uninstall_k8s
init_env
check_linux_system
set_hostname $MASTER_HOST
config_yum_source
install_docker $DOCKER_VERSION
install_kubeadm_kubectl_kubelet $K8S_VERSION
pull_k8s_images $K8S_VERSION
init_k8s_master $MASTER_IP
config_network
