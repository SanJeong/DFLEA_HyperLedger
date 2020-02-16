#!/bin/bash

set -ev

docker exec cli peer chaincode install -n dflea -v 1.0 -p github.com/chaincode/dflea
docker exec cli peer chaincode instantiate -o orderer0.example.com:7050 -C mychannel -n dflea -v 1.0 -c '{"Args":["init"]}' -P "OR('Org1MSP.member','Org2MSP.member')" --collections-config  /opt/gopath/src/github.com/chaincode/dflea/collections_config.json

sleep 3

docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n dflea -c '{"Args":["addUser","21400684@handong.edu","정산"]}'
sleep 3

docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n dflea -c '{"Args":["addDataset","2020san","21400684@handong.edu","foodimg","20200207","c/docu/"]}'
sleep 3

docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n dflea -c '{"Args":["addPurchase","2020san","21400684@handong.edu","foodimg","20200208"]}'
sleep 3

docker exec cli peer chaincode query -C mychannel -n dflea -c '{"Args":["readUser","21400684@handong.edu"]}'
docker exec cli peer chaincode query -C mychannel -n dflea -c '{"Args":["readDataset","2020san"]}'
docker exec cli peer chaincode query -C mychannel -n dflea -c '{"Args":["readDatasetPrivateDetails","2020san"]}'
docker exec cli peer chaincode query -C mychannel -n dflea -c '{"Args":["readPurchase","2020san"]}'
docker exec cli peer chaincode invoke -C mychannel -n dflea -c '{"Args":["transferPurchase","2020san","Success"]}'
sleep 3
docker exec cli peer chaincode query -C mychannel -n dflea -c '{"Args":["readPurchase","2020san"]}'