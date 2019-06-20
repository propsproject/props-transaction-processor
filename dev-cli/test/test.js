const fs = require('fs');
const pendingProps = require('../pending-props');
const { exec } = require('child_process');
const chai = require('chai');
const chaiAsPromised = require('chai-as-promised');
const waitUntil = require('async-wait-until');
chai.use(chaiAsPromised);
const expect = chai.expect;
const Web3 = require('web3');
const web3 = new Web3();
const BigNumber = require('bignumber.js');
BigNumber.config({ EXPONENTIAL_AT: 1e+9 })

function add(accumulator, a) {
    return accumulator + a;
}

let startTime = Math.floor(Date.now() / 1000);
let earningAddresses = [];
const secondsInDay = 20; // should match development.json "seconds_in_day"
const amounts = [125, 50.5];
const descriptions = ["Broadcasting", "Watching"];
const waitTimeUntilOnChain = 750; // miliseconds
const longerTestWaitMultiplier = 6;

// data about the ethereum transactions we're testing with
const walletAddress =  "0x2d4dcf292bc5bd8d7246099052dfc76b3cdd3524";
const pk = "759b603832da1100ab47c0f4aa6d445637eb5873d25cadd40484c48970b814d7";
const balanceAtBlock = "428521654000000000000000";
const txHash = "0x0d4d80b54378376131e1ec60ee804fa58f0c33151cd340c8a971cca0a4033834";
const blockNum = "3961643";
const timestamp = "1551553747";

const balanceAtBlock2 = "428513884000000000000000";
const txHash2 = "0x9ef12357191c917cbc3c8102c36948dc731b650852448c51f4705d0f30119100";
const blockNum2 = "3966915";
const timestamp2 = "1551632827";

const settlementTxHash = "0xd0dae165cd740518faf212781e4a707a738970c030d7a3b27f04109ca607447e";
const settlementAmount = 1; // which is 1e18 = 1000000000000000000
const settlementBalanceAtBlock = "428511433000000000000000";
const settlementTimestamp = "1553107867";
const settlementBlockNum = "3967331";


// remove pending-props.js logs
pendingProps.setLoggerType(-1);

