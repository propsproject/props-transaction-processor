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
const waitTimeUntilOnChain = 1300; // miliseconds
const longerTestWaitMultiplier = 6;
let balanceUpdateIndex = 0;

// data about the ethereum transactions we're testing with
const sawtoothPk1 = "5895c973a69c4fe662fcda172900a98bb918c0c31bf374f1b781bc34531cce3f";
const sawtoothPk2 = "af37d6a745b32ef52c80b4b6b18560dfd085e5f3e9ee819478b795733a19257c";
const walletAddress =  "0x2d4dcf292bc5bd8d7246099052dfc76b3cdd3524";
const pk = "5895c973a69c4fe662fcda172900a98bb918c0c31bf374f1b781bc34531cce3f";
const balanceAtBlock = "428521654000000000000000";
const txHash = "0x0d4d80b54378376131e1ec60ee804fa58f0c33151cd340c8a971cca0a4033834";
const blockNum = "3961643";
const timestamp = "1551553747";

const balanceAtBlock2 = "428513884000000000000000";
const txHash2 = "0x9ef12357191c917cbc3c8102c36948dc731b650852448c51f4705d0f30119100";
const blockNum2 = "3966915";
const timestamp2 = "1551632827";



// remove pending-props.js logs
pendingProps.setLoggerType(-1);


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

