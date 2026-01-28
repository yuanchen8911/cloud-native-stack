# CNCF AI Conformance v1.34: Minimum Component Set

## Target: Kubernetes 1.34 with Open Source, Helm-Based Components

**Last Updated:** January 15, 2026

---

## Executive Summary

Based on thorough analysis of the CNCF AI Conformance v1.34 specification and NVIDIA's open source stack, here is the **minimum set** of components needed to satisfy all 9 MUST requirements.

### The Minimum Stack (4 Components)

| Component | Helm Chart | Version | Requirements Satisfied |
|-----------|------------|---------|----------------------|
| **NVIDIA GPU Operator** | `nvidia/gpu-operator` | v25.10.1 | `dra_support`*, `secure_accelerator_access`, `accelerator_metrics`, `pod_autoscaling`** |
| **Dynamo Platform*** | `nvidia/ai-dynamo/dynamo-platform` | v0.8.0 | `robust_controller`, `ai_inference`, `gang_scheduling` |
| **kube-prometheus-stack** | `prometheus-community/kube-prometheus-stack` | Chart 80.14.4 | `ai_service_metrics` |
| **Cluster Autoscaler** | `autoscaler/cluster-autoscaler` | App 1.34.2 | `cluster_autoscaling` |

*`dra_support` requires separate `nvidia-dra-driver-gpu` chart (Technology Preview)
**`pod_autoscaling` uses standard Kubernetes HPA with DCGM Exporter metrics
***With `--set grove.enabled=true --set kai-scheduler.enabled=true` for gang scheduling

> **⚠️ Important:** Grove and KAI Scheduler are NOT bundled by default in Dynamo Platform - they must be explicitly enabled.

**Total: 4 Helm releases**

---

## Requirements Coverage Matrix

| Requirement | Component | Version | Notes |
|-------------|-----------|---------|-------|
| **dra_support** | GPU Operator + DRA Driver | v25.10.1 + 25.3.0 | Device plugin included; DRA driver is **separate chart** (Tech Preview) |
| **ai_inference** | Dynamo Platform | v0.8.0 | OpenAI-compatible frontend + Gateway API via GAIE |
| **gang_scheduling** | Dynamo + KAI Scheduler | v0.12.6 | `--set kai-scheduler.enabled=true` (NOT bundled by default) |
| **cluster_autoscaling** | Cluster Autoscaler | 1.34.2 | **Separate install** - not part of any NVIDIA component |
| **pod_autoscaling** | GPU Operator + HPA | v25.10.1 | Standard HPA with DCGM Exporter metrics |
| **accelerator_metrics** | GPU Operator | v25.10.1 | DCGM Exporter (v4.4.2-1) included |
| **ai_service_metrics** | kube-prometheus-stack | Chart 80.14.4 | Prometheus + ServiceMonitor/PodMonitor CRDs |
| **secure_accelerator_access** | GPU Operator | v25.10.1 | Device plugin isolation + CDI (enabled by default) |
| **robust_controller** | Dynamo Operator | v0.8.0 | DynamoGraphDeployment + 4 other CRDs |

---

## Component Details

### 1. NVIDIA GPU Operator

**Version:** v25.10.1 (December 2024)

**What's bundled:**
- NVIDIA GPU Drivers (v580.105.08)
- Container toolkit (nvidia-container-runtime v1.18.1)
- Device plugin (GPU allocation & isolation)
- DCGM Exporter (v4.4.2-1, GPU metrics for Prometheus)
- GPU Feature Discovery (node labeling)
- MIG Manager (for GPU partitioning)
- Container Device Interface (CDI) - enabled by default
- GDS (GPUDirect Storage) Driver

**What's NOT bundled (separate chart):**
- DRA Driver (`nvidia-dra-driver-gpu`) - required for `dra_support` conformance

