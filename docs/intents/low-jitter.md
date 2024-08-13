# Low Jitter Intent
The low jitter intent identifies the best path between a given source and destination pair based on the lowest jitter. This intent is ideal for applications requiring a stable connection between two points.

## Network Topology
The system is designed that the higher the delay the higher the jitter. Thus, the network topology below, with the associated latency values, is used to demonstrate this intent. 
The latency impairments were applied using the [Lab Impairment Script](https://github.com/hawkv6/network/blob/main/docs/network.md#lab-impairments-scripts).

However, for manual testing an additional jitter can also be set via the `clab-telemetry-linker -n <node> -i <interface> -j <jitter value>` command.

## Example Scenario

In this example, Host-A (acting as a client) requests a low-latency path to Host-B (acting as a server). The HawkEye controller calculates the optimal path based on the lowest latency between the two hosts.

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
          - intent: low-jitter
```

### API Request
The `low jitter` request is saved `single intent` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:b::10",
    "intents": [
        {
            "type": "INTENT_TYPE_LOW_JITTER"
        }
    ]
}
```

### Result 
The result includes the following SID List, ensuring the lowest jitter path between the two hosts:
- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:4:0:1::`
- `fc00:0:5:0:1::`
- `fc00:0:6:0:1::`
- `fc00:0:b:0:1::` 

The packets are routed through the following devices:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-4 -> XR-5 -> XR-6 -> SITE-B -> HOST-B

![Low Latency Path](../images/hawkv6-low-latency-intent.drawio.svg)

