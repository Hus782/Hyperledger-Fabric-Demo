/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./CAUtil.js');
const FabricCAServices = require('fabric-ca-client');


const path = require('path');
const fs = require('fs');
const walletPath = path.join(process.cwd(), 'wallet');
const mspOrg1 = 'Org1MSP';
const org1UserId = 'appUser';
const channelName = 'mychannel';
const chaincodeName = 'voting';

function printJSON(result){
    if (`${result}` !== '') {
        console.log(JSON.parse(result.toString()));
    }
}
async function main() {
    try {
        // load the network configuration

        const ccpPath = path.resolve(__dirname, '..', '..', 'fabric-samples','test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new CA client for interacting with the CA.
        const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');


        // Create a new file system based wallet for managing identities.
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // in a real application this would be done on an administrative flow, and only once
        await enrollAdmin(caClient, wallet, mspOrg1);

        // in a real application this would be done only when a new user was required to be added
        // and would be part of an administrative flow
        await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: org1UserId, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(channelName);

        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);

        try{
            
            let result = await contract.submitTransaction('CreateVoter');
            console.log(result.toString());
 
            console.log('-------Query all items-------');
            result = await contract.evaluateTransaction('QueryAllItems');
            printJSON(result)

            console.log('-------Voting for item with ID 1-------');
            result =  await contract.submitTransaction('Vote','1');
            console.log(result.toString());
            if (`${result}` === 'Voted successfully') {
                result = await contract.submitTransaction('UpdateVoter','1');
                console.log(result.toString());
            }

            console.log('-------Query all items-------');
            result = await contract.evaluateTransaction('QueryAllItems');
            printJSON(result)

            console.log('-------Voting for item with ID 1-------');
            result =  await contract.submitTransaction('Vote','1');
            console.log(result.toString());
            if (`${result}` === 'Voted successfully') {
                result = await contract.submitTransaction('UpdateVoter','1');
                console.log(result.toString());
            }

            console.log('-------Query all items-------');
            result = await contract.evaluateTransaction('QueryAllItems');
            printJSON(result)
            
        }
        catch (error) {
            console.error(`Failed to evaluate transaction: ${error}`);
            process.exit(1);
        }

        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

main();
