const pendingProps = require('./pending-props');
const cli = require('caporal');
const figlet = require('figlet');

// console.log(`blockUpdateAddress=`+pendingProps.CONFIG.earnings.namespaces.blockUpdateAddress());
// console.log(`activityLogAddress=`+pendingProps.CONFIG.earnings.namespaces.activityLogAddress('2','1','1'));
// console.log(`balanceUpdateAddress=`+pendingProps.CONFIG.earnings.namespaces.balanceUpdateAddress('xx','0x00'));
// console.log(`settlementAddress=`+pendingProps.CONFIG.earnings.namespaces.settlementAddress('tt'));
// console.log(`balanceAddress=`+pendingProps.CONFIG.earnings.namespaces.balanceAddress('1','1'));
// console.log(`walletLinkAddress=`+pendingProps.CONFIG.earnings.namespaces.walletLinkAddress('1'));
// console.log(`transactionAddress=`+pendingProps.CONFIG.earnings.namespaces.transactionAddress('1','1','1','1'));
// process.exit(1);
// const sigTest = async () => {
//     const sig = await pendingProps.signMessage(`app1user1`, '0x2d4dcf292bc5bd8d7246099052dfc76b3cdd3524', '759b603832da1100ab47c0f4aa6d445637eb5873d25cadd40484c48970b814d7'); // "signature11";
//     console.log("sig="+sig);
//     process.exit(0);
// }
// sigTest();

