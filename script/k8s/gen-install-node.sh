#!/bin/bash
function gen_install_node() {
  set -e
  join_token=$(kubeadm token create --print-join-command)
  if [ ! -f "count.txt" ]; then
    echo 1 >count.txt
  fi
  num=$(cat count.txt)
  echo $((num + 1)) >count.txt
  hostname="k8s-node"$num

  \cp install-node-template.tpl install-node#genrated.sh
  chmod +x install-node#genrated.sh
  sed -i "s#?HOSTNAME#${hostname}#g" install-node#genrated.sh
  sed -i "s#?JOIN_COMMAND#${join_token}#g" install-node#genrated.sh
}
