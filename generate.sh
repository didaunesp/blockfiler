#!/bin/sh
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=mychannel

# remove previous crypto material and config transactions
rm -fr config/*
rm -fr crypto-config/*

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# generate genesis block for orderer
configtxgen -profile OrdererGenesis -outputBlock ./config/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

# generate channel configuration transaction
configtxgen -profile Channel -outputCreateChannelTx ./config/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# generate anchor peer empresa transaction
configtxgen -profile Channel -outputAnchorPeersUpdate ./config/EmpresaMSPanchors.tx -channelID $CHANNEL_NAME -asOrg EmpresaMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for Empresa..."
  exit 1
fi

# generate anchor peer callativo transaction
configtxgen -profile Channel -outputAnchorPeersUpdate ./config/CallAtivoMSPanchors.tx -channelID $CHANNEL_NAME -asOrg CallAtivoMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for CallAtivoMSP..."
  exit 1
fi

# generate anchor peer callreativo transaction
configtxgen -profile Channel -outputAnchorPeersUpdate ./config/CallReativoMSPanchors.tx -channelID $CHANNEL_NAME -asOrg CallReativoMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for CallReativoMSP..."
  exit 1
fi
