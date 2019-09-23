!/bin/bash

set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
export CHANNEL_NAME=mychannel

export COMPANY_DOMAIN=blockfiler.com
export PEER_NUMBER=0
export ORGANIZATION_NAME=Blockfiler
export ORGANIZATION_NAME2=blockfiler.com

export COMPANY2_DOMAIN=blockfiler2.com
export ORGANIZATION2_NAME=Blockfiler2
export ORGANIZATION2_NAME2=blockfiler2.com

export COMPANY3_DOMAIN=blockfiler3.com
export ORGANIZATION3_NAME=Blockfiler3
export ORGANIZATION3_NAME2=blockfiler3.com

rm -rf ./client/controllers/users/*

docker-compose -f docker-compose.yml down

docker-compose -p network -f docker-compose.yml up -d 

#export FABRIC_START_TIMEOUT=2

sleep 2
# Create the channel
docker exec cli.${ORGANIZATION_NAME2} peer channel create -o orderer.${ORGANIZATION_NAME2}:7050 -c ${CHANNEL_NAME} -f /etc/hyperledger/configtx/channel.tx

sleep 2

# Join peer0.${ORGANIZATION_NAME2} to the channel.
docker exec cli.${ORGANIZATION_NAME2} peer channel join -b ${CHANNEL_NAME}.block 

sleep 2

#fetch channel block for org 2
docker exec  cli.${ORGANIZATION2_NAME2} peer channel fetch 0 ${CHANNEL_NAME}.block -o orderer.${ORGANIZATION_NAME2}:7050 -c ${CHANNEL_NAME}

sleep 2

# Join peer0.${ORGANIZATION2_NAME2} to the channel.
docker exec cli.${ORGANIZATION2_NAME2} peer channel join -b ${CHANNEL_NAME}.block

sleep 5

#fetch channel block for org 3
docker exec  cli.${ORGANIZATION3_NAME2} peer channel fetch 0 ${CHANNEL_NAME}.block -o orderer.${ORGANIZATION_NAME2}:7050 -c ${CHANNEL_NAME}

sleep 2

# Join peer0.${ORGANIZATION3_NAME2} to the channel.
docker exec cli.${ORGANIZATION3_NAME2} peer channel join -b ${CHANNEL_NAME}.block

sleep 5

#update org1 anchor peer
docker exec cli.${ORGANIZATION_NAME2} peer channel -c ${CHANNEL_NAME} update -o orderer.${ORGANIZATION_NAME2}:7050 -f /etc/hyperledger/configtx/Org1MSPanchors.tx 
#update org2 anchor peer
docker exec cli.${ORGANIZATION2_NAME2} peer channel -c ${CHANNEL_NAME} update -o orderer.${ORGANIZATION_NAME2}:7050 -f /etc/hyperledger/configtx/Org2MSPanchors.tx
#update org3 anchor peer
docker exec cli.${ORGANIZATION3_NAME2} peer channel -c ${CHANNEL_NAME} update -o orderer.${ORGANIZATION_NAME2}:7050 -f /etc/hyperledger/configtx/Org3MSPanchors.tx


#install chaincode on peer
docker exec cli.${ORGANIZATION_NAME2} peer chaincode install -n chaincode -v 1.0 -p github.com/chaincode -l golang

sleep 2

#instantiate chaincode on peer
docker exec cli.${ORGANIZATION_NAME2} peer chaincode instantiate -o orderer.${ORGANIZATION_NAME2}:7050 -C ${CHANNEL_NAME} -n chaincode -l golang -v 1.0 -c '{"Args":[]}' -P "OR('BlockfilerMSP.member')" 