cli
    .version('0.0.1')
    .command('issue', 'issue props to a recipient')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .argument('<amount>', 'amount of props to issue in earning')
    .argument('[description]', 'reason for this earning (optional)')
    .argument('[timestamp]', 'timestamp of the transaction (will default to now if not passed')
    .action(async (args, options, logger) => {
        logger.info(`issuing props of amount ${args.amount} to application ${args.application} user ${args.user} for ${args.description} (${args.timestamp})`);
        try {
            await pendingProps.transaction(pendingProps.transactionTypes.ISSUE, args.application, args.user, args.amount, args.description ? args.description : '', {}, args.timestamp ? args.timestamp : 0);
        } catch (e) {
            logger.error(`error issuing earning: ${e}`)
        }
    })
    .command('revoke', 'revoke props from recipient')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .argument('<amount>', 'amount of props to issue in earning')
    .argument('[description]', 'reason for this earning (optional)')
    .argument('[timestamp]', 'timestamp of the transaction (will default to now if not passed')
    .action(async (args, options, logger) => {
        logger.info(`revoking props of amount ${args.amount} to application ${args.application} user ${args.user} for ${args.description} (${args.timestamp})`);
        try {
            await pendingProps.transaction(pendingProps.transactionTypes.REVOKE, args.application, args.user, args.amount, args.description ? args.description : '', {}, args.timestamp ? args.timestamp : 0);
        } catch (e) {
            logger.error(`error issuing earning: ${e}`)
        }
    })
    .command('updateLastEthBlockId', 'Update last eth block id for which events were added')
    .argument('<blockid>', 'last block id')
    .argument('<timestamp>', 'last block timestamp')
    .action(async (args, options, logger) => {
        logger.info(`updating last eth block id with: ${args.blockid}, ${args.timestamp}`);
        try {
            await pendingProps.updateLastBlockId(args.blockid, args.timestamp);
        } catch (e) {
            logger.error(`error updating last block id: ${e}`)
        }
    })
    .command('linkWallet', 'Link wallet to an application user')
    .argument('<address>', 'ethereum address of the recipient')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .argument('<signature>', 'signature by linked wallet of application.user')
    .argument('[timestamp]', 'timestamp of the transaction (will default to now if not passed')
    .action(async (args, options, logger) => {
        logger.info(`linking wallet ${args.address} to application ${args.application} and user ${args.user}, 
        using signature ${args.signature} (${args.timestamp})`);
        try {
            await pendingProps.linkWallet(args.address, args.application, args.user, args.signature, args.timestamp ? args.timestamp : 0);
        } catch (e) {
            logger.error(`error linking wallet: ${e}`)
        }
    })
    .command('externalBalanceUpdate', 'Update balance of address from main chain')
    .argument('<address>', 'ethereum address of the recipient')
    .argument('<balance>', 'new balance on main chain')
    .argument('<ethtransactionhash>', 'new balance on main chain')
    .argument('<blockid>', 'block id of the transaction')
    .argument('<timestamp>', 'eth transaction timestamp')
    .action(async (args, options, logger) => {
        logger.info(`sending updated main chain balance ${args.address} for ${args.balance}, 
        with txHash ${args.ethtransactionhash} (blockId: ${args.blockid}, timestamp: ${args.timestamp})`);
        try {
            await pendingProps.externalBalanceUpdate(args.address, args.balance, args.ethtransactionhash, args.blockid, args.timestamp);
        } catch (e) {
            logger.error(`error updating external balance: ${e}`)
        }
    })
    .command('settle', 'settle props to a recipient')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .argument('<amount>', 'amount of props that was settled')
    .argument('<toaddress>', 'users wallet who received settlement')
    .argument('<fromaddress>', 'application rewards address')
    .argument('<ethtransactionhash>', 'settlement transaction hash')
    .argument('<blockid>', 'block id of the settlement transaction hash')
    .argument('<timestamp>', 'timestamp of the settlement transaction hash')
    .argument('<tobalance>', 'timestamp of the settlement transaction hash')
    .action(async (args, options, logger) => {
        logger.info(`settling props of amount ${args.amount} to application ${args.application} user ${args.user} for ${args.description}, 
        with txHash ${args.ethtransactionhash} (blockId: ${args.blockid}, timestamp: ${args.timestamp}, to: ${args.toaddress}, from: ${args.fromaddress}, toBalance: ${args.tobalance})`);
        try {
            await pendingProps.settle(args.application, args.user, args.amount, args.toaddress, args.fromaddress, args.ethtransactionhash, args.blockid, args.timestamp, args.tobalance);
        } catch (e) {
            logger.error(`error settling: ${e}`)
        }
    })
    .command('state-query', 'get transaction(s) or balance(s) etc. from the state')
    .argument('<stateaddress>', 'state address for query')
    .argument('<t>', 'state type')
    .action(async (args, options, logger) => {
        try {
            await pendingProps.queryState(args.stateaddress, args.t);
        } catch (e) {
            logger.error(`error making state query: ${e}`)
        }
    })
    .command('logActivity', 'Log activity to the state')
    .argument('<userId>', 'User id')
    .argument('<appId>', 'Application id')
    .argument('<timestamp>', 'The timestamp')
    .argument('<date>', 'The date in YYYYMMDD format')
    .action(async (args, options, logger) => {
       try {
           await pendingProps.logActivity(args.userId, args.appId, args.timestamp, args.date);
       } catch (e) {
           logger.error(`error logging activity: ${e}`);
       }
    })
    .command('activityAddress', 'Get activity address')
    .argument('<userId>', 'User id')
    .argument('<appId>', 'Application id')
    .argument('<date>', 'The timestamp')
    .action(async (args, options, logger) => {
        try {
            logger.info(pendingProps.CONFIG.earnings.namespaces.activityLogAddress(args.date, args.appId, args.userId));
        } catch (e) {
            logger.error(`error logging activity: ${e}`);
        }
    })
    .command('balance-address', 'Get balance address')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .action((args, options, logger) => {
        const application = args.application == 0 ? '': args.application;
        logger.info(pendingProps.CONFIG.earnings.namespaces.balanceAddress(application, args.user));
    })
    .command('walletlink-address', 'Get wallet link address')
    .argument('<wallet>', 'wallet address')
    .action((args, options, logger) => {
        logger.info(pendingProps.CONFIG.earnings.namespaces.walletLinkAddress(args.wallet));
    })
    .command('balance-update-tx-address', 'Get balance update tx address')
    .argument('<ethtransactionhash>', 'transaction hash')
    .argument('<address>', 'ethereum address for which an update took place')
    .action((args, options, logger) => {
        logger.info(pendingProps.CONFIG.earnings.namespaces.balanceAddress(args.ethtransactionhash, args.address));
    })
    .command('settlement-address', 'Get balance update tx address')
    .argument('<ethtransactionhash>', 'transaction hash')
    .action((args, options, logger) => {
        logger.info(pendingProps.CONFIG.earnings.namespaces.settlementAddress(args.ethtransactionhash));
    });

const banner = figlet.textSync('props-chain-cli', {
    font: 'Slant',
    horizontalLayout: 'fitted',
    verticalLayout: 'default'
});

cli.description(`\n\n${banner}`);
cli.parse(process.argv);