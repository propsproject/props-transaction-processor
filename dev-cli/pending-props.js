// import { default as Web3} from 'web3';
const Web3 = require('web3');
const web3 = new Web3();
const {
    createContext,
    CryptoFactory
} = require('sawtooth-sdk/signing');
const {
    Secp256k1PrivateKey,
    Secp256k1PublicKey
} = require('sawtooth-sdk/signing/secp256k1');
const {
    createHash
} = require('crypto');
const {
    protobuf
} = require('sawtooth-sdk');
const request = require('request');
const colors = require('colors');
const proto = require('google-protobuf');
const any = require('google-protobuf/google/protobuf/any_pb.js');
const payloads_pb = require('./proto/payload_pb');
const earnings_pb = require('./proto/earning_pb');
const balance_pb = require('./proto/balance_pb');
const activity_pb = require('./proto/activity_pb');
const users_pb = require('./proto/users_pb');
const ethUtil = require('ethereumjs-util');
const ethWallet = require('ethereumjs-wallet');
const context = createContext('secp256k1');
const opn = require('opn');
const axios = require('axios');
const prettyjson = require('prettyjson');
const moment = require('moment');
const BigNumber = require('bignumber.js');
BigNumber.config({ EXPONENTIAL_AT: 1e+9 })
const restAPIHost = process.env.REST_API_URL != undefined ? process.env.REST_API_URL : 'http://127.0.0.1:8008';

// hard coded example private key
// const pk = Buffer.from("196749ed808372060eaeffe10e56de82a48829fcf52199847e1e1db4b780ced0", 'hex');
const pk = Buffer.from("fd899d64b5209b53e6b6380dbe195500d988b2184d3a7076681370d5d1c58408", 'hex');

const priv = new Secp256k1PrivateKey(pk);
const signer = new CryptoFactory(context).newSigner(priv);

let loggerType = 0; // 0 console, -1 none

const setLoggerType = (type) => {
    loggerType = type;
}

const signMessage = async (msg, address, pk) => {
    const privateKey = pk;
    const account = web3.eth.accounts.privateKeyToAccount('0x' + privateKey);
    web3.eth.accounts.wallet.add(account);
    web3.eth.defaultAccount = account.address;
    return web3.eth.sign(msg, address);
}

const recoverFromSignature = async (msg, sig) => {

    //web3.eth.accounts.wallet.add(account);
    //web3.eth.defaultAccount = account.address;
    return web3.eth.accounts.recover(msg, sig);
}

const log = (str, type="info") => {
    if (loggerType>=0) {
        switch (type) {
            case "info":
                console.log(str);
                break;
            case "warn":
                console.warn(str);
                break;
            case "err":
                console.error(str);
                break;
        }
    }
}

