# Constrained Intents
In addition to combining intents, path selection can be refined using constraints. These constraints allow you to specify the maximum tolerable packet loss, latency, and jitter, as well as the minimum required bandwidth.

This section demonstrates how to apply these constraints to optimize the path between two hosts. As an example, the combination of low latency and low packet loss intents is shown with constraints on packet loss, latency, and bandwidth.

## Network Topology

The network topology below, with the associated latency, packet loss values and bandwidth impairments, is used to demonstrate this intent. The impairments were applied using the [Lab Impairment Script](https://github.com/hawkv6/network/blob/main/docs/network.md#lab-impairments-scripts).


![Hawkv6 Network with Latency and Packet Loss and BW](../../images/hawkv6-network-packet-loss-delay-bw.drawio.svg)

## Example Scenarios

The following scenarios illustrate how combined intents with constraints can be utilized to optimize the path between two hosts:

1. **Scenario 1**: Prioritizing low latency, with packet loss as a secondary factor, and applying maximum constraints for the path between Host-A and Host-C.
2. **Scenario 2**: Prioritizing low latency, with packet loss as a secondary factor, while incorporating both maximum and minimum constraints for the path between Host-A and Host-C.


### Scenario 1: Low Latency and Low Packet Loss with max Constraints

In this scenario, HOST-A (acting as client) requests a path to HOST-C ( acting as server) that prioritizes low latency and low packet loss. Additionally, the request specifies constraints: latency must not exceed 25ms, and packet loss must remain below 1%. The HawkEye controller calculates the optimal path, ensuring it meets these constraints while achieving the lowest possible latency and packet loss between the two hosts.


#### HawkWing Configuration
```yaml
---
client_ipv6_address: 2001:db8:a::10
hawkeye:
  enabled: true
  address: 2001:db8:e5::e
  port: 10000
services:
  webserver-c:
    ipv6_addresses:
      - 2001:db8:c::10
    applications:
      - port: 80
        intents:
          - intent: low-latency
            max_value: 25000 #us
          - intent: low-packet-loss
            max_value: 1 #%
```

#### API Request
The `low latency and low packet loss with max constraints` request is saved `combined intents` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:c::10",
    "intents": [
        {
            "type": "INTENT_TYPE_LOW_LATENCY",
            "values": [
                {
                    "type": "VALUE_TYPE_MAX_VALUE",
                    "number_value": 25000
                }
            ] 
        },
        {
            "type": "INTENT_TYPE_LOW_PACKET_LOSS",
            "values": [
                {
                    "type": "VALUE_TYPE_MAX_VALUE",
                    "number_value": 1
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
- `fc00:0:3:0:1::`
- `fc00:0:7:0:1::`
- `fc00:0:c:0:1::`

The packets are routed through the following sequence of devices:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-3 -> XR-7 -> SITE-C -> HOST-C

Compared to the [example without constraints](two-intents.md#scenario-1--low-latency-and-low-packet-loss), this path successfully avoids the link between XR-5 and XR-7, which exceeds the 1% packet loss constraint. As a result, the path calculation selects the best alternative route that meets the specified requirements:


![Low Latency and Low Packet Loss Path with constraints](../../images/hawkv6-low-latency-low-loss-with-constraints-intent.drawio.svg)


#### Scenario 2: Low Packet Loss and Low Latency with max and min constraints

In this scenario, the path between HOST-A and HOST-C is optimized with a focus on minimizing packet loss and latency, while also enforcing additional constraints. The minimum required bandwidth is set to 200 Mbit/s, the maximum allowable packet loss is 3%, and the maximum latency is capped at 25ms.

The HawkEye controller computes the optimal path by considering all these constraints.


#### HawkWing Configuration
```yaml
---
client_ipv6_address: 2001:db8:a::10
hawkeye:
  enabled: true
  address: 2001:db8:e5::e
  port: 10000
services:
  webserver-c:
    ipv6_addresses:
      - 2001:db8:c::10
    applications:
      - port: 80
        intents:
          - intent: low-latency
            max_value: 25000 #us
          - intent: low-packet-loss
            max_value: 3 #%
          - intent: high-bandwidth
            min_value: 200000 #200 Mbit/s
```

#### API Request
The `low packet loss and low latency` request is saved `combined intents` folder. The JSON request format is as follows:
```
{
    "ipv6_destination_address": "2001:db8:c::10",
    "ipv6_source_address": "2001:db8:a::10",
    "intents": [
                {
            "type": "INTENT_TYPE_LOW_LATENCY",
            "values": [
                {
                    "type": "VALUE_TYPE_MAX_VALUE",
                    "number_value": 25000
                }
            ]
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
- `fc00:0:6:0:1::`
- `fc00:0:8:0:1::`
- `fc00:0:c:0:1::`

The packets are routed through the devices in this sequence:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-6 -> XR-8 -> SITE-C -> HOST-C

This path is optimized to prioritize minimizing latency while considering packet loss, along with the additional constraints of a minimum required bandwidth of 200 Mbit/s, a maximum latency of 25ms, and a maximum packet loss of 3%. The result indicates that a completely new path was selected, as many links failed to meet the bandwidth requirement.


![Low Low Packet Loss and Low Latency Path](../../images/hawkv6-low-latency-low-loss-with-min-constraints-intent.drawio.svg)

