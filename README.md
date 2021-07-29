# Hyperledger-Fabric-Demo
## A simple voting demo app using Hyperledger Fabric 
Chaincode written in Golang, client side is using the node.js SDK
### To run the project 
1. [Install](https://hyperledger-fabric.readthedocs.io/en/latest/install.html) Fabric
2. Copy the repository to the same directory as fabric-samples
3. Create the network using the test-network scripts in fabric-samples
    ```
    ./startFabric.sh
    ```
4. Build the project
    ```
    cd javascript
    npm install
    ```
##### Supported commands are:
1. node enrollAdmin.js – register and admin with name “admin”
2. node registerUser.js – register an User with name “appUser”
3. node query.js – get all voting items available
4. node queryUsers.js – get all voters and their vote information
5. node vote.js itemID – vote for an item with ID itemID
6. node index.js – demonstrate all the above-mentioned commands
