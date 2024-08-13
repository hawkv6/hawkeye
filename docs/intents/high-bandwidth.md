# High Bandwidth Intent
The high bandwidth intent identifies the best path between a given source and destination pair based on the highest available bandwidth. The available bandwidth equals the maximum link bandwidth - the consumed bandwidth. This intent is ideal for applications requiring high bandwidth.

## Network Topology
The network topology below, with the associated rate limits, is used to demonstrate this intent. The bandwidth/rate impairments were applied using the [Lab Impairment Script](https://github.com/hawkv6/network/blob/main/docs/network.md#lab-impairments-scripts). Since, the there is almost no traffic in the network, the path should follow the path including the highest bandwidth links.

![Hawkv6 Network with Latency Impairments](../images/hawkv6-network-bw.drawio.svg)

## Example Scenario
In this example, Host-A (acting as a client) requests a high bandwidth path to Host-B (acting as a server). The HawkEye controller calculates the optimal path based on the highest available bandwidth between the two hosts.

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
          - intent: high-bandwidth
```

### API Request
The `high bandwidth` request is saved `single intent` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:c::10",
    "intents": [
        {
            "type": "INTENT_TYPE_HIGH_BANDWIDTH"
        }
    ]
}
```

### Result 
The result includes the following SID List, ensuring the lowest latency path between the two hosts:
- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:6:0:1::`
- `fc00:0:8:0:1::`
- `fc00:0:b:0:1::`


The packets are routed through the following devices:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-6 -> SITE-B -> HOST-B

![High Bandwidth Path](../images/hawkv6-high-bw-intent.drawio.svg)

