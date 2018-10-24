Current build process:
#Software install
sudo yum install golang

#Setup source
mkdir -p $HOME/go/src/redhat.com/consulting
cd $HOME/go/src/redhat.com/consulting
git clone https://github.com/tristianc/stager.git

#Setup dependencies
tar xf ~/Downloads/stager_deps.tar.gz -C $HOME/go/src

#Build
cd cd $HOME/go/src/redhat.com/consulting/stager
go build

#Run
./stager

