# CORD Infrastructure Initializer

The CORD Infrastructure Initializer is a [Kubernetes initializer](https://kubernetes.io/docs/admin/extensible-admission-controllers/#what-are-initializers) that injects a sidecar container into a pod based on policy.

## Usage

```
cord-infra-initializer -h
```
```
Usage of cord-infra-initializer:
  -annotation string
    	The annotation to trigger initialization (default "initializer.onf.io/cord-infra")
  -configmap string
    	The cord infrastructure initializer configuration configmap (default "cord-infra-initializer")
  -initializer-name string
    	The initializer name (default "cord-infra.initializer.onf.io")
  -namespace string
    	The configuration namespace (default "utility")
  -require-annotation
    	Require annotation for initialization
```
