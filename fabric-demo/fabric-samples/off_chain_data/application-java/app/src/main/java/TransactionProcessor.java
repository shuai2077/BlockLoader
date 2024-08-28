/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import com.google.protobuf.InvalidProtocolBufferException;
import parser.Transaction;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;

public final class TransactionProcessor {
    // Typically we should ignore read/write sets that apply to system chaincode namespaces.
    private static final Set<String> SYSTEM_CHAINCODE_NAMES = Set.of(
            "_lifecycle",
            "cscc",
            "escc",
            "lscc",
            "qscc",
            "vscc"
    );

    private final long blockNumber;
    private final Transaction transaction;
    private final Store store;

    public TransactionProcessor(final Transaction transaction, final long blockNumber, final Store store) {
        this.blockNumber = blockNumber;
        this.transaction = transaction;
        this.store = store;
    }

    private static boolean isSystemChaincode(final String chaincodeName) {
        return SYSTEM_CHAINCODE_NAMES.contains(chaincodeName);
    }

    public void process() throws IOException {
        var transactionId = transaction.getChannelHeader().getTxId();

        var writes = getWrites();
        if (writes.isEmpty()) {
            System.out.println("Skipping read-only or system transaction " + transactionId);
            return;
        }

        System.out.println("Process transaction " + transactionId);
        store.store(blockNumber, transactionId, writes);
    }

    private List<Write> getWrites() throws InvalidProtocolBufferException {
        var channelName = transaction.getChannelHeader().getChannelId();

        var writes = new ArrayList<Write>();
        for (var readWriteSet : transaction.getNamespaceReadWriteSets()) {
            var namespace = readWriteSet.getNamespace();
            if (isSystemChaincode(namespace)) {
                continue;
            }

            readWriteSet.getReadWriteSet().getWritesList().stream()
                    .map(write -> new Write(channelName, namespace, write))
                    .forEach(writes::add);
        }

        return writes;
    }
}
