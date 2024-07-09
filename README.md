# Kaas

Kaas (Kubernetes as a Service) is a self-service platform written with Golang that allows users to deploy, monitor, and
manage their applications and databases effortlessly.

## Technologies used

- [Golang](https://golang.org/), Programming language.
- [Echo](https://echo.labstack.com/), HTTP web framework.
- [PostgreSQL](https://www.postgresql.org/), Database management system.
- [Minikube](https://minikube.sigs.k8s.io/docs/), Tool for running Kubernetes locally.
- [Helm](https://helm.sh/), Kubernetes package manager.
- [Docker](https://www.docker.com/), Containerization platform.

## Features

- **Self-Service Deployment**: Users can deploy their applications and databases with a simple request.

- **Database Management**: Automates the creation and management of PostgreSQL databases.

- **Health Monitoring**: Regularly checks the status of deployed applications and records their health metrics.

## Getting Started

### Prerequisites

- Golang
- Docker
- Kubernetes cluster

### Installation

Build the Docker image:

```sh
docker build -t <image> -f build/kaas/Dockerfile .
 ```

Deploy to your Kubernetes cluster using Helm:

```sh
helm package kaas-api
 ```

```sh
helm install kaas-api-release ./kaas-api-0.1.0.tgz
 ```

### Configuration

Modify the `values.yaml` file to set your desired configurations, such as replica counts, image repositories, and
database settings.

## Monitoring

To monitor the health of an application, I use a CronJob that regularly checks the application's health status every 5
minutes. You can still change this interval as per your requirements.

## Examples

### Deploying an Application

To deploy an application, send a POST request to the `/create` endpoint with the following JSON payload:

```json
{
  "AppName": "example-application",
  "Replicas": 3,
  "ImageAddress": "your-docker-repo/your-app-image",
  "ImageTag": "latest",
  "ServicePort": 8080,
  "Resources": {
    "CPU": "500m",
    "RAM": "256Mi"
  },
  "Envs": [
    {
      "Key": "ENV_VAR1",
      "Value": "value1",
      "IsSecret": true
    },
    {
      "Key": "ENV_VAR2",
      "Value": "value2",
      "IsSecret": true
    }
  ]
}
```

### Check Application Status

To check the status of a specific application, send a GET request to the `/status/{appName}` endpoint,
replacing `{appName}` with the name of your application. Example:

```json
{
  "DeploymentName": "example-application",
  "Replicas": 3,
  "ReadyReplicas": 2,
  "PodStatuses": [
    {
      "Name": "example-application-pod-1",
      "Phase": "Running",
      "HostIP": "192.168.49.2",
      "PodIP": "10.244.1.4",
      "StartTime": "2024-06-07T12:30:00Z"
    },
    {
      "Name": "example-application-pod-2",
      "Phase": "Running",
      "HostIP": "192.168.49.2",
      "PodIP": "10.244.1.5",
      "StartTime": "2024-06-07T12:32:00Z"
    },
    {
      "Name": "example-application-pod-3",
      "Phase": "Pending",
      "HostIP": "",
      "PodIP": "",
      "StartTime": ""
    }
  ]
}

```

### Monitoring an Application

To monitor the health of an application, send a GET request to the `/health/{appName}` endpoint, replacing `{appName}`
with the name of your application. Example:

```json
{
  "AppName": "example-application",
  "FailureCount": 0,
  "SuccessCount": 3,
  "LastSuccess": "2024-07-01T12:30:00Z",
  "LastFailure": "2024-07-01T14:00:00Z"
}

```
