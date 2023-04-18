FROM golang:1.20.2-bullseye

RUN apt-get update && \
    apt-get install lsb-release -y

RUN go version

RUN apt-get update

################## Tooling prerequisites  ######################
ARG AZURE_CLI_VERSION=2.46.0
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list
RUN curl -L https://packages.microsoft.com/keys/microsoft.asc | apt-key add -
RUN apt-get install apt-transport-https
RUN apt-get update && apt-get install -y azure-cli=${AZURE_CLI_VERSION}-1~bullseye

# Install terraform
ARG TERRAFORM_VERSION=1.4.0
RUN curl -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip --output terraform.zip
RUN apt-get install unzip
RUN unzip terraform.zip
RUN mv terraform /usr/local/bin

# k8s CLI
# You must use a kubectl version that is within one minor version difference of your cluster.
# For example, a v1.24 client can communicate with v1.23, v1.24, and v1.25 control planes.
ARG KUBECTL_VERSION=v1.25.7
RUN curl -LO https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
RUN kubectl version --client=true

# kubelogin CLI
ARG KUBELOGIN_VERSION=v0.0.28
RUN curl -L https://github.com/Azure/kubelogin/releases/download/${KUBELOGIN_VERSION}/kubelogin-linux-amd64.zip --output kubelogin.zip
RUN unzip kubelogin.zip -d kubelogin
RUN chmod +x kubelogin
RUN mv kubelogin/bin/linux_amd64/kubelogin /usr/local/bin
RUN kubelogin --version

# Helm
ENV HELM_VERSION v3.11.3
RUN curl -LO https://get.helm.sh/helm-${HELM_VERSION}-linux-386.tar.gz
RUN tar xvzf helm-${HELM_VERSION}-linux-386.tar.gz
RUN mv linux-386/helm $GOPATH/bin
RUN helm version

# copsctl
ENV COPSCTL_VERSION 0.9.0
RUN curl -LO https://github.com/conplementAG/copsctl/releases/download/v${COPSCTL_VERSION}/copsctl_${COPSCTL_VERSION}_Linux_x86_64.tar.gz
RUN tar xvzf copsctl_${COPSCTL_VERSION}_Linux_x86_64.tar.gz
RUN mv copsctl $GOPATH/bin
RUN copsctl --version

# sops for configuration management
ARG SOPS_VERSION=3.7.3
RUN curl -LO https://github.com/mozilla/sops/releases/download/v${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux
RUN mv sops-v${SOPS_VERSION}.linux $GOPATH/bin
RUN chmod +x $GOPATH/bin/sops-v${SOPS_VERSION}.linux
RUN mv $GOPATH/bin/sops-v${SOPS_VERSION}.linux $GOPATH/bin/sops # rename the file
RUN sops --version

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