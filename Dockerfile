FROM golang:1.22.2-bullseye

RUN apt-get update && \
    apt-get install lsb-release unzip -y

RUN go version

RUN apt-get update

################## Tooling prerequisites  ######################
# Azure Cli
# https://github.com/Azure/azure-cli/releases
ARG AZURE_CLI_VERSION=2.59.0
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list && \
    curl -L https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
    apt-get install apt-transport-https  && \
    apt-get update && apt-get install -y azure-cli=${AZURE_CLI_VERSION}-1~bullseye  && \
    az version

# terraform
# https://releases.hashicorp.com/terraform/
ARG TERRAFORM_VERSION=1.7.5
RUN curl -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip --output terraform.zip && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin && \
    terraform version

# k8s CLI
# You must use a kubectl version that is within one minor version difference of your cluster.
# For example, a v1.24 client can communicate with v1.23, v1.24, and v1.25 control planes.
ARG KUBECTL_VERSION=v1.28.5
RUN curl -LO https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl  && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl  && \
    kubectl version --client=true -o=json

# kubelogin CLI
ARG KUBELOGIN_VERSION=v0.1.1
ARG KUBELOGIN_SHA256=51938ead06738bce6bb8fa4418fe043a5106a7d8bad037a65b57170c4f3fba0e
RUN curl -LO https://github.com/Azure/kubelogin/releases/download/${KUBELOGIN_VERSION}/kubelogin-linux-amd64.zip  && \
    echo "${KUBELOGIN_SHA256} kubelogin-linux-amd64.zip" | sha256sum -c && \
    unzip kubelogin-linux-amd64.zip -d kubelogin && \
    chmod +x kubelogin && \
    mv kubelogin/bin/linux_amd64/kubelogin /usr/local/bin && \
    kubelogin --version

# helm
ARG HELM_VERSION=3.14.4
RUN curl -L  https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz --output helm.tar.gz && \
    tar xvzf helm.tar.gz && \
    mv linux-amd64/helm /usr/local/bin && \
    helm version

# copsctl
ARG COPSCTL_VERSION=0.12.3
RUN curl -LO https://github.com/conplementAG/copsctl/releases/download/v${COPSCTL_VERSION}/copsctl_Linux_x86_64.tar.gz && \
    tar xvzf copsctl_Linux_x86_64.tar.gz && \
    mv copsctl $GOPATH/bin && \
    copsctl --version

# sops
ARG SOPS_VERSION=v3.8.1
ARG SOPS_SHA256=d6bf07fb61972127c9e0d622523124c2d81caf9f7971fb123228961021811697
RUN curl -LO https://github.com/mozilla/sops/releases/download/${SOPS_VERSION}/sops-${SOPS_VERSION}.linux.amd64 && \
    echo "${SOPS_SHA256} sops-${SOPS_VERSION}.linux.amd64" | sha256sum -c && \
    mv sops-${SOPS_VERSION}.linux.amd64 sops && \
    chmod +x sops && \
    mv sops /usr/local/bin && \
    sops --version

######################  Compile  ##########################
# to have something to compile, we will use the example-infra CLI in cmd/example-infra directory. This also makes this
# Dockefile a good boilerplate for other infra projects, where cops-hq would be used as a library.
ADD . /src
WORKDIR /src/cmd/example-infra

# we build a binary called example-infra which we can use later on without always using "go run ."
RUN go build -o example-infra
RUN cp example-infra $GOPATH/bin # example-infra cli will now be globally available, although it should always be called from cmd/example-infra due to paths!
RUN example-infra --version # semantic check if it works

# as a check if all the required tooling is installed, we can use the in-built cops-hq CLI method
RUN example-infra hq check-dependencies

# change back to root working directory to do some more checks
WORKDIR /src

RUN go vet ./...
RUN go test ./...

# set back to command line tool root, because this is from where we could execute the commands in the future
WORKDIR /src/cmd/example-infra
CMD [ "/bin/bash" ]