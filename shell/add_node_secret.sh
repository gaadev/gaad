#!/bin/bash
if test ! -d ~/.devops; then
  mkdir ~/.devops
fi
if test ! -f ~/.devops/deploy-target; then
touch ~/.devops/deploy-target
fi
exist=$(cat ~/.devops/deploy-target | grep "$1=")
if test -z "$exist" ; then
 echo $2 >> ~/.devops/deploy-target;
fi