describe('Sawtooth side chain test', () => {
    describe('Successfully update last eth block Id', () => {
        let lastEthBlockAddress, ethBlockOnChain, ethTimestamp;
        before(async () => {
            ethTimestamp =  Math.floor(Date.now() / 1000);
            await pendingProps.updateLastBlockId(blockNum, ethTimestamp);
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            lastEthBlockAddress = pendingProps.CONFIG.earnings.namespaces.blockUpdateAddress();
            ethBlockOnChain = await pendingProps.queryState(lastEthBlockAddress, 'lastblockid');
        });
        it('Last Eth Synched Block details are correct', () => {
            expect(ethBlockOnChain[0].id).to.be.equal(parseInt(blockNum, 10));
            expect(ethBlockOnChain[0].timestamp).to.be.equal(parseInt(ethTimestamp, 10));
        });
    });
    describe('Successfully issue an earning', () => {
        const addresses = {};
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user = "user1";
        let earningOnChain, balanceAddress, balanceOnChain, balanceObj, balanceDetails, earningPropsAmount,
            balancePendingAmount, balanceTotalPendingAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app, user, amounts[0], descriptions[0], addresses);
            balanceUpdateIndex += 1;
            // console.log(JSON.stringify(addresses));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
            earningAddresses.push(addresses['stateAddress']);
            // console.log(`earningOnChain=${JSON.stringify(earningOnChain)}`);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
        });


        it('Transaction details are correct', () => {
            expect(earningPropsAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(earningOnChain[0].transaction.userId).to.be.equal(user);
            expect(earningOnChain[0].transaction.applicationId).to.be.equal(app);
            expect(earningOnChain[0].transaction.description).to.be.equal(descriptions[0]);
            expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.ISSUE);
        });
        it('Balance details are correct', () => {

            expect(balancePendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(balanceObj.userId).to.be.equal(user);
            expect(balanceObj.applicationId).to.be.equal(app);
            expect(balanceDetails.lastUpdateType).to.be.equal(0);
            expect(balanceObj.type).to.be.equal(0);
            expect(balanceObj.linkedWallet).to.be.equal("");
            expect(balanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
    });
    describe('Successfully log an activity', () => {
        const addresses = {};
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user = "user1";
        const timestamp = Math.floor(Date.now() / 1000);
        const rewardsDay = pendingProps.calcRewardsDay(timestamp);
        let activityOnChain, balanceAddress, balanceOnChain, balanceObj, activityObj, balanceDetails,
            balancePendingAmount, balanceTotalPendingAmount;
        before(async () => {
            await pendingProps.logActivity(user, app, timestamp, rewardsDay, addresses);
            // console.log(JSON.stringify(addresses));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            activityOnChain = await pendingProps.queryState(addresses['stateAddress'], 'activity');
            // console.log(`activityOnChain=${JSON.stringify(activityOnChain)}`);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            activityObj = activityOnChain[0];
            balanceObj = activityObj.balance;
            balanceDetails = balanceObj.balanceDetails;
            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
        });
        it('Activity details are correct', () => {
            expect(activityObj.userId).to.be.equal(user);
            expect(activityObj.applicationId).to.be.equal(app);
            expect(activityObj.date).to.be.equal(rewardsDay);
            expect(activityObj.timestamp).to.be.equal(timestamp);
        });
        it('Activity balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(balanceObj.userId).to.be.equal(user);
            expect(balanceObj.applicationId).to.be.equal(app);
            expect(balanceDetails.lastUpdateType).to.be.equal(0);
            expect(balanceObj.type).to.be.equal(0);
            expect(balanceObj.linkedWallet).to.be.equal("");
            expect(balanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
    });
    describe('Successfully issue another earning', () => {
        const addresses = {};
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user = "user1";
        const timestamp = Math.floor(Date.now() / 1000);
        const rewardsDay = pendingProps.calcRewardsDay(timestamp);
        let earningOnChain, balanceAddress, balanceOnChain, balanceObj, balanceDetails, earningPropsAmount,
            balancePendingAmount, balanceTotalPendingAmount, activityAddress, activityOnChain, activityBalanceObj, activityBalanceDetails, activityBalancePendingAmount, activityBalanceTotalPendingAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app, user, amounts[1], descriptions[1], addresses);
            balanceUpdateIndex += 1;
            // console.log(JSON.stringify(addresses));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain*2)
            }, 10000, 100);
            earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
            earningAddresses.push(addresses['stateAddress']);
            // console.log(`earningOnChain=${JSON.stringify(earningOnChain)}`);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            activityAddress = pendingProps.CONFIG.earnings.namespaces.activityLogAddress(rewardsDay.toString(), app, user);
            activityOnChain = await pendingProps.queryState(activityAddress, 'activity');

            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);

            activityBalanceObj = activityOnChain[0].balance;
            activityBalanceDetails = activityBalanceObj.balanceDetails;
            activityBalancePendingAmount = new BigNumber(activityBalanceDetails.pending, 10);
            activityBalanceTotalPendingAmount = new BigNumber(activityBalanceDetails.totalPending, 10);
        });
        it('Transaction details are correct', () => {
            expect(earningPropsAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(earningOnChain[0].transaction.userId).to.be.equal(user);
            expect(earningOnChain[0].transaction.applicationId).to.be.equal(app);
            expect(earningOnChain[0].transaction.description).to.be.equal(descriptions[1]);
            expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.ISSUE);
        });
        it('Balance details are correct', () => {
            const sum = amounts.slice(0, 2).reduce(add).toString();
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal(sum);
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(sum);
            expect(balanceObj.userId).to.be.equal(user);
            expect(balanceObj.applicationId).to.be.equal(app);
            expect(balanceDetails.lastUpdateType).to.be.equal(0);
            expect(balanceObj.type).to.be.equal(0);
            expect(balanceObj.linkedWallet).to.be.equal("");
            expect(balanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
        it('Activity balance details are correctly updated', () => {
            const sum = amounts.slice(0, 2).reduce(add).toString();
            expect(activityBalancePendingAmount.div(1e18).toString()).to.be.equal(sum);
            expect(activityBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(sum);
            expect(activityBalanceObj.userId).to.be.equal(user);
            expect(activityBalanceObj.applicationId).to.be.equal(app);
            expect(activityBalanceDetails.lastUpdateType).to.be.equal(0);
            expect(activityBalanceObj.type).to.be.equal(0);
            expect(activityBalanceObj.linkedWallet).to.be.equal("");
            expect(activityBalanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
    });
    describe('Successfully revoke an amount', () => {
        const addresses = {};
        const user = "user1";
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        let earningOnChain, balanceAddress, balanceOnChain, balanceObj, balanceDetails, earningPropsAmount,
            balancePendingAmount, balanceTotalPendingAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, app, user, amounts[0], descriptions[0] + "-revoke", addresses);
            balanceUpdateIndex += 1;
            // console.log(JSON.stringify(revokeAddress));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            earningOnChain = await pendingProps.queryState(addresses['stateAddress'], 'transaction');
            // console.log(`earningOnChain=${JSON.stringify(earningOnChain)}`);

            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
        });
        it('Transaction details are correct', () => {

            expect(earningPropsAmount.div(1e18).toString()).to.be.equal((amounts[0]).toString());
            expect(earningOnChain[0].transaction.userId).to.be.equal(user);
            expect(earningOnChain[0].transaction.applicationId).to.be.equal(app);
            expect(earningOnChain[0].transaction.description).to.be.equal(descriptions[0] + "-revoke");
            expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.REVOKE);
        });

        it('Balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(balanceObj.userId).to.be.equal(user);
            expect(balanceObj.applicationId).to.be.equal(app);
            expect(balanceDetails.lastUpdateType).to.be.equal(0);
            expect(balanceObj.type).to.be.equal(0);
            expect(balanceObj.linkedWallet).to.be.equal("");
            expect(balanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
    });

    describe('Successfully update balance from mainchain transaction', () => {
        let balanceAddress, balanceOnChain, balanceObj, balanceDetails;
        let txAddress, txData;
        before(async () => {
            await pendingProps.externalBalanceUpdate(walletAddress, balanceAtBlock, txHash, blockNum, timestamp);
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > (waitTimeUntilOnChain * longerTestWaitMultiplier))
            }, 20000, 100);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
            // console.log(`test balanceAddress=${balanceAddress}`);
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            txAddress = pendingProps.CONFIG.earnings.namespaces.balanceUpdateAddress(txHash, walletAddress);
            txData = await pendingProps.queryState(txAddress, 'transfer_tx');
        });

        it('Balance details are correct', async () => {
            expect(balanceDetails.pending).to.be.equal('0');
            expect(balanceDetails.totalPending).to.be.equal('0');
            expect(balanceDetails.transferable).to.be.equal(balanceAtBlock);
            expect(balanceObj.userId).to.be.equal(walletAddress);
            expect(balanceObj.applicationId).to.be.equal('');
            expect(balanceDetails.lastUpdateType).to.be.equal(1);
            expect(balanceObj.type).to.be.equal(1);
            expect(balanceObj.linkedWallet).to.be.equal("");
        });

        it('Balance update transaction details are correct', async () => {
            expect(txData[0].publicAddress).to.be.equal(walletAddress);
            expect(txData[0].txHash).to.be.equal(txHash);
            expect(txData[0].blockId.toString()).to.be.equal(blockNum);
            expect(txData[0].timestamp.toString()).to.be.equal(timestamp);
            expect(txData[0].onchainBalance).to.be.equal(balanceAtBlock);
        });

    });

    describe('Successfully link app user to wallet', () => {
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user = "user1";
        const walletPk = "759b603832da1100ab47c0f4aa6d445637eb5873d25cadd40484c48970b814d7";
        let sig, userBalanceAddress, balanceAddress, walletLinkAddress, walletLinkOnChain, balanceOnChain,
            userBalanceOnChain,
            balanceObj, balanceDetails, userBalanceObj, userBalanceDetails, walletApplicationUsers,
            balancePendingAmount, balanceTotalPendingAmount,
            userBalancePendingAmount, userBalanceTotalPendingAmount;
        before(async () => {
            sig = await pendingProps.signMessage(`${app}_${user}`, walletAddress, walletPk);
            await pendingProps.linkWallet(walletAddress, app, user, sig);
            balanceUpdateIndex += 1;
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress);
            walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
            walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
            // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);

            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

            userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
            // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);

            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            userBalanceObj = userBalanceOnChain[0];
            userBalanceDetails = userBalanceObj.balanceDetails;
            walletApplicationUsers = walletLinkOnChain[0].usersList;

            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
        });
        it('Wallet link data is correct', () => {
            expect(walletApplicationUsers.length).to.be.equal(1);
            expect(walletLinkOnChain[0].address).to.be.equal(walletAddress);
            expect(walletApplicationUsers[0].applicationId).to.be.equal(app);
            expect(walletApplicationUsers[0].userId).to.be.equal(user);
            expect(walletApplicationUsers[0].signature).to.be.equal(sig);
        });
        it('Wallet balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceDetails.transferable).to.be.equal(balanceAtBlock);
            expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp, 10));
            expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum, 10));
            expect(balanceObj.userId).to.be.equal(walletAddress);
            expect(balanceObj.applicationId).to.be.equal('');
            expect(balanceDetails.lastUpdateType).to.be.equal(1);
            expect(balanceObj.type).to.be.equal(1);
            expect(balanceObj.linkedWallet).to.be.equal("");
        });
        it('User balance details are correct', () => {

            userBalancePendingAmount = new BigNumber(userBalanceDetails.pending, 10);
            userBalanceTotalPendingAmount = new BigNumber(userBalanceDetails.totalPending, 10);

            expect(userBalancePendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(userBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock);
            expect(userBalanceObj.userId).to.be.equal(user);
            expect(userBalanceObj.applicationId).to.be.equal(app);
            expect(userBalanceObj.type).to.be.equal(0);
            expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);
            expect(userBalanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
    });
    describe('Successfully issue an earning to user with linked wallet', () => {
        const addresses = {};
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user = "user1";
        let balanceAddress, balanceOnChain, walletBalanceAddress, walletBalanceOnChain, balanceObj, balanceDetails,
            walletBalanceObj, walletBalanceDetails,
            balancePendingAmount, balanceTotalPendingAmount, balanceTransferableAmount, walletBalancePendingAmount,
            walletBalanceTotalPendingAmount, walletBalanceTransferableAmount, expectedTransferableAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app, user, amounts[0], descriptions[0], addresses);
            balanceUpdateIndex += 1;
            // console.log(JSON.stringify(addresses));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');

            walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
            walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');

            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            earningAddresses.push(addresses['stateAddress']);

            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;

            walletBalanceObj = walletBalanceOnChain[0];
            walletBalanceDetails = walletBalanceObj.balanceDetails;

            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
            balanceTransferableAmount = new BigNumber(balanceDetails.transferable, 10);

            walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
            walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
            walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);

            expectedTransferableAmount = new BigNumber(balanceAtBlock, 10);
        });
        it('User balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(balanceTransferableAmount.toString()).to.be.equal(expectedTransferableAmount.toString());
            expect(balanceObj.userId).to.be.equal(user);
            expect(balanceObj.applicationId).to.be.equal(app);
            expect(balanceObj.linkedWallet).to.be.equal(walletAddress);
            expect(balanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
        it('Wallet balance details are correct', () => {
            expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTransferableAmount.toString()).to.be.equal(expectedTransferableAmount.toString());
        });
    });
    describe('Successfully revoke an earning of a user with a linked wallet', () => {
        const addresses = {};
        const user = "user1";
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        let balanceAddress, balanceOnChain, walletBalanceAddress, walletBalanceOnChain, balanceObj, balanceDetails,
            walletBalanceObj, walletBalanceDetails,
            balancePendingAmount, balanceTotalPendingAmount, balanceTotalTransferable, walletBalancePendingAmount,
            walletBalanceTotalPendingAmount, walletBalanceTransferableAmount, expectedTransferableAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, app, user, amounts[0], descriptions[0], addresses);
            balanceUpdateIndex += 1;
            // console.log(JSON.stringify(revokeAddress));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);

            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');

            walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
            walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;

            walletBalanceObj = walletBalanceOnChain[0];
            walletBalanceDetails = walletBalanceObj.balanceDetails;

            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
            balanceTotalTransferable = new BigNumber(balanceDetails.transferable, 10);

            walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
            walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
            walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);


            expectedTransferableAmount = new BigNumber(balanceAtBlock, 10);
        });
        it('User balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(balanceTotalTransferable.toString()).to.be.equal(expectedTransferableAmount.toString());
            expect(balanceObj.userId).to.be.equal(user);
            expect(balanceObj.applicationId).to.be.equal(app);
            expect(balanceObj.linkedWallet).to.be.equal(walletAddress);
            expect(balanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
        it('Wallet balance details are correct', () => {
            expect(walletBalancePendingAmount.toString()).to.be.equal('0');
            expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTransferableAmount.toString()).to.be.equal(expectedTransferableAmount.toString());
        });
    });
    describe('Successfully update mainchain balance of a linked wallet with an activity object (2nd update)', () => {
        const user = "user1";
        const app = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const rewardsDay = pendingProps.calcRewardsDay(timestamp2);
        let walletLinkAddress, walletLinkOnChain, balanceAddress, balanceOnChain, balanceObj, balanceDetails,
            userBalanceAddress, userBalanceOnChain, userBalanceObj, userBalanceDetails,
            balanceTotalPendingAmount, userBalancePendingAmount, userBalanceTotalPendingAmount,
            activityAddress, activityOnChain, activityBalanceObj, activityBalanceDetails, activityBalancePendingAmount, activityBalanceTotalPendingAmount;
        before(async () => {
            // create an activity that will get updated by this external transaction
            await pendingProps.logActivity(user, app, timestamp2, rewardsDay);
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain * 2)
            }, 10000, 100);
            await pendingProps.externalBalanceUpdate(walletAddress, balanceAtBlock2, txHash2, blockNum2, timestamp2);
            balanceUpdateIndex += 1;
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > (waitTimeUntilOnChain * longerTestWaitMultiplier))
            }, 300000, 100);

            walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
            walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);
            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;

            userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
            userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
            // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);
            userBalanceObj = userBalanceOnChain[0];
            userBalanceDetails = userBalanceObj.balanceDetails;

            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);

            userBalancePendingAmount = new BigNumber(userBalanceDetails.pending, 10);
            userBalanceTotalPendingAmount = new BigNumber(userBalanceDetails.totalPending, 10);

            activityAddress = pendingProps.CONFIG.earnings.namespaces.activityLogAddress(rewardsDay.toString(), app, user);
            activityOnChain = await pendingProps.queryState(activityAddress, 'activity');
            activityBalanceObj = activityOnChain[0].balance;
            activityBalanceDetails = activityBalanceObj.balanceDetails;
            activityBalancePendingAmount = new BigNumber(activityBalanceDetails.pending, 10);
            activityBalanceTotalPendingAmount = new BigNumber(activityBalanceDetails.totalPending, 10);

        });
        it('Wallet balance details are correct', () => {
            expect(balanceDetails.pending).to.be.equal('0');
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp2, 10));
            expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum2, 10));
            expect(balanceObj.userId).to.be.equal(walletAddress);
            expect(balanceObj.applicationId).to.be.equal('');
            expect(balanceDetails.lastUpdateType).to.be.equal(1);
            expect(balanceObj.type).to.be.equal(1);
            expect(balanceObj.linkedWallet).to.be.equal("");
        });
        it('Linked user balance details are correct', () => {
            expect(userBalancePendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(userBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(userBalanceObj.userId).to.be.equal(user);
            expect(userBalanceObj.applicationId).to.be.equal(app);
            expect(userBalanceObj.type).to.be.equal(0);
            expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);
            expect(userBalanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
        it('Activity balance details are correct', () => {
            expect(activityBalancePendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(activityBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(activityBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(activityBalanceObj.userId).to.be.equal(user);
            expect(activityBalanceObj.applicationId).to.be.equal(app);
            expect(activityBalanceObj.type).to.be.equal(0);
            expect(activityBalanceObj.linkedWallet).to.be.equal(walletAddress);
            expect(activityBalanceObj.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
    });
    describe('Successfully update balances with linked wallet based on settlement transaction', () => {
        const addresses = {};
        const app1 = "0xa80a6946f8af393d422cd6feee9040c25121a3b8";
        const user1 = "32f2be121e8b2efc4b04a45511412f60";
        const user1Wallet = "0x2755ef71ec620348570bd8866045f97aca250ce1";
        const app1user1Sig = "0xbf18182fc0f4c148795fb03cc30bf172f8becd5a8b5d6be4ed3b9d2370bb361d4bee6a12911cd38fab2d5b582192e7bbf5fd156db99ebd1dbd0a7b0c2879d1961c";

        const settlementTxHash = "0x5645a41ccc7c7a757831369677dc1bc39d9c58b1cd8541e2306cfbbb32da0054";
        const settlementAmount = 5; // which is 1e18 = 5000000000000000000
        const settlementTimestamp = "1563692060";
        const settlementBlockNum = "4770812";
        const settlementApplicationRewardsAddress = "0xd8186f92ba7cc1991f6e3ab842cb50c29bbfdc6a";
        const settleTransactionAddress = pendingProps.CONFIG
            .earnings
            .namespaces
            .transactionAddress(pendingProps.transactionTypes.SETTLE, app1, user1, settlementTimestamp);
        let earningOnChain, balanceAddress, balanceOnChain, walletBalanceAddress, walletBalanceOnChain, balanceObj,
            balanceDetails, walletBalanceObj, walletBalanceDetails, earningPropsAmount,
            balancePendingAmount, balanceTotalPendingAmount, balanceTransferableAmount, walletBalancePendingAmount,
            walletBalanceTotalPendingAmount, walletBalanceTransferableAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app1, user1, amounts[0], descriptions[0], addresses);
            await pendingProps.linkWallet(user1Wallet, app1, user1, app1user1Sig);
            await pendingProps.settle(app1, user1, settlementAmount, user1Wallet, settlementApplicationRewardsAddress, settlementTxHash, settlementBlockNum, settlementTimestamp, '25000000000000000000', addresses);


            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > (waitTimeUntilOnChain * longerTestWaitMultiplier * 2))
            }, 300000, 100);

            earningOnChain = await pendingProps.queryState(settleTransactionAddress, 'transaction');
            // console.log(JSON.stringify(earningOnChain));
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", user1Wallet)
            walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
            //     // console.log(`balanceOnChain1=${JSON.stringify(balanceOnChain1)}`);
            //     // console.log(`balanceOnChain2=${JSON.stringify(balanceOnChain2)}`);
            //
            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            walletBalanceObj = walletBalanceOnChain[0];
            walletBalanceDetails = walletBalanceObj.balanceDetails;

            earningPropsAmount = new BigNumber(earningOnChain[0].transaction.amount, 10);
            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
            balanceTransferableAmount = new BigNumber(balanceDetails.transferable, 10);
            walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
            walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
            walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);
        });
        it('Transaction details are correct', () => {
            expect(earningPropsAmount.div(1e18).toString()).to.be.equal(settlementAmount.toString());
            expect(earningOnChain[0].transaction.userId).to.be.equal(user1);
            expect(earningOnChain[0].transaction.applicationId).to.be.equal(app1);
            expect(earningOnChain[0].transaction.description).to.be.equal('Settlement');
            expect(earningOnChain[0].transaction.type).to.be.equal(pendingProps.transactionTypes.SETTLE);
            expect(earningOnChain[0].transaction.txHash).to.be.equal(settlementTxHash);
            expect(earningOnChain[0].transaction.wallet).to.be.equal(user1Wallet);
        });
        it('User balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal((amounts[0] - settlementAmount).toString());
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal((amounts[0] - settlementAmount).toString());
            expect(balanceTransferableAmount.toString()).to.be.equal("25000000000000000000");
            expect(balanceObj.balanceUpdateIndex).to.be.equal(3);
        });
        it('Wallet balance details are correct', () => {
            expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTransferableAmount.toString()).to.be.equal("25000000000000000000");
        });
    });
    describe('Successfully link another app user to same wallet', () => {
        const app = "0x39dbb8ddeb0d0e86f17aa23d9dac4eeb69b76511";
        const user = "user1";
        const linkedApp = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const linkedUser = "user1";
        const walletPk = "759b603832da1100ab47c0f4aa6d445637eb5873d25cadd40484c48970b814d7";
        let sig, userBalanceAddress, balanceAddress, walletLinkAddress, walletLinkOnChain, balanceOnChain,
            userBalanceOnChain, balanceObj, balanceDetails, userBalanceObj, userBalanceDetails,
            walletApplicationUsers, balanceTotalPendingAmount, userBalanceTotalPendingAmount;
        before(async () => {
            sig = await pendingProps.signMessage(`${app}_${user}`, walletAddress, walletPk); // "signature21";
            pendingProps.newSigner(sawtoothPk2);

            // issue
            await pendingProps.linkWallet(walletAddress, app, user, sig);
            balanceUpdateIndex += 1;
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress);
            walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
            walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
            // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);

            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

            userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
            // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);

            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            userBalanceObj = userBalanceOnChain[0];
            userBalanceDetails = userBalanceObj.balanceDetails;
            walletApplicationUsers = walletLinkOnChain[0].usersList;
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
            userBalanceTotalPendingAmount = new BigNumber(userBalanceDetails.totalPending, 10);
        });
        it('Wallet link data is correct', () => {
            expect(walletApplicationUsers.length).to.be.equal(2);
            expect(walletLinkOnChain[0].address).to.be.equal(walletAddress);
            expect(walletApplicationUsers[0].applicationId).to.be.equal(linkedApp);
            expect(walletApplicationUsers[0].userId).to.be.equal(linkedUser);
            expect(walletApplicationUsers[1].applicationId).to.be.equal(app);
            expect(walletApplicationUsers[1].userId).to.be.equal(user);
            expect(walletApplicationUsers[1].signature).to.be.equal(sig);
        });
        it('Wallet balance details are correct', () => {
            expect(balanceDetails.pending).to.be.equal('0');
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp2, 10));
            expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum2, 10));
            expect(balanceObj.userId).to.be.equal(walletAddress);
            expect(balanceObj.applicationId).to.be.equal('');
            expect(balanceDetails.lastUpdateType).to.be.equal(1);
            expect(balanceObj.type).to.be.equal(1);
            expect(balanceObj.linkedWallet).to.be.equal("");
        });
        it('User balance details are correct', () => {
            expect(userBalanceDetails.pending).to.be.equal('0');
            expect(userBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[1].toString());
            expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(userBalanceObj.userId).to.be.equal(user);
            expect(userBalanceObj.applicationId).to.be.equal(app);
            expect(userBalanceObj.type).to.be.equal(0);
            expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);
            expect(userBalanceObj.balanceUpdateIndex).to.be.equal(1);
        });
    });
    describe('Successfully issue an earning to user with linked wallet with another user', () => {
        const addresses = {};
        const app1 = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user1 = "user1";
        const app2 = "0x39dbb8ddeb0d0e86f17aa23d9dac4eeb69b76511";
        const user2 = "user1";
        let balanceAddress1, balanceOnChain1, balanceAddress2, balanceOnChain2, walletBalanceAddress,
            walletBalanceOnChain, balanceObj1, balanceDetails1, balanceObj2, balanceDetails2,
            walletBalanceObj, walletBalanceDetails, balancePendingAmount1, balancePendingAmount2,
            balanceTotalPendingAmount1, balanceTotalPendingAmount2, balanceTransferableAmount1,
            balanceTransferableAmount2,
            walletBalancePendingAmount, walletBalanceTotalPendingAmount, walletBalanceTransferableAmount;
        before(async () => {
            pendingProps.newSigner(sawtoothPk1);
            // issue
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, app1, user1, amounts[0] * 2, descriptions[0], addresses);
            balanceUpdateIndex += 1;
            // console.log(JSON.stringify(addresses));
            earningAddresses = [];
            earningAddresses.push(addresses['stateAddress']);
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            balanceAddress1 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
            balanceOnChain1 = await pendingProps.queryState(balanceAddress1, 'balance');
            balanceAddress2 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app2, user2)
            balanceOnChain2 = await pendingProps.queryState(balanceAddress2, 'balance');
            walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
            walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
            // console.log(`balanceOnChain1=${JSON.stringify(balanceOnChain1)}`);
            // console.log(`balanceOnChain2=${JSON.stringify(balanceOnChain2)}`);

            balanceObj1 = balanceOnChain1[0];
            balanceDetails1 = balanceObj1.balanceDetails;
            balanceObj2 = balanceOnChain2[0];
            balanceDetails2 = balanceObj2.balanceDetails;
            walletBalanceObj = walletBalanceOnChain[0];
            walletBalanceDetails = walletBalanceObj.balanceDetails;

            balancePendingAmount1 = new BigNumber(balanceDetails1.pending, 10);
            balancePendingAmount2 = new BigNumber(balanceDetails2.pending, 10);
            balanceTotalPendingAmount1 = new BigNumber(balanceDetails1.totalPending, 10);
            balanceTotalPendingAmount2 = new BigNumber(balanceDetails2.totalPending, 10);
            balanceTransferableAmount1 = new BigNumber(balanceDetails1.transferable, 10);
            balanceTransferableAmount2 = new BigNumber(balanceDetails2.transferable, 10);
            walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
            walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
            walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);
        });
        it('User1 Balance details are correct', () => {
            expect(balancePendingAmount1.div(1e18).toString()).to.be.equal((amounts[0] * 2 + amounts[1]).toString());
            expect(balanceTotalPendingAmount1.div(1e18).toString()).to.be.equal((amounts[0] * 2 + amounts[1]).toString());
            expect(balanceTransferableAmount1.toString()).to.be.equal(balanceAtBlock2.toString());
            expect(balanceObj1.balanceUpdateIndex).to.be.equal(balanceUpdateIndex);
        });
        it('User2 Balance details are correct', () => {
            expect(balancePendingAmount2.div(1e18).toString()).to.be.equal('0');
            expect(balanceTotalPendingAmount2.div(1e18).toString()).to.be.equal((amounts[0] * 2 + amounts[1]).toString());
            expect(balanceTransferableAmount2.toString()).to.be.equal(balanceAtBlock2.toString());
            expect(balanceObj2.balanceUpdateIndex).to.be.equal(2);
        });

        it('Wallet Balance details are correct', () => {
            expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTransferableAmount.toString()).to.be.equal(balanceAtBlock2.toString());
        });
    });
    describe('Successfully revoke an earning of a user with a linked wallet with another user', () => {
        const addresses = {};
        const app1 = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const user1 = "user1";
        const app2 = "0x39dbb8ddeb0d0e86f17aa23d9dac4eeb69b76511";
        const user2 = "user1";
        let balanceAddress1, balanceOnChain1, balanceAddress2, balanceOnChain2, walletBalanceAddress,
            walletBalanceOnChain, balanceObj1, balanceDetails1, balanceObj2, balanceDetails2,
            walletBalanceObj, walletBalanceDetails, balancePendingAmount1, balancePendingAmount2,
            balanceTotalPendingAmount1, balanceTotalPendingAmount2, balanceTransferableAmount1,
            balanceTransferableAmount2,
            walletBalancePendingAmount, walletBalanceTotalPendingAmount, walletBalanceTransferableAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, app1, user1, amounts[0], descriptions[0], addresses);
            // console.log(JSON.stringify(addresses));
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            balanceAddress1 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app1, user1)
            balanceOnChain1 = await pendingProps.queryState(balanceAddress1, 'balance');
            balanceAddress2 = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app2, user2)
            balanceOnChain2 = await pendingProps.queryState(balanceAddress2, 'balance');
            walletBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress)
            walletBalanceOnChain = await pendingProps.queryState(walletBalanceAddress, 'balance');
            // console.log(`balanceOnChain1=${JSON.stringify(balanceOnChain1)}`);
            // console.log(`balanceOnChain2=${JSON.stringify(balanceOnChain2)}`);

            balanceObj1 = balanceOnChain1[0];
            balanceDetails1 = balanceObj1.balanceDetails;
            balanceObj2 = balanceOnChain2[0];
            balanceDetails2 = balanceObj2.balanceDetails;
            walletBalanceObj = walletBalanceOnChain[0];
            walletBalanceDetails = walletBalanceObj.balanceDetails;

            balancePendingAmount1 = new BigNumber(balanceDetails1.pending, 10);
            balancePendingAmount2 = new BigNumber(balanceDetails2.pending, 10);
            balanceTotalPendingAmount1 = new BigNumber(balanceDetails1.totalPending, 10);
            balanceTotalPendingAmount2 = new BigNumber(balanceDetails2.totalPending, 10);
            balanceTransferableAmount1 = new BigNumber(balanceDetails1.transferable, 10);
            balanceTransferableAmount2 = new BigNumber(balanceDetails2.transferable, 10);
            walletBalancePendingAmount = new BigNumber(walletBalanceDetails.pending, 10);
            walletBalanceTotalPendingAmount = new BigNumber(walletBalanceDetails.totalPending, 10);
            walletBalanceTransferableAmount = new BigNumber(walletBalanceDetails.transferable, 10);
        });
        it('User1 Balance details are correct', () => {
            expect(balancePendingAmount1.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(balanceTotalPendingAmount1.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(balanceTransferableAmount1.toString()).to.be.equal(balanceAtBlock2.toString());
        });
        it('User2 Balance details are correct', () => {
            expect(balancePendingAmount2.div(1e18).toString()).to.be.equal('0');
            expect(balanceTotalPendingAmount2.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(balanceTransferableAmount2.toString()).to.be.equal(balanceAtBlock2.toString());
        });

        it('Wallet Balance details are correct', () => {
            expect(walletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(walletBalanceTransferableAmount.toString()).to.be.equal(balanceAtBlock2.toString());
        });
    });
    describe('Successfully link same wallet to another user with the same app after issue to new user', () => {
        const app = "0x39dbb8ddeb0d0e86f17aa23d9dac4eeb69b76511";
        const user = "user1";
        const user2 = "user2";
        const linkedApp = "0x96c41cfd601a477e80fd9fbf256e767e92ac4278";
        const linkedUser = "user1";
        const walletPk = "759b603832da1100ab47c0f4aa6d445637eb5873d25cadd40484c48970b814d7";
        let sig, userBalanceAddress, oldUserBalanceAddress, otherAppUserBalanceAddress, balanceAddress, walletLinkAddress, walletLinkOnChain, balanceOnChain, userBalanceOnChain,
            oldUserBalanceOnChain, otherAppUserBalanceOnChain, balanceObj, balanceDetails, userBalanceObj, userBalanceDetails, oldUserBalanceObj, oldUserBalanceDetails, otherAppUserBalanceObj,
            otherAppUserBalanceDetails, walletApplicationUsers, balanceTotalPendingAmount, userBalanceTotalPendingAmount, userBalancePendingAmount, oldUserBalanceTotalPendingAmount, oldUserBalancePendingAmount, otherAppUserBalanceTotalPendingAmount;
        before(async () => {
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, linkedApp, user2, amounts[0], descriptions[0], {});

            sig = await pendingProps.signMessage(`${linkedApp}_${user2}`, walletAddress, walletPk); // "signature21";
            pendingProps.newSigner(sawtoothPk2);

            await pendingProps.linkWallet(walletAddress, linkedApp, user2, sig);
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                // console.log(`waiting for transaction ${ Math.floor(Date.now() / 1000) - global.timeOfStart}...`);
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(linkedApp, user2);
            oldUserBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(linkedApp, linkedUser);
            otherAppUserBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress);
            walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddress);
            walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
            // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);

            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

            userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
            oldUserBalanceOnChain = await pendingProps.queryState(oldUserBalanceAddress, 'balance');
            otherAppUserBalanceOnChain = await pendingProps.queryState(otherAppUserBalanceAddress, 'balance');
            // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);

            balanceObj = balanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            userBalanceObj = userBalanceOnChain[0];
            userBalanceDetails = userBalanceObj.balanceDetails;
            oldUserBalanceObj = oldUserBalanceOnChain[0];
            oldUserBalanceDetails = oldUserBalanceObj.balanceDetails;
            otherAppUserBalanceObj = otherAppUserBalanceOnChain[0];
            otherAppUserBalanceDetails = otherAppUserBalanceObj.balanceDetails;
            walletApplicationUsers = walletLinkOnChain[0].usersList;
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
            userBalanceTotalPendingAmount = new BigNumber(userBalanceDetails.totalPending, 10);
            userBalancePendingAmount = new BigNumber(userBalanceDetails.pending, 10);
            oldUserBalanceTotalPendingAmount = new BigNumber(oldUserBalanceDetails.totalPending, 10);
            oldUserBalancePendingAmount = new BigNumber(oldUserBalanceDetails.pending, 10);
            otherAppUserBalanceTotalPendingAmount = new BigNumber(otherAppUserBalanceDetails.totalPending, 10);

        });
        it('Unlinked user balance details are correct', () => {
            expect(oldUserBalancePendingAmount.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(oldUserBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal((amounts[0] + amounts[1]).toString());
            expect(oldUserBalanceDetails.transferable).to.be.equal('0');
            expect(oldUserBalanceObj.userId).to.be.equal(linkedUser);
            expect(oldUserBalanceObj.applicationId).to.be.equal(linkedApp);
            expect(oldUserBalanceObj.type).to.be.equal(0);
            expect(oldUserBalanceObj.linkedWallet).to.be.equal('');
        });
        it('Newly linked user balance details are correct', () => {
            expect(userBalancePendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(userBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(userBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(userBalanceObj.userId).to.be.equal(user2);
            expect(userBalanceObj.applicationId).to.be.equal(linkedApp);
            expect(userBalanceObj.type).to.be.equal(0);
            expect(userBalanceObj.linkedWallet).to.be.equal(walletAddress);
        });
        it('Other linked user balance details are correct', () => {
            expect(otherAppUserBalanceDetails.pending).to.be.equal('0');
            expect(otherAppUserBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal(amounts[0].toString());
            expect(otherAppUserBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(otherAppUserBalanceObj.userId).to.be.equal(user);
            expect(otherAppUserBalanceObj.applicationId).to.be.equal(app);
            expect(otherAppUserBalanceObj.type).to.be.equal(0);
            expect(otherAppUserBalanceObj.linkedWallet).to.be.equal(walletAddress);
        });
        it('Wallet balance details are correct', () => {
            expect(balanceDetails.pending).to.be.equal('0');
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(balanceDetails.timestamp).to.be.equal(parseInt(timestamp2, 10));
            expect(balanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum2, 10));
            expect(balanceObj.userId).to.be.equal(walletAddress);
            expect(balanceObj.applicationId).to.be.equal('');
            expect(balanceDetails.lastUpdateType).to.be.equal(1);
            expect(balanceObj.type).to.be.equal(1);
            expect(balanceObj.linkedWallet).to.be.equal("");
        });
        it('Wallet link data is correct', () => {
            expect(walletApplicationUsers.length).to.be.equal(2);
            expect(walletLinkOnChain[0].address).to.be.equal(walletAddress);
            expect(walletApplicationUsers[1].applicationId).to.be.equal(linkedApp);
            expect(walletApplicationUsers[1].userId).to.be.equal(user2);
            expect(walletApplicationUsers[0].applicationId).to.be.equal(app);
            expect(walletApplicationUsers[0].userId).to.be.equal(user);
            expect(walletApplicationUsers[1].signature).to.be.equal(sig);
        });
    });

    describe('Successfully link another wallet with no props to an app user with an existing wallet', () => {
        const app = "0x39dbb8ddeb0d0e86f17aa23d9dac4eeb69b76511";
        const user = "user1";
        const walletPk = "54fea0e281b1f5add50c7b87bb8ac5feb6d84d8b5d7f4a1b643e0c7ed9924114";
        const walletAddr = "0xf731762520db9f451c7f883475fb7a27afba5516";
        const oldWalletAddr = walletAddress;
        let sig, userBalanceAddress, balanceAddress, walletLinkAddress, walletLinkOnChain, balanceOnChain,
            userBalanceOnChain, oldWalletLinkAddress, oldWalletLinkOnChain,
            balanceObj, balanceDetails, userBalanceObj, userBalanceDetails, walletApplicationUsers, oldWalletApplicationUsers,
            balancePendingAmount, balanceTotalPendingAmount, oldWalletBalanceDetails, oldWalletBalancePendingAmount, oldWalletBalanceTotalPendingAmount,
            userBalancePendingAmount, userBalanceTotalPendingAmount, oldWalletBalanceObj, oldWalletbalanceAddress, oldWalletbalanceOnChain;
        before(async () => {
            sig = await pendingProps.signMessage(`${app}_${user}`, walletAddr, walletPk);
            await pendingProps.linkWallet(walletAddr, app, user, sig);
            balanceUpdateIndex += 1;
            global.timeOfStart = Math.floor(Date.now());
            // wait a bit for it to be on chain
            await waitUntil(() => {
                const timePassed = Math.floor(Date.now()) - global.timeOfStart;
                return (timePassed > waitTimeUntilOnChain)
            }, 10000, 100);
            userBalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress(app, user);
            balanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddr);
            oldWalletbalanceAddress = pendingProps.CONFIG.earnings.namespaces.balanceAddress("", walletAddress);
            walletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(walletAddr);
            oldWalletLinkAddress = pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(oldWalletAddr);
            walletLinkOnChain = await pendingProps.queryState(walletLinkAddress, 'walletlink');
            oldWalletLinkOnChain = await pendingProps.queryState(oldWalletLinkAddress, 'walletlink');
            // console.log(`walletLinkOnChain=${JSON.stringify(walletLinkOnChain)}`);
            // console.log(`oldWalletLinkOnChain=${JSON.stringify(oldWalletLinkOnChain)}`);

            balanceOnChain = await pendingProps.queryState(balanceAddress, 'balance');
            // console.log(`balanceOnChain=${JSON.stringify(balanceOnChain)}`);

            oldWalletbalanceOnChain = await pendingProps.queryState(oldWalletbalanceAddress, 'balance');
            // console.log(`oldWalletbalanceOnChain=${JSON.stringify(oldWalletbalanceOnChain)}`);

            userBalanceOnChain = await pendingProps.queryState(userBalanceAddress, 'balance');
            // console.log(`userBalanceOnChain=${JSON.stringify(userBalanceOnChain)}`);

            balanceObj = balanceOnChain[0];
            oldWalletBalanceObj = oldWalletbalanceOnChain[0];
            balanceDetails = balanceObj.balanceDetails;
            oldWalletBalanceDetails = oldWalletBalanceObj.balanceDetails;
            userBalanceObj = userBalanceOnChain[0];
            userBalanceDetails = userBalanceObj.balanceDetails;
            walletApplicationUsers = walletLinkOnChain[0].usersList;
            oldWalletApplicationUsers = oldWalletLinkOnChain[0].usersList;
            balancePendingAmount = new BigNumber(balanceDetails.pending, 10);
            balanceTotalPendingAmount = new BigNumber(balanceDetails.totalPending, 10);
            oldWalletBalancePendingAmount = new BigNumber(oldWalletBalanceDetails.pending, 10);
            oldWalletBalanceTotalPendingAmount = new BigNumber(oldWalletBalanceDetails.totalPending, 10);
        });
        it('Old Wallet link data is correct', () => {
            expect(oldWalletApplicationUsers.length).to.be.equal(1);
            expect(oldWalletLinkOnChain[0].address).to.be.equal(oldWalletAddr);
            expect(oldWalletApplicationUsers[0].applicationId).to.be.equal('0x96c41cfd601a477e80fd9fbf256e767e92ac4278');
            expect(oldWalletApplicationUsers[0].userId).to.be.equal("user2");
        });
        it('New Wallet link data is correct', () => {
            expect(walletApplicationUsers.length).to.be.equal(1);
            expect(walletLinkOnChain[0].address).to.be.equal(walletAddr);
            expect(walletApplicationUsers[0].applicationId).to.be.equal(app);
            expect(walletApplicationUsers[0].userId).to.be.equal(user);
            expect(walletApplicationUsers[0].signature).to.be.equal(sig);
        });
        it('Wallet balance details are correct', () => {
            expect(balancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(balanceDetails.transferable).to.be.equal('0');
            expect(balanceObj.userId).to.be.equal(walletAddr);
            expect(balanceObj.applicationId).to.be.equal('');
            expect(balanceObj.linkedWallet).to.be.equal("");
        });
        it('Old Wallet balance details are correct', () => {
            expect(oldWalletBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(oldWalletBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(oldWalletBalanceDetails.transferable).to.be.equal(balanceAtBlock2);
            expect(oldWalletBalanceDetails.timestamp).to.be.equal(parseInt(timestamp2, 10));
            expect(oldWalletBalanceDetails.lastEthBlockId).to.be.equal(parseInt(blockNum2, 10));
            expect(oldWalletBalanceObj.userId).to.be.equal(oldWalletAddr);
            expect(oldWalletBalanceObj.applicationId).to.be.equal('');
            expect(oldWalletBalanceDetails.lastUpdateType).to.be.equal(1);
            expect(oldWalletBalanceObj.type).to.be.equal(1);
            expect(oldWalletBalanceObj.linkedWallet).to.be.equal("");
        });
        it('User balance details are correct', () => {

            userBalancePendingAmount = new BigNumber(userBalanceDetails.pending, 10);
            userBalanceTotalPendingAmount = new BigNumber(userBalanceDetails.totalPending, 10);

            expect(userBalancePendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(userBalanceTotalPendingAmount.div(1e18).toString()).to.be.equal('0');
            expect(userBalanceDetails.transferable).to.be.equal('0');
            expect(userBalanceObj.userId).to.be.equal(user);
            expect(userBalanceObj.applicationId).to.be.equal(app);
            expect(userBalanceObj.linkedWallet).to.be.equal(walletAddr);
        });
    });

    // TODO - add more test for error scenarios such as replaying the same transaction, last eth block smaller than current, bad signatures, etc.

    after(async () => {

    });

});
