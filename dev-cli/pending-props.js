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
const transaction_pb = require('./proto/transaction_pb');
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

const transactionTypes = {
    ISSUE: payloads_pb.Method.ISSUE,
    REVOKE: payloads_pb.Method.REVOKE,
    SETTLE: payloads_pb.Method.SETTLE
}
// hard coded example private key
// const pk = Buffer.from("196749ed808372060eaeffe10e56de82a48829fcf52199847e1e1db4b780ced0", 'hex');
let pk = Buffer.from("5895c973a69c4fe662fcda172900a98bb918c0c31bf374f1b781bc34531cce3f", 'hex');

let priv = new Secp256k1PrivateKey(pk);
let signer = new CryptoFactory(context).newSigner(priv);

let loggerType = 0; // 0 console, -1 none

const newSigner = (pk) => {
    priv = new Secp256k1PrivateKey(Buffer.from(pk, 'hex'));
    signer = new CryptoFactory(context).newSigner(priv);
}

const setLoggerType = (type) => {
    loggerType = type;
}

const signMessage = async (msg, address, _pk) => {
    const privateKey = _pk;
    const account = web3.eth.accounts.privateKeyToAccount('0x' + privateKey);
    // web3.eth.accounts.wallet.add(account);
    // web3.eth.defaultAccount = account.address;
    // return web3.eth.sign(msg, address);
    const signed = account.sign(msg);
    return signed.signature;
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
                "settlement": createHash('sha512')
                    .update("pending-props:earnings:settlements")
                    .digest('hex')
                    .substring(0, 6),
                "balance": createHash('sha512')
                    .update("pending-props:earnings:balance")
                    .digest('hex')
                    .substring(0, 6),
                "transaction": createHash('sha512')
                    .update("pending-props:earnings:transaction")
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
                "rewardEntity": createHash('sha512')
                    .update("pending-props:earnings:rewardentity")
                    .digest('hex')
                    .substring(0,6),
            },
            rewardEntityAddressBySidechainAddress(address, type) {
                const prefix = this.prefixes.rewardEntity;
                const part1 = createHash('sha512')
                    .update(`${normalizeAddress(address)}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 60);
                const part2 = createHash('sha512')
                    .update(`${type}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 4);
                return `${prefix}${part1}${part2}`
            },
            rewardEntityAddressByRewardsAddress(address, type) {
                const prefix = this.prefixes.rewardEntity;
                const part1 = createHash('sha512')
                    .update(`${normalizeAddress(address)}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 60);
                const part2 = createHash('sha512')
                    .update(`${type}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 4);
                return `${prefix}${part1}${part2}`
            },
            transactionAddress(type, applicationId, userId, timestamp) {
                const prefix = this.prefixes.transaction
                const part1 = createHash('sha512')
                    .update(`${type}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 2);
                const part2 = createHash('sha512')
                    .update(`${applicationId}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 10);
                const part3 = createHash('sha512')
                    .update(`${userId}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 42);
                const part4 = createHash('sha512')
                    .update(`${timestamp}`)
                    .digest('hex')
                    .toLowerCase()
                    .substring(0, 10);
                return `${prefix}${part1}${part2}${part3}${part4}`
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
            settlementAddress(txHash) {
                const prefix = this.prefixes.settlement;
                const body = createHash('sha512')
                    .update(`${txHash}`)
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

console.log(JSON.stringify(CONFIG.earnings.namespaces));

const newRPCRequest = (params, method) => {
    const reqParams = new payloads_pb.Params();
    reqParams.setData(params);
    const payload = new payloads_pb.RPCRequest();
    payload.setMethod(method);
    payload.setParams(reqParams);

    return payload;
};

const externalBalanceUpdate = async (address, balance, ethTransactionHash, blockId, timestamp, addresses = {}) => {
    address = normalizeAddress(address);
    const txHash = normalizeAddress(ethTransactionHash);
    const balanceUpdate = new balance_pb.BalanceUpdate();
    balanceUpdate.setPublicAddress(address);
    balanceUpdate.setOnchainBalance(balance);
    balanceUpdate.setTxHash(txHash);
    balanceUpdate.setBlockId(blockId);
    balanceUpdate.setTimestamp(timestamp);
    //setup RPC request
    const params = new any.Any();
    params.setValue(balanceUpdate.serializeBinary());
    params.setTypeUrl('github.com/propsproject/props-transaction-processor/protos/pending_props_pb.BalanceUpdate');

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
            // include the settle transaction address in case this transfer it a settlement
            const transactionAddress = CONFIG
                .earnings
                .namespaces
                .transactionAddress(transactionTypes.SETTLE, linkedApplicationUsers[0].usersList[i].applicationId, linkedApplicationUsers[0].usersList[i].userId, timestamp)
            linkedApplicationUserAddresses.push(transactionAddress);
            addresses['stateAddress'] = transactionAddress;
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
    params.setTypeUrl('github.com/propsproject/props-transaction-processor/protos/pending_props_pb.WalletToUser');

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
// await pendingProps.settle(args.application, args.user, args.amount, args.toaddress, args.fromaddress, args.ethtransactionhash, args.blockid, args.timestamp);
const settle = async (applicationId, userId, amount, toAddress, fromAddress, txHash, blockId, timestamp, addresses = {}) => {

    const settlementData = new payloads_pb.SettlementData();
    settlementData.setApplicationId(applicationId);
    settlementData.setUserId(userId);
    const propsAmount = new BigNumber(amount, 10);
    const tokensAmount = propsAmount.times(1e18);
    settlementData.setAmount(tokensAmount.toString());
    settlementData.setToAddress(toAddress);
    settlementData.setFromAddress(fromAddress);
    settlementData.setTxHash(txHash);
    settlementData.setBlockId(blockId);
    settlementData.setTimestamp(timestamp);

    //setup RPC request
    const params = new any.Any();
    params.setValue(settlementData.serializeBinary());
    params.setTypeUrl('github.com/propsproject/props-transaction-processor/protos/pending_props_pb.SettlementData');

    const request = newRPCRequest(params, payloads_pb.Method.SETTLEMENT);
    const requestBytes = request.serializeBinary();

    // compute state address
    const stateAddress = CONFIG
        .earnings
        .namespaces
        .transactionAddress(transactionTypes.SETTLE, applicationId, userId, timestamp);

    // compute balance and balaneTimestamp addresses for outputs
    const balanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress(applicationId, userId);

    const settlementAddress = CONFIG
        .earnings
        .namespaces
        .settlementAddress(txHash);


    const walletBalanceAddress = CONFIG
        .earnings
        .namespaces
        .balanceAddress("",toAddress);

    const walletLinkAddress = CONFIG
        .earnings
        .namespaces
        .walletLinkAddress(toAddress);

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
    log("Settlement Address = "+settlementAddress)
    log("Transaction Address = "+stateAddress);
    log("Balance Address = "+balanceAddress);
    log("Wallet Balance Address = "+walletBalanceAddress);
    log("Wallet Link Address = "+walletLinkAddress);
    const inputs = [stateAddress, balanceAddress, settlementAddress, walletBalanceAddress, walletLinkAddress,  ...linkedApplicationUserAddresses];
    const outputs = [stateAddress, balanceAddress, settlementAddress, walletBalanceAddress, walletLinkAddress, ...linkedApplicationUserAddresses];
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

    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const transaction = async (transactionType, applicationId, userId, amount, description = '', addresses = {}) => {


    const transactionData = new transaction_pb.Transaction();
    const timestamp = moment().unix();
    transactionData.setType(transactionType);
    transactionData.setTimestamp(timestamp);
    transactionData.setApplicationId(applicationId);
    transactionData.setUserId(userId);
    const propsAmount = new BigNumber(amount, 10);
    const tokensAmount = propsAmount.times(1e18);
    transactionData.setAmount(tokensAmount.toString());
    transactionData.setDescription(description);

    // console.log('transactionData:',transactionType, timestamp, applicationId, userId);
    //setup RPC request
    const params = new any.Any();
    params.setValue(transactionData.serializeBinary());
    params.setTypeUrl('github.com/propsproject/props-transaction-processor/protos/pending_props_pb.Transaction');

    const request = newRPCRequest(params, transactionType);
    const requestBytes = request.serializeBinary();

    // compute state address
    const stateAddress = CONFIG
        .earnings
        .namespaces
        .transactionAddress(transactionType, applicationId, userId, timestamp);

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
    log("Transaction Address = "+stateAddress);
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

    return await submitTransaction(transactionHeaderBytes, requestBytes);
};

const updateLastBlockId = async (blockId) => {
    const blockUpdate = new payloads_pb.LastEthBlock();
    blockUpdate.setId(blockId);
    //setup RPC request
    const params = new any.Any();
    params.setValue(blockUpdate.serializeBinary());
    params.setTypeUrl('github.com/propsproject/props-transaction-processor/protos/pending_props_pb.LastEthBlock');

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

        if (t === "transaction") {
            return deserializeTransactions(data);
        } else if (t === "balance") {
            return deserializeBalance(data);
        } else if (t === "lastblockid") {
            return deserializeLastEthBlockId(data);
        } else if (t === "walletlink") {
            return deserializeWalletLink(data);
        } else if (t === "activity") {
            return deserializeActivity(data)
        } else if (t === "settlement") {
            return deserializeSettlement(data)
        } else if (t === "transfer_tx") {
            return deserializeTransferTx(data)
        }else {
            log(`unknown state type ${t} expected earning or settlements`);
        }


    } catch (e) {
        throw e;
    }
};

const deserializeTransactions = (data) => {
    const recipients = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const transaction = new transaction_pb
            .Transaction
            .deserializeBinary(bytes);

        const output = {
            'state-address': entry.address,
            'transaction-data': transaction.toObject(),
            'human-readable': {
                timestamp: moment(new Date(transaction.getTimestamp() * 1000)).format('L'),
            }
        };
        log(prettyjson.render(output));
        recipients.push({ applicationId: transaction.getApplicationId(), userId: transaction.getUserId(), transaction:transaction.toObject()})
    });
    return recipients
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

const deserializeTransferTx = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const transferTxData = new payloads_pb
            .BalanceUpdate
            .deserializeBinary(bytes);
        retData.push(transferTxData.toObject());
        log(prettyjson.render(transferTxData.toObject()));
    });
    return retData;
}

const deserializeSettlement = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const settlementData = new payloads_pb
            .SettlementData
            .deserializeBinary(bytes);
        retData.push(settlementData.toObject());
        log(prettyjson.render(settlementData.toObject()));
    });
    return retData;
}

const deserializeLastEthBlockId = (data) => {
    const retData = []
    data.forEach(entry => {
        const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
        const lastBlock = new payloads_pb
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
    params.setTypeUrl('github.com/propsproject/props-transaction-processor/protos/pending_props_pb.ActivityLog');
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

const calcDay = function(secondsInDay) {
    const currentTimestamp = Math.floor(Date.now()/1000);
    const secondsLeft = secondsInDay - (currentTimestamp % secondsInDay);
    const ret =  {
        rewardsDay: ((currentTimestamp - (currentTimestamp % secondsInDay)) / secondsInDay),
        secondsLeft
    }
    return ret;
}


module.exports = {
    transaction,
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
    transactionTypes,
    newSigner,
    calcDay,
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