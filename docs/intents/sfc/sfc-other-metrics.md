# SFC Intent with Additional Metrics

SFC intents can be enhanced by incorporating metrics such as low latency, low jitter, low packet loss, or a combination thereof. Additionally, constraints like maximum packet loss, maximum latency, maximum jitter, and minimum bandwidth can be applied, similar to those used in [combined metrics](../combined-intents/constrained-intents.md).

## Network Topology

The network topology shown below illustrates the test environment. It includes two segment-routing-aware firewalls (SERA-1 and SERA-2) connected to XR-2 and XR-3, respectively, which manage traffic by blocking or allowing specific types. Additionally, two SNORT IDS instances (SNORT-1 and SNORT-2) are connected to XR-6 and XR-7, respectively, and are responsible for detecting and preventing network attacks. 

The topology also shows the packet loss, latency, and bandwidth impairments on the links between the routers. These impairments were applied using the [Lab Impairment Script](https://github.com/hawkv6/network/blob/main/docs/network.md#lab-impairments-scripts).


![Hawkv6 Network with Services and Impairments](../../images/hawkv6-network-sfc-with-impairments.drawio.svg)


## Example Scenarios

The following scenarios demonstrate how SFC intents can optimize the path between two hosts:

1. **Scenario 1**: Path between HOST-A and HOST-B using the packet loss metric, routed through a firewall and an IDS.
2. **Scenario 2**: Path between HOST-A and HOST-C using a combination of low latency and low packet loss, with additional constraints, routed through both a firewall and an IDS.

### Scenario 1: SFC with FW and IDS Optimized for Low Packet Loss

In this scenario, HOST-A (acting as the client) requests a path to HOST-B (acting as the server) that must traverse the optimal firewall and IDS instances based on the packet loss metric. The path is determined by evaluating the packet loss on each link and considering the status of healthy firewall and IDS instances. The firewall instances considered are SERA-1 and SERA-2, while the IDS instances are SNORT-1 and SNORT-2.


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
          - intent: sfc
            functions:
              - fw
              - ids
          - intent: low-packet-loss
```

#### API Request
The `sfc with firwall and IDS low packet loss ` request is saved `sfc` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:b::10", 
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
            "type": "INTENT_TYPE_LOW_PACKET_LOSS"
        }
    ]
}
```

#### Result
The following SID List is generated to ensure an optimized path between the two hosts:

- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:2f::`
- `fc00:0:6:0:1::`
- `fc00:0:6f::`
- `fc00:0:8:0:1::`
- `fc00:0:b:0:1::`

The packets are routed through the following sequence of devices:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> SERA-1 -> XR-6 -> SNORT-1 -> XR-8 -> SITE-B -> HOST-B

This path routes the traffic from HOST-A to HOST-B through the most optimal firewall and IDS instances, prioritized based on the packet loss metric. In this scenario, SERA-1 is identified as the optimal firewall, and SNORT-1 as the best IDS instance. The path avoids the high-packet-loss link between XR-6 and SITE-B by rerouting traffic through XR-8.


![SFC with only FW](../../images/hawkv6-sfc-fw-low-loss-intent.drawio.svg)


#### Scenario 2: SFC with FW and IDS, Low Latency, and Low Packet Loss with Max Constraints

This scenario demonstrates the capability to combine multiple metrics and constraints. HOST-A requests a path to HOST-C that must traverse the most optimal firewall and IDS instances, considering both low latency and low packet loss metrics. Additionally, there are constraints set to a maximum of 3% packet loss and a minimum required bandwidth of 200 Mbit/s, along with the requirement for healthy firewall and IDS instances.

In this scenario, the firewall instances considered are SERA-1 and SERA-2, while the IDS instances are SNORT-1 and SNORT-2.



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
          - intent: low-latency
          - intent: low-packet-loss
            max_value: 3  #%
          - intent: high-bandwidth
            min_value: 200000 #200 Mbit/s
```

#### API Request
The `sfc with firewall and IDS low latency and low packet loss with max and min constraints` request is saved `sfc` folder. The JSON request format is as follows:
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
            "type": "INTENT_TYPE_LOW_LATENCY"
        },
        {
            "type": "INTENT_TYPE_LOW_PACKET_LOSS",
             "values": [
                {
                    "type": "VALUE_TYPE_MAX_VALUE",
                    "number_value": 3
                }
            ]
        },
        {
            "type": "INTENT_TYPE_HIGH_BANDWIDTH",
             "values": [
                {
                    "type": "VALUE_TYPE_MIN_VALUE",
                    "number_value": 200000
                }
            ]
        }
    ]
}
```

#### Result

The following SID List is generated to ensure an optimized path between the two hosts:
- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:2f::`
- `fc00:0:6:0:1::`
- `fc00:0:6f::`
- `fc00:0:8:0:1::`
- `fc00:0:c:0:1::`

The packets are routed through the devices in this sequence:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> SERA-1 -> XR-6 -> SNORT-1 -> XR-8 -> SITE-C -> HOST-C

This path leads from HOST-A to HOST-C through the most optimal firewall and IDS instances, considering the IGP metric. In this scenario, SERA-1 is identified as the best firewall instance, and SNORT-1 as the best IDS instance. The path effectively avoids routing through SERA-2 and SNORT-2, as the paths through this service function chain do not meet the set constraints.

![SFC with FW and IDS](../../images/hawkv6-sfc-fw-low-latency-low-loss-min-intent.drawio.svg)

