# AutoScaler

AutoScaler is an application that dynamically adjusts the number of replicas based on CPU utilization of a target application. It fetches status from an upstream service, calculates new replica counts, and updates them as needed.

## Features

- Automatic scaling of application replicas based on CPU utilization.
- Configuration via environment variables or configuration file (using `viper`).
- Flexible logging with adjustable log levels.
- HTTP client for making API requests.
- Structured logging for easy monitoring and debugging.

## Config parameters exposed
| Name                     | env var to set                 | default                                | description                                                                                                                |
|--------------------------|--------------------------------|----------------------------------------|----------------------------------------------------------------------------------------------------------------------------|
| TargetCPU                | `APP_TARGETCPU`                | `0.80`                                 | TargetCPU is the target CPU utilization percentage to maintain. (e.g., 0.80 for 80%)                                       |
| StatusEndpoint           | `APP_STATUSENDPOINT`           | `http://localhost:8123/app/status`     | StatusEndpoint is the URL endpoint for fetching current application status.                                                |
| ReplicasEndpoint         | `APP_REPLICAENDPOINT`          | `http://localhost:8123/app/replicas`   | ReplicasEndpoint is the URL endpoint for updating the number of application replicas.                                      |
| PollInterval             | `APP_POLLINTERVAL`             | `10s`                                  | PollInterval is the interval between consecutive status checks. (3s)                                                       |
| CoolDownPeriod           | `APP_COOLDOWNPERIOD`           | `20`                                   | CoolDownPeriod is the duration to wait after scaling replicas before making another scaling decision.                      |
| ApiTimeOut               | `APP_APITIMEOUT`               | `2s`                                   | ApiTimeOut is the timeout duration for API requests.                                                                       |
| LogLevel                 | `APP_LOGLEVEL`                 | `DEBUG`                                | LogLevel sets the logging level for the application (e.g., "debug", "info", "warn").                                       |
| DownscaleAfterAttempts   | `APP_DOWNSCALEAFTERATTEMPTS`   | `3`                                    | DownscaleAfterAttempts specifies the number of retry attempts after which the system initiates downscaling of resources.   |


## Main logic.
This section describes the logic for dynamically adjusting the number of replicas in response to CPU utilization.
#### Scaling Up:
When the current CPU utilization exceeds the target CPU utilization, the system increases the number of replicas to handle the higher load. 
The new number of replicas is calculated by taking the current number of replicas and scaling it proportionally to the CPU overage. The formula used is:
```
newReplicas=ceil(⌈currentReplicas×(currentCPU/TargetCPU)])
```
#### Scaling Down:
When the current CPU utilization is at or below the target CPU utilization, the system attempts to decrease the number of replicas. However, to avoid rapid scaling down and potential instability, the system only reduces the replicas after a certain number of consecutive downscale attempts. 
If the downscale attempts threshold is reached, the number of replicas is scaled down proportionally. The formula used is:
```
newReplicas=floor(⌈currentReplicas×(currentCPU/TargetCPU)])
```

## Edge cases
- If `TargetCPU` is set to zero, the logic will attempt to divide by zero, causing a runtime error. (as per the current code, it will set new replicas to 1)
- If currentCPU is zero, the downscaling logic may produce unexpected results. For example, using the downscale formula with zero CPU would result in zero replicas, which might not be desirable.
- If currentReplicas is zero, any calculations based on the current number of replicas will result in zero. This could prevent scaling up when needed.
- If `currentCPU` is negative due to a misconfiguration or erroneous input, the scaling logic will produce nonsensical results.
- If `downscaleAfterAttempts` is very high, it might take too long to scale down, leading to unnecessary resource usage and costs.
- There is no upper limit in the logic. If the currentCPU is significantly higher than the `TargetCPU`, the logic could request an unreasonably high number of replicas.
- Similar to the maximum limit, there should be a lower bound to prevent the number of replicas from dropping below a safe operational threshold.

## Build and Deployment Instructions

This project uses a `Makefile` to manage common tasks such as building, running, and pushing the application. Below are the commands you can use with `make` to perform these tasks.


## Unit test
```
❯  go clean -i -cache && go test -v  ./...                                                                                                                                                                                           ─╯
?       github.com/suyog1pathak/autoscaler/api/v1/model [no test files]
?       github.com/suyog1pathak/autoscaler/cmd  [no test files]
?       github.com/suyog1pathak/autoscaler/pkg/config   [no test files]
?       github.com/suyog1pathak/autoscaler/pkg/util     [no test files]
=== RUN   TestCalculateNewReplicas
--- PASS: TestCalculateNewReplicas (0.00s)
PASS
ok      github.com/suyog1pathak/autoscaler/pkg/autoscaler       1.624s
=== RUN   TestClient
=== RUN   TestClient/GET_/test
=== RUN   TestClient/POST_/test
--- PASS: TestClient (0.02s)
    --- PASS: TestClient/GET_/test (0.02s)
    --- PASS: TestClient/POST_/test (0.00s)
PASS
ok      github.com/suyog1pathak/autoscaler/pkg/rest     1.456s

```


### Prerequisites

Ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (1.22)
- [Docker](https://docs.docker.com/get-docker/) (optional)

### Makefile Commands

#### `start`

This command runs the application.

```
make start
```
#### `build`
This command compiles the application and creates an executable binary.
```
make build
```

#### `docker-build`
This command builds a Docker image for the application. The image is tagged with the specified artifactURL and TAG. By default, artifactURL is set to services and TAG is set to latest.
```
make docker-build artifactURL=my-repo/my-service TAG=v1.0.0
```

#### `docker-push`
This command pushes the Docker image to the specified repository. Ensure you have the necessary permissions to push to the repository.
`default: services:latest`
```
make docker-push artifactURL=my-repo/my-service TAG=v1.0.0
```


