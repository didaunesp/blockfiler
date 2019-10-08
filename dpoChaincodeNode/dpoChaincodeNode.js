/*
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
*/

const shim = require('fabric-shim');
const util = require('util');

var Chaincode = class {

    // Initialize the chaincode
    async Init(stub) {
        return shim.success();
    }

    async Invoke(stub) {
        let ret = stub.getFunctionAndParameters();
        console.info(ret);
        let method = this[ret.fcn];
        if (!method) {
            console.log('no method of name:' + ret.fcn + ' found');
            return shim.success();
        }
        try {
            let payload = await method(stub, ret.params);
            return shim.success(payload);
        } catch (err) {
            console.log(err);
            return shim.error(err);
        }
    }

    async invoke(stub, args) {
        if (args.length != 3) {
            throw new Error('Incorrect number of arguments. Expecting 3');
        }

        let A = args[0];
        let B = args[1];
        if (!A || !B) {
            throw new Error('asset holding must not be empty');
        }

        // Get the state from the ledger
        let Avalbytes = await stub.getState(A);
        if (!Avalbytes) {
            throw new Error('Failed to get state of asset holder A');
        }
        let Aval = parseInt(Avalbytes.toString());

        let Bvalbytes = await stub.getState(B);
        if (!Bvalbytes) {
            throw new Error('Failed to get state of asset holder B');
        }

        let Bval = parseInt(Bvalbytes.toString());
        // Perform the execution
        let amount = parseInt(args[2]);
        if (typeof amount !== 'number') {
            throw new Error('Expecting integer value for amount to be transaferred');
        }

        Aval = Aval - amount;
        Bval = Bval + amount;
        console.info(util.format('Aval = %d, Bval = %d\n', Aval, Bval));

        // Write the states back to the ledger
        await stub.putState(A, Buffer.from(Aval.toString()));
        await stub.putState(B, Buffer.from(Bval.toString()));

    }

    // Deletes an entity from state
    async delete(stub, args) {
        if (args.length != 1) {
            throw new Error('Incorrect number of arguments. Expecting 1');
        }

        let A = args[0];

        // Delete the key from the state in ledger
        await stub.deleteState(A);
    }

    // query callback representing the query of a chaincode
    async query(stub, args) {
        if (args.length != 2) {
            throw new Error('Incorrect number of arguments. Expecting 2 get: ' + args.length)
        }

        let jsonResp = {};
        let key = args[0];
        let user = args[1];

        // Get the state from the ledger
        let Avalbytes = await stub.getState(A);
        if (!Avalbytes) {
            jsonResp.error = 'Failed to get state for ' + A;
            throw new Error(JSON.stringify(jsonResp));
        }

        jsonResp.name = A;
        jsonResp.amount = Avalbytes.toString();
        console.info('Query Response:');
        console.info(jsonResp);
        return Avalbytes;

        /* #################################################################### */

        collection:= "collectionPublico"
        QueryAsBytes, _ := APIstub.GetPrivateData(collection, key)
        if len(QueryAsBytes) == 0 {
            collection = "collectionReativo"
            QueryAsBytes, _ = APIstub.GetPrivateData(collection, key)
            if len(QueryAsBytes) == 0 {
                collection = "collectionAtivo"
                QueryAsBytes, _ = APIstub.GetPrivateData(collection, key)
                if len(QueryAsBytes) == 0 {
                    collection = "collectionEmpresa"
                    QueryAsBytes, _ = APIstub.GetPrivateData(collection, key)
                }
            }
        }
        fmt.Println(string(QueryAsBytes))

        if len(QueryAsBytes) > 0 {

            if !s.updateRegister(APIstub, QueryAsBytes, user, collection, key) {
                return shim.Error("erro ao atualizar registro")
            }
            // var args2 []string
            // args2[0] = register.Key
            // args2[1] = register.Content
            // s.create(APIstub, args2)
        }

        return shim.Success(QueryAsBytes)
    }
};

shim.start(new Chaincode());
