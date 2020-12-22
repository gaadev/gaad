#!/bin/bash
git_download_url_prefix=https://mirrors.edge.kernel.org/pub/software/scm/git/
git_tar_name=git-2.29.2.tar.xz
git_dir_name=git-2.29.2

function download() {
  wget "$git_download_url_prefix$git_tar_name"
  xz -d $git_tar_name
  tar -xvf $git_tar_name
  cd $git_dir_name
  make
}

download
