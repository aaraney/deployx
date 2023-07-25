#!/bin/bash

install_with_plugin=false
install_without_plugin=false

while getopts ":yn" opt; do
  case $opt in
  n)
    install_without_plugin=true
    ;;
  y)
    install_with_plugin=true
    ;;
  \?)
    echo "Invalid option: -$OPTARG" >&2
    ;;
  esac
done

url="https://api.github.com/repos/rei-x/deployx/releases/latest"

os=$(uname -s)
arch=$(uname -m)

if [ "$arch" = "x86_64" ]; then
  arch="amd64"
elif [ "$arch" = "aarch64" ]; then
  arch="arm64"
elif [ "$arch" = "armv7l" ]; then
  arch="arm"
elif [ "$arch" = "armv6l" ]; then
  arch="arm"
elif [ "$arch" = "armv5l" ]; then
  arch="arm"
elif [ "$arch" = "i386" ]; then
  arch="386"
else
  echo "Unsupported architecture: $arch"
  exit 1
fi

if [ "$os" = "Linux" ]; then
  os="linux"
elif [ "$os" = "Darwin" ]; then
  os="darwin"
elif [ "$os" = "Windows" ]; then
  os="windows"
else
  echo "Unsupported OS: $os"
  exit 1
fi

echo "Downloading deployx binary for $os $arch"

curl --progress-bar --show-error -sL $url | grep "browser_download_url.*deployx.*$os.*$arch" | cut -d '"' -f 4 | wget -qi -

# check md5sum

echo "Checking md5sum of file"

md5sum=$(md5sum deployx*.tar.gz | cut -d ' ' -f 1)
md5sum_file=$(cat deployx*.tar.gz.md5 | cut -d ' ' -f 1)

if [ "$md5sum" != "$md5sum_file" ]; then
  echo "md5sum does not match"
  exit 1
fi

echo "MD5sum matches"

tar -xzf deployx*.tar.gz

# ask if user wants to install as docker plugin

if [ "$install_with_plugin" = false ] && [ "$install_without_plugin" = false ]; then
  echo "Do you want to install deployx as a docker plugin? (y/n)"
  read install_as_plugin
else
  install_as_plugin="y"
fi

if [ "$install_with_plugin" = true ]; then
  install_as_plugin="y"
fi

if [ "$install_without_plugin" = true ]; then
  install_as_plugin="n"
fi

if [ "$install_as_plugin" = "y" ]; then
  install_path="${HOME}"/.docker/cli-plugins

  echo "Installing deployx to $install_path"

  mkdir -p "${install_path}" 2>/dev/null || sudo mkdir -p "${install_path}"

  if [ $? -ne 0 ]; then
    echo "Failed to create directory $install_path"
    exit 1
  fi

  # move deployx to installation path
  mv deployx "${install_path}"/docker-deployx 2>/dev/null || sudo mv deployx "${install_path}"/docker-deployx

  if [ $? -ne 0 ]; then
    echo "Failed to move deployx to $install_path"
    exit 1
  fi

  echo "You can now use docker deployx"
else
  echo "You can now use ./deployx"
fi

# clean

rm deployx*.tar.gz*
