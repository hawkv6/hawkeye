# hawkeye
Controller for Enabling Intent-Driven End-to-End SRv6 Networking

## Configuration

Here's how you can document these environment variables in your README file:

## Configuration

The application can be configured via command-line flags or environment variables. The following options are available:

- `--jagw-service-address` or `-j`: This sets the JAGW Service Address. It can be a hostname like `localhost` or an IP address like `127.0.0.1`. The corresponding environment variable is `HAWKEYE_JAGW_SERVICE_ADDRESS`.

- `--jagw-request-port` or `-r`: This sets the JAGW Request Port, for example `9903`. The corresponding environment variable is `HAWKEYE_JAGW_REQUEST_PORT`.

- `--jagw-subscription-port` or `-s`: This sets the JAGW Subscription Port, for example `9902`. The corresponding environment variable is `HAWKEYE_JAGW_SUBSCRIPTION_PORT`.

- `--grpc-port` or `-p`: This sets the gRPC Port, for example `10000`. The corresponding environment variable is `HAWKEYE_GRPC_PORT`.

- `HAWKEYE_FLAPPING_THRESHOLD`: This variable sets the flapping threshold. It should be a float value. If not provided, the default value is `0.1`, which means that paths will only change if the alternative path is 10% better.

- `HAWKEYE_TWO_FACTOR_WEIGHTS`: This variable sets the weights when the request includes two factors. It should be a comma-separated string of float values. If not provided, the default value is `0.7,0.3`.

- `HAWKEYE_THREE_FACTOR_WEIGHTS`: This variable sets the weights when the request includes three factors. It should be a comma-separated string of float values. If not provided, the default value is `0.5,0.3,0.2`.


You can set the environment variables in your shell before running the application, like this:

```bash
export HAWKEYE_JAGW_SERVICE_ADDRESS=localhost
export HAWKEYE_JAGW_REQUEST_PORT=9903
export HAWKEYE_JAGW_SUBSCRIPTION_PORT=9902
export HAWKEYE_GRPC_PORT=10000
export HAWKEYE_FLAPPING_THRESHOLD=0.1
export HAWKEYE_TWO_FACTOR_WEIGHTS=0.7,0.3
export HAWKEYE_THREE_FACTOR_WEIGHTS=0.5,0.3,0.2

Or you can pass the command-line flags when running the application, like this:

```bash
./bin/hawkeye --jagw-service-address localhost --jagw-request-port 9903 --jagw-subscription-port 9902 --grpc-port 10000
```