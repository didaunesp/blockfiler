!/bin/bash

set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

export COMPANY_DOMAIN=blockfiler.com
export PEER_NUMBER=0
export ORGANIZATION_NAME=Blockfiler
export ORGANIZATION_NAME2=blockfiler.com

export COMPANY2_DOMAIN=blockfiler2.com
export ORGANIZATION2_NAME=Blockfiler2
export ORGANIZATION2_NAME2=blockfiler2.com

rm -rf ./client/controllers/users/*

docker-compose -f docker-compose.yml down

docker-compose -p network -f docker-compose.yml up -d 

#export FABRIC_START_TIMEOUT=2

sleep 2
# Create the channel
docker exec cli.${ORGANIZATION_NAME2} peer channel create -o orderer.${ORGANIZATION_NAME2}:7050 -c mychannel -f /etc/hyperledger/configtx/channel.tx

# Join peer0.${ORGANIZATION_NAME2} to the channel.
docker exec cli.${ORGANIZATION_NAME2} peer channel join -b mychannel.block 

sleep 5

# Join peer0.${ORGANIZATION2_NAME2} to the channel.
#docker exec -e "CORE_PEER_ADDRESS=peer0.${ORGANIZATION2_NAME2}:7051" -e "CORE_PEER_LOCALMSPID=Blockfiler2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/peer/crypto/peerOrganizations/${COMPANY2_DOMAIN}/users/Admin@${COMPANY2_DOMAIN}/msp" cli.${ORGANIZATION_NAME2} peer channel join -b mychannel.block
docker exec cli.${ORGANIZATION2_NAME2} peer channel join -b mychannel.block

sleep 5

#install chaincode on peer
docker exec cli.${ORGANIZATION_NAME2} peer chaincode install -n chaincode -v 1.0 -p github.com/chaincode -l golang

sleep 2

#instantiate chaincode on peer
docker exec cli.${ORGANIZATION_NAME2} peer chaincode instantiate -o orderer.${ORGANIZATION_NAME2}:7050 -C mychannel -n chaincode -l golang -v 1.0 -c '{"Args":[]}' -P "OR('BlockfilerMSP.member')" 


