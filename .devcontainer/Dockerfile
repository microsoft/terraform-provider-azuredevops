FROM ubuntu:bionic
ARG GO_VERSION=1.13
ARG TERRAFORM_VERSION=0.12.13
RUN apt-get update && \
    apt-get install wget unzip gcc git vim -y && \
    wget "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -xvf "go${GO_VERSION}.linux-amd64.tar.gz" && \
    mv go /usr/local && \
    wget "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" && \
    unzip "./terraform_${TERRAFORM_VERSION}_linux_amd64.zip" -d /usr/local/bin/ && \
    rm "go${GO_VERSION}.linux-amd64.tar.gz" && \
    rm "terraform_${TERRAFORM_VERSION}_linux_amd64.zip" && \
    rm -rf /var/lib/apt/lists/*
ENV PATH="${PATH}":/usr/local/go/bin
RUN git clone https://github.com/microsoft/terraform-provider-azuredevops.git && \
    cd terraform-provider-azuredevops && \
    ./scripts/build.sh && \
    ./scripts/local-install.sh && \
    rm -rf ../terrform-provider-azuredevops
WORKDIR /workspaces/terraform-provider-azuredevops