const CONFIG = {
    earnings: {
        familyName: "pending-earnings",
        familyVersion: "1.0",
        namespaces: {
            prefixes: {
                "pending": createHash('sha512')
                    .update("pending-props:earnings:pending")
                    .digest('hex')
                    .substring(0, 6),
                "revoked": createHash('sha512')
                    .update("pending-props:earnings:revoked")
                    .digest('hex')
                    .substring(0, 6),
                "settled": createHash('sha512')
                    .update("pending-props:earnings:settled")
                    .digest('hex')
                    .substring(0, 6),
                "settlements": createHash('sha512')
                    .update("pending-props:earnings:settlements")
                    .digest('hex')
                    .substring(0, 6),
                "balance": createHash('sha512')
                    .update("pending-props:earnings:balance")
                    .digest('hex')
                    .substring(0, 6),
                "balanceUpdate": createHash('sha512')
                    .update("pending-props:earnings:bal-rtx")
                    .digest('hex')
                    .substring(0, 6),
                "blockIdUpdate": createHash('sha512')
                    .update("pending-props:earnings:lastethblock")
                    .digest('hex')
                    .substring(0, 6),
                "walletLink": createHash('sha512')
                    .update("pending-props:earnings:walletl")
                    .digest('hex')
                    .substring(0, 6),
                "activityLog": createHash('sha512')
                    .update("pending-props:earnings:activity_log")
                    .digest('hex')
                    .substring(0,6),
            },
            earningAddress(status, args) {
                const prefix = this.prefixes[status];
                let address = prefix;
                args.forEach(a => {
                    address = address.concat(createHash('sha512').update(`${a.data}`).digest('hex').substring(a.start, a.end));
                });
                // const recID =
                // createHash('sha512').update(recipient).digest('hex').substring(0, 4); const
                // appID = createHash('sha512').update(application).digest('hex').substring(0,
                // 4); const postfix =
                // createHash('sha512').update(`${recipient}${application}${signature}`).digest('
                // hex').toLowerCase().substring(0, 56);
                return address
            },
            settlementAddress(ethereumTxtHash) {
                const prefix = this.prefixes.settlements;
                const postfix = createHash('sha512')
                    .update(`${normalizeAddress(ethereumTxtHash)}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 64);
                return `${prefix}${postfix}`
            },
            balanceAddress(applicationId, userId) {
                const prefix = this.prefixes.balance;
                const part1 = createHash('sha512')
                    .update(`${applicationId}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 10);
                const part2 = createHash('sha512')
                    .update(`${userId}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 54);
                return `${prefix}${part1}${part2}`
            },
            walletLinkAddress(address) {
                const prefix = this.prefixes.walletLink;
                const body = createHash('sha512')
                    .update(`${address}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 64);
                return `${prefix}${body}`
            },
            balanceUpdateAddress(txHash, address) {
                txHash = normalizeAddress(txHash);
                const prefix = this.prefixes.balanceUpdate;
                const body = createHash('sha512')
                    .update(`${txHash}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 40);
                const postfix = createHash('sha512')
                    .update(`${address}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 24);
                return `${prefix}${body}${postfix}`
            }
            ,
            blockUpdateAddress() {
                const prefix = this.prefixes.blockIdUpdate;
                const body = createHash('sha512')
                    .update('LastEthBlockAddress')
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 64);
                return `${prefix}${body}`
            },
            activityLogAddress(date, appId, userId) {
                const prefix = this.prefixes.activityLog;
                console.log(prefix);
                const part1 = createHash('sha512')
                    .update(date)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 8);
                const part2 = createHash('sha512')
                    .update(appId)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 10);
                const part3 = createHash('sha512')
                    .update(userId)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0,46);
                return `${prefix}${part1}${part2}${part3}`
            }
        }
    }
};

const newRPCRequest = (params, method) => {
    const reqParams = new payloads_pb.Params();
    reqParams.setData(params);
    const payload = new payloads_pb.RPCRequest();
    payload.setMethod(method);
    payload.setParams(reqParams);

    return payload;
};

const newEarningsDetails = (amount, applicationId, userId, description = '') => {
    const details = new earnings_pb.EarningDetails();
    details.setTimestamp(moment().unix());
    details.setUserId(userId);
    details.setApplicationId(applicationId);
    details.setDescription(description);
    BigNumber.config({ EXPONENTIAL_AT: 1e+9 });
    const propsAmount = new BigNumber(amount, 10);
    const tokensAmount = propsAmount.times(1e18);
    const zero = new BigNumber(0, 10);
    details.setAmountEarned(tokensAmount.toString());
    details.setAmountSettled(zero.toString());

    return details;
};

