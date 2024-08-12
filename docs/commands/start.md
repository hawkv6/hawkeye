# Start the Controller

## Overview
The `start` command starts the HawkEye controller.

## Command Syntax
```bash
hawkeye start -j <jagw-service-address> -s <jagw-subscription-port> -r <jagw-request-port> -p <grpc-port> -c <consul-server-address>
```

- `-j` or `--jagw-service-address`: The address of the Jalapeno API Gateway (JAGW) service if not set via the environment variable `HAWKEYE_JAGW_SERVICE_ADDRESS`.
- `-s` or `--jagw-subscription-port`: The port number for JAGW subscription if not set via the environment variable `HAWKEYE_JAGW_SUBSCRIPTION_PORT`.
- `-r` or `--jagw-request-port`: The port number for JAGW request if not set via the environment variable `HAWKEYE_JAGW_REQUEST_PORT`.
- `-p` or `--grpc-port`: The port number for the gRPC API if not set via the environment variable `HAWKEYE_GRPC_PORT`.
- `-c` or `--consul-server-address`: The address of the Consul server if not set via the environment variable `HAWKEYE_CONSUL_SERVER_ADDRESS`.

## Example
```bash
hawkeye start -j 10.8.39.69 -s 9902 -r 9903 -p 10000 -c consul-hawkv6.stud.network.garden
INFO[2024-08-12T13:01:30+02:00] Config created successfully                   subsystem=cmd
INFO[2024-08-12T13:01:30+02:00] Got 24 SRv6 SIDs from JAGW                    subsystem=jagw
INFO[2024-08-12T13:01:30+02:00] Got 84 LsPrefixes from JAGW                   subsystem=jagw
INFO[2024-08-12T13:01:30+02:00] Got 11 LsNodes from JAGW                      subsystem=jagw
INFO[2024-08-12T13:01:30+02:00] Got 36 LsLinks from JAGW                      subsystem=jagw
INFO[2024-08-12T13:01:30+02:00] Closing JAGW Request Service                  subsystem=jagw
INFO[2024-08-12T13:01:30+02:00] Starting controller                           subsystem=controller
INFO[2024-08-12T13:01:30+02:00] Starting JAGW Subscription Service            subsystem=jagw
INFO[2024-08-12T13:01:30+02:00] Starting processing network updates with hold time 1s  subsystem=processor
INFO[2024-08-12T13:01:30+02:00] Starting monitoring services                  subsystem=service
INFO[2024-08-12T13:01:30+02:00] Listening on :10000                           subsystem=messaging
INFO[2024-08-12T13:01:30+02:00] Service SERA-1 from type fw with sid fc00:0:2f:: created - healthy: true  subsystem=service
INFO[2024-08-12T13:01:30+02:00] Service SERA-2 from type fw with sid fc00:0:3f:: created - healthy: true  subsystem=service
INFO[2024-08-12T13:01:30+02:00] Services or service health changed - sending update message  subsystem=service
INFO[2024-08-12T13:01:30+02:00] Service SNORT-1 from type ids with sid fc00:0:6f:: created - healthy: true  subsystem=service
INFO[2024-08-12T13:01:30+02:00] Service SNORT-2 from type ids with sid fc00:0:7f:: created - healthy: true  subsystem=service
INFO[2024-08-12T13:01:30+02:00] Services or service health changed - sending update message  subsystem=service
^C
INFO[2024-08-12T13:01:33+02:00] Received interrupt signal, shutting down      subsystem=cmd
INFO[2024-08-12T13:01:33+02:00] Stopping the gRPC server                      subsystem=messaging
INFO[2024-08-12T13:01:33+02:00] Stopping JAGW Subscription Service            subsystem=jagw
INFO[2024-08-12T13:01:35+02:00] Stopping monitoring services, can take up to 10s  subsystem=service
INFO[2024-08-12T13:01:40+02:00] Stopping monitoring health state for service type: ids  subsystem=service
INFO[2024-08-12T13:01:40+02:00] Stopping monitoring health state for service type: fw  subsystem=service
INFO[2024-08-12T13:01:40+02:00] Stopping network processor                    subsystem=processor
INFO[2024-08-12T13:01:40+02:00] Stopping session controller                   subsystem=controller
INFO[2024-08-12T13:01:40+02:00] All services stopped successfully             subsystem=cmd
```