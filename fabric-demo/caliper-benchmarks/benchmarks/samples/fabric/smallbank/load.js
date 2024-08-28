'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');



const initial_balance = 1000000;
let account_array = [];
let prefix;

/**
 * Workload module for the benchmark round.
 */
class CreateCarWorkload extends WorkloadModuleBase {
    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        // this.txIndex = 0;
        this.txnPerBatch = 1;

    }
      /**
     * Initialize the workload module with the given parameters.
     * @param {number} workerIndex The 0-based index of the worker instantiating the workload module.
     * @param {number} totalWorkers The total number of workers participating in the round.
     * @param {number} roundIndex The 0-based index of the currently executing round.
     * @param {Object} roundArguments The user-provided arguments for the round from the benchmark configuration file.
     * @param {BlockchainInterface} sutAdapter The adapter of the underlying SUT.
     * @param {Object} sutContext The custom context object provided by the SUT adapter.
     * @async
     */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.txnPerBatch = 1;
    }
    generateAccount() {
        // should be [a-z]{1,9}
        if(typeof prefix === 'undefined') {
            prefix = process.pid;
        }
        let count = account_array.length+1;
        let num = prefix.toString() + count.toString();
        return num;
    }


    random_string() {
        let text = '';
        const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
        for (let i = 0; i < 12; i++) {
            text += possible.charAt(Math.floor(Math.random() * possible.length));
        }
        return text;
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        let workload = [];
        for(let i= 0; i < this.txnPerBatch; i++) {
            let acc_id = this.generateAccount();
            account_array.push(acc_id);
            let args = {
                contractId: 'smallbank',
                contractVersion: 'v1',
                contractFunction: 'create_account',
                contractArguments: [acc_id, this.random_string(),initial_balance,initial_balance],
                timeout: 30
            };
            workload.push(args);
        }
        await this.sutAdapter.sendRequests(workload);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new CreateCarWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
module.exports.account_array = account_array;
