# Intents Overview

This section provides a detailed overview of all available intents within the platform. Each intent serves as a unique identifier for a specific action that can be initiated, enabling the platform to interpret and execute user requests accordingly.

For each intent, an example is provided that includes the relevant HawkWing configuration and the corresponding API request. The examples are organized in the [`/api`](/api/) folder, where [Kreya API Files](https://kreya.app/) can be easily imported for practical use. Each example has been thoroughly tested within the provided test network environment, offering practical references for implementing these intents in various configurations.

## Intents Categories

The intents are categorized into four main groups: **Single Intents**, **Combined Intents**, **Service Function Chaining (SFC)**, and **Flexible Algorithm (Flex Algo)**. Each category includes detailed information and examples to help implement the intents effectively.

### Single Intents

Single intents are used to define a specific path between two hosts based on a single network metric. Below are the available single intents:

- **Low Latency**: [Learn more](single-intent/low-latency.md)
- **Low Jitter**: [Learn more](single-intent/low-jitter.md)
- **Low Packet Loss**: [Learn more](single-intent/low-packet-loss.md)
- **High Bandwidth**: [Learn more](single-intent/high-bandwidth.md)
- **Low Bandwidth**: [Learn more](single-intent/low-bandwidth.md)
- **Low Utilization**: [Learn more](single-intent/low-utilization.md)

### Combined Intents

Combined intents allow optimization of paths between hosts using multiple metrics. Below are the available combined intents:

- **Two Intents**: [Learn more](combined-intents/two-intents.md)
- **Three Intents**: [Learn more](combined-intents/three-intents.md)
- **Constrained Intents**: [Learn more](combined-intents/constrained-intents.md)

### Service Function Chaining (SFC)

SFC intents are used to define paths between two hosts that include specific service functions. Below are the available SFC intents:

- **SFC with IGP Metric**: [Learn more](sfc/sfc-igp-metric.md)
- **SFC with Other Metrics**: [Learn more](sfc/sfc-other-metrics.md)

### Flexible Algorithm (Flex Algo)

Flex Algo intents allow calculation of paths on specific subgraphs of the network topology, enabling the exclusion of certain links or nodes. Below are the available Flex Algo intents:

- **Flex Algo**: [Learn more](flex-algo/flex-algo-overview.md)
