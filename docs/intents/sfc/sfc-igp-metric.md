# SFC Intent with IGP Metric

Service Function Chain (SFC) intents are used to establish a path where packets are routed through a sequence of service nodes before reaching their destination. When no additional intents are specified, the optimal service function chain is determined using the IGP metric. The order of services in the path request dictates the sequence in which they are traversed. For instance, if a firewall is specified before an IDS, the path will first route through the firewall, followed by the IDS.

## Network Topology

The network topology below illustrates the test environment. Two segment-routing-aware firewalls (SERA-1 and SERA-2) are connected to XR-2 and XR-3, respectively. These firewalls manage traffic by blocking or allowing specific types. Additionally, two SNORT IDS instances (SNORT-1, SNORT-2) are connected to XR-6 and XR-7, respectively, and are used to detect and prevent network attacks.


![Hawkv6 Network with Services](../../images/hawkv6-network-sfc-overview.drawio.svg)

## Example Scenarios

The following scenarios demonstrate how SFC intents can optimize the path between two hosts:

1. **Scenario 1**: Path between HOST-A and HOST-C using the IGP metric, routed through a firewall.
2. **Scenario 2**: Path between HOST-A and HOST-C using the IGP metric, routed through both a firewall and an IDS.

### Scenario 1: SFC with Firewall and IGP Metric

In this scenario, HOST-A (acting as the client) requests a path to HOST-B (acting as the server) that must pass through the most optimally placed firewall. The path is determined by evaluating the IGP metric and the status of healthy firewall instances. In this scenario, the firewall instances considered are SERA-1 and SERA-2.



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
```

#### API Request
The `sfc with only firwall IGP metric` request is saved `sfc` folder. The JSON request format is as follows:
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
- `fc00:0:8:0:1::`

The packets are routed through the following sequence of devices:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> SERA-1 -> XR-6 -> XR-8 -> SITE-B -> HOST-B

This path directs the traffic from HOST-A to HOST-B through the best-placed firewall, considering the IGP metric. In this scenario, SERA-1 is identified as the most optimal firewall instance for the path.


![SFC with only FW](../../images/hawkv6-sfc-fw-intent.drawio.svg)



#### Scenario 2: SFC with FW, IDS, and IGP Metric

This scenario builds upon the previous one by adding an IDS instance to the path. The optimal path is calculated by considering the IGP metric alongside the healthy firewall and IDS instances. In this scenario, SERA-1 and SERA-2 serve as the firewall instances, while SNORT-1 and SNORT-2 are the IDS instances.


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
```

#### API Request
The `sfc with firwall and IDS IGP metric` request is saved `sfc` folder. The JSON request format is as follows:
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
- `fc00:0:b:0:1::`

The packets are routed through the devices in this sequence:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> SERA-1 -> XR-6 -> SNORT-1 -> SITE-B -> HOST-B

This path leads from HOST-A to HOST-B through the best placed firewall and IDS instances, optimized according to the IGP metric. In this scenario, SERA-1 is identified as the best placed firewall instance, and SNORT-1 as the best placed IDS instance.


![SFC with FW and IDS](../../images/hawkv6-sfc-fw-ids-intent.drawio.svg)

