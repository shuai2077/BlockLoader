# BlockLoader
Blockchain distributed testing

---

### 1. Network Structure

We are setting up a multi-machine, multi-node network with the following structure. The network consists of four organizations, `org1`,`org2`,`org3` and `org4`, each with one `peer` nodes, and one `orderer` nodes, and there is also a separate node for Caliper to test the performance of the network.

| **Name**                  | **IP Address** | **Hosts**                 | **Organization** |
|---------------------------|----------------|---------------------------|------------------|
| peer0.org1.example.com    | 172.19.127.240 | peer0.org1.example.com    | org1             |
| peer0.org2.example.com    | 172.19.127.234 | peer0.org2.example.com    | org2             |
| peer0.org3.example.com    | 172.19.127.233 | peer0.org3.example.com    | org3             |
| peer0.org4.example.com    | 172.19.127.236 | peer0.org4.example.com    | org4             |
| orderer1.org1.example.com | 172.19.127.238 | orderer1.org1.example.com | org1             |
| orderer2.org2.example.com | 172.19.127.239 | orderer2.org2.example.com | org2             |
| orderer3.org3.example.com | 172.19.127.232 | orderer3.org3.example.com | org3             |
| orderer4.org4.example.com | 172.19.127.235 | orderer4.org4.example.com | org4             |

### 2. Setting Up Network Hosts

First, open the terminal by right-clicking, then use the following command to check the IP address of the current virtual machine on each of the eight virtual machines. The last line will display the IP address of the machine.

```bash
ifconfig
```

Next, configure the network hosts on all servers by performing the following steps on each of the three virtual machines.

```bash
gedit /etc/hosts
```

Add the following lines at the end of the file (the IP addresses and hosts can be chosen freely, but they must remain consistent):

```plaintext
172.19.127.240  peer0.org1.example.com
172.19.127.234  peer0.org2.example.com
172.19.127.233  peer0.org3.example.com
172.19.127.236  peer0.org4.example.com
172.19.127.238  orderer1.example.com
172.19.127.239  orderer2.example.com
172.19.127.232  orderer3.example.com
172.19.127.235  orderer4.example.com
```

### 3. Caliper Config File

We write caliper configuration files based on different test benchmarks, add monitor to the configuration file to monitor the operation of nodes on different servers, we here take the linear-rate test benchmark as an example:

```yaml
test:
  workers:
    number: 4
  rounds:
    - label: load.
      txNumber: 100000
      rateControl:
        type: linear-rate
        opts:
          startingTps: 250
          finishingTps: 1000
      workload:
        module: benchmarks/samples/fabric/smallbank/load.js
        arguments:
          txnPerBatch: 1
    - label: mixed.
      txNumber: 100000
      rateControl:
        type: linear-rate
        opts:
          startingTps: 250
          finishingTps: 1000
      workload:
        module: benchmarks/samples/fabric/smallbank/mixed.js
        arguments:
          queryRatio: 0.1
          txnPerBatch: 1
monitors:
  resource:
    - module: docker
      options:
        interval: 1
        containers:
          - peer0.org1.example.com
          - http://172.19.127.240:2375/peer0.org1.example.com
          - peer0.org2.example.com
          - http://172.19.127.234:2375/peer0.org2.example.com
          - peer0.org3.example.com
          - http://172.19.127.233:2375/peer0.org3.example.com
          - peer0.org4.example.com
          - http://172.19.127.236:2375/peer0.org4.example.com
          - orderer1.example.com
          - http://172.19.127.238:2375/orderer1.example.com
          - orderer2.example.com
          - http://172.19.127.239:2375/orderer2.example.com
          - orderer3.example.com
          - http://172.19.127.232:2375/orderer3.example.com
          - orderer4.example.com
          - http://172.19.127.235:2375/orderer4.example.com
      charting:
        bar:
          metrics: [all]
```

In order for caliper to successfully detect the operation of nodes on different servers, we need to modify the docker configuration as follows:

```bash
cp /usr/lib/systemd/system/docker.service /usr/lib/systemd/system/docker.service.bk

vi /usr/lib/systemd/system/docker.service
#ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2375 -H fd:// --containerd=/run/containerd/containerd.sock

systemctl daemon-reload
systemctl restart docker
```

### 4. Test

After completing the above configuration, we can start the test.

First, run the following command on a different server to turn on the node docker service.

```bash
cd ~/testwork/fixtures/

docker-compose up -d
```

We then use the fabric sdk to set up the blockchain.

```bash
cd ~/testwork/

./testwork
```

After the blockchain is deployed, we run the following command to turn on caliper testing

```bash
cd ~/fabric-demo/caliper-benchmarks/

npx caliper launch manager \
--caliper-workspace ./ \
--caliper-networkconfig networks/fabric/test-network-sdk.yaml \
--caliper-benchconfig benchmarks/samples/fabric/fabric-sdk-go/config.yaml \
--caliper-flow-only-test \
--caliper-fabric-gateway-enabled
```

---