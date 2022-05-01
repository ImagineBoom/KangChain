export BYFN_CA1_PRIVATE_KEY=$(cd demo0/crypto-config/peerOrganizations/org1.example.com/ca && ls *_sk)
export BYFN_CA2_PRIVATE_KEY=$(cd demo0/crypto-config/peerOrganizations/org2.example.com/ca && ls *_sk)
export GO111MODULE=on
.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### Prerequisites
pre:
	@echo "make crypto-materials"
	$(shell mkdir -p $(GOPATH)/src/github.com/KangChain.com/KangChain/demo0/channel-artifacts)
	@echo "cryptogen ..."
	@if [ -d "demo0/crypto-config" ]; then \
        rm -rf demo0/crypto-config; \
		echo "removed crypto-config"; \
		echo "recreate crypto-config"; \
    fi
	cd demo0 && cryptogen generate --config=crypto-config.yaml --output="crypto-config"
	@echo "configtxgen ..."
	@if [ -d "demo0/channel-artifacts" ]; then \
        rm -rf demo0/channel-artifacts; \
		echo "removed channel-artifacts"; \
		mkdir -p $(GOPATH)/src/github.com/KangChain.com/KangChain/demo0/channel-artifacts; \
		echo "recreate channel-artifacts"; \
	fi
	cd demo0 && configtxgen -profile SampleMultiNodeEtcdRaft -channelID byfn-sys-channel  -outputBlock ./channel-artifacts/genesis.block
	cd demo0 && configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
	cd demo0 && configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
	cd demo0 && configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
##### BUILD
build:
ifneq (go.mod,$(wildcard go.mod))
	@export GO111MODULE=on && go mod init && go get github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl@v1.0.0-beta2
#	&&go get github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl@master
endif
	#export GO111MODULE=on && go get github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl@master
	export GO111MODULE=on && go get github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl@v1.0.0-beta2
	@echo "Build ..."
	go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	cd demo0 && docker-compose -f docker-compose-cli.yaml up -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	cd demo0 && docker-compose -f docker-compose-cli.yaml down
	@docker volume prune
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./KangChain

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/msp
	@rm -rf /tmp/state-store
	@docker rm -f -v `docker ps -a --no-trunc | grep "demo0" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "demo0" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"