describe('Sawtooth side chain test', async () => {
    before(async () => {
        // const msgToSign = "some message";
        // const signed = await pendingProps.signMessage(msgToSign, walletAddress, pk);
        // console.log(`signed=${signed}`);
        // const recovered = await pendingProps.recoverFromSignature(msgToSign, signed);
        // console.log(`recovered=${recovered}`);
        // process.exit(0);
        console.log(`will wait for sawtooth to be ready...`);
        let REGEX = /Now building on top of block.*/g;
        // const
        await waitUntil(() => {
            console.log(`still waiting ${ Math.floor(Date.now() / 1000) - startTime}...`);
            const fileContents = fs.readFileSync('/tmp/out.log', 'utf8');
            const results = fileContents.match(REGEX);
            if (results!=null && results.length > 0) {
                return true;
            } else {
                return false;
            }
        }, 90000, 1000);
        // wait 5 mote seconds to make sure everything is ready to go
        // execute tp now that sawtooth is ready
        console.log(`will wait for tp to be ready...`);
        exec('cd ../ && go run cmd/main.go -c -f ./configs/development.json  >> /tmp/out.log 2>> /tmp/out.log && cd dev-cli', (err, stdout, stderr) => {
            if (err) {
                console.log(`node couldn't execute the command: ${err}`);
                return;
            }

            // the *entire* stdout and stderr (buffered)
            console.log(`stdout: ${stdout}`);
            console.log(`stderr: ${stderr}`);
        });
        // execSync('cd ../ && go run cmd/main.go -c -f ./configs/development.json  >> /tmp/out.log 2>> /tmp/out.log && cd dev-cli');
        REGEX = /registered transaction processor.*pending/g;
        startTime = Math.floor(Date.now() / 1000);
        await waitUntil(() => {
            console.log(`still waiting for tp ${ Math.floor(Date.now() / 1000) - startTime}...`);
            const fileContents = fs.readFileSync('/tmp/out.log', 'utf8');
            const results = fileContents.match(REGEX);
            if (results!=null && results.length > 0) {
                return true;
            } else {
                return false;
            }
        }, 90000, 1000);
        global.timeOfStart = Math.floor(Date.now() / 1000);
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now() / 1000) - global.timeOfStart;
            console.log(`waiting for one more second before testing ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > 1)
        }, 10000, 1000);
    });

    it('Successfully issue an earning', async() => {
        const addresses = {};
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user = "user1";
        // issue
        await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app, user, amounts[0], descriptions[0], addresses);
        // console.log(JSON.stringify(addresses));
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
        earningAddresses.push(addresses['stateAddress']);
        // console.log(`earningOnChain=${JSON.stringify(earningOnChain)}`);
        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;
        const earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
        const balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
        const balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);

        // expect earning details to be correct
        expect(earningPropsAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(earningOnChain[0].transaction.userId).to.be.equal(user);
        expect(earningOnChain[0].transaction.applicationId).to.be.equal(app);
        expect(earningOnChain[0].transaction.description).to.be.equal(descriptions[0]);
        expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.ISSUE);

        // expect balance details to be correct
        expect(balancePendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceObj.userId).to.be.equal(user);
        expect(balanceObj.applicationId).to.be.equal(app);
        expect(balanceDetails.lastUpdateType).to.be.equal(0);
        expect(balanceObj.type).to.be.equal(0);
        expect(balanceObj.linkedWallet).to.be.equal("");
    });

    it('Successfully issue another earning', async() => {
        const addresses = {};
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user = "user1";
        // issue
        await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app, user, amounts[1], descriptions[1], addresses);
        // console.log(JSON.stringify(addresses));
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
        earningAddresses.push(addresses['stateAddress']);
        // console.log(`earningOnChain=${JSON.stringify(earningOnChain)}`);
        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;
        const earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
        const balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
        const balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);

        // expect earning details to be correct
        expect(earningPropsAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
        expect(earningOnChain[0].transaction.userId).to.be.equal(user);
        expect(earningOnChain[0].transaction.applicationId).to.be.equal(app);
        expect(earningOnChain[0].transaction.description).to.be.equal(descriptions[1]);
        expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.ISSUE);

        // expect balance details to be correct
        const sum = amounts.slice(0,2).reduce(add).toString();
        expect(balancePendingAmount.div(1e18).toString()).to.be.equal(sum);
        expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(sum);
        expect(balanceObj.userId).to.be.equal('user1');
        expect(balanceObj.applicationId).to.be.equal('0xa80a6946f8af393d422cd6feee9040c25121a3b8');
        expect(balanceDetails.lastUpdateType).to.be.equal(0);
        expect(balanceObj.type).to.be.equal(0);
        expect(balanceObj.linkedWallet).to.be.equal("");
    });

    it('Successfully revoke an amount', async() => {
        const addresses = {};
        const user = "user1";
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, app, user, amounts[0]+amounts[1], descriptions[0]+"-revoke", addresses);
        // console.log(JSON.stringify(revokeAddress));
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
        // console.log(`earningOnChain=${JSON.stringify(earningOnChain)}`);

        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;
        const earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
        const balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
        const balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);

        expect(earningPropsAmount.div(1e18).toString()).to.be.equal((amounts[0]+amounts[1]).toString());
        expect(earningOnChain[0].transaction.userId).to.be.equal(user);
        expect(earningOnChain[0].transaction.applicationId).to.be.equal(app);
        expect(earningOnChain[0].transaction.description).to.be.equal(descriptions[0]+"-revoke");
        expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.REVOKE);

        // expect balance details to be correct
        expect(balancePendingAmount.div(1e18).toString()).to.be.equal('0');
        expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
        expect(balanceObj.userId).to.be.equal(user);
        expect(balanceObj.applicationId).to.be.equal(app);
        expect(balanceDetails.lastUpdateType).to.be.equal(0);
        expect(balanceObj.type).to.be.equal(0);
        expect(balanceObj.linkedWallet).to.be.equal("");
    });


    it('Successfully update mainchain balance', async() => {
        await pendingProps.externalBalanceUpdate(walletAddress, balanceAtBlock, txHash, blockNum, timestamp);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > (waitTimeUntilOnChain*longerTestWaitMultiplier))
        }, 10000, 100);
        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;

        // expect balance details to be correct
        expect(balanceDetails.pending).to.be.equal('0');
        expect(balanceDetails.totalPending).to.be.equal('0');
        expect(balanceDetails.transferable).to.be.equal(balanceAtBlock);
        expect(balanceObj.userId).to.be.equal(walletAddress);
        expect(balanceObj.applicationId).to.be.equal('');
        expect(balanceDetails.lastUpdateType).to.be.equal(1);
        expect(balanceObj.type).to.be.equal(1);
        expect(balanceObj.linkedWallet).to.be.equal("");
    });

    it('Successfully link app user to wallet', async() => {
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user = "user1";
        const sig = await pendingProps.signMessage(`${app}_${user}`, walletAddress, pk); // "signature11";
        const testSig =  await pendingProps.signMessage(`${app}_8195af8336c01e8014348a906b6adfcf`, walletAddress, pk); // "signature11";
        // console.log(`testSig=${testSig}`);
        // issue
        await pendingProps.linkWallet(walletAddress, app, user, sig);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress);
        const walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
        const walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
        // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);

        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

        const userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
        // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);



        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;
        const userBalanceObj = userBalanceOnChain[0];
        const userBalanceDetails = userBalanceObj.balanceDetails;
        const walletApplicationUsers = walletLinkOnChain[0].usersList;

        // expect balance details to be correct
        expect(balanceDetails.pending).to.be.equal('0');
        expect(balanceDetails.totalPending).to.be.equal('0');
        expect(balanceDetails.transferable).to.be.equal(balanceAtBlock);
        expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp,10));
        expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum,10));
        expect(balanceObj.userId).to.be.equal(walletAddress);
        expect(balanceObj.applicationId).to.be.equal('');
        expect(balanceDetails.lastUpdateType).to.be.equal(1);
        expect(balanceObj.type).to.be.equal(1);
        expect(balanceObj.linkedWallet).to.be.equal("");

        // expect user balance details to be correct
        expect(userBalanceDetails.pending).to.be.equal('0');
        expect(userBalanceDetails.totalPending).to.be.equal('0');
        expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock);
        expect(userBalanceObj.userId).to.be.equal(user);
        expect(userBalanceObj.applicationId).to.be.equal(app);
        expect(userBalanceObj.type).to.be.equal(0);
        expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);

        // expect wallet link is correctly set up
        expect(walletApplicationUsers.length).to.be.equal(1);
        expect(walletLinkOnChain[0].address).to.be.equal(walletAddress);
        expect(walletApplicationUsers[0].applicationId).to.be.equal(app);
        expect(walletApplicationUsers[0].userId).to.be.equal(user);
        expect(walletApplicationUsers[0].signature).to.be.equal(sig);

    });

    it('Successfully update last eth block Id', async() => {
        // issue
        await pendingProps.updateLastBlockId(blockNum);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const lastEthBlockAddress = pendingProps.CONFIG.earnings.namespaces.blockUpdateAddress();
        const ethBlockOnChain = await pendingProps.queryState(lastEthBlockAddress, 'lastblockid');
        // console.log(`ethBlockOnChain=${JSON.stringify(ethBlockOnChain)}`);

        // expect last eth block to be correct
        expect(ethBlockOnChain[0].id).to.be.equal(parseInt(blockNum, 10));
    });

    it('Successfully issue an earning to user with linked wallet', async() => {
        const addresses = {};
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user = "user1";
        // issue
        await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app, user, amounts[0], descriptions[0], addresses);
        // console.log(JSON.stringify(addresses));
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');

        const walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');

        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
        earningAddresses.push(addresses['stateAddress']);

        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;

        const walletBalanceObj = walletBalanceOnChain[0];
        const walletBalanceDetails = walletBalanceObj.balanceDetails;

        const balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
        const balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
        const balanceTransferableAmount = new BigNumber(balanceDetails.transferable, 10);

        const walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
        const walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
        const walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);

        const expectedTransferableAmount = new BigNumber(balanceAtBlock, 10);

        // expect balance details to be correct considering the linked wallet
        expect(balancePendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceTransferableAmount.toString()).to.be.equal(expectedTransferableAmount.toString());

        expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
        expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(walletBalanceTransferableAmount.toString()).to.be.equal(expectedTransferableAmount.toString());
        expect(balanceObj.userId).to.be.equal(user);
        expect(balanceObj.applicationId).to.be.equal(app);
        expect(balanceObj.linkedWallet).to.be.equal(walletAddress);
    });

    it('Successfully revoke an earning of a user with a linked wallet', async() => {
        const addresses = {};
        const user = "user1";
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, app, user, amounts[0], descriptions[0], addresses);
        // console.log(JSON.stringify(revokeAddress));
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);

        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');

        const walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;

        const walletBalanceObj = walletBalanceOnChain[0];
        const walletBalanceDetails = walletBalanceObj.balanceDetails;

        const balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
        const balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
        const balanceTotalTransferable = new BigNumber(balanceDetails.transferable, 10);

        const walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
        const walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
        const walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);


        const expectedTransferableAmount = new BigNumber(balanceAtBlock, 10);

        // expect balance details to be correct considering the linked wallet
        expect(balancePendingAmount.toString()).to.be.equal('0');
        expect(balanceTotalPendingAmount.toString()).to.be.equal('0');
        expect(balanceTotalTransferable.toString()).to.be.equal(expectedTransferableAmount.toString());

        expect(walletBalancePendingAmount.toString()).to.be.equal('0');
        expect(walletBalanceTotalPendingAmount.toString()).to.be.equal('0');
        expect(walletBalanceTransferableAmount.toString()).to.be.equal(expectedTransferableAmount.toString());

        expect(balanceObj.userId).to.be.equal(user);
        expect(balanceObj.applicationId).to.be.equal(app);
        expect(balanceObj.linkedWallet).to.be.equal(walletAddress);
    });

    it('Successfully update mainchain balance of a linked wallet (2nd update)', async() => {
        const user = "user1";
        const app = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";

        await pendingProps.externalBalanceUpdate(walletAddress, balanceAtBlock2, txHash2, blockNum2, timestamp2);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > (waitTimeUntilOnChain*longerTestWaitMultiplier))
        }, 300000, 100);

        const walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
        const walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
        // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);

        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;


        const userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
        const userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
        // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);
        const userBalanceObj = userBalanceOnChain[0];
        const userBalanceDetails = userBalanceObj.balanceDetails;

        // expect balance details to be correct
        expect(balanceDetails.pending).to.be.equal('0');
        expect(balanceDetails.totalPending).to.be.equal('0');
        expect(balanceDetails.transferable).to.be.equal(balanceAtBlock2);
        expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp2,10));
        expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum2,10));
        expect(balanceObj.userId).to.be.equal(walletAddress);
        expect(balanceObj.applicationId).to.be.equal('');
        expect(balanceDetails.lastUpdateType).to.be.equal(1);
        expect(balanceObj.type).to.be.equal(1);
        expect(balanceObj.linkedWallet).to.be.equal("");

        // expect user balance details to be correct
        expect(userBalanceDetails.pending).to.be.equal('0');
        expect(userBalanceDetails.totalPending).to.be.equal('0');
        expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
        expect(userBalanceObj.userId).to.be.equal(user);
        expect(userBalanceObj.applicationId).to.be.equal(app);
        expect(userBalanceObj.type).to.be.equal(0);
        expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);



    });

    it('Successfully link another app user to same wallet', async() => {
        const app = "app2";
        const user = "user1";
        const linkedApp = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const linkedUser = "user1";
        const sig =  await pendingProps.signMessage(`${app}_${user}`, walletAddress, pk); // "signature21";

        // issue
        await pendingProps.linkWallet(walletAddress, app, user, sig);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
        const balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress);
        const walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
        const walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
        // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);

        const balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
        // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

        const userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
        // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);

        const balanceObj = balanceOnChain[0];
        const balanceDetails = balanceObj.balanceDetails;
        const userBalanceObj = userBalanceOnChain[0];
        const userBalanceDetails = userBalanceObj.balanceDetails;
        const walletApplicationUsers = walletLinkOnChain[0].usersList;

        // expect balance details to be correct
        expect(balanceDetails.pending).to.be.equal('0');
        expect(balanceDetails.totalPending).to.be.equal('0');
        expect(balanceDetails.transferable).to.be.equal(balanceAtBlock2);
        expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp2,10));
        expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum2,10));
        expect(balanceObj.userId).to.be.equal(walletAddress);
        expect(balanceObj.applicationId).to.be.equal('');
        expect(balanceDetails.lastUpdateType).to.be.equal(1);
        expect(balanceObj.type).to.be.equal(1);
        expect(balanceObj.linkedWallet).to.be.equal("");

        // expect user balance details to be correct
        expect(userBalanceDetails.pending).to.be.equal('0');
        expect(userBalanceDetails.totalPending).to.be.equal('0');
        expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
        expect(userBalanceObj.userId).to.be.equal(user);
        expect(userBalanceObj.applicationId).to.be.equal(app);
        expect(userBalanceObj.type).to.be.equal(0);
        expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);

        // expect wallet link is correctly set up
        expect(walletApplicationUsers.length).to.be.equal(2);
        expect(walletLinkOnChain[0].address).to.be.equal(walletAddress);
        expect(walletApplicationUsers[1].applicationId).to.be.equal(linkedApp);
        expect(walletApplicationUsers[1].userId).to.be.equal(linkedUser);
        expect(walletApplicationUsers[0].applicationId).to.be.equal(app);
        expect(walletApplicationUsers[0].userId).to.be.equal(user);
        expect(walletApplicationUsers[0].signature).to.be.equal(sig);
    });

    it('Successfully issue an earning to user with linked wallet with another user', async() => {
        const addresses = {};
        const app1 = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user1 = "user1";
        const app2 = "app2";
        const user2 = "user1";
        // issue
        await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app1, user1, amounts[0]*2, descriptions[0], addresses);
        // console.log(JSON.stringify(addresses));
        earningAddresses = [];
        earningAddresses.push(addresses['stateAddress']);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const balanceAddress1 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
        const balanceOnChain1 = await pendingProps.queryState(balanceAddress1, 'balance');
        const balanceAddress2 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app2, user2)
        const balanceOnChain2 = await pendingProps.queryState(balanceAddress2, 'balance');
        const walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
        // console.log(`balanceOnChain1=${JSON.stringify(balanceOnChain1)}`);
        // console.log(`balanceOnChain2=${JSON.stringify(balanceOnChain2)}`);

        const balanceObj1 = balanceOnChain1[0];
        const balanceDetails1 = balanceObj1.balanceDetails;
        const balanceObj2 = balanceOnChain2[0];
        const balanceDetails2 = balanceObj2.balanceDetails;
        const walletBalanceObj = walletBalanceOnChain[0];
        const walletBalanceDetails = walletBalanceObj.balanceDetails;

        const balancePendingAmount1 = new BigNumber(balanceDetails1.pending, 10);
        const balancePendingAmount2 = new BigNumber(balanceDetails2.pending, 10);
        const balanceTotalPendingAmount1 = new BigNumber(balanceDetails1.totalPending, 10);
        const balanceTotalPendingAmount2 = new BigNumber(balanceDetails2.totalPending, 10);
        const balanceTransferableAmount1 = new BigNumber(balanceDetails1.transferable, 10);
        const balanceTransferableAmount2 = new BigNumber(balanceDetails2.transferable, 10);
        const walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
        const walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
        const walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);

        // expect balance details to be correct considering the linked wallet
        expect(balancePendingAmount1.div(1e18).toString()).to.be.equal((amounts[0]*2).toString());
        expect(balancePendingAmount2.div(1e18).toString()).to.be.equal('0');
        expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
        expect(balanceTotalPendingAmount1.div(1e18).toString()).to.be.equal((amounts[0]*2).toString());
        expect(balanceTotalPendingAmount2.div(1e18).toString()).to.be.equal((amounts[0]*2).toString());
        expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal((amounts[0]*2).toString());
        expect(balanceTransferableAmount1.toString()).to.be.equal(balanceAtBlock2.toString());
        expect(balanceTransferableAmount2.toString()).to.be.equal(balanceAtBlock2.toString());
        expect(walletBalanceTransferableAmount.toString()).to.be.equal(balanceAtBlock2.toString());
    });

    it('Successfully revoke an earning of a user with a linked wallet with another user', async() => {
        const addresses = {};
        const app1 = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user1 = "user1";
        const app2 = "app2";
        const user2 = "user1";
        await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, app1, user1, amounts[0], descriptions[0], addresses);
        // console.log(JSON.stringify(addresses));
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const balanceAddress1 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
        const balanceOnChain1 = await pendingProps.queryState(balanceAddress1, 'balance');
        const balanceAddress2 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app2, user2)
        const balanceOnChain2 = await pendingProps.queryState(balanceAddress2, 'balance');
        const walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
        // console.log(`balanceOnChain1=${JSON.stringify(balanceOnChain1)}`);
        // console.log(`balanceOnChain2=${JSON.stringify(balanceOnChain2)}`);

        const balanceObj1 = balanceOnChain1[0];
        const balanceDetails1 = balanceObj1.balanceDetails;
        const balanceObj2 = balanceOnChain2[0];
        const balanceDetails2 = balanceObj2.balanceDetails;
        const walletBalanceObj = walletBalanceOnChain[0];
        const walletBalanceDetails = walletBalanceObj.balanceDetails;

        const balancePendingAmount1 = new BigNumber(balanceDetails1.pending, 10);
        const balancePendingAmount2 = new BigNumber(balanceDetails2.pending, 10);
        const balanceTotalPendingAmount1 = new BigNumber(balanceDetails1.totalPending, 10);
        const balanceTotalPendingAmount2 = new BigNumber(balanceDetails2.totalPending, 10);
        const balanceTransferableAmount1 = new BigNumber(balanceDetails1.transferable, 10);
        const balanceTransferableAmount2 = new BigNumber(balanceDetails2.transferable, 10);
        const walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
        const walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
        const walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);

        // expect balance details to be correct considering the linked wallet
        expect(balancePendingAmount1.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balancePendingAmount2.div(1e18).toString()).to.be.equal('0');
        expect(balanceTotalPendingAmount1.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceTotalPendingAmount2.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceTransferableAmount1.toString()).to.be.equal(balanceAtBlock2.toString());
        expect(balanceTransferableAmount2.toString()).to.be.equal(balanceAtBlock2.toString());
        expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
        expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(walletBalanceTransferableAmount.toString()).to.be.equal(balanceAtBlock2.toString());
    });

    it('Issue an earning before and after day change properly update precut and current balance', async() => {
        const addresses = {};
        const app1 = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user1 = "user1";
        const app2 = "app2";
        const user2 = "user1";
        // issue
        // wait for new day to begin
        let currentDay = pendingProps.calcDay(secondsInDay);
        console.log(`waiting for new "day" to start ${currentDay.secondsLeft} seconds...`);
        await waitUntil(() => {
            const day = pendingProps.calcDay(secondsInDay);

            return (day.rewardsDay == (currentDay.rewardsDay + 1))
        }, 30000, 100);
        await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app1, user1, amounts[0], descriptions[0], addresses);

        currentDay = pendingProps.calcDay(secondsInDay);
        console.log(`waiting for new "day" to start ${currentDay.secondsLeft} seconds...`);
        await waitUntil(() => {
            const day = pendingProps.calcDay(secondsInDay);
            return (day.rewardsDay == (currentDay.rewardsDay + 1))
        }, 30000, 100);

        await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app1, user1, amounts[0], descriptions[0], addresses);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > waitTimeUntilOnChain)
        }, 10000, 100);
        const balanceAddress1 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
        const balanceOnChain1 = await pendingProps.queryState(balanceAddress1, 'balance');
        const balanceAddress2 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app2, user2)
        const balanceOnChain2 = await pendingProps.queryState(balanceAddress2, 'balance');
        const walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');

        const balanceObj1 = balanceOnChain1[0];
        const balanceDetails1 = balanceObj1.balanceDetails;
        const balancePreCutoffDetails1 = balanceObj1.preCutoffDetails;
        const balanceObj2 = balanceOnChain2[0];
        const balanceDetails2 = balanceObj2.balanceDetails;
        const balancePreCutoffDetails2 = balanceObj2.preCutoffDetails;
        const walletBalanceObj = walletBalanceOnChain[0];
        const walletBalanceDetails = walletBalanceObj.balanceDetails;
        const walletPreCutoffDetails = walletBalanceObj.preCutoffDetails;

        expect(balancePreCutoffDetails1.timestamp).to.be.lt(balanceDetails1.timestamp);
        expect(balancePreCutoffDetails2.timestamp).to.be.lt(balanceDetails2.timestamp);
        expect(walletPreCutoffDetails.timestamp).to.be.lt(walletBalanceDetails.timestamp);

        const balanceTotalPendingAmount1 = new BigNumber(balanceDetails1.totalPending, 10);
        const balancePreCutoffTotalPendingAmount1 = new BigNumber(balancePreCutoffDetails1.totalPending, 10);
        const balanceTotalPendingAmount2 = new BigNumber(balanceDetails2.totalPending, 10);
        const balancePreCutoffTotalPendingAmount2 = new BigNumber(balancePreCutoffDetails2.totalPending, 10);
        const walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
        const walletPreCutoffBalanceTotalPendingAmount = new BigNumber(walletPreCutoffDetails.totalPending, 10);

        expect(balanceTotalPendingAmount1.minus(balancePreCutoffTotalPendingAmount1).div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(balanceTotalPendingAmount2.minus(balancePreCutoffTotalPendingAmount2).div(1e18).toString()).to.be.equal(amounts[0].toString());
        expect(walletBalanceTotalPendingAmount.minus(walletPreCutoffBalanceTotalPendingAmount).div(1e18).toString()).to.be.equal(amounts[0].toString());

    });

    it('Successfully settle a linked wallet with an external balance update with tx from being the rewardsAddress of the app (hardcoded in TP)', async() => {
        const addresses = {};
        const app1 = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user1 = "user1";
        const app2 = "app2";
        const user2 = "user1";

        /*
        const settlementTxHash = "0xd0dae165cd740518faf212781e4a707a738970c030d7a3b27f04109ca607447e";
        const settlementAmount = 1; // which is 1e18 = 1000000000000000000
        const settlementBalanceAtBlock = "428511433000000000000000";
        const settlementTimestamp = "1553107867";
        const settlementBlockNum = "3967331";

         */

        await pendingProps.externalBalanceUpdate(walletAddress, settlementBalanceAtBlock, settlementTxHash, settlementBlockNum, settlementTimestamp, addresses);
        global.timeOfStart = Math.floor(Date.now());
        // wait a bit for it to be on chain
        await waitUntil(() => {
            const timePassed =  Math.floor(Date.now()) - global.timeOfStart;
            // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
            return (timePassed > (waitTimeUntilOnChain*longerTestWaitMultiplier*2))
        }, 300000, 100);

        const earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
            // console.log(JSON.stringify(earningOnChain));
        const balanceAddress1 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
        const balanceOnChain1 = await pendingProps.queryState(balanceAddress1, 'balance');
        const balanceAddress2 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app2, user2)
        const balanceOnChain2 = await pendingProps.queryState(balanceAddress2, 'balance');
        const walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
        const walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
        //     // console.log(`balanceOnChain1=${JSON.stringify(balanceOnChain1)}`);
        //     // console.log(`balanceOnChain2=${JSON.stringify(balanceOnChain2)}`);
        //
        const balanceObj1 = balanceOnChain1[0];
        const balanceDetails1 = balanceObj1.balanceDetails;
        const balanceObj2 = balanceOnChain2[0];
        const balanceDetails2 = balanceObj2.balanceDetails;
        const walletBalanceObj = walletBalanceOnChain[0];
        const walletBalanceDetails = walletBalanceObj.balanceDetails;


        const earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
        // expect earning details to be correct
        expect(earningPropsAmount.div(1e18).toString()).to.be.equal(settlementAmount.toString());
        expect(earningOnChain[0].transaction.userId).to.be.equal(user1);
        expect(earningOnChain[0].transaction.applicationId).to.be.equal(app1);
        expect(earningOnChain[0].transaction.description).to.be.equal('Settlement');
        expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.SETTLE);
        expect(earningOnChain[0].transaction.txHash).to.be.equal(settlementTxHash);
        expect(earningOnChain[0].transaction.wallet).to.be.equal(walletAddress);

        const balancePendingAmount1 = new BigNumber(balanceDetails1.pending, 10);
        const balancePendingAmount2 = new BigNumber(balanceDetails2.pending, 10);
        const balanceTotalPendingAmount1 = new BigNumber(balanceDetails1.totalPending, 10);
        const balanceTotalPendingAmount2 = new BigNumber(balanceDetails2.totalPending, 10);
        const balanceTransferableAmount1 = new BigNumber(balanceDetails1.transferable, 10);
        const balanceTransferableAmount2 = new BigNumber(balanceDetails2.transferable, 10);
        const walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
        const walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
        const walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);

        // expect balance details to be correct considering the linked wallet
        expect(balancePendingAmount1.div(1e18).toString()).to.be.equal(((amounts[0]*3)-settlementAmount).toString());
        expect(balancePendingAmount2.div(1e18).toString()).to.be.equal('0');
        expect(balanceTotalPendingAmount1.div(1e18).toString()).to.be.equal(((amounts[0]*3)-settlementAmount).toString());
        expect(balanceTotalPendingAmount2.div(1e18).toString()).to.be.equal(((amounts[0]*3)-settlementAmount).toString());
        expect(balanceTransferableAmount1.toString()).to.be.equal(settlementBalanceAtBlock.toString());
        expect(balanceTransferableAmount2.toString()).to.be.equal(settlementBalanceAtBlock.toString());
        expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
        expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(((amounts[0]*3)-settlementAmount).toString());
        expect(walletBalanceTransferableAmount.toString()).to.be.equal(settlementBalanceAtBlock.toString());
    });



    // TODO - add more test for error scenarios such as replaying the same transaction, last eth block smaller than current, bad signatures, etc.

    after(async () => {

    });

});
