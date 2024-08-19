<h1 align="center">HawkEye</h1>
<p align="center">
    <br>
    <img alt="GitHub Release" src="https://img.shields.io/github/v/release/hawkv6/hawkeye?display_name=release&style=flat-square">
    <img src="https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat-square">
    <img src="https://img.shields.io/github/actions/workflow/status/hawkv6/hawkeye/testing.yaml?style=flat-square&label=tests">
    <img src="https://img.shields.io/codecov/c/github/hawkv6/hawkeye?style=flat-square">
    <img src="https://img.shields.io/github/actions/workflow/status/hawkv6/hawkeye/golangci-lint.yaml?style=flat-square&label=checks">
</p>

<p align="center">
</p>


## Overview
HawkEye is an advanced controller designed to facilitate Intent-Based Networking (IBN) and Segment Routing over IPv6 (SRv6). It integrates with open-source projects like [Jalapeno](https://github.com/cisco-open/jalapeno) and the [Jalapeno API Gateway](https://github.com/jalapeno-api-gateway), and utilizes standardized protocols such as [BMP](https://datatracker.ietf.org/doc/rfc7854/) and [YANG-Push](https://datatracker.ietf.org/doc/rfc8641/). Additionally, HawkEye leverages [Consul](https://www.consul.io/) as a service registry to gather critical information about network services and perform health checks.

When receiving path requests from clients via gRPC, HawkEye calculates optimal paths based on the specified intents and maps these paths to segment lists, which are then returned to the clients for packet encapsulation. The application continuously monitors the network, adjusting paths as network conditions evolve. If a change is necessary to maintain the intent, an updated segment list is automatically provided to the client, ensuring that the network meets the desired requirements.

![HawkEye Architecture](docs/images/Hawkv6-High-Level-Architecture-high.drawio.svg)

## Key Features

- **Intent Fulfillment**: HawkEye processes user-defined intents, such as ensuring low latency, high bandwidth, or specific service function chaining, and calculates optimal paths in the SRv6 network to meet these requirements. It returns the segment list that must be applied to the packet to achieve the desired behavior.

- **Real-Time Network Monitoring**: HawkEye continuously monitors network performance metrics, including latency, jitter, and packet loss. It responds to network changes such as link removals or performance metric variations, ensuring that the network adheres to the specified intents. This task leverages components like [Jalapeno](https://github.com/cisco-open/jalapeno), [JAGW](https://github.com/jalapeno-api-gateway), [Consul](https://www.consul.io/), the [generic-processor](https://github.com/hawkv6/generic-processor), and the [clab-telemetry-linker](https://github.com/hawkv6/clab-telemetry-linker).

- **Event-Driven Architecture**: HawkEye operates in an event-driven manner, automatically recalculating paths and updating segment lists when network conditions change or when service health checks indicate issues. Operators can define thresholds that determine when a path should be switched, preventing unnecessary path flapping. This feature ensures that the network remains responsive and reliable.

- **Service Registry Integration**: HawkEye integrates with [Consul](https://www.consul.io/) to monitor service availability and health. This integration ensures that service function chaining remains intact and that paths are adjusted if a service instance becomes unavailable.

- **Session Management**: HawkEye tracks ongoing sessions for each client, maintaining a record of path requests and updates. This session management capability allows for seamless adjustments to paths as network conditions change, ensuring continuous intent compliance.

- **Interoperability**: HawkEye leverages standardized technologies like YANG-Push for telemetry and BMP for performance measurement, ensuring compatibility with existing network hardware and software.

## Design Considerations

The design principles and detailed information about the HawkEye implementation are available in the [design documentation](docs/design.md).


## Intent Overview
An overview about the intents and examples can be found in the [intent documentation](docs/intents/overview.md).


## Usage
```
hawkeye [command]
```
### Commands
- Start the controller: [`start`](docs/commands/start.md)
- Check Version: `version`

## Installation

### Using Package Manager
For Debian-based systems, install the package using apt:
```
sudo apt install ./hawkeye_{version}_amd64.deb
```

### Using Docker 
```
docker run --rm  -e HAWKEYE_JAGW_SERVICE_ADDRESS=10.8.39.69 -e HAWKEYE_JAGW_REQUEST_PORT=9903 -e HAWKEYE_JAGW_SUBSCRIPTION_PORT=9902 -e HAWKEYE_GRPC_PORT=10000 -e HAWKEYE_CONSUL_SERVER_ADDRESS=consul-hawkv6.stud.network.garden -e HAWKEYE_SKIP_TLS_VERIFICATION=TRUE ghcr.io/hawkv6/hawkeye:latest start
```

### Using Binary
```
git clone https://github.com/hawkv6/hawkeye
cd hawkeye && make binary
sudo ./bin/hawkeye
```

## Getting Started

1. Deploy all necessary Kubernetes resources.
   - For more details, refer to the [hawkv6 deployment documentation](https://github.com/hawkv6/deployment).

2. Ensure the network is properly configured and operational.
   - Additional information can be found in the [hawkv6 testnetwork documentation](https://github.com/hawkv6/network).

3. Confirm that `clab-telemetry-linker` is active and running.
   - Detailed instructions are available in the [clab-telemetry-linker documentation](https://github.com/hawkv6/clab-telemetry-linker).

4. Confirm that the `generic-processor` is active and running.
   - Detailed instructions are available in the [generic processor documentation](https://github.com/hawkv6/generic-processor).

5. Install the HawkEye controller using one of the methods described above.

5. Start the HawkEye controller.
   - For more information, refer to the [start command documentation](docs/commands/start.md).

## Additional Information
- Environment variables are documented in the [env documentation](docs/env.md).
- The proto/API definiton is included via submodule and can be found [here](https://github.com/hawkv6/proto/blob/main/intent.proto).
- Limitations are documented in the [limitations documentation](docs/limitations.md).
- Unit tests are documented in the [unit tests documentation](docs/unit-tests.md).