**Installation:**
```bash
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update

# GPU Operator v25.10.1
helm install gpu-operator nvidia/gpu-operator \
  --namespace gpu-operator \
  --create-namespace \
  --version v25.10.1 \
  --set dcgmExporter.serviceMonitor.enabled=true

# Optional: DRA Driver (Technology Preview for GPU allocation)
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu \
  --version 25.3.0 \
  --namespace nvidia-dra-driver-gpu \
  --create-namespace
```

**Satisfies:** `dra_support`*, `secure_accelerator_access`, `accelerator_metrics`, `pod_autoscaling`

*With separate DRA driver installation

---

### 2. Dynamo Platform (with Grove + KAI Scheduler)

**Version:** v0.8.0 (January 15, 2025)

**What's bundled by default:**
- **Dynamo Operator** - Manages DynamoGraphDeployment and related CRDs
- **NATS** - Messaging layer
- **etcd** - Configuration store
- **OpenAI-compatible frontend** - HTTP API server
- **Kubernetes-native service discovery** - Uses EndpointSlices

**What's NOT bundled (must enable with flags):**
- **Grove** - Enable with `--set grove.enabled=true`
- **KAI Scheduler** - Enable with `--set kai-scheduler.enabled=true`

**CRDs Provided:**
- `DynamoGraphDeployment (DGD)` - Core deployment resource
- `DynamoGraphDeploymentRequest (DGDR)` - SLA-driven deployments
- `DynamoComponentDeployment` - Individual component specs
- `DynamoGraphDeploymentScalingAdapter (DGDSA)` - HPA/KEDA integration
- `DynamoModel` - Model references and endpoint discovery

**Installation:**
```bash
export DYNAMO_VERSION=0.8.0  # https://github.com/ai-dynamo/dynamo/releases

# Install CRDs
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-crds-${DYNAMO_VERSION}.tgz
helm install dynamo-crds dynamo-crds-${DYNAMO_VERSION}.tgz --namespace default

# Install Platform with Grove + KAI Scheduler (required for gang_scheduling)
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-platform-${DYNAMO_VERSION}.tgz
helm install dynamo-platform dynamo-platform-${DYNAMO_VERSION}.tgz \
  --namespace dynamo-system \
  --create-namespace \
  --set grove.enabled=true \
  --set kai-scheduler.enabled=true
```

> **⚠️ Note:** For `gang_scheduling` conformance, you MUST enable KAI Scheduler with `--set kai-scheduler.enabled=true`

**Satisfies:** `robust_controller`, `ai_inference`, `gang_scheduling`*

*Only with `--set kai-scheduler.enabled=true`

---

### 3. kube-prometheus-stack

**Version:** Chart 80.14.4 / App v0.87.1 (January 2025)

**What's included:**
- Prometheus Operator (v0.87.1)
- Prometheus (metrics collection)
- Alertmanager (alerting)
- Grafana (visualization)
- kube-state-metrics
- Node Exporter (host metrics)
- ServiceMonitor/PodMonitor CRDs

**Installation:**
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false
```

**Satisfies:** `ai_service_metrics`

**Note:** This is an independent component - NOT bundled with any NVIDIA product.

---

### 4. Cluster Autoscaler

**Version:** App 1.34.2 / Chart 9.54.1 (January 2025)

> **⚠️ Important:** Cluster Autoscaler is NOT bundled with Dynamo Platform or any NVIDIA component. It must be installed separately.

**For AWS:**
```bash
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm repo update

helm install cluster-autoscaler autoscaler/cluster-autoscaler \
  --namespace kube-system \
  --set cloudProvider=aws \
  --set autoDiscovery.clusterName=${CLUSTER_NAME} \
  --set awsRegion=${AWS_REGION}
```

**For GCP:**
```bash
helm install cluster-autoscaler autoscaler/cluster-autoscaler \
  --namespace kube-system \
  --set cloudProvider=gce \
  --set autoscalingGroupsnamePrefix=${MIG_PREFIX}
```

**For Azure:**
```bash
helm install cluster-autoscaler autoscaler/cluster-autoscaler \
  --namespace kube-system \
  --set cloudProvider=azure \
  --set azureResourceGroup=${RESOURCE_GROUP} \
  --set azureVMType=vmss