const externalBalanceUpdate = async (address, balance, ethTransactionHash, blockId, timestamp) => {
    address = normalizeAddress(address);
    const txHash = normalizeAddress(ethTransactionHash);
    const balanceUpdate = new earnings_pb.BalanceUpdate();
    balanceUpdate.setPublicAddress(address);
    balanceUpdate.setOnchainBalance(balance);
    balanceUpdate.setTxHash(txHash);
    balanceUpdate.setBlockId(blockId);
    balanceUpdate.setTimestamp(timestamp);
    //setup RPC request
    const params = new any.Any();
    params.setValue(balanceUpdate.serializeBinary());
    params.setTypeUrl('github.com/propsproject/pending-props/protos/pending_props_pb.BalanceUpdate');

    const request = newRPCRequest(params, payloads_pb.Method.BALANCE_UPDATE);
    const requestBytes = request.serializeBinary();

    //compute balance and balaneTimestamp addresses for outputs
    const balanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress("", address);

    const balanceUpdateAddress = CONFIG
        .earnings
        .namespaces
        .balanceUpdateAddress(txHash, address);

    const walletLinkAddress = CONFIG
        .earnings
        .namespaces
        .walletLinkAddress(address);

    const linkedApplicationUsers = await getLinkedUsersFromWalletLinkAddress(walletLinkAddress);
    const linkedApplicationUserAddresses = [];
    log(`linkedApplicationUsers ${JSON.stringify(linkedApplicationUsers)}`);
    if (linkedApplicationUsers.length > 0) {
        for (let i = 0; i < linkedApplicationUsers[0].usersList.length; ++i) {
            const linkedBalanceAddress = CONFIG
                .earnings
                .namespaces
                .balanceAddress(linkedApplicationUsers[0].usersList[i].applicationId, linkedApplicationUsers[0].usersList[i].userId);
            if (linkedBalanceAddress != balanceAddress) {
                linkedApplicationUserAddresses.push(linkedBalanceAddress);
            }
        }
        log(`linkedApplicationUsersAddresses ${JSON.stringify(linkedApplicationUserAddresses)}`);
    }

    log("Balance Address = "+balanceAddress);
    log("BalanceUpdate Address = "+balanceUpdateAddress);
    log("WalletLink Address = "+walletLinkAddress);
    const inputs = [balanceAddress, balanceUpdateAddress, walletLinkAddress, ...linkedApplicationUserAddresses];
    const outputs = [balanceAddress, balanceUpdateAddress, walletLinkAddress, ...linkedApplicationUserAddresses];

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();

    // console.log(colors.yellow(`balance state addresses: recipientBalanceAddress:${recipientBalanceAddress} fromBalanceAddress:${fromBalanceAddress} ${balanceTimestampAddressPrefix}`));
    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const linkWallet = async (address, applicationId, userId, signature) => {
    address = normalizeAddress(address);
    const walletToUser = new users_pb.WalletToUser();
    walletToUser.setAddress(address);
    const applicationUser = new users_pb.ApplicationUser();
    applicationUser.setUserId(userId);
    applicationUser.setApplicationId(applicationId);
    applicationUser.setSignature(signature);
    applicationUser.setTimestamp(moment().unix());
    walletToUser.addUsers(applicationUser);

    //setup RPC request
    const params = new any.Any();
    params.setValue(walletToUser.serializeBinary());
    params.setTypeUrl('github.com/propsproject/pending-props/protos/pending_props_pb.WalletToUser');

    const request = newRPCRequest(params, payloads_pb.Method.WALLET_LINK);
    const requestBytes = request.serializeBinary();

    const walletLinkAddress = CONFIG
        .earnings
        .namespaces
        .walletLinkAddress(address);


    //add new application user balance address for reading and writing
    const balanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress(applicationId, userId);

    //add walletBalance address in case it needs to be created (0 props scenario)
    const walletBalanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress("", address);

    // read walletLinkAddress existing data and add to inputs/outputs so it can read/write from/to it
    const linkedApplicationUsers = await getLinkedUsersFromWalletLinkAddress(walletLinkAddress);
    const linkedApplicationUserAddresses = [];
    log(`linkedApplicationUsers ${JSON.stringify(linkedApplicationUsers)}`);
    if (linkedApplicationUsers.length > 0) {
        for (let i = 0; i < linkedApplicationUsers[0].usersList.length; ++i) {
            linkedApplicationUserAddresses.push(
                CONFIG
                    .earnings
                    .namespaces
                    .balanceAddress(linkedApplicationUsers[0].usersList[i].applicationId, linkedApplicationUsers[0].usersList[i].userId)
            )
        }
        log(`linkedApplicationUsersAddresses ${JSON.stringify(linkedApplicationUserAddresses)}`);
    }

    log("walletLinkAddress = "+walletLinkAddress);
    log("balanceAddress = "+balanceAddress);
    log("walletBalanceAddress = "+walletBalanceAddress);
    log("linkedApplicationUserAddresses = "+linkedApplicationUserAddresses.join(","));
    const inputs = [walletLinkAddress, balanceAddress, walletBalanceAddress, ...linkedApplicationUserAddresses];
    const outputs = [walletLinkAddress, balanceAddress, walletBalanceAddress, ...linkedApplicationUserAddresses];

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();

    // console.log(colors.yellow(`balance state addresses: recipientBalanceAddress:${recipientBalanceAddress} fromBalanceAddress:${fromBalanceAddress} ${balanceTimestampAddressPrefix}`));
    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const issue = async (applicationId, userId, amount, description = '', addresses = {}) => {
    //setup details
    const details = newEarningsDetails(amount, applicationId, userId, description);
    const hashToSign = createHash('sha512')
        .update(details.serializeBinary())
        .digest('hex')
        .toLowerCase()

    const earningsSignature = signer.sign(Buffer.from(hashToSign));

    const earning = new earnings_pb.Earning();
    earning.setDetails(details);
    earning.setSignature(earningsSignature);

    //setup RPC request
    const params = new any.Any();
    params.setValue(earning.serializeBinary());
    params.setTypeUrl('github.com/propsproject/pending-props/protos/pending_props_pb.Earning');

    const request = newRPCRequest(params, payloads_pb.Method.ISSUE);
    const requestBytes = request.serializeBinary();

    const addressArgs = [{
        data: details.getApplicationId(),
        start: 0,
        end: 4
    }, {
        data: details.getUserId(),
        start: 0,
        end: 4
    }, {
        data: `${applicationId}${userId}${earningsSignature}`,
        start: 0,
        end: 56
    }];

    // compute state address
    const stateAddress = CONFIG
        .earnings
        .namespaces
        .earningAddress("pending", addressArgs);

    // compute balance and balaneTimestamp addresses for outputs
    const balanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress(applicationId, userId);

    // check if user balance is linked to a wallet
    const balance = await getLinkedWalletFromBalanceAddress(balanceAddress);
    const linkedWalletAddress = (balance[0]===undefined || !('linkedWallet' in balance[0])) ? "" : balance[0].linkedWallet;
    const linkedApplicationUserAddresses = []
    if (linkedWalletAddress.length > 0) {
        const walletLinkAddress = CONFIG
            .earnings
            .namespaces
            .walletLinkAddress(linkedWalletAddress);

        const walletBalanceAddress = CONFIG
            .earnings
            .namespaces
            .balanceAddress("", linkedWalletAddress);
        const linkedApplicationUsers = await getLinkedUsersFromWalletLinkAddress(walletLinkAddress);
        linkedApplicationUserAddresses.push(walletLinkAddress);
        linkedApplicationUserAddresses.push(walletBalanceAddress);
        log(`linkedApplicationUsers ${JSON.stringify(linkedApplicationUsers)}`);
        if (linkedApplicationUsers.length > 0) {
            for (let i = 0; i < linkedApplicationUsers[0].usersList.length; ++i) {
                const linkedBalanceAddress = CONFIG
                    .earnings
                    .namespaces
                    .balanceAddress(linkedApplicationUsers[0].usersList[i].applicationId, linkedApplicationUsers[0].usersList[i].userId);
                if (linkedBalanceAddress != balanceAddress) {
                    linkedApplicationUserAddresses.push(linkedBalanceAddress);
                }
            }
            log(`linkedApplicationUsersAddresses ${JSON.stringify(linkedApplicationUserAddresses)}`);
        }
    }
    log(JSON.stringify(linkedApplicationUserAddresses));
    log("Earnings (pending) Address = "+stateAddress);
    log("Balance Address = "+balanceAddress);
    const inputs = [stateAddress, balanceAddress, ...linkedApplicationUserAddresses];
    const outputs = [stateAddress, balanceAddress, ...linkedApplicationUserAddresses];
    addresses['stateAddress'] = stateAddress;
    addresses['balanceAddress'] = balanceAddress;
    addresses['linkedApplicationUserAddresses'] = linkedApplicationUserAddresses;
    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();

    log(`state addresses: ${stateAddress}`);
    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const updateLastBlockId = async (blockId) => {
    const blockUpdate = new earnings_pb.LastEthBlock();
    blockUpdate.setId(blockId);
    //setup RPC request
    const params = new any.Any();
    params.setValue(blockUpdate.serializeBinary());
    params.setTypeUrl('github.com/propsproject/pending-props/protos/pending_props_pb.LastEthBlock');

    const request = newRPCRequest(params, payloads_pb.Method.LAST_ETH_BLOCK_UPDATE);
    const requestBytes = request.serializeBinary();
    //compute balance and balaneTimestamp addresses for outputs
    const lastBlockAddress = CONFIG
        .earnings
        .namespaces
        .blockUpdateAddress();


    log("Last Eth Block Id Address = "+lastBlockAddress);
    const inputs = [lastBlockAddress];
    const outputs = [lastBlockAddress];

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();

    // console.log(colors.yellow(`balance state addresses: recipientBalanceAddress:${recipientBalanceAddress} fromBalanceAddress:${fromBalanceAddress} ${balanceTimestampAddressPrefix}`));
    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

// settle can only be updated to new schema once wallet linking exists and will require a read before write
const settle = async (transactionHash, recipient) => {
    return true; // TODO settle requires app address look up to appId will deal with it when dealing with settlements
    recipient = normalizeAddress(recipient);
    transactionHash = normalizeAddress(transactionHash);
    const appAddr = normalizeAddress(ethUtil.pubToAddress(signer.getPublicKey().asBytes(), true).toString('hex'));
    const pendingAddresses = await getPendingEarningsAddress(recipient, appAddr);

    const settleAddresses = [];
    pendingAddresses.forEach(address => {
        settleAddresses.push(`${CONFIG.earnings.namespaces.prefixes.settled}${address.substring(6)}`);
    });

    const settlementAddress = CONFIG
        .earnings
        .namespaces
        .settlementAddress(transactionHash);
    //setup RPC request
    const paramData = JSON.stringify({
        "eth_transaction_hash": transactionHash,
        recipient: recipient,
        "pending_addresses": pendingAddresses,
        "timestamp": moment().unix()
    });
    const params = new any.Any();
    params.setValue(Buffer.from(paramData));
    const request = newRPCRequest(params, payloads_pb.Method.SETTLE);
    const requestBytes = request.serializeBinary();

    //compute balance and balaneTimestamp addresses for outputs
    const balanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress(recipient);

    log("Balance Address = "+balanceAddress);

    let inputs,
        outputs;
    inputs = outputs = [
        balanceAddress, settlementAddress, ...pendingAddresses,
        ...settleAddresses
    ];
    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();
    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const revoke = async (addresses, revokeAddress = {}) => {
    //setup RPC request
    log(`debug:: addresses ${JSON.stringify(addresses)}`);
    const paramData = JSON.stringify({
        addresses: addresses,
        timestamp: moment().unix()
    });
    const params = new any.Any();
    params.setValue(Buffer.from(paramData));
    const request = newRPCRequest(params, payloads_pb.Method.REVOKE);
    const requestBytes = request.serializeBinary();

    const revokeAddresses = [];
    const linkedApplicationUserAddresses = [];
    let recipientsAddresses = [];
    for (const address of addresses) {
        revokeAddresses.push(`${CONFIG.earnings.namespaces.prefixes.revoked}${address.substring(6)}`);
        recipientsAddresses = recipientsAddresses.concat(await getRecipientFromStateAddress(address));
    }
    //compute balance and balaneTimestamp addresses for outputs
    const balanceAddresses = [];
    for (let i = 0; i < recipientsAddresses.length; i = i + 1)
    {
        const recipient = recipientsAddresses[i];
        const balanceAddress = CONFIG
            .earnings
            .namespaces
            .balanceAddress(recipient.applicationId, recipient.userId)
        balanceAddresses.push(balanceAddress);

        // check if user balance is linked to a wallet
        const balance = await getLinkedWalletFromBalanceAddress(balanceAddress);
        const linkedWalletAddress = (balance[0]===undefined || !('linkedWallet' in balance[0])) ? "" : balance[0].linkedWallet;
        if (linkedWalletAddress.length > 0) {
            const walletLinkAddress = CONFIG
                .earnings
                .namespaces
                .walletLinkAddress(linkedWalletAddress);

            const walletBalanceAddress = CONFIG
                .earnings
                .namespaces
                .balanceAddress("", linkedWalletAddress);
            const linkedApplicationUsers = await getLinkedUsersFromWalletLinkAddress(walletLinkAddress);
            linkedApplicationUserAddresses.push(walletLinkAddress);
            linkedApplicationUserAddresses.push(walletBalanceAddress);
            log(`linkedApplicationUsers ${JSON.stringify(linkedApplicationUsers)}`);
            if (linkedApplicationUsers.length > 0) {
                for (let i = 0; i < linkedApplicationUsers[0].usersList.length; ++i) {
                    const linkedBalanceAddress = CONFIG
                        .earnings
                        .namespaces
                        .balanceAddress(linkedApplicationUsers[0].usersList[i].applicationId, linkedApplicationUsers[0].usersList[i].userId);
                    if (linkedBalanceAddress != balanceAddress) {
                        linkedApplicationUserAddresses.push(linkedBalanceAddress);
                    }
                }
                log(`linkedApplicationUsersAddresses ${JSON.stringify(linkedApplicationUserAddresses)}`);
            }
        }
    }

    log(JSON.stringify(linkedApplicationUserAddresses));
    const inputs = [...balanceAddresses, ...addresses, ...revokeAddresses, ...linkedApplicationUserAddresses];
    const outputs = [...balanceAddresses, ...addresses, ...revokeAddresses, ...linkedApplicationUserAddresses];
    revokeAddress['balanceAddresses'] = balanceAddresses;
    revokeAddress['addresses'] = addresses;
    revokeAddress['revokeAddresses'] = revokeAddresses;
    revokeAddress['linkedApplicationUserAddresses'] = linkedApplicationUserAddresses;

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();

    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const submitTransaction = async (transactionHeaderBytes, requestBytes) => {
    try {
        const signature = signer.sign(transactionHeaderBytes);

        const transaction = protobuf
            .Transaction
            .create({
                header: transactionHeaderBytes,
                headerSignature: signature,
                payload: requestBytes
            });

        const transactions = [transaction];

        const batchHeaderBytes = protobuf
            .BatchHeader
            .encode({
                signerPublicKey: signer
                    .getPublicKey()
                    .asHex(),
                transactionIds: transactions.map((txn) => txn.headerSignature)
            })
            .finish();

        const signature1 = signer.sign(batchHeaderBytes);

        const batch = protobuf
            .Batch
            .create({
                header: batchHeaderBytes,
                headerSignature: signature1,
                transactions: transactions
            });

        const batchListBytes = protobuf
            .BatchList
            .encode({
                batches: [batch]
            })
            .finish();

        const reqConfig = {
            method: 'POST',
            url: restAPIHost+'/batches',
            data: batchListBytes,
            headers: {
                'Content-Type': 'application/octet-stream'
            }
        };

        const response = await axios(reqConfig);

        const link = response.data.link;
        log(`transaction submitted successfully`);
        log(`status: ${link}`);
        //opn(link);
        return response;
    } catch (e) {
        throw e;
    }
};

const getRecipientFromStateAddress = async (address) => {
    try {
        const reqConfig = {
            method: 'GET',
            url: `${restAPIHost}/state?address=${address}`,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        const response = await axios(reqConfig);
        const data = response.data.data;

        return deserializeEarnings(data);
    } catch (e) {
        throw e;
    }
}

const getLinkedUsersFromWalletLinkAddress = async (address) => {
    try {
        const reqConfig = {
            method: 'GET',
            url: `${restAPIHost}/state?address=${address}`,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        const response = await axios(reqConfig);
        const data = response.data.data;

        return deserializeWalletLink(data);
    } catch (e) {
        throw e;
    }
}

const getLinkedWalletFromBalanceAddress = async (address) => {
    try {
        const reqConfig = {
            method: 'GET',
            url: `${restAPIHost}/state?address=${address}`,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        const response = await axios(reqConfig);
        const data = response.data.data;

        return deserializeBalance(data);
    } catch (e) {
        throw e;
    }
}


const queryState = async (address, t) => {
    try {
        const reqConfig = {
            method: 'GET',
            url: `${restAPIHost}/state?address=${address}`,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        const response = await axios(reqConfig);
        const data = response.data.data;

        if (t === "earning") {
            return deserializeEarnings(data);
        } else if (t === "settlement") {
            return deserializeSettlements(data);
        } else if (t === "balance") {
            return deserializeBalance(data);
        } else if (t === "lastblockid") {
            return deserializeLastEthBlockId(data);
        } else if (t === "walletlink") {
            return deserializeWalletLink(data);
        } else if (t === "activity") {
            return deserializeActivity(data)
        } else {
            log(`unknown state type ${t} expected earning or settlements`);
        }


    } catch (e) {
        throw e;
    }
};

const deserializeEarnings = (data) => {
    const recipients = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const earning = new earnings_pb
            .Earning
            .deserializeBinary(bytes);

        const output = {
            'state-address': entry.address,
            'earning-data': earning.toObject(),
            'human-readable': {
                timestamp: moment(new Date(earning.getDetails().getTimestamp() * 1000)).format('L'),
                status: Object
                    .keys(earnings_pb.Status)
                    .filter((k, i) => i === earning.getDetails().getStatus())
            }
        };
        log(prettyjson.render(output));
        recipients.push({ applicationId: earning.getDetails().getApplicationId(), userId: earning.getDetails().getUserId(), earning:earning.toObject()})
    });
    return recipients
};

const deserializeSettlements = (data) => {
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const settlement = new earnings_pb
            .Settlements
            .deserializeBinary(bytes);
        log(prettyjson.render(settlement.toObject()));
    });
};

const deserializeBalance = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const balance = new balance_pb
            .Balance
            .deserializeBinary(bytes);
        retData.push(balance.toObject());
        log(prettyjson.render(balance.toObject()));
    });
    return retData;
};

const deserializeWalletLink = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const walletLink = new users_pb
            .WalletToUser
            .deserializeBinary(bytes);
        retData.push(walletLink.toObject());
        log(prettyjson.render(walletLink.toObject()));
    });
    return retData;
};

