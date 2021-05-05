#!/bin/bash

#set -ev

main(){
   exportVariables
   
   ./clearEnvironment.sh

   ./generate.sh

   startNetwork

   createAndJoinChannel

   updateAnchorPeers

   installAndIntantiateChaincodes

   createUsers

}

exportVariables(){
   # don't rewrite paths for Windows Git Bash users
   export MSYS_NO_PATHCONV=1
   export CHANNEL_NAME=mychannel

   export COMPANY_DOMAIN=empresa.com
   export PEER_NUMBER=0
   export ORGANIZATION_NAME=Empresa
   export ORGANIZATION_NAME2=empresa.com

   export COMPANY2_DOMAIN=callativo.com
   export ORGANIZATION2_NAME=CallAtivo
   export ORGANIZATION2_NAME2=callativo.com

   export COMPANY3_DOMAIN=callreativo.com
   export ORGANIZATION3_NAME=CallReativo
   export ORGANIZATION3_NAME2=callreativo.com
}

startNetwork(){
   docker-compose -p network -f docker-compose.yml up -d 
   sleep 2
}

createAndJoinChannel(){
   
   createChannel

   joinOrg1OnChannel

   joinOrg2OnChannel

   joinOrg3OnChannel
}

createChannel(){
   docker exec cli.${ORGANIZATION_NAME2} peer channel create -o orderer.${ORGANIZATION_NAME2}:7050 -c ${CHANNEL_NAME} -f /etc/hyperledger/configtx/channel.tx
   sleep 2
}

joinOrg1OnChannel(){
   docker exec cli.${ORGANIZATION_NAME2} peer channel join -b ${CHANNEL_NAME}.block 
   sleep 2
}

joinOrg2OnChannel(){
   docker exec  cli.${ORGANIZATION2_NAME2} peer channel fetch 0 ${CHANNEL_NAME}.block -o orderer.${ORGANIZATION_NAME2}:7050 -c ${CHANNEL_NAME}
   sleep 2
   docker exec cli.${ORGANIZATION2_NAME2} peer channel join -b ${CHANNEL_NAME}.block
   sleep 2
}

joinOrg3OnChannel(){
   docker exec  cli.${ORGANIZATION3_NAME2} peer channel fetch 0 ${CHANNEL_NAME}.block -o orderer.${ORGANIZATION_NAME2}:7050 -c ${CHANNEL_NAME}
   sleep 2
   docker exec cli.${ORGANIZATION3_NAME2} peer channel join -b ${CHANNEL_NAME}.block
   sleep 2
}

updateAnchorPeers(){
   docker exec cli.${ORGANIZATION_NAME2} peer channel -c ${CHANNEL_NAME} update -o orderer.${ORGANIZATION_NAME2}:7050 -f /etc/hyperledger/configtx/EmpresaMSPanchors.tx 
   docker exec cli.${ORGANIZATION2_NAME2} peer channel -c ${CHANNEL_NAME} update -o orderer.${ORGANIZATION_NAME2}:7050 -f /etc/hyperledger/configtx/CallAtivoMSPanchors.tx
   docker exec cli.${ORGANIZATION3_NAME2} peer channel -c ${CHANNEL_NAME} update -o orderer.${ORGANIZATION_NAME2}:7050 -f /etc/hyperledger/configtx/CallReativoMSPanchors.tx
   sleep 2
}

installAndIntantiateChaincodes(){
   #install chaincode on peer 1
   docker exec cli.${ORGANIZATION_NAME2} peer chaincode install -n chaincode -v 1.0 -p github.com/chaincode -l golang

   #install dpoChaincode on peer 1
   docker exec cli.${ORGANIZATION_NAME2} peer chaincode install -n dpoChaincode -v 1.0 -p github.com/dpoChaincode -l golang

   #install chaincode on peer 2
   docker exec cli.${ORGANIZATION2_NAME2} peer chaincode install -n chaincode -v 1.0 -p github.com/chaincode -l golang

   #install chaincode on peer 3
   docker exec cli.${ORGANIZATION3_NAME2} peer chaincode install -n chaincode -v 1.0 -p github.com/chaincode -l golang

   sleep 2

   #instantiate chaincode on peer 1
   docker exec cli.${ORGANIZATION_NAME2} peer chaincode instantiate -o orderer.${ORGANIZATION_NAME2}:7050 -C ${CHANNEL_NAME} -n chaincode -l golang -v 1.0 -c '{"Args":[]}' -P "OR('EmpresaMSP.member', 'CallAtivoMSP.member', 'CallReativoMSP.member')" 

   #instantiate dpoChaincode on peer 1
   docker exec cli.${ORGANIZATION_NAME2} peer chaincode instantiate -o orderer.${ORGANIZATION_NAME2}:7050 -C ${CHANNEL_NAME} -n dpoChaincode -l golang -v 1.0 -c '{"Args":[]}' -P "OR('EmpresaMSP.member', 'CallAtivoMSP.member', 'CallReativoMSP.member')" --collections-config  /opt/gopath/src/github.com/dpoChaincode/collections_config.json

   sleep 2
}

createUsers(){
   cd scripts
   npm install
   node empresaEnrollAdmin.js
   node empresaRegisterUser.js
   node callAtivoEnrollAdmin.js
   node callAtivoRegisterUser.js
   node callReativoEnrollAdmin.js
   node callReativoRegisterUser.js
}

main