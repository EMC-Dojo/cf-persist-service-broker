set -x

add-apt-repository -y ppa:ubuntu-lxc/lxd-stable
apt-get -y update
apt-get -y install golang git

mkdir -p gocode
export GOPATH=/home/vcap/gocode
export PATH=$PATH:$GOPATH/bin

go get --insecure -f -u gopkg.in/yaml.v2
go get --insecure -f -u github.com/onsi/ginkgo/ginkgo
go get --insecure -f -u github.com/onsi/gomega
go get --insecure -f -u github.com/cloudfoundry-incubator/volume_driver_cert
go get --insecure -f -u code.cloudfoundry.org/clock
go get --insecure -f -u code.cloudfoundry.org/cfhttp
go get --insecure -f -u code.cloudfoundry.org/cfhttp/handlers
go get --insecure -f -u github.com/cloudfoundry/gunk/http_wrap
cd $GOPATH/src/github.com/cloudfoundry-incubator/volume_driver_cert

export FIXTURE_FILENAME=/home/vcap/rexray_config.json
ginkgo -r