const deserializeLastEthBlockId = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const lastBlock = new earnings_pb
            .LastEthBlock
            .deserializeBinary(bytes);
        retData.push(lastBlock.toObject());
        log(prettyjson.render(lastBlock.toObject()));
    });
    return retData;
};

const deserializeActivity = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const activityLog = new activity_pb
            .ActivityLog
            .deserializeBinary(bytes);
        retData.push(activityLog.toObject());
        log(prettyjson.render(activityLog.toObject()));
    });
    return retData;
};

const logActivity = async(userId, appId, timestamp, date) => {

    //setup RPC request
    const paramData = JSON.stringify({
        userId,
        appId,
        timestamp,
        date,
    });

    activityLog = new activity_pb.ActivityLog();
    activityLog.setUserId(userId);
    activityLog.setApplicationId(appId);
    activityLog.setTimestamp(timestamp);
    activityLog.setDate(date);

    const params = new any.Any();
    params.setValue(activityLog.serializeBinary());
    params.setTypeUrl('github.com/propsproject/pending-props/protos/pending_props_pb.ActivityLog');
    const request = newRPCRequest(params, payloads_pb.Method.ACTIVITY_LOG);
    const requestBytes = request.serializeBinary();

    const activityAddress = CONFIG
        .earnings
        .namespaces
        .activityLogAddress(date, appId, userId);

    log(activityAddress);

    const inputs = [activityAddress];
    const outputs = [activityAddress];

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf
        .TransactionHeader
        .encode({
            familyName: CONFIG.earnings.familyName,
            familyVersion: CONFIG.earnings.familyVersion,
            inputs: inputs,
            outputs: outputs,
            signerPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, we're signing the batch with the same private key, but the
            // batch can be signed by another party, in which case, the public key will need
            // to be associated with that key.
            batcherPublicKey: signer
                .getPublicKey()
                .asHex(),
            // In this example, there are no dependencies.  This list should include an
            // previous transaction header signatures that must be applied for this
            // transaction to successfully commit. For example, dependencies:
            // ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c
            // 122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
            dependencies: [],
            payloadSha512: createHash('sha512')
                .update(requestBytes)
                .digest('hex')
        })
        .finish();

    return await submitTransaction(transactionHeaderBytes, requestBytes);
}

