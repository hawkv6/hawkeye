# Low Packet Loss
The low packet loss Intent identifies the best path between a given source and destination pair based on the lowest packet loss. This intent is ideal for applications requiring a reliable connection between two points.

## Network Topology
The network topology below, with the associated packet loss values, is used to demonstrate this intent. The packet loss impairments were applied using the [Lab Impairment Script](https://github.com/hawkv6/network/blob/main/docs/network.md#lab-impairments-scripts)

![Hawkv6 Network with Packet Loss Impairments](../images/hawkv6-network-packet-loss.drawio.svg)

## Example Scenario

In this example, Host-A (acting as a client) requests a low-packet-loss path to Host-B (acting as a server). The HawkEye controller calculates the optimal path based on the lowest packet loss between the two hosts.

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
          - intent: low-packet-loss
```

### API Request
The `low packet loss` request is saved `single intent` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:b::10",
    "intents": [
        {
            "type": "INTENT_TYPE_LOW_PACKET_LOSS"
        }
    ]
}
```

### Result 
The result includes the following SID List, ensuring the lowest packet loss path between the two hosts:
- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:6:0:1::`
- `fc00:0:8:0:1::`
- `fc00:0:b:0:1::`

The packets are routed through the following devices:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-6 -> XR-8 -> SITE-B -> HOST-B

![Low Packet Loss Path](../images/hawkv6-low-packet-loss-intent.drawio.svg)

