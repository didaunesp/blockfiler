!/bin/bash

set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

export COMPANY_DOMAIN=blockfiler.com
export PEER_NUMBER=0
export ORGANIZATION_NAME=Blockfiler
export ORGANIZATION_NAME2=blockfiler.com

rm -rf ./client/controllers/users/*

docker-compose -f docker-compose.yml down

docker-compose -p network -f docker-compose.yml up -d 

#export FABRIC_START_TIMEOUT=2

sleep 5
# Create the channel
docker exec cli.${ORGANIZATION_NAME2} peer channel create -o orderer.${ORGANIZATION_NAME2}:7050  -c mychannel -f /etc/hyperledger/configtx/channel.tx

# Join peer0.${ORGANIZATION_NAME2} to the channel.
docker exec cli.${ORGANIZATION_NAME2} peer channel join -b mychannel.block 

sleep 5

docker exec cli.${ORGANIZATION_NAME2} peer chaincode install -n chaincode -v 1.0 -p github.com/chaincode -l golang

sleep 2
docker exec cli.${ORGANIZATION_NAME2} peer chaincode instantiate -o orderer.${ORGANIZATION_NAME2}:7050 -C mychannel -n chaincode -l golang -v 1.0 -c '{"Args":[]}' -P "OR('BlockfilerMSP.member')" 


