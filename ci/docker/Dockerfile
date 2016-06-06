FROM ubuntu:14.04
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN apt-get -y update
RUN apt-get install -y git wget curl

RUN wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
RUN tar zxvf go1.6.2.linux-amd64.tar.gz -C /usr/local/

ENV GOROOT /usr/local/go
ENV PATH $GOROOT/bin:$PATH

RUN apt-get install -y autoconf bison build-essential libssl-dev libyaml-dev libreadline6-dev zlib1g-dev libncurses5-dev libffi-dev libgdbm3 libgdbm-dev

RUN git clone https://github.com/sstephenson/rbenv.git /root/.rbenv
ENV PATH /root/.rbenv/bin:$PATH
ENV PATH /root/.rbenv/shims:$PATH
ENV RBENV_SHELL sh

RUN git clone https://github.com/sstephenson/ruby-build.git /root/.rbenv/plugins/ruby-build
ENV PATH /root/.rbenv/plugins/ruby-build/bin:$PATH

RUN rbenv install 2.2.4
RUN rbenv global 2.2.4

RUN echo eval "$(rbenv init -)" > /etc/profile.d/rbenv.sh && \
    chmod +x /etc/profile.d/rbenv.sh && \
    source /etc/profile.d/rbenv.sh && \
    gem install bundler

RUN wget -q -O - "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx -O > /usr/bin/cf
RUN chmod +x /usr/bin/cf

ENV GOBIN /usr/bin
ENV GOPATH /root/gocode

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go install github.com/onsi/ginkgo/ginkgo
RUN go get github.com/tools/godep
RUN go install github.com/tools/godep
