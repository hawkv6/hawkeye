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
HawkEye is an advanced controller designed to bring Intent-Based Networking (IBN) to life within Segment Routing over IPv6 (SRv6) networks. The project showcases how HawkEye translates high-level user intents into precise network instructions (segments), ensuring seamless and efficient network management.

![HawkEye Architecture](docs/images/Hawkv6-High-Level-Architecture.drawio.svg)

HawkEye continuously maintains an up-to-date view of the network by requesting and subscribing to data from the Jalapeno API Gateway (JAGW). This real-time data allows HawkEye to monitor network changes and respond proactively. Additionally, HawkEye integrates with Consul as a service registry, obtaining vital information about network services and conducting health checks to ensure optimal performance.

HawkEye includes a gRPC API that enables clients within the SRv6 network to send intents directly to the controller. When an intent is received, HawkEye computes the optimal segment list to ensure the network behavior aligns with the specified intent.


## Key Features
- **Intent Fulfillment**: HawkEye processes user-defined intents, such as ensuring low latency, high bandwidth, or specific service chaining, and calculates optimal paths in the SRv6 network to meet these requirements. It returns the segment list that must be applied to the packet to achieve the desired behavior.

- **Real-Time Network Monitoring**: HawkEye continuously monitors network performance metrics, including latency, jitter, and packet loss. It responds to network changes in real-time, ensuring that the network adheres to the specified intents.

- **Event-Driven Architecture**: HawkEye operates in an event-driven manner, automatically recalculating paths and updating segment lists when network conditions change or when service health checks indicate issues. The operator can define a threshold which defines when the path should be switched. This feature ensures that the network remains responsive and reliable and does not flap between paths unnecessarily.

- **Interoperability**: HawkEye leverages standardized technologies like YANG-Push for telemetry and BMP for performance measurement, ensuring compatibility with existing network hardware and software.


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
   - Detailed instructions are available in the [clab-telemetry-linker documentation](https://github.com/hawkv6/generic-processor).

4. Confirm that the `generic-processor` is active and running.
   - Detailed instructions are available in the [generic processor documentation](https://github.com/hawkv6/generic-processor).

5. Install the HawkEye controller using one of the methods described above.

5. Start the HawkEye controller.
   - For more information, refer to the [start command documentation](docs/commands/start.md).

## Additional Information
- The design considerations are documented in the [design file](docs/design.md).
- Environment variables are documented in the [env file](docs/env.md).
- The proto definiton is included via submodule and can be found [here](https://github.com/hawkv6/proto/blob/main/intent.proto).
- Limitations are documented in the [limitations file](docs/limitations.md).