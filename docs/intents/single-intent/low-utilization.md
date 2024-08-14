# Low Utilization Intent
The `low utilization` intent identifies the best path between a given source and destination pair based on the lowest utilization. This intent is ideal for applications requiring the least congested path. Also it can be used to balance the traffic in the network.

## Network Topology
The network topology below, shows a high utilization between shortest path (based on the IGP metric) between Host-A and Host-B. 
Such a situation can be established by using the following command on HOST-A:
```
nping -6 2001:db8:b::10 -e eth1 -udp -data-length 1000 -delay 1ms -count 100000
```

![Hawkv6 Network with Utilization](../../images/hawkv6-network-utilization.drawio.svg)

## Example Scenario
In this example scenario, Host-A (acting as a client) requests a low-utilization path to Host-B (acting as a server). The HawkEye controller calculates the optimal path based on the lowest utilization between the two hosts.


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
          - intent: low-utilization
```

### API Request
The `low utilization` request is saved `single intent` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:b::10",
    "intents": [
        {
            "type": "INTENT_TYPE_LOW_UTILIZATION"
        }
    ]
}
```

### Result 
The result includes the following SID List, ensuring the lowest utilization path between the two hosts:
- `fc00:0:1:0:1::`
- `fc00:0:3:0:1::`
- `fc00:0:7:0:1::`
- `fc00:0:8:0:1::`
- `fc00:0:b:0:1::`


The packets are routed through the following devices:
HOST-A -> SITE-A -> XR-1 -> XR-3 -> XR-7 -> XR-8 -> SITE-B -> HOST-B

![Low Utilization Path](../../images/hawkv6-low-utilization-intent.drawio.svg)

