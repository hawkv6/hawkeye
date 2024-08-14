# Flex Algo Intent

The `flex algo` intent allows users to calculate a path on a specific subgraph of the network topology, enabling the exclusion of certain links or nodes. This intent is particularly useful for applications requiring additional constraints, such as enforcing encryption along the path or incorporating geo-based information into the calculation.

When only the `flex algo` intent is provided, the HawkEye controller determines the optimal path based on the IGP cost between the two hosts. However, users can also define additional intents to be included in the calculation, allowing for the use of other metrics such as packet loss, latency, and jitter. The concept of constraints is also supported, providing further flexibility in path selection.

## Network Topology

The network topology below illustrates the pre-defined Flex Algo 128, highlighted in blue. All path requests specifying Flex Algo 128 will be calculated on this blue subgraph.


![Flex Algo 128](../../images/hawkv6-network-flex-algo-128-with-details.drawio.svg) 

## Example Scenarios

The following scenarios illustrate how Flex Algo intents can optimize the path between two hosts:

1. **Scenario 1**: Path between HOST-A and HOST-B using Flex Algo 128 and the IGP metric.
2. **Scenario 2**: Path between HOST-A and HOST-C using Flex Algo 128 with a combination of low latency and low packet loss.
3. **Scenario 3**: Path between HOST-A and HOST-C using Flex Algo 128 and a Service Function Chain (SFC) with low packet loss.

### Scenario 1: Flex Algo 128 with IGP Metric

In this scenario, HOST-A (acting as the client) requests a path to HOST-B (acting as the server) using Flex Algo 128. The HawkEye controller calculates the optimal path based on the IGP cost between the two hosts.


#### HawkWing Configuration
```yaml
---
client_ipv6_address: 2001:db8:a::10
hawkeye:
  enabled: true
  address: 2001:db8:e5::e
  port: 10000
services:
  webserver-b:
    ipv6_addresses:
      - 2001:db8:b::10
    applications:
      - port: 80
        intents:
          - intent: flex-algo
            flex_algo_number: 128
```

#### API Request
The `flex-algo128 igp metric` request is saved `flex-algo` folder. The JSON request format is as follows:
```
{
    "ipv6_destination_address": "2001:db8:b::10",
    "ipv6_source_address": "2001:db8:a::10",
    "intents": [
        {
            "type": "INTENT_TYPE_FLEX_ALGO",
            "values": [
        {
            "type": "VALUE_TYPE_FLEX_ALGO_NR",
            "number_value": 128
        }
            ]
        }
    ]
}
```
#### Result

The following SID list is generated to ensure the optimal path with the lowest IGP metric using Flex Algo 128 between the two hosts:

- `fc00:0:1:128:1::`
- `fc00:0:3:128:1::`
- `fc00:0:7:128:1::`
- `fc00:0:8:128:1::`
- `fc00:0:b:128:1::`

The packets are routed through the following devices:
HOST-A -> SITE-A -> XR-1 -> XR-3 -> XR-7 -> XR-8 -> SITE-B -> HOST-B


![Flex Algo 128 Path](../../images/hawkv6-flex-algo-128-intent.drawio.svg)


### Scenario 2: Flex Algo 128 with Low Latency and Low Packet Loss

In this scenario, Host-A (acting as the client) requests a path to Host-B (acting as the server) using Flex Algo 128. The HawkEye controller calculates the optimal path by prioritizing both low latency and low packet loss intents to ensure an efficient and reliable connection between the two hosts.

#### HawkWing Configuration
```yaml
---
client_ipv6_address: 2001:db8:a::10
hawkeye:
  enabled: true
  address: 2001:db8:e5::e
  port: 10000
services:
  webserver-b:
    ipv6_addresses:
      - 2001:db8:b::10
    applications:
      - port: 80
        intents:
          - intent: flex-algo
            flex_algo_number: 128
          - intent: low-latency
          - intent: low-packet-loss
```