const getPendingEarningsAddress = async (recipient, owner) => {
    recipient = normalizeAddress(recipient);
    owner = normalizeAddress(owner);
    try {
        const addressArgs = [{
            data: recipient,
            start: 0,
            end: 4
        }, {
            data: owner,
            start: 0,
            end: 4
        }];

        //compute state address
        const queryAddress = CONFIG
            .earnings
            .namespaces
            .earningAddress("pending", addressArgs);
        const reqConfig = {
            method: 'GET',
            url: `${restAPIHost}/state?address=${queryAddress}`,
            headers: {
                'Content-Type': 'application/json'
            }
        };
        const response = await axios(reqConfig);
        const data = response.data.data;
        return data.map(d => d.address);
    } catch (e) {
        throw e;
    }
};

const normalizeAddress = (str) => {
    if (str.length > 0) {
        if (str.substr(0,2) === '0x') {
            return str.toLowerCase();
        } else {
            return `0x${str.toLowerCase()}`
        }
    }
    return str;
};

// const stripPrefix = (str) => {
//     return str.substr(0, 2) === "0x" ?
//         str.substr(2) :
//         str
// };

const padWithZeros = (address) => {
    for (let i = 0; address.length + i < 70; i++) {
        address = address.concat(`${ 0}`);
    }

    return address;
};

