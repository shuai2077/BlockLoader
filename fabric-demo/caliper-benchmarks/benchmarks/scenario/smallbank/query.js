/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

const OperationBase = require('./utils/operation-base');
const Smallbank = require('./utils/smallbank');

/**
 * Workload module for Smallbank queries.
 */
class Query extends OperationBase {
    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
    }

    /**
     * Creates a Smallbank instance initialized for the configured number of accounts.
     * @return {Smallbank} The instance.
     * @protected
     */
    createSmallbank() {
        return new Smallbank(this.accounts);
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        const queryArgs = this.smallbank.getQueryArguments();
        const request = this.createConnectorRequest('query', queryArgs);
        await this.sutAdapter.sendRequests(request);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new Query();
}

module.exports.createWorkloadModule = createWorkloadModule;