#### API Request
The `flex-algo 128 with low latency and low packet loss` request is saved `flex-algo` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:b::10",
    "intents": [
        {
            "type": "INTENT_TYPE_FLEX_ALGO",
            "values": [
        {
            "type": "VALUE_TYPE_FLEX_ALGO_NR",
            "number_value": 128
        }
            ]
        },
        {
            "type": "INTENT_TYPE_LOW_LATENCY"
        },
        {
            "type": "INTENT_TYPE_LOW_PACKET_LOSS"
        },
    ]
}
``` 

#### Result

The following SID List is generated to ensure the optimal path with the lowest latency and packet loss using Flex Algo 128 between the two hosts:
- `fc00:0:1:128:1::`
- `fc00:0:3:128:1::`
- `fc00:0:7:128:1::`
- `fc00:0:8:128:1::`
- `fc00:0:b:128:1::`

The packets traverse the devices in this order:
HOST-A -> SITE-A -> XR-1 -> XR-3 -> XR-7 -> XR-8 -> SITE-B -> HOST-B

The path remains unchanged from the previous scenario because, within this subgraph, only one viable alternative exists. The alternative route through SITE-C would introduce higher latency and packet loss, resulting in a higher overall cost.

![Flex Algo 128 Low Latency Low Loss Path](../../images/hawkv6-flex-algo-128-low-latency-low-loss-intent.drawio.svg)


### 
### Scenario 3: Flex Algo 128 with SFC and Low Packet Loss
In this scenario, Host-A (acting as the client) requests a path to Host-C (acting as the server) using Flex Algo 128. The HawkEye controller calculates the optimal route based on the low packet loss intent, incorporating a Service Function Chain (SFC) that includes both a firewall and an IDS, while adhering to the Flex Algo 128 constraints. The available service nodes in this scenario are the SERA-2 firewall and the SNORT-2 IDS.


#### HawkWing Configuration
```yaml
---
client_ipv6_address: 2001:db8:a::10
hawkeye:
  enabled: true
  address: 2001:db8:e5::e
  port: 10000
services:
  webserver-b:
    ipv6_addresses:
      - 2001:db8:c::10
    applications:
      - port: 80
        intents:
          - intent: sfc
            functions:
              - fw
              - ids
          - intent: flex-algo
            flex_algo_number: 128
          - intent: low-packet-loss
```

#### API Request
The `flex-algo 128 sfc with firwall and ids low packet loss` request is saved `flex-algo` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:c::10",
    "intents": [
        {
            "type": "INTENT_TYPE_SFC",
            "values": [
                {
                    "type": "VALUE_TYPE_SFC",
                    "string_value": "fw"
                },
                {
                    "type": "VALUE_TYPE_SFC",
                    "string_value": "ids"
                }
            ]
        },
        {
            "type": "INTENT_TYPE_FLEX_ALGO",
            "values": [
                {
                    "type": "VALUE_TYPE_FLEX_ALGO_NR",
                    "number_value": 128
                }
            ]
        },
        {
            "type": "INTENT_TYPE_LOW_PACKET_LOSS"
        }
    ]
}
```

#### Result
The following SID List is generated to ensure the optimal path between the two hosts:

- `fc00:0:1:128:1::`
- `fc00:0:3:128:1::`
- `fc00:0:3f::`
- `fc00:0:7:128:1::`
- `fc00:0:7f::`
- `fc00:0:c:128:1::`

The packets traverse the devices in this order:
HOST-A -> SITE-A -> XR-1 -> XR-3 -> SERA-2 -> XR-7 -> SNORT-2 -> SITE-C -> HOST-C

Since the only available Service Function Chain (SFC) on the Flex Algo 128 subgraph involves the SERA-2 firewall and SNORT-2 IDS, the path is determined by this service sequence. 

![Flex Algo 128 Low Latency Low Loss SFC Path](../../images/hawkv6-flex-algo-128-sfc-intent.drawio.svg)