```

**For On-Premise:** Use Cluster API provider or skip if nodes are statically provisioned.

**Satisfies:** `cluster_autoscaling`

---

## Complete Installation Script

```bash
#!/bin/bash
# CNCF AI Conformance v1.34 - Minimum Stack
# Updated: January 2026
set -e

export DYNAMO_VERSION=0.8.0  # https://github.com/ai-dynamo/dynamo/releases
export CLUSTER_NAME=my-cluster

# Add all Helm repos
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm repo update

# 1. GPU Operator v25.10.1
echo "Installing GPU Operator v25.10.1..."
helm install gpu-operator nvidia/gpu-operator \
  --namespace gpu-operator --create-namespace \
  --version v25.10.1 \
  --set dcgmExporter.serviceMonitor.enabled=true

# 2. Dynamo Platform v0.8.0 (with Grove + KAI Scheduler for gang_scheduling)
echo "Installing Dynamo Platform v0.8.0..."
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-crds-${DYNAMO_VERSION}.tgz
helm install dynamo-crds dynamo-crds-${DYNAMO_VERSION}.tgz --namespace default

helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-platform-${DYNAMO_VERSION}.tgz
helm install dynamo-platform dynamo-platform-${DYNAMO_VERSION}.tgz \
  --namespace dynamo-system --create-namespace \
  --set grove.enabled=true \
  --set kai-scheduler.enabled=true

# 3. Prometheus Stack (Chart 80.14.4)
echo "Installing kube-prometheus-stack..."
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false

# 4. Cluster Autoscaler (App 1.34.2) - cloud only, NOT bundled with Dynamo
# Uncomment for your cloud provider:
# echo "Installing Cluster Autoscaler..."
# helm install cluster-autoscaler autoscaler/cluster-autoscaler \
#   --namespace kube-system \
#   --set autoDiscovery.clusterName=${CLUSTER_NAME} \
#   --set cloudProvider=aws

echo "Minimum AI Conformance stack installed!"
echo ""
echo "Component versions installed:"
echo "  - GPU Operator: v25.10.1"
echo "  - Dynamo Platform: v0.8.0 (with Grove + KAI Scheduler)"
echo "  - kube-prometheus-stack: Chart 80.14.4"
echo "  - Cluster Autoscaler: (uncomment to install)"
```

---

## What Each NVIDIA Component Provides

### GPU Operator vs Dynamo vs KAI Scheduler vs Grove

| Component | Layer | Purpose |
|-----------|-------|---------|
| **GPU Operator** | Infrastructure | GPU drivers, device plugin, metrics, isolation |
| **Dynamo Platform** | Application | LLM inference orchestration, CRDs, frontend |
| **KAI Scheduler** | Scheduling | GPU-aware scheduling, fractional GPUs, gang |
| **Grove** | Orchestration | Disaggregated workloads, PodCliqueSet, autoscaling |

### Integration Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     USER WORKLOADS                          │
│         DynamoGraphDeployment / PodCliqueSet CRs            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              DYNAMO PLATFORM                         │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────────┐  │   │
│  │  │ Dynamo  │ │  Grove  │ │  KAI    │ │ Frontend  │  │   │
│  │  │Operator │ │         │ │Scheduler│ │ (OpenAI)  │  │   │
│  │  └─────────┘ └─────────┘ └─────────┘ └───────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              GPU OPERATOR                            │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────────┐  │   │
│  │  │ Driver  │ │ Device  │ │  DCGM   │ │ Container │  │   │
│  │  │         │ │ Plugin  │ │Exporter │ │ Toolkit   │  │   │
│  │  └─────────┘ └─────────┘ └─────────┘ └───────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │            KUBERNETES + PROMETHEUS                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Alternative Configurations

### Option A: Minimal (No Multi-node/Disaggregated)

If you only need single-node inference (no disaggregated prefill/decode):

```bash
# Skip Grove and KAI Scheduler
helm install dynamo-platform dynamo-platform-${RELEASE_VERSION}.tgz \
  --namespace dynamo-system --create-namespace
  # grove.enabled and kai-scheduler.enabled default to false

