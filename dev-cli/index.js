const pendingProps = require('./pending-props');
const cli = require('caporal');
const figlet = require('figlet');

cli
    .version('0.0.1')
    .command('issue', 'issue an earning to a recipient')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .argument('<amount>', 'amount of props to issue in earning')
    .argument('<description>', 'reason for this earning (optional)')
    .action(async (args, options, logger) => {
        logger.info(`issuing earning of amount ${args.amount} to application ${args.application} user ${args.user} for ${args.description}`);
        try {
            await pendingProps.issue(args.application, args.user, args.amount, args.description ? args.description : '');
        } catch (e) {
            logger.error(`error issuing earning: ${e}`)
        }
    })
    .command('revoke', 'revoke earning(s) by their state address')
    .argument('<earnings...>', 'earnings addresses', cli.LIST)
    .action(async (args, options, logger) => {
        logger.info(`revoking earning(s): ${args.earnings}`);
        try {
            await pendingProps.revoke(args.earnings);
        } catch (e) {
            logger.error(`error revoking earning(s): ${e}`)
        }
    })
    .command('settle', 'settle some earnings that have been paid on the ethereum chain')
    .argument('<ethtransactionhash>', 'hash of the ethereum transaction that contains settlement for side chain earning')
    .argument('<recipient>', 'recipient address')
    .action(async (args, options, logger) => {
        logger.info(`settling earning(s) with hash: ${args.ethtransactionhash}`);
        try {
            await pendingProps.settle(args.ethtransactionhash, args.recipient);
        } catch (e) {
            logger.error(`error settling earning(s): ${e}`)
        }
    })
    .command('updateLastEthBlockId', 'Update last eth block id for which events were added')
    .argument('<blockid>', 'last block id')
    .action(async (args, options, logger) => {
        logger.info(`updating last eth block id with: ${args.blockid}`);
        try {
            await pendingProps.updateLastBlockId(args.blockid);
        } catch (e) {
            logger.error(`error updating last block id: ${e}`)
        }
    })
    .command('linkWallet', 'Link wallet to an application user')
    .argument('<address>', 'ethereum address of the recipient')
    .argument('<application>', 'UDID of an authorized application')
    .argument('<user>', 'UDID of an authorized application user')
    .argument('<signature>', 'signature by linked wallet of application.user')
    .action(async (args, options, logger) => {
        logger.info(`linking wallet ${args.address} to application ${args.application} and user ${args.user}, 
        using signature ${args.signature}`);
        try {
            await pendingProps.linkWallet(args.address, args.application, args.user, args.signature);
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
    .command('state-query', 'get earning(s) or settlement(s) from the state')
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
    });

const banner = figlet.textSync('props-chain-cli', {
    font: 'Slant',
    horizontalLayout: 'fitted',
    verticalLayout: 'default'
});

cli.description(`\n\n${banner}`);
cli.parse(process.argv);