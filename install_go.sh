#!/bin/bash

# Define the version of Go to install
GO_VERSION="1.23.3"

# Determine the operating system and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Adjust architecture naming for Go
if [ "$ARCH" == "x86_64" ]; then
  ARCH="amd64"
elif [[ "$ARCH" == "armv"* ]]; then
  ARCH="armv6l"
elif [ "$ARCH" == "aarch64" ]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Construct the download URL
GO_TARFILE="go${GO_VERSION}.${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://go.dev/dl/${GO_TARFILE}"

# Download the Go tarball
echo "Downloading Go ${GO_VERSION} for ${OS}/${ARCH}..."
wget -q $DOWNLOAD_URL -O /tmp/$GO_TARFILE

if [ $? -ne 0 ]; then
  echo "Failed to download Go tarball."
  exit 1
fi

# Remove any previous Go installation
sudo rm -rf /usr/local/go

# Extract the tarball to /usr/local
echo "Installing Go to /usr/local..."
sudo tar -C /usr/local -xzf /tmp/$GO_TARFILE

# Clean up
rm /tmp/$GO_TARFILE

# Set up Go environment variables
echo "Setting up Go environment variables..."
if ! grep -q 'export PATH=$PATH:/usr/local/go/bin' ~/.profile; then
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
fi

if ! grep -q 'export GOPATH=$HOME/go' ~/.profile; then
  echo 'export GOPATH=$HOME/go' >> ~/.profile
fi

# Apply the changes to the current session
source ~/.profile

# Verify the installation
echo "Verifying Go installation..."
go version

if [ $? -ne 0 ]; then
  echo "Go installation failed."
  exit 1
else
  echo "Go ${GO_VERSION} installed successfully."
fi

snap install go --classic