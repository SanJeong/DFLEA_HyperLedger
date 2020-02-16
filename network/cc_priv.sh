#chaincode insall
docker exec cli peer chaincode install -n mym -v 1.0 -p github.com/mymarbles
#chaincode instatiate
docker exec cli peer chaincode instantiate -n mym -v 1.0 -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")' --collections-config  /opt/gopath/src/github.com/mymarbles/collections_config.json 
sleep 5

#chaincode invoke init marble1
export MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)
docker exec cli peer chaincode invoke -o orderer.example.com:7050 -n mym -C mychannel -c '{"Args":["initMarble"]}'  --transient "{\"marble\":\"$MARBLE\"}"
sleep 5

#chaincode query
docker exec cli peer chaincode query -C mychannel -n mym -c '{"Args":["readMarble","marble1"]}'

docker exec cli peer chaincode query -C mychannel -n mym -c '{"Args":["readMarblePrivateDetails","marble1"]}'

#export CORE_PEER_LOCALMSPID=Org2MSP
#export CORE_PEER_ADDRESS=peer0.org2.example.com:7051
#export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

echo '-------------------------------------END-------------------------------------'
