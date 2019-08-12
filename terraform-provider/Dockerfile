ARG GO_VERSION=1.12
FROM golang:${GO_VERSION}-stretch

ARG TF_VERSION=0.12.5

# Configure Go
ENV GOLANG_VERSION=$GO_VERSION
ENV PATH /usr/local/go/bin:/usr/local/go:$PATH
ENV GOPATH /go
ENV GOBIN /usr/local/go


# Install Terraform
RUN apt-get update && \
    apt-get install -y curl zip
RUN wget --quiet https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip \
  && unzip terraform_${TF_VERSION}_linux_amd64.zip \
  && mv terraform /usr/bin \
  && rm terraform_${TF_VERSION}_linux_amd64.zip

WORKDIR /init
ARG PROVIDER_SRC_DIR=provider-src/
# install project dependencies ahead of time to prevent slow builds for code changes
ADD $PROVIDER_SRC_DIR/src/go.mod .
ADD $PROVIDER_SRC_DIR/src/go.sum .
RUN go mod download


# setup project workspace
ENV WORKSPACE_ROOT=/workspace/
WORKDIR $WORKSPACE_ROOT

# copy terraform provider code
WORKDIR $PROVIDER_SRC_DIR
ADD $PROVIDER_SRC_DIR .

# Run Unit Tests
RUN cd src      && \
    go test -v

# Build plugins
RUN mkdir build/                                    && \
    cd src/                                         && \
    VERSION=$(cat ../PROVIDER_VERSION.txt)          && \
    NAME=$(cat ../PROVIDER_NAME.txt)                && \
    ARTIFACT=terraform-provider-${NAME}_v${VERSION} && \
    echo "Building terraform provider $ARTIFACT"    && \
    go build -o ../build/$ARTIFACT

# Install plugins
ENV PLUGINS_DIR=/root/.terraform.d/plugins/linux_amd64/
RUN mkdir -p $PLUGINS_DIR      && \ 
    mv ./build/* $PLUGINS_DIR

# Copy terraform files
ARG TERRAFORM_SRC_DIR=terraform-src/
WORKDIR $WORKSPACE_ROOT/$TERRAFORM_SRC_DIR
ADD $TERRAFORM_SRC_DIR .

CMD bash