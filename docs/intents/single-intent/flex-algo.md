# Flex Algo Intent
The `flex algo` intent enables users to calculate a path on a specific subgraph of the network topology, allowing the exclusion of certain links or nodes. This intent is ideal for applications that require additional constraints, such as ensuring encryption on the path or including geo-based information in the calculation.

## Network Topology

The network topology below, shows to pre-defined flex algos of the test network. The flex algo 128 and 129 are used to demonstrate this intent.


## Flex Algos Comparison

| Flex Algo 128 | Flex Algo 129 |
|---------------|---------------|
| ![Flex Algo 128](../../images/hawkv6-network-flex-algo-128.drawio.svg) | ![Flex Algo 129](../../images/hawkv6-network-flex-algo-129.drawio.svg)

## Example Scenario
In this example, Host-A (acting as a client) requests a path to Host-B (acting as a server) using flex algo 128. The HawkEye controller calculates the optimal path based on the IGP cost between the two hosts.

### HawkWing Configuration
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

### API Request
The `flex-algo128` (and `flex-algo-129`) request is saved `single intent` folder. The JSON request format is as follows:
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

### Result 
The result includes the following SID List, ensuring the lowest latency path between the two hosts:
- `fc00:0:1:128:1::`
- `fc00:0:3:128:1::`
- `fc00:0:7:128:1::`
- `fc00:0:8:128:1::`
- `fc00:0:b:128:1::`


The packets are routed through the following devices:
HOST-A -> SITE-A -> XR-1 -> XR-3 -> XR-7 -> XR-8 -> SITE-B -> HOST-B

![Flex Algo 128 Path](../../images/hawkv6-flex-algo-128-intent.drawio.svg)