# Use Kueue for gang scheduling instead
kubectl apply --server-side -f \
  https://github.com/kubernetes-sigs/kueue/releases/download/v0.14.2/manifests.yaml
```

### Option B: Maximum Performance (Full NVIDIA Stack)

For production LLM inference at scale:

```bash
helm install dynamo-platform dynamo-platform-${RELEASE_VERSION}.tgz \
  --namespace dynamo-system --create-namespace \
  --set grove.enabled=true \
  --set kai-scheduler.enabled=true

# Plus: Deploy Inference Gateway for Gateway API
cd /path/to/dynamo/deploy/inference-gateway
./install_gaie_crd_kgateway.sh
```

---

## Verification Tests

### Test 1: GPU Isolation (secure_accelerator_access)
```bash
# Pod without GPU request should NOT see GPUs
kubectl run no-gpu-test --image=busybox --rm -it --restart=Never -- \
  sh -c "ls /dev/nvidia* 2>&1 || echo 'PASS: No GPU access'"
```

### Test 2: GPU Metrics (accelerator_metrics)
```bash
kubectl port-forward -n monitoring svc/kube-prometheus-stack-prometheus 9090:9090 &
curl 'localhost:9090/api/v1/query?query=DCGM_FI_DEV_GPU_UTIL'
```

### Test 3: AI Operator (robust_controller)
```bash
kubectl apply -f examples/backends/vllm/deploy/agg.yaml -n dynamo-system
kubectl get dynamographdeployment -n dynamo-system
```

### Test 4: Gang Scheduling
```bash
# With KAI Scheduler
kubectl get pods -n kai-scheduler
# Verify jobs with insufficient resources stay pending as a gang
```

---

## Summary

| Requirement | Component | Version | Installation |
|-------------|-----------|---------|--------------|
| `dra_support` | GPU Operator + DRA Driver | v25.10.1 + 25.3.0 | `nvidia/gpu-operator` + `nvidia/nvidia-dra-driver-gpu` |
| `ai_inference` | Dynamo Platform | v0.8.0 | `nvidia/ai-dynamo/dynamo-platform` |
| `gang_scheduling` | KAI Scheduler (via Dynamo) | v0.12.6 | `--set kai-scheduler.enabled=true` |
| `cluster_autoscaling` | Cluster Autoscaler | 1.34.2 | `autoscaler/cluster-autoscaler` (**separate**) |
| `pod_autoscaling` | GPU Operator + HPA | v25.10.1 | DCGM metrics for HPA scaling |
| `accelerator_metrics` | GPU Operator | v25.10.1 | DCGM Exporter v4.4.2-1 included |
| `ai_service_metrics` | kube-prometheus-stack | Chart 80.14.4 | `prometheus-community/kube-prometheus-stack` |
| `secure_accelerator_access` | GPU Operator | v25.10.1 | Device Plugin + CDI included |
| `robust_controller` | Dynamo Operator | v0.8.0 | Included with Dynamo Platform |

**Total Helm releases: 4** (GPU Operator, Dynamo Platform, Prometheus Stack, Cluster Autoscaler)

> **Key Points:**
> - Grove and KAI Scheduler are NOT bundled by default in Dynamo - must enable with `--set` flags
> - DRA Driver is separate from GPU Operator - install `nvidia-dra-driver-gpu` chart
> - Cluster Autoscaler is NOT part of any NVIDIA component - separate installation required

---

## References

- [NVIDIA GPU Operator](https://github.com/NVIDIA/gpu-operator)
- [NVIDIA Dynamo](https://github.com/ai-dynamo/dynamo)
- [NVIDIA KAI Scheduler](https://github.com/NVIDIA/KAI-Scheduler)
- [NVIDIA Grove](https://github.com/NVIDIA/grove)
- [CNCF k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance)
- [WG AI Conformance](https://github.com/kubernetes-sigs/wg-ai-conformance)
