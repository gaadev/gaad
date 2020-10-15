#!/bin/bash
source log.sh
source config
echo "xxxj"

function uninstall_k8s() {
        echo "ww"
  #hash判断命令是否存在
  hash docker >/dev/null 2>&1 && systemctl restart docker
  hash kubectl >/dev/null 2>&1 && kubectl delete cm kubeadm-config -n kube-system
  hash kubeadm >/dev/null 2>&1 && kubeadm reset -f
  modprobe -r ipip
  lsmoda 
#  rm -rf ~/.kube/
#  rm -rf /etc/kubernetes/
#  rm -rf /etc/systemd/system/kubelet.service.d
#  rm -rf /etc/systemd/system/kubelet.service
#  rm -rf /usr/bin/kube*
#  rm -rf /var/lib/cni/
#  rm -rf /var/lib/kubelet/*
#  rm -rf /etc/cni
        echo "eee"
}
uninstall_k8s
