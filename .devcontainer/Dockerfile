ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-bullseye

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

ARG TERRAFORM_VERSION=1.1.8
ARG TFLINT_VERSION=0.35.0
ARG TERRAFORM_PROVIDER_AZUREDEVOPS_NAME=terraform-provider-azuredevops

RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils 2>&1 \
    # Verify git, process tools, lsb-release (common in install instructions for CLIs), unzip installed
    && apt-get -y install git procps lsb-release unzip \
    # Install Editor
    && apt-get install vim -y \
    # Install missing icu-devtools to be able to use liveshare. More info: https://stackoverflow.com/questions/58736399/how-to-get-vscode-liveshare-extension-working-when-running-inside-vscode-remote
    && apt update && apt install icu-devtools -y \
    # Install Terraform, tflint
    && mkdir -p /tmp/docker-downloads \
    && curl -sSL -o /tmp/docker-downloads/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip /tmp/docker-downloads/terraform.zip \
    && mv terraform /usr/local/bin \
    && curl -sSL -o /tmp/docker-downloads/tflint.zip https://github.com/wata727/tflint/releases/download/v${TFLINT_VERSION}/tflint_linux_amd64.zip \
    && unzip /tmp/docker-downloads/tflint.zip \
    && mv tflint /usr/local/bin \
    && cd ~ \
    && rm -rf /tmp/docker-downloads \
    # Install goreleaser
    && go install github.com/goreleaser/goreleaser@latest \
    #
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# Enable go modules
ENV GO111MODULE=on

# Install Go tools
RUN \
    # --> Delve for debugging
    go install github.com/go-delve/delve/cmd/dlv@v1.7.3 \
    # --> Go language server
    && go install golang.org/x/tools/gopls@latest \
    # --> GolangCI-lint
    && curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sed 's/tar -/tar --no-same-owner -/g' | sh -s -- -b $(go env GOPATH)/bin \
    #
    # Clean up
    && rm -rf /go/src/ && rm -rf /go/pkg

# Customize bash
RUN \
    # Git command prompt
    git clone https://github.com/magicmonty/bash-git-prompt.git ~/.bash-git-prompt --depth=1 \
    && echo "if [ -f \"$HOME/.bash-git-prompt/gitprompt.sh\" ]; then GIT_PROMPT_ONLY_IN_REPO=1 && source $HOME/.bash-git-prompt/gitprompt.sh; fi" >> "/root/.bashrc" \
    # Add alias
    && echo "alias tf=terraform" >> "/root/.bashrc"

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog
