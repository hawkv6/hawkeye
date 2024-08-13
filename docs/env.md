# Environment Variables

## Configuration

HawkEye can be configured using the following environment variables:

- **`HAWKEYE_JAGW_SERVICE_ADDRESS`**: Sets the JAGW Service Address. Accepts a hostname (e.g., `localhost`) or an IP address (e.g., `127.0.0.1`).

- **`HAWKEYE_JAGW_REQUEST_PORT`**: Specifies the JAGW Request Port, for example, `9903`.

- **`HAWKEYE_JAGW_SUBSCRIPTION_PORT`**: Specifies the JAGW Subscription Port, for example, `9902`.

- **`HAWKEYE_GRPC_PORT`**: Defines the gRPC Port, for example, `10000`.

- **`HAWKEYE_FLAPPING_THRESHOLD`**: Sets the flapping threshold as a float value. The default is `0.1`, meaning paths change only if the alternative path is 10% better.

- **`HAWKEYE_TWO_FACTOR_WEIGHTS`**: Sets the weights for requests involving two factors. Accepts a comma-separated string of float values. Default is `0.7,0.3`.

- **`HAWKEYE_THREE_FACTOR_WEIGHTS`**: Sets the weights for requests involving three factors. Accepts a comma-separated string of float values. Default is `0.7,0.2,0.1`.

- **`HAWKEYE_SKIP_TLS_VERIFICATION`**: Skips TLS verification when set to `true` or `TRUE`. The default is `false`.

- **`HAWKEYE_CONSUL_QUERY_WAIT_TIME`**: Sets the wait time for Consul long-polling queries. The default is `5s`.

- **`HAWKEYE_NETWORK_PROCESSOR_HOLD_TIME`**: Sets the hold time for the network processor. The default is `1s`. Meaning the network processor will trigger a recalculation if no updates are received within x seconds.
