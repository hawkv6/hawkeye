# Combining Three Intents 
It is possible to combine the `low latency`, `low jitter`, and `low packet loss` intents to find a specific path between two hosts which is optimized based on these intent combination. The order of these intents influences their weighting, with the first intent receiving a higher weight than the second and third. By default, the first intent has a weight of 0.7, the second has a weight of 0.2, and the third has a weight of 0.1. These weights can be adjusted by setting the `HAWKEYE_THREE_FACTOR_WEIGHTS` environment variable.

This page demonstrates the functionality using a combination of `low packet loss`, `low latency`, and `low jitter`. Other combinations are also possible.

## Network Topology

The network topology below, with the associated latency and packet loss values, is used to demonstrate this intent. The latency impairments were applied using the [Lab Impairment Script](https://github.com/hawkv6/network/blob/main/docs/network.md#lab-impairments-scripts).


![Hawkv6 Network with Latency and Packet Loss Impairments](../../images/hawkv6-network-packet-loss-delay.drawio.png)

## Example Scenarios

These two example scenarios demonstrate how combined intents can be used to optimize the path between two hosts:

1. **Scenario 1**: Prioritizing low packet loss over low latency and low jitter with the default weights for the path between Host-A and Host-C.
2. **Scenario 2**: Prioritizing low latency over low packet loss and low jitter with the adjusted weights for the path between Host-A and Host-C.

Note that the order in which the intents are listed influences their weighting, with the first intent receiving higher priority.


### Scenario 1:  Low Packet Loss, Low Latency, and Low Jitter with default weights
In this example, Host-A (acting as a client) requests a path to Host-C (acting as a server) that prioritizes low packet loss, low latency, and low jitter. In this scenario, packet loss is given higher importance than latency and jitter, with default parameters assigning a weight of 0.7 to packet loss, 0.2 to latency, and 0.1 to jitter. The HawkEye controller calculates the optimal path based on these weights to achieve the lowest packet loss, latency, and jitter between the two hosts.


#### HawkWing Configuration
```yaml
---
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
          - intent: low-packet-loss
          - intent: low-latency
          - intent: low-jitter
```

#### API Request
The `low packet loss and low latency and low jitter` request is saved `combined intents` folder. The JSON request format is as follows:
```
{
    "ipv6_source_address": "2001:db8:a::10",
    "ipv6_destination_address": "2001:db8:c::10",
    "intents": [
        {
            "type": "INTENT_TYPE_LOW_PACKET_LOSS"
        },
        {
            "type": "INTENT_TYPE_LOW_LATENCY"
        }, 
        {
            "type": "INTENT_TYPE_LOW_JITTER"
        }
    ]
}
```

#### Result

The following SID List is generated to optimize the path between the two hosts:

- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:3:0:1::`
- `fc00:0:7:0:1::`
- `fc00:0:c:0:1::`

The packets traverse the network in the following sequence:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-3 -> XR-7 -> SITE-C -> HOST-C

This path is optimized to prioritize minimizing packet loss while also considering latency and jitter. The selected route avoids segments with higher packet loss to ensure the best possible performance. However, it does include links with higher delay and jitter as a trade-off to minimize packet loss.


![Low Latency and Low Packet Loss Path](../../images/hawkv6-low-loss-low-latency-low-jitter-default-weights-intent.drawio.svg)


### Scenario 2:  Low Packet Loss, Low Latency, and Low Jitter with non-default weights

In this scenario, the path optimization focuses on minimizing packet loss, with secondary considerations for latency and jitter. The `HAWKEYE_THREE_FACTOR_WEIGHTS` environment variable is configured with weights of 0.5 for packet loss, 0.3 for latency, and 0.2 for jitter. This setup ensures that the HawkEye controller prioritizes low packet loss, followed by low latency and jitter, when calculating the optimal path between the two hosts.


#### HawkWing Configuration
The HawkWing configuration is identical to the previous scenario.


#### API Request
The api request is identical to the previous scenario.

#### Result

The following SID List is generated to optimize the path between the two hosts:

- `fc00:0:1:0:1::`
- `fc00:0:2:0:1::`
- `fc00:0:3:0:1::`
- `fc00:0:7:0:1::`
- `fc00:0:c:0:1::`

Packets are routed through the devices in this sequence:
HOST-A -> SITE-A -> XR-1 -> XR-2 -> XR-3 -> XR-7 -> SITE-C -> HOST-C

This path is optimized to prioritize minimizing packet loss while also considering latency and jitter. The result illustrates that, in this scenario, the influence of low latency and low jitter led to the selection of a path that includes the higher packet loss link between XR-5 and XR-7. This outcome highlights the significant impact that weight selection can have on the path calculation.


![Low Low Packet Loss and Low Latency Path](../../images/hawkv6-low-loss-low-latency-low-jitter-non-default-weights-intent.drawio.svg)

