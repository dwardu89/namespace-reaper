# namespace-reaper

This is a controller that will monitors namespaces in Kubernetes and destroy the ones that have expired.
It can be used when multiple namespaces are being created to test out tools, or even debug and they can be destroyed automatically after a time period has elapsed.

## Deploying to k8s

You can deploy this to your cluster using the helm command below:

``` bash

helm upgrade --install namespace-reaper helm/namespace-reaper
```