module.exports = {
    issue,
    revoke,
    settle,
    externalBalanceUpdate,
    queryState,
    updateLastBlockId,
    linkWallet,
    logActivity,
    CONFIG,
    setLoggerType,
    signMessage,
    recoverFromSignature,
};

// const secp256k1 = require('secp256k1');
//
// const isValidSig = (msg, sig, pubK) => {     const hash = hashIt(msg);
//
//     console.log(Buffer.from(pubK).length);     return
// secp256k1.verify(Buffer.from(hash, 'hex'), Buffer.from(hash),
// Buffer.from(pubK)) };
//
// const hashIt = (data) => {     const hash =
// createHash('sha256').update(data).digest('hex');     console.log("HASH:",
// hash);     return hash };
//
// const doTheThing = (msg, sig, pubK) => {     if(isValidSig(msg, sig, pubK)) {
//         return ethUtil.pubToAddress(pubk)     }
//
//     throw new Error("invalid signature") };
//
// const s =
// "67ddf55291cd9ab1e8179254086386c107edf94767324bbd8820bc3c71d9009a1c01656dcc1ec
// fbc4a63dd1446c2e71460abe59ce1a758e0f1180151744c55561b"; const p =
// "b95a3c633d02b64a59acaf603f4e1776be4416f68b8259eec2b992902e4c1a2a52a80e362ed95
// 61b795da4c126df1d6e230c6f593f5c922aba877501a96c5d9"; const msg = "newuser";
//
// isValidSig(msg, s, p); console.log(signer.getPublicKey().asBytes().length);
// console.log(signer.getPublicKey().asHex().length);
// console.log(signer.getPublicKey().asHex());
// console.log(ethUtil.pubToAddress(signer.getPublicKey().asBytes(),
// true).toString('hex'));