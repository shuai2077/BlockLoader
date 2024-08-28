'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');


const operation_type = ['transact_savings','deposit_checking','send_payment','write_check', 'amalgamate'];
let prefix;
//let account_array = [];
//let acc = require('./load.js');
//account_array = acc.account_array;


/**
 * Workload module for the benchmark round.
 */
class CreateCarWorkload extends WorkloadModuleBase {
	/**
	 * Initializes the workload module instance.
	 */
	constructor() {
		super();
		this.txnPerBatch;
		this.queryRatio;
		this.account_array=[];
        
    }

    getAccount() {
        return Math.floor(Math.random()*Math.floor(this.account_array.length));
    }

    /**
     * Get two accounts
     * @return {Array} index of two accounts
     */
    get2Accounts() {
        let idx1 = this.getAccount();
        let idx2 = this.getAccount();
        if(idx2 === idx1) {
            idx2 = this.getAccount();
        }
        return [idx1, idx2];
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        // if (!this.roundArguments.hasOwnProperty('txnPerBatch')) {
        //     return Promise.reject(new Error('ycsb.txn - \'txnPerBatch\' is missed in the arguments'));
        // }
        this.txnPerBatch = this.roundArguments.txnPerBatch;

        // if (!this.roundArguments.hasOwnProperty('readRatio') || !inRange(this.roundArguments.readRatio)) {
        //     return Promise.reject(new Error('ycsb.txn - \'readRatio\' is missed in the arguments or not in the correct range.'));
        // }
	let acc = require('./load.js');
        this.account_array = acc.account_array;
	console.log(this.account_array);
        this.queryRatio = this.roundArguments.queryRatio;
    }
    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        let prob = Math.random();
        if (Math.random() < this.queryRatio) {
            let query_promises = [];
            let acc_num  = this.account_array[Math.floor(Math.random()*(this.account_array.length))];
            for (let i=0; i < this.txnPerBatch; ++i) {
                let args = {
                    contractId: 'smallbank',
                    contractVersion: 'v1',
                    contractFunction: 'query',
                    contractArguments: [acc_num],
                    timeout: 30,
                    readOnly: true
                };
                query_promises.push(args);
            }
            await this.sutAdapter.sendRequests(query_promises);
        } else {
            let workload = [];
            for(let j= 0; j<this.txnPerBatch; j++) {
                let op_index =  Math.floor(Math.random() * Math.floor(operation_type.length));
                let acc_index = this.getAccount();
                let random_op = operation_type[op_index];
                let random_acc = this.account_array[acc_index];
                let amount = Math.floor(Math.random() * 200);
                let op_payload;
                switch(random_op) {
                    case 'transact_savings': {
                        op_payload = {
                            contractId: 'smallbank',
                            contractVersion: 'v1',
                            contractFunction: random_op,
                            contractArguments: [amount, random_acc],
                            timeout: 30
                        };
                        break;
                    }
                    case 'deposit_checking': {
                        op_payload = {
                            contractId: 'smallbank',
                            contractVersion: 'v1',
                            contractFunction: random_op,
                            contractArguments: [amount, random_acc],
                            timeout: 30
                        };
                        break;
                    }
                    case 'send_payment': {
                        let accounts = this.get2Accounts();
                        op_payload = {
                            contractId: 'smallbank',
                            contractVersion: 'v1',
                            contractFunction: random_op,
                            contractArguments: [amount,this.account_array[accounts[0]],this.account_array[accounts[1]]],
                            timeout: 30
                        };
                        break;
                    }
                    case 'write_check': {
                        op_payload = {
                            contractId: 'smallbank',
                            contractVersion: 'v1',
                            contractFunction: random_op,
                            contractArguments: [amount, random_acc],
                            timeout: 30
                        };
                        break;
                    }
                    case 'amalgamate': {
                        let accounts = this.get2Accounts();
                        op_payload = {
                            contractId: 'smallbank',
                            contractVersion: 'v1',
                            contractFunction: random_op,
                            contractArguments: [this.account_array[accounts[0]], this.account_array[accounts[1]]],
                            timeout: 30
                        };
                        break;
                    }
                    default: {
                        throw new Error('Invalid operation!!!');
                    }
                }
                workload.push(op_payload);
            }
            await this.sutAdapter.sendRequests(workload);
        }
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
