# CNCF Kubernetes AI Conformance v1.34: Implementation Analysis

## A Comprehensive Guide to Achieving AI Conformance Using Open Source Components

**Target Platform:** Kubernetes 1.34
**Document Version:** 1.1
**Date:** January 2026
**Last Updated:** January 15, 2026

---

## Table of Contents

1. [Executive Summary](#executive-summary) - **Includes Minimum Stack (4 Components)**
2. [Understanding CNCF AI Conformance](#understanding-cncf-ai-conformance)
3. [Conformance Requirements Analysis](#conformance-requirements-analysis)
4. [Provider Implementation Review](#provider-implementation-review)
5. [Open Source Component Stack](#open-source-component-stack) - **Minimum vs Extended Options**
6. [Detailed Component Analysis](#detailed-component-analysis)
7. [Recommended Implementation Architecture](#recommended-implementation-architecture)
8. [Helm-Based Deployment Strategy](#helm-based-deployment-strategy)
9. [Testing and Validation](#testing-and-validation)
10. [Conclusion and Recommendations](#conclusion-and-recommendations) - **Complete Installation Script**

---

## 1. Executive Summary

The CNCF Kubernetes AI Conformance program establishes a standardized set of requirements that Kubernetes platforms must meet to effectively support AI/ML workloads. This document provides a thorough analysis of the v1.34 conformance specification, examines how 12 major cloud providers have implemented these requirements, and presents a comprehensive strategy for achieving conformance using exclusively open source, Helm-deployable components.

### Key Findings

- **9 MUST-level requirements** across 6 categories must be satisfied for conformance
- All requirements can be met using **100% open source components**
- **Helm charts are available** for all recommended components
- **Only 4 Helm releases needed** to satisfy all requirements (see Minimum Stack below)
- **NVIDIA Dynamo Platform v0.8.0**: Dynamo Operator + optional Grove/KAI Scheduler (enabled via `--set` flags)
- **NVIDIA GPU Operator v25.10.1**: drivers, device plugin, DCGM Exporter, container toolkit
- **DRA support requires separate nvidia-dra-driver-gpu chart** (not bundled in GPU Operator)

### Minimum Stack (4 Components)

| Component | Helm Chart | Requirements Satisfied |
|-----------|------------|----------------------|
| **NVIDIA GPU Operator** | `nvidia/gpu-operator` | `dra_support`, `secure_accelerator_access`, `accelerator_metrics`, `pod_autoscaling`** |
| **Dynamo Platform** | `nvidia/ai-dynamo/dynamo-platform` | `robust_controller`, `ai_inference`, `gang_scheduling`* |
| **kube-prometheus-stack** | `prometheus-community/kube-prometheus-stack` | `ai_service_metrics` |
| **Cluster Autoscaler** | `autoscaler/cluster-autoscaler` | `cluster_autoscaling` |

*With `--set grove.enabled=true --set kai-scheduler.enabled=true`
**`pod_autoscaling` uses standard Kubernetes HPA with DCGM Exporter metrics from GPU Operator

> **⚠️ Note:** Cluster Autoscaler is a **separate installation** - it is NOT bundled with Dynamo Platform.

### Quick Install
```bash
export DYNAMO_VERSION=0.8.0  # Latest as of Jan 2026

# GPU Operator v25.10.1
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia && helm repo update
helm install gpu-operator nvidia/gpu-operator -n gpu-operator --create-namespace --version=v25.10.1

# Dynamo Platform v0.8.0 (with Grove + KAI Scheduler for gang scheduling)
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-crds-${DYNAMO_VERSION}.tgz
helm install dynamo-crds dynamo-crds-${DYNAMO_VERSION}.tgz -n default
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-platform-${DYNAMO_VERSION}.tgz
helm install dynamo-platform dynamo-platform-${DYNAMO_VERSION}.tgz -n dynamo-system --create-namespace \
  --set grove.enabled=true --set kai-scheduler.enabled=true

# Prometheus Stack (Chart 80.14.4)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Cluster Autoscaler (App 1.34.2)
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm install cluster-autoscaler autoscaler/cluster-autoscaler -n kube-system \
  --set autoDiscovery.clusterName=${CLUSTER_NAME} --set cloudProvider=aws
```

---

## 2. Understanding CNCF AI Conformance

### What is AI Conformance?

The CNCF AI Conformance program, managed by the [WG AI Conformance Working Group](https://github.com/kubernetes-sigs/wg-ai-conformance), defines a set of requirements that Kubernetes platforms must satisfy to be certified as AI-ready. The program ensures that:

- AI/ML workloads can run reliably across conformant platforms
- Portability is maintained between different Kubernetes distributions
- Best practices for GPU/accelerator management are standardized
- Observability and security requirements are consistently implemented

### Conformance Levels

Requirements are categorized into two levels:
- **MUST**: Mandatory for conformance certification
- **SHOULD**: Recommended but not required (graduating to MUST in future versions)

### Version Alignment

The AI conformance version aligns with Kubernetes releases. Version 1.34 conformance:
- Requires Kubernetes v1.34 as the base platform
- Assumes DRA (Dynamic Resource Allocation) APIs are available (GA in 1.32+)
- Targets the v1 stable APIs for core features

---

## 3. Conformance Requirements Analysis

### Complete Requirements Matrix (v1.34)

| Category | ID | Requirement | Level | Complexity |
|----------|-----|-------------|-------|------------|
| **Accelerators** | dra_support | Dynamic Resource Allocation APIs | MUST | High |
| **Networking** | ai_inference | Gateway API with advanced traffic management | MUST | Medium |
| **Scheduling** | gang_scheduling | All-or-nothing scheduling (Kueue, Volcano) | MUST | Medium |
| **Scheduling** | cluster_autoscaling | Autoscale nodes with accelerators | MUST | Medium |
| **Scheduling** | pod_autoscaling | HPA with custom metrics for accelerators | MUST | Medium |
| **Observability** | accelerator_metrics | Fine-grained GPU/accelerator metrics | MUST | Medium |
| **Observability** | ai_service_metrics | Prometheus-format metrics collection | MUST | Low |
| **Security** | secure_accelerator_access | Isolated accelerator access via device plugin/DRA | MUST | High |
| **Operator** | robust_controller | AI operator (Ray, Kubeflow) with CRD reconciliation | MUST | Medium |

### Detailed Requirement Descriptions

#### 3.1 DRA Support (dra_support)
> "Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts."

**Technical Requirements:**
- DRA feature gate enabled (GA in Kubernetes 1.32+)
- DRA driver installed for accelerator type (NVIDIA, AMD, Intel)
- Support for ResourceClaims and ResourceClaimTemplates
- Proper device plugin or DRA driver integration

#### 3.2 AI Inference Networking (ai_inference)
> "Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes."

**Technical Requirements:**
- Gateway API CRDs installed (v1.0+)
- GatewayClass and Gateway resources functional
- HTTPRoute support with:
  - Weighted traffic splitting
  - Header-based routing
  - Request/response transformation (for OpenAI protocol headers)

#### 3.3 Gang Scheduling (gang_scheduling)
> "The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads."

**Technical Requirements:**
- Support for all-or-nothing pod scheduling
- Queue management for pending workloads
- Integration with native Kubernetes scheduler or replacement scheduler
- Support for common AI workload patterns (RayJob, PyTorchJob, etc.)

#### 3.4 Cluster Autoscaling (cluster_autoscaling)
> "If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators."

**Technical Requirements:**
- Automatic node provisioning based on pending pod requests
- Recognition of accelerator resource requests (nvidia.com/gpu, etc.)
- Support for heterogeneous node pools (different GPU types)
- Scale-down when resources are idle

#### 3.5 Pod Autoscaling (pod_autoscaling)
> "If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads."

**Technical Requirements:**
- Standard HPA functionality with accelerator workloads
- Custom metrics API implementation (metrics.k8s.io)
- External metrics support for GPU utilization metrics
- Integration with metrics adapters (Prometheus Adapter, KEDA)

**Note:** This requirement is satisfied by standard Kubernetes HPA combined with GPU metrics from DCGM Exporter (included in GPU Operator). KAI Scheduler provides GPU-aware scheduling for autoscaled workloads but does not replace HPA.

#### 3.6 Accelerator Metrics (accelerator_metrics)
> "For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint."

**Required Metrics:**
- Per-accelerator utilization (%)
- Memory usage (bytes/percentage)
- Temperature (if available)
- Power draw (if available)
- Interconnect bandwidth (if available - NVLink, etc.)

#### 3.7 AI Service Metrics (ai_service_metrics)
> "Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format)."

**Technical Requirements:**
- Prometheus-compatible metrics endpoint
- Service discovery for metric endpoints
- Support for PodMonitor/ServiceMonitor CRDs
- Scalable metrics collection and storage

#### 3.8 Secure Accelerator Access (secure_accelerator_access)
> "Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads."

**Technical Requirements:**
- Pods without GPU requests cannot access GPU devices
- GPU isolation between pods (each pod sees only allocated GPUs)
- Proper device file isolation (/dev/nvidia* not exposed without request)
- Container runtime enforcement (runtimeClass for NVIDIA)

#### 3.9 Robust Controller (robust_controller)
> "The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably."

**Technical Requirements:**
- AI operator pods running correctly
- Operator webhooks operational (validating/mutating)
- Custom resources reconcile successfully
- Workload lifecycle management functional

---

## 4. Provider Implementation Review

### Providers Analyzed

| Provider | Platform | Kubernetes Version | Notable Approach |
|----------|----------|-------------------|------------------|
| **GKE** | Google Kubernetes Engine | 1.34.0-gke.1662000 | Native DRA, Kueue |
| **EKS** | Amazon Elastic Kubernetes Service | 1.34.1-eks.4 | Karpenter, Volcano options |
| **AKS** | Azure Kubernetes Service | v1.34 | KEDA, Kueue |
| **Kubermatic** | Kubermatic Kubernetes Platform | v2.29 | KubeLB, Kubeflow |
| **Talos** | Talos Linux (Sidero Labs) | v1.11.3 | Minimal, KubeRay |
| **ACK** | Alibaba Container Service | 1.34.1-aliyun.1 | Native AI extensions |
| **OKE** | Oracle Kubernetes Engine | v1.34 | Istio Gateway API |
| **LKE** | Linode Kubernetes Engine | v1.34 | Prometheus Operator |
| **VKS** | vSphere Kubernetes Service | v3.5.0 | Istio, on-prem focus |
| **CCE** | Baidu Cloud Container Engine | 1.34 | Native schedulers |
| **CKS** | CoreWeave Kubernetes Service | v1.34 | K-Gateway, SUNK |

### Common Implementation Patterns

#### Gang Scheduling Solutions Adopted

| Solution | Providers Using It |
|----------|-------------------|
| **Kueue** | GKE, AKS, Kubermatic, Talos, OKE, VKS, CCE, CKS, LKE |
| **Volcano** | EKS, ACK |
| **Yunikorn** | EKS |
| **AWS Batch** | EKS |
| **Custom** | CoreWeave (SUNK/Slurm) |

**Key Insight:** Kueue has become the de facto standard for gang scheduling in Kubernetes AI workloads, with 9 out of 12 providers listing it as their primary or only solution.

#### Gateway API Implementations

| Implementation | Providers Using It |
|----------------|-------------------|
| **Istio** | OKE, VKS, Kubermatic |
| **Traefik** | Talos |
| **Cilium** | Talos (alternative) |
| **AWS Load Balancer Controller** | EKS |
| **Cloud-native** | GKE, AKS, ACK |
| **K-Gateway** | CKS |

#### AI Operators Validated

| Operator | Providers |
|----------|-----------|
| **KubeRay** | GKE, AKS, ACK, Talos, OKE, EKS |
| **Kubeflow** | GKE, Kubermatic, OKE, CKS |
| **vLLM** | EKS, CKS |
| **NVIDIA Triton** | EKS |
| **AIBrix** | EKS |

---

## 5. Open Source Component Stack

### Minimum vs Extended Component Selection

Based on analysis of all 12 provider implementations and the NVIDIA open source stack, there are two approaches:

#### Option A: Minimum Stack (Recommended)

Using NVIDIA's integrated platform which bundles multiple capabilities:

| Requirement | Component | Included In | Helm Chart |
|-------------|-----------|-------------|------------|
| **DRA Support** | Device Plugin + DRA Driver | GPU Operator | `nvidia/gpu-operator` |
| **Gateway API** | GAIE Integration + Frontend | Dynamo Platform | `nvidia/ai-dynamo/dynamo-platform` |
| **Gang Scheduling** | KAI Scheduler | Dynamo Platform* | `--set kai-scheduler.enabled=true` |
| **Cluster Autoscaling** | Cluster Autoscaler | **Separate** (NOT in Dynamo) | `autoscaler/cluster-autoscaler` |
| **Pod Autoscaling** | HPA + DCGM metrics | Kubernetes + GPU Operator | Prometheus Adapter optional |
| **Accelerator Metrics** | DCGM Exporter | GPU Operator | (included) |
| **AI Service Metrics** | Prometheus | kube-prometheus-stack | `prometheus-community/kube-prometheus-stack` |
| **Secure Access** | Device Plugin | GPU Operator | (included) |
| **AI Operator** | Dynamo Operator | Dynamo Platform | (included) |

*Requires `--set grove.enabled=true --set kai-scheduler.enabled=true`

**Total: 4 Helm releases**

> **⚠️ Important:** Cluster Autoscaler is NOT bundled with Dynamo Platform. It must be installed separately using `helm install autoscaler/cluster-autoscaler` or enabled via cloud provider native autoscaling (GKE Autopilot, AKS auto-scale, EKS managed node groups).

#### Option B: Extended Stack (À la carte)

For environments that prefer individual components:

| Requirement | Primary Component | Alternative | Helm Chart |
|-------------|------------------|-------------|------------|
| **DRA Support** | NVIDIA GPU Operator | AMD/Intel equivalents | ✅ |
| **Gateway API** | Envoy Gateway | Istio, Traefik, Cilium | ✅ |
| **Gang Scheduling** | Kueue | Volcano, KAI Scheduler | ✅ |
| **Cluster Autoscaling** | Karpenter / Cluster Autoscaler | N/A | ✅ |
| **Pod Autoscaling** | KEDA + Prometheus Adapter | metrics-server + HPA | ✅ |
| **Accelerator Metrics** | NVIDIA DCGM Exporter | GPU Operator built-in | ✅ |
| **AI Service Metrics** | kube-prometheus-stack | OpenTelemetry Collector | ✅ |
| **Secure Access** | NVIDIA GPU Operator | Device plugin only | ✅ |
| **AI Operator** | KubeRay | Kubeflow Trainer, Dynamo | ✅ |

**Total: 8-10 Helm releases**

### What's Included in Each NVIDIA Component

#### NVIDIA GPU Operator
The GPU Operator is a single Helm release that manages:
- **NVIDIA drivers** - Automatic installation and updates
- **Container toolkit** - nvidia-container-runtime configuration
- **Device plugin** - GPU allocation and isolation (`secure_accelerator_access`)
- **DCGM Exporter** - GPU metrics for Prometheus (`accelerator_metrics`)
- **GPU Feature Discovery** - Node labeling with GPU capabilities
- **MIG Manager** - Multi-Instance GPU partitioning (optional)
- **MPS Control Daemon** - GPU sharing via MPS (optional)

#### NVIDIA Dynamo Platform
The Dynamo Platform is a comprehensive AI inference platform that includes:
- **Dynamo Operator** - Manages `DynamoGraphDeployment` CRDs (`robust_controller`)
- **OpenAI-compatible Frontend** - High-performance HTTP API server
- **Service Discovery** - Kubernetes-native (no external etcd/NATS required)
- **Grove** (optional) - Disaggregated inference orchestration with `PodCliqueSet`
- **KAI Scheduler** (optional) - GPU-native scheduling with fractional GPUs, gang scheduling
- **Inference Gateway Integration** - Gateway API via GAIE (`ai_inference`)
- **Observability** - Prometheus configs and Grafana dashboards

---

## 6. Detailed Component Analysis

### 6.1 NVIDIA GPU Operator + DRA Driver

**Purpose:** Unified GPU management including drivers, container toolkit, device plugin. DRA support is a separate component.

**Latest Version:** v25.10.1 (December 2024)

**Bundled Components:**
- NVIDIA GPU Drivers (v580.105.08)
- Kubernetes Device Plugin for GPUs
- NVIDIA Container Toolkit (v1.18.1)
- GPU Feature Discovery (GFD)
- DCGM (v4.4.2-1) and DCGM Exporter
- MIG Manager
- GDS (GPUDirect Storage) Driver
- Container Device Interface (CDI) - enabled by default

**Key Features:**
- Automated driver installation and updates
- Container runtime configuration (nvidia-container-runtime)
- Device plugin deployment
- MIG (Multi-Instance GPU) support
- GPU sharing via time-slicing

**Helm Installation:**
```bash
# Add NVIDIA Helm repo
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update

# Install GPU Operator v25.10.1
helm install gpu-operator nvidia/gpu-operator \
  --namespace gpu-operator \
  --create-namespace \
  --version v25.10.1

# For DRA support (SEPARATE chart - required for dra_support conformance)
# Note: GPU allocation via DRA is Technology Preview; ComputeDomains is fully supported
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu \
  --version 25.3.0 \
  --namespace nvidia-dra-driver-gpu \
  --create-namespace \
  --set resources.gpus.enabled=false \
  --set nvidiaDriverRoot=/run/nvidia/driver
```

**Configuration Highlights:**
```yaml
# values.yaml for GPU Operator
driver:
  enabled: true
  version: "570.76.02"
toolkit:
  enabled: true
devicePlugin:
  enabled: true
dcgmExporter:
  enabled: true
mig:
  strategy: single  # or mixed
```

**Conformance Coverage:** `dra_support`, `secure_accelerator_access`, `accelerator_metrics` (partial)

---

### 6.2 Envoy Gateway

**Purpose:** Kubernetes-native API Gateway implementing Gateway API with advanced traffic management.

**Key Features:**
- Full Gateway API support (HTTPRoute, GRPCRoute, TCPRoute)
- Weighted traffic splitting for canary/blue-green deployments
- Header-based routing (critical for OpenAI protocol headers)
- Request/response transformation
- Rate limiting and authentication
- Native Envoy proxy performance

**Helm Installation:**
```bash
# Install Envoy Gateway
helm install eg oci://docker.io/envoyproxy/gateway-helm \
  --version v1.4.6 \
  -n envoy-gateway-system \
  --create-namespace
```

**Example Gateway Configuration:**
```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: ai-inference-gateway
  namespace: inference
spec:
  gatewayClassName: eg
  listeners:
  - name: http
    port: 80
    protocol: HTTP
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: vllm-route
spec:
  parentRefs:
  - name: ai-inference-gateway
  rules:
  - matches:
    - headers:
      - name: x-model-id
        value: llama-3-70b
    backendRefs:
    - name: vllm-llama-70b
      port: 8000
      weight: 90
    - name: vllm-llama-70b-canary
      port: 8000
      weight: 10
```

**Alternatives Comparison:**

| Feature | Envoy Gateway | Istio | Traefik | Cilium |
|---------|---------------|-------|---------|--------|
| Gateway API | ✅ Full | ✅ Full | ✅ Full | ✅ Full |
| Complexity | Medium | High | Low | Medium |
| Performance | High | High | Medium | High |
| Resource Usage | Medium | High | Low | High |
| Service Mesh | Optional | Built-in | No | Built-in |

**Conformance Coverage:** `ai_inference`

---

### 6.3 Kueue

**Purpose:** Kubernetes-native job queueing system for batch, HPC, and AI/ML workloads.

**Key Features:**
- Gang scheduling (all-or-nothing admission)
- Multi-tenant quota management with ClusterQueues
- Priority-based preemption
- Borrowing between queues (Cohorts)
- Integration with JobSet, RayJob, PyTorchJob
- Topology-aware scheduling

**Helm Installation:**
```bash
# Using kubectl (recommended by upstream)
VERSION=v0.14.2
kubectl apply --server-side -f \
  https://github.com/kubernetes-sigs/kueue/releases/download/$VERSION/manifests.yaml

# Or via Helm (community charts available)
helm repo add kueue https://kubernetes-sigs.github.io/kueue/charts
helm install kueue kueue/kueue -n kueue-system --create-namespace
```

**Example Configuration:**
```yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: ai-training
spec:
  namespaceSelector: {}
  resourceGroups:
  - coveredResources: ["cpu", "memory", "nvidia.com/gpu"]
    flavors:
    - name: gpu-a100
      resources:
      - name: "nvidia.com/gpu"
        nominalQuota: 8
      - name: "cpu"
        nominalQuota: 64
      - name: "memory"
        nominalQuota: 512Gi
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  name: team-ml
  namespace: ml-workloads
spec:
  clusterQueue: ai-training
```

**Volcano Comparison:**

| Feature | Kueue | Volcano |
|---------|-------|---------|
| Kubernetes Native | ✅ Uses native scheduler | ❌ Custom scheduler |
| Complexity | Lower | Higher |
| Job Types | JobSet, RayJob, PyTorchJob | VcJob, custom types |
| DRA Support | Emerging | Yes (with feature gate) |
| CNCF Status | kubernetes-sigs | CNCF Incubating |
| Community Adoption | Growing rapidly | Mature |

**Conformance Coverage:** `gang_scheduling`

---

### 6.4 Cluster Autoscaling

#### Option A: Karpenter (Cloud environments)

**Purpose:** High-performance node autoscaler with fast provisioning and intelligent bin-packing.

**Key Features:**
- Sub-minute node provisioning
- Intelligent instance selection based on workload requirements
- GPU/accelerator-aware scheduling
- Spot instance support with interruption handling
- Consolidation and defragmentation

**Helm Installation (AWS):**
```bash
helm upgrade --install karpenter oci://public.ecr.aws/karpenter/karpenter \
  --version "${KARPENTER_VERSION}" \
  --namespace "${KARPENTER_NAMESPACE}" \
  --set "settings.clusterName=${CLUSTER_NAME}" \
  --set "settings.interruptionQueue=${CLUSTER_NAME}" \
  --set "serviceAccount.annotations.eks\.amazonaws\.com/role-arn=arn:aws:iam::${AWS_ACCOUNT_ID}:role/KarpenterControllerRole-${CLUSTER_NAME}"
```

#### Option B: Cluster Autoscaler (Multi-cloud/On-premise)

**Purpose:** Standard Kubernetes cluster autoscaler supporting multiple cloud providers and Cluster API.

**Helm Installation:**
```bash
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm install cluster-autoscaler autoscaler/cluster-autoscaler \
  -n kube-system \
  --set autoDiscovery.clusterName=${CLUSTER_NAME} \
  --set cloudProvider=aws  # or gce, azure, clusterapi
```

**Conformance Coverage:** `cluster_autoscaling`

---

### 6.5 Pod Autoscaling with KEDA

**Purpose:** Event-driven autoscaling including scale-to-zero capability for GPU workloads.

**Key Features:**
- Scale based on custom metrics (GPU utilization, queue depth)
- Scale to zero when idle (cost savings for GPU workloads)
- Integration with Prometheus metrics
- 60+ built-in scalers (Kafka, RabbitMQ, etc.)

**Helm Installation:**
```bash
helm repo add kedacore https://kedacore.github.io/charts
helm repo update
helm install keda kedacore/keda \
  --namespace keda \
  --create-namespace
```

**Example ScaledObject for GPU Workloads:**
```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: vllm-scaledobject
spec:
  scaleTargetRef:
    name: vllm-deployment
  minReplicaCount: 0
  maxReplicaCount: 10
  triggers:
  - type: prometheus
    metadata:
      serverAddress: http://prometheus.monitoring:9090
      metricName: DCGM_FI_DEV_GPU_UTIL
      threshold: '80'
      query: |
        avg(DCGM_FI_DEV_GPU_UTIL{pod=~"vllm-.*"})
```

**Supporting Components:**

```bash
# metrics-server (required for basic HPA)
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm install metrics-server metrics-server/metrics-server -n kube-system

# Prometheus Adapter (for custom metrics)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus-adapter prometheus-community/prometheus-adapter \
  -n monitoring \
  --set prometheus.url=http://prometheus.monitoring.svc
```

**Conformance Coverage:** `pod_autoscaling`

---

### 6.6 Observability Stack

#### NVIDIA DCGM Exporter

**Purpose:** GPU metrics exporter for Prometheus.

**Metrics Exposed:**
- `DCGM_FI_DEV_GPU_UTIL` - GPU utilization (%)
- `DCGM_FI_DEV_MEM_COPY_UTIL` - Memory utilization (%)
- `DCGM_FI_DEV_GPU_TEMP` - Temperature (°C)
- `DCGM_FI_DEV_POWER_USAGE` - Power (W)
- `DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL` - NVLink bandwidth

**Helm Installation:**
```bash
helm repo add gpu-helm-charts https://nvidia.github.io/dcgm-exporter/helm-charts
helm repo update
helm install dcgm-exporter gpu-helm-charts/dcgm-exporter \
  -n gpu-operator \
  --set serviceMonitor.enabled=true
```

**Note:** If using NVIDIA GPU Operator, DCGM Exporter is included and managed automatically.

#### kube-prometheus-stack

**Purpose:** Complete monitoring stack with Prometheus, Alertmanager, and Grafana.

**Helm Installation:**
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false
```

#### OpenTelemetry Collector (Alternative/Complementary)

**Purpose:** Vendor-neutral telemetry collection supporting metrics, logs, and traces.

**Helm Installation:**
```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm install otel-collector open-telemetry/opentelemetry-collector \
  --namespace monitoring \
  --set mode=daemonset \
  --set presets.hostMetrics.enabled=true \
  --set presets.kubernetesAttributes.enabled=true
```

**Conformance Coverage:** `accelerator_metrics`, `ai_service_metrics`

---

### 6.7 KubeRay (AI Operator)

**Purpose:** Kubernetes operator for Ray clusters, enabling distributed AI/ML workloads.

**Key Features:**
- RayCluster management
- RayJob for batch inference/training
- RayService for serving with autoscaling
- Integration with Kueue for gang scheduling
- Zero-downtime upgrades

**Helm Installation:**
```bash
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update

# Install KubeRay Operator
helm install kuberay-operator kuberay/kuberay-operator \
  --namespace ray-system \
  --create-namespace \
  --version 1.5.1

# Deploy a RayCluster
helm install raycluster kuberay/ray-cluster \
  --namespace ray-workloads \
  --create-namespace
```

**Example RayCluster with GPU:**
```yaml
apiVersion: ray.io/v1
kind: RayCluster
metadata:
  name: gpu-cluster
spec:
  rayVersion: '2.53.0'
  headGroupSpec:
    rayStartParams:
      dashboard-host: '0.0.0.0'
    template:
      spec:
        containers:
        - name: ray-head
          image: rayproject/ray-ml:2.53.0-gpu
          resources:
            limits:
              cpu: "4"
              memory: "16Gi"
  workerGroupSpecs:
  - groupName: gpu-workers
    replicas: 2
    template:
      spec:
        containers:
        - name: ray-worker
          image: rayproject/ray-ml:2.53.0-gpu
          resources:
            limits:
              cpu: "8"
              memory: "32Gi"
              nvidia.com/gpu: "1"
        runtimeClassName: nvidia
```

**Alternative: Kubeflow Trainer**

```bash
helm install kubeflow-trainer oci://ghcr.io/kubeflow/charts/kubeflow-trainer \
  --namespace kubeflow-system \
  --create-namespace
```

**Conformance Coverage:** `robust_controller`

---

### 6.8 NVIDIA Advanced AI Infrastructure (Open Source)

NVIDIA has open-sourced several advanced components that significantly enhance Kubernetes AI capabilities beyond the basic conformance requirements. These components represent production-grade solutions used internally and with enterprise customers.

#### 6.8.1 NVIDIA Dynamo - Distributed Inference Framework

**Purpose:** High-throughput, low-latency inference serving framework for generative AI and reasoning models at datacenter scale.

**Latest Version:** v0.8.0 (January 15, 2025)

**GitHub:** [ai-dynamo/dynamo](https://github.com/ai-dynamo/dynamo)

**Key Features:**
- **Disaggregated Prefill/Decode**: Separates inference stages for optimal GPU utilization
- **Dynamic GPU Scheduling**: Intelligent resource allocation based on workload demands
- **LLM-Aware Request Routing**: Smart routing based on model characteristics and load
- **Multi-Node NVLink Support**: Optimized for NVIDIA GB200 and similar systems
- **Backend Agnostic**: Supports vLLM, TensorRT-LLM, SGLang, and PyTorch
- **Kubernetes-Native Service Discovery**: Uses EndpointSlices (no external etcd/NATS required)
- **Multimodal Support**: Audio and video inference (new in v0.8.0)
- **Agentic Workflows**: Extended tool calling support

**Performance:**
- Up to **30x throughput increase** on DeepSeek-R1 671B with GB200 NVL72
- **2x+ throughput improvement** for Llama 70B on Hopper GPUs

**What's Bundled vs Optional:**
- **Bundled**: Dynamo Operator, NATS, etcd
- **Optional (must enable)**: Grove (`--set grove.enabled=true`), KAI Scheduler (`--set kai-scheduler.enabled=true`)

**CRDs Provided:**
- `DynamoGraphDeployment (DGD)` - Core deployment resource
- `DynamoGraphDeploymentRequest (DGDR)` - SLA-driven deployments with auto-profiling
- `DynamoComponentDeployment` - Individual component specs
- `DynamoGraphDeploymentScalingAdapter (DGDSA)` - HPA/KEDA integration
- `DynamoModel` - Model references and endpoint discovery

**Installation:**
```bash
export DYNAMO_VERSION=0.8.0

# Install CRDs
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-crds-${DYNAMO_VERSION}.tgz
helm install dynamo-crds dynamo-crds-${DYNAMO_VERSION}.tgz --namespace default

# Install Platform (basic - single node)
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-platform-${DYNAMO_VERSION}.tgz
helm install dynamo-platform dynamo-platform-${DYNAMO_VERSION}.tgz \
  --namespace dynamo-system \
  --create-namespace

# For multi-node with Grove + KAI Scheduler
helm install dynamo-platform dynamo-platform-${DYNAMO_VERSION}.tgz \
  --namespace dynamo-system \
  --create-namespace \
  --set grove.enabled=true \
  --set kai-scheduler.enabled=true
```

**Example Dynamo Service:**
```yaml
apiVersion: dynamo.nvidia.com/v1
kind: DynamoService
metadata:
  name: llama-70b-service
spec:
  model: meta-llama/Llama-2-70b-chat-hf
  replicas: 2
  prefillWorkers: 4
  decodeWorkers: 8
  routingStrategy: load-aware
  resources:
    gpu: nvidia.com/gpu
    gpuCount: 8
```

**Conformance Relevance:** Enhances `ai_inference`, `robust_controller`

---

#### 6.8.2 NVIDIA KAI Scheduler - GPU-Native Kubernetes Scheduler

**Purpose:** Kubernetes-native GPU scheduling solution designed for large-scale AI workloads, open-sourced from Run:ai platform.

**Latest Version:** v0.12.6 (January 14, 2025)

**GitHub:** [NVIDIA/KAI-Scheduler](https://github.com/NVIDIA/KAI-Scheduler)

**Key Features:**
- **GPU Sharing**: Multiple workloads efficiently share GPUs (time-slicing, MPS)
- **Fractional GPU Allocation**: Request portions of GPUs via `gpu-fraction` annotation
- **Gang/Batch Scheduling**: All-or-nothing scheduling with bin-packing and spread options
- **Hierarchical Queues**: Two-level queue hierarchy for organizational control
- **DRA Support**: Native Dynamic Resource Allocation integration
- **Topology-Aware Scheduling**: Optimizes placement based on NVLink, GPU interconnects (block, rack, hostname levels)
- **Fair-Sharing**: Time-aware resource distribution across tenants

**Architecture:**
- Runs alongside default Kubernetes scheduler (does not replace it)
- Captures cluster state, computes optimal distribution
- Applies scheduling actions: allocate, consolidate, reclaim, preempt

**Deployment Options:**
1. **Standalone**: Install independently via Helm
2. **With Dynamo**: Enable via `--set kai-scheduler.enabled=true`

**Standalone Installation:**
```bash
# Install via Helm (standalone)
helm upgrade -i kai-scheduler \
  oci://ghcr.io/nvidia/kai-scheduler/kai-scheduler \
  -n kai-scheduler \
  --create-namespace \
  --version v0.12.6

# For GPU sharing support
helm upgrade -i kai-scheduler \
  oci://ghcr.io/nvidia/kai-scheduler/kai-scheduler \
  -n kai-scheduler \
  --create-namespace \
  --version v0.12.6 \
  --set "global.gpuSharing=true"
```

**Example GPU Sharing Configuration:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: shared-gpu-pod
  annotations:
    kai.scheduler.nvidia.com/gpu-fraction: "0.5"  # Request 50% of GPU
spec:
  schedulerName: kai-scheduler
  containers:
  - name: inference
    image: nvidia/cuda:12.0-runtime
    resources:
      limits:
        nvidia.com/gpu: "1"
```

**Comparison with Other Schedulers:**

| Feature | KAI Scheduler | Kueue | Volcano |
|---------|---------------|-------|---------|
| GPU Sharing | ✅ Native | ❌ Limited | ✅ Via plugins |
| Fractional GPUs | ✅ Yes | ❌ No | ✅ Yes |
| Gang Scheduling | ✅ Yes | ✅ Yes | ✅ Yes |
| DRA Support | ✅ Yes | ✅ Emerging | ✅ Yes |
| Topology-Aware | ✅ Advanced | ✅ Basic | ✅ Basic |
| Multi-tenant | ✅ Hierarchical | ✅ Cohorts | ✅ Queues |
| Run alongside default | ✅ Yes | ✅ Yes | ❌ Replaces |

**Important Note:** KAI Scheduler does not enforce GPU memory isolation. Applications must respect their allocated memory limits.

**Conformance Relevance:** Enhances `gang_scheduling`, `cluster_autoscaling`, `pod_autoscaling`

---

#### 6.8.3 NVIDIA Grove - Kubernetes AI Orchestration API

**Purpose:** Single declarative interface for orchestrating complex AI inference workloads from single-pod to multi-node, datacenter-scale deployments.

**Latest Version:** v0.1.0-alpha.3 (October 2024) - Active development

**GitHub:** [NVIDIA/grove](https://github.com/NVIDIA/grove)

**Key Features:**
- **Unified Custom Resource**: Describe entire inference systems (prefill, decode, routing) as single CR
- **PodCliqueScalingGroups**: Bundle and scale tightly coupled pods as single entity
- **Network Topology-Aware Placement**: Gang scheduling with GPU cluster topology optimization
- **Hierarchical Autoscaling**: Scale from modules to full service replicas
- **Startup Ordering**: Explicit coordination of multi-component AI stacks
- **Scales to tens of thousands of GPUs**

**How It Works:**
1. Define inference system in single Grove CR
2. Grove scheduler applies gang-scheduling logic
3. Ensures all components scheduled together or not at all
4. Optimizes placement based on GPU topology (NVLink, interconnects)
5. Coordinates startup order (prefill → decode → routing)

**Installation:**
```bash
# Clone Grove repository
git clone https://github.com/NVIDIA/grove.git
cd grove

# Install Grove CRDs and operator
kubectl apply -f deploy/crds/
helm install grove-operator ./charts/grove-operator \
  --namespace grove-system \
  --create-namespace
```

**Example Grove Deployment:**
```yaml
apiVersion: grove.nvidia.com/v1alpha1
kind: GroveService
metadata:
  name: disaggregated-llm
spec:
  components:
    - name: prefill
      replicas: 4
      template:
        spec:
          containers:
          - name: prefill-worker
            image: nvidia/dynamo-prefill:latest
            resources:
              limits:
                nvidia.com/gpu: "2"
    - name: decode
      replicas: 8
      template:
        spec:
          containers:
          - name: decode-worker
            image: nvidia/dynamo-decode:latest
            resources:
              limits:
                nvidia.com/gpu: "1"
    - name: router
      replicas: 2
      template:
        spec:
          containers:
          - name: router
            image: nvidia/dynamo-router:latest
  startupOrder:
    - prefill
    - decode
    - router
  scalingPolicy:
    minReplicas: 1
    maxReplicas: 100
    metrics:
    - type: Resource
      resource:
        name: nvidia.com/gpu
        target:
          utilization: 80
```

**Conformance Relevance:** Enhances `gang_scheduling`, `ai_inference`, `robust_controller`, `cluster_autoscaling`

---

### 6.9 Advanced Stack Comparison: When to Use What

| Use Case | Recommended Components |
|----------|----------------------|
| **Basic AI Conformance** | GPU Operator + Kueue + KubeRay + Envoy Gateway |
| **High-throughput LLM Inference** | GPU Operator + Dynamo + Grove + KAI Scheduler |
| **Multi-tenant GPU Sharing** | GPU Operator + KAI Scheduler + Kueue |
| **Distributed Training** | GPU Operator + KubeRay + Kueue |
| **Disaggregated Inference (Prefill/Decode)** | GPU Operator + Dynamo + Grove |
| **Cost-optimized Batch Processing** | GPU Operator + Kueue + KEDA |

### Integration Example: Full NVIDIA Stack

```bash
#!/bin/bash
# Advanced NVIDIA AI Stack Installation

# 1. GPU Operator (foundation)
helm install gpu-operator nvidia/gpu-operator \
  --namespace gpu-operator --create-namespace

# 2. KAI Scheduler (advanced GPU scheduling)
helm upgrade -i kai-scheduler \
  oci://ghcr.io/nvidia/kai-scheduler/kai-scheduler \
  -n kai-scheduler --create-namespace

# 3. Grove (orchestration API)
kubectl apply -f https://raw.githubusercontent.com/NVIDIA/grove/main/deploy/crds/
helm install grove-operator grove/grove-operator \
  -n grove-system --create-namespace

# 4. Dynamo Operator (inference framework)
helm install dynamo-operator dynamo/dynamo-operator \
  -n dynamo-system --create-namespace

# 5. Supporting infrastructure
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace
```

---

## 7. Recommended Implementation Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          KUBERNETES CLUSTER (v1.34)                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  Control Plane  │  │   GPU Nodes     │  │  CPU Nodes      │             │
│  │                 │  │                 │  │                 │             │
│  │  - API Server   │  │  - nvidia.com/  │  │  - Standard     │             │
│  │  - Controller   │  │    gpu: N       │  │    workloads    │             │
│  │  - Scheduler    │  │  - DRA Driver   │  │                 │             │
│  │  - etcd         │  │  - GPU Operator │  │                 │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                           PLATFORM LAYER                                    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    GPU Management (gpu-operator namespace)           │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────────┐ │   │
│  │  │GPU Operator  │ │Device Plugin │ │DCGM Exporter │ │DRA Driver   │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Networking (envoy-gateway-system namespace)       │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────────────┐ │   │
│  │  │Envoy Gateway │ │Gateway API   │ │HTTPRoutes / GRPCRoutes       │ │   │
│  │  │Controller    │ │CRDs          │ │(AI Inference routing)        │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Scheduling (kueue-system namespace)               │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────────────┐ │   │
│  │  │Kueue         │ │ClusterQueues │ │LocalQueues                   │ │   │
│  │  │Controller    │ │              │ │(per namespace/team)          │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Autoscaling                                       │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────────┐ │   │
│  │  │Karpenter/CA  │ │KEDA          │ │metrics-server│ │Prom Adapter │ │   │
│  │  │(cluster)     │ │(pod scaling) │ │              │ │             │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Observability (monitoring namespace)              │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────────┐ │   │
│  │  │Prometheus    │ │Alertmanager  │ │Grafana       │ │ServiceMon.  │ │   │
│  │  │              │ │              │ │              │ │PodMonitors  │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    AI Operators (ray-system namespace)               │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────────────┐ │   │
│  │  │KubeRay       │ │RayCluster    │ │RayJob / RayService           │ │   │
│  │  │Operator      │ │CRD           │ │CRDs                          │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Helm-Based Deployment Strategy

### Complete Installation Script

```bash
#!/bin/bash
# CNCF AI Conformance Stack Installation
# Target: Kubernetes 1.34

set -e

# Add all Helm repositories
echo "Adding Helm repositories..."
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo add kedacore https://kedacore.github.io/charts
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm repo add gpu-helm-charts https://nvidia.github.io/dcgm-exporter/helm-charts
helm repo update

# 1. GPU Operator (dra_support, secure_accelerator_access, accelerator_metrics)
echo "Installing NVIDIA GPU Operator..."
helm install gpu-operator nvidia/gpu-operator \
  --namespace gpu-operator \
  --create-namespace \
  --version v25.10.1 \
  --set driver.enabled=true \
  --set toolkit.enabled=true \
  --set devicePlugin.enabled=true \
  --set dcgmExporter.enabled=true \
  --set dcgmExporter.serviceMonitor.enabled=true

# 2. Envoy Gateway (ai_inference)
echo "Installing Envoy Gateway..."
helm install eg oci://docker.io/envoyproxy/gateway-helm \
  --version v1.4.6 \
  -n envoy-gateway-system \
  --create-namespace

# 3. Kueue (gang_scheduling)
echo "Installing Kueue..."
kubectl apply --server-side -f \
  https://github.com/kubernetes-sigs/kueue/releases/download/v0.14.2/manifests.yaml

# 4. Metrics Server (pod_autoscaling - basic)
echo "Installing metrics-server..."
helm install metrics-server metrics-server/metrics-server \
  -n kube-system \
  --set args={--kubelet-insecure-tls}

# 5. kube-prometheus-stack (ai_service_metrics, accelerator_metrics collection)
echo "Installing kube-prometheus-stack..."
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false

# 6. Prometheus Adapter (pod_autoscaling - custom metrics)
echo "Installing Prometheus Adapter..."
helm install prometheus-adapter prometheus-community/prometheus-adapter \
  -n monitoring \
  --set prometheus.url=http://kube-prometheus-stack-prometheus.monitoring.svc \
  --set prometheus.port=9090

# 7. KEDA (pod_autoscaling - advanced)
echo "Installing KEDA..."
helm install keda kedacore/keda \
  --namespace keda \
  --create-namespace

# 8. KubeRay Operator (robust_controller)
echo "Installing KubeRay Operator..."
helm install kuberay-operator kuberay/kuberay-operator \
  --namespace ray-system \
  --create-namespace \
  --version 1.5.1

# Verify installations
echo "Verifying installations..."
kubectl get pods -n gpu-operator
kubectl get pods -n envoy-gateway-system
kubectl get pods -n kueue-system
kubectl get pods -n monitoring
kubectl get pods -n keda
kubectl get pods -n ray-system

echo "AI Conformance stack installation complete!"
```

### Deployment Order and Dependencies

```
┌────────────────────────────────────────────────────────────────┐
│                    DEPLOYMENT ORDER                            │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│  Phase 1: Foundation (No Dependencies)                         │
│  ├── GPU Operator                                              │
│  ├── Envoy Gateway                                             │
│  └── metrics-server                                            │
│                                                                │
│  Phase 2: Core Platform (Depends on Phase 1)                   │
│  ├── kube-prometheus-stack (needs GPU metrics from Phase 1)    │
│  └── Kueue                                                     │
│                                                                │
│  Phase 3: Advanced Features (Depends on Phase 2)               │
│  ├── Prometheus Adapter (needs Prometheus from Phase 2)        │
│  ├── KEDA (needs Prometheus for GPU metrics)                   │
│  └── KubeRay Operator                                          │
│                                                                │
│  Phase 4: Workloads                                            │
│  ├── ClusterQueues / LocalQueues                               │
│  ├── Gateway / HTTPRoutes                                      │
│  └── RayClusters / RayJobs                                     │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

---

## 9. Testing and Validation

### Conformance Test Suite

Based on the provider implementations, here are the tests to validate each requirement:

#### Test 1: DRA Support
```yaml
# test-dra-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dra-test
spec:
  containers:
  - name: cuda-test
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda11.7.1-ubuntu20.04
    resources:
      limits:
        nvidia.com/gpu: "1"
  runtimeClassName: nvidia
```

```bash
# Verify DRA is working
kubectl apply -f test-dra-pod.yaml
kubectl exec dra-test -- nvidia-smi
# Expected: GPU information displayed
```

#### Test 2: Secure Accelerator Access (Comprehensive)

The secure_accelerator_access requirement has multiple sub-tests:

**Test 2a: Unauthorized Access Prevention**
```yaml
# Pod WITHOUT GPU request should NOT see GPU devices
apiVersion: v1
kind: Pod
metadata:
  name: gpu-unauthorized
spec:
  containers:
  - name: test-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda11.7.1-ubuntu20.04
    command: ["sleep", "3600"]
    # Note: No GPU resource request, no runtimeClassName
  restartPolicy: Never
```

```bash
kubectl apply -f gpu-unauthorized.yaml
kubectl exec gpu-unauthorized -- nvidia-smi
# Expected: Command fails - nvidia-smi not found or cannot communicate with driver

kubectl exec gpu-unauthorized -- ls -la /dev/nvidia* 2>&1
# Expected: No such file or directory
```

**Test 2b: Privileged Container Bypass Prevention**
```yaml
# Even privileged containers should NOT access GPUs without resource request
apiVersion: v1
kind: Pod
metadata:
  name: gpu-privileged-test
spec:
  containers:
  - name: privileged-container
    image: ubuntu:22.04
    command: ["sleep", "3600"]
    securityContext:
      privileged: true
    # Note: No GPU resource request
  restartPolicy: Never
```

```bash
kubectl apply -f gpu-privileged-test.yaml
kubectl exec gpu-privileged-test -- ls -la /dev/nvidia* 2>&1
# Expected: Without nvidia.com/gpu request and runtimeClassName, GPU not usable
```

**Test 2c: RuntimeClass Enforcement**
```yaml
# Pod with GPU request but without runtimeClassName
apiVersion: v1
kind: Pod
metadata:
  name: gpu-wrong-runtime
spec:
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda11.7.1-ubuntu20.04
    command: ["sleep", "3600"]
    resources:
      limits:
        nvidia.com/gpu: 1
  # Note: No runtimeClassName specified
  restartPolicy: Never
```

Best practice: Always specify `runtimeClassName: nvidia` for GPU workloads.

#### Test 3: GPU Isolation Between Pods

**Test 3a: Different GPU UUIDs**
```bash
# Deploy two pods requesting 1 GPU each on a multi-GPU node
kubectl exec gpu-pod-1 -- nvidia-smi -L
# Example: GPU 0: NVIDIA A100-SXM4-40GB (UUID: GPU-11111111-...)

kubectl exec gpu-pod-2 -- nvidia-smi -L
# Example: GPU 0: NVIDIA A100-SXM4-40GB (UUID: GPU-22222222-...)

# Expected: Each pod sees exactly 1 GPU, with DIFFERENT UUIDs
```

**Test 3b: Resource Exhaustion Behavior**
```bash
# On single-GPU node, second pod should remain pending
kubectl get pods gpu-pod-1 gpu-pod-2
# Expected on 1-GPU node:
# gpu-pod-1   Running
# gpu-pod-2   Pending

kubectl describe pod gpu-pod-2 | grep -A5 Events
# Expected: "Insufficient nvidia.com/gpu" message
```

#### Test 4: Gang Scheduling with Kueue
```yaml
# test-gang-scheduling.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: gang-test
  labels:
    kueue.x-k8s.io/queue-name: test-queue
spec:
  parallelism: 4
  completions: 4
  template:
    spec:
      containers:
      - name: worker
        image: busybox
        command: ["sleep", "30"]
        resources:
          requests:
            nvidia.com/gpu: "1"
      restartPolicy: Never
```

```bash
# With insufficient GPUs, all pods should remain pending (gang)
kubectl apply -f test-gang-scheduling.yaml
kubectl get pods -l job-name=gang-test
# Expected: All pods pending until sufficient resources
```

#### Test 5: Gateway API Routing
```bash
# Test weighted traffic split
for i in {1..100}; do
  curl -H "Host: inference.example.com" http://gateway-ip/v1/chat/completions
done
# Expected: ~90% to primary, ~10% to canary
```

#### Test 6: Accelerator Metrics
```bash
# Verify GPU metrics in Prometheus
kubectl port-forward -n monitoring svc/kube-prometheus-stack-prometheus 9090:9090
curl 'http://localhost:9090/api/v1/query?query=DCGM_FI_DEV_GPU_UTIL'
# Expected: Per-GPU utilization metrics
```

#### Test 7: HPA with GPU Metrics
```bash
# Verify custom metrics API
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/DCGM_FI_DEV_GPU_UTIL"
# Expected: GPU metrics available via custom metrics API
```

#### Test 8: AI Operator Validation
```bash
# Deploy RayCluster and verify reconciliation
helm install test-cluster kuberay/ray-cluster -n ray-workloads
kubectl get rayclusters -n ray-workloads
kubectl get pods -n ray-workloads
# Expected: Ray head and worker pods running
```

---

## 10. Conclusion and Recommendations

### Recommended Minimum Stack (4 Components)

| # | Component | Requirements Covered | Helm Chart |
|---|-----------|---------------------|------------|
| 1 | **NVIDIA GPU Operator** | `dra_support`, `secure_accelerator_access`, `accelerator_metrics`, `pod_autoscaling`* | `nvidia/gpu-operator` |
| 2 | **Dynamo Platform**** | `robust_controller`, `ai_inference`, `gang_scheduling` | `nvidia/ai-dynamo/dynamo-platform` |
| 3 | **kube-prometheus-stack** | `ai_service_metrics` | `prometheus-community/kube-prometheus-stack` |
| 4 | **Cluster Autoscaler** | `cluster_autoscaling` | `autoscaler/cluster-autoscaler` |

*`pod_autoscaling` uses standard Kubernetes HPA with DCGM Exporter metrics from GPU Operator
**With `--set grove.enabled=true --set kai-scheduler.enabled=true`

> **⚠️ Important:** Cluster Autoscaler is NOT bundled with Dynamo Platform - it must be installed separately.

### Complete Installation (Minimum Stack)

```bash
#!/bin/bash
# CNCF AI Conformance v1.34 - Minimum Stack Installation
# Updated: January 2026

export DYNAMO_VERSION=0.8.0  # https://github.com/ai-dynamo/dynamo/releases

# Add repos
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm repo update

# 1. GPU Operator v25.10.1 (drivers, device plugin, DCGM metrics)
helm install gpu-operator nvidia/gpu-operator \
  -n gpu-operator --create-namespace \
  --version v25.10.1 \
  --set dcgmExporter.serviceMonitor.enabled=true

# 2. Dynamo Platform v0.8.0 (operator + optional Grove/KAI Scheduler)
# Note: Grove and KAI Scheduler are NOT bundled by default - must enable
helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-crds-${DYNAMO_VERSION}.tgz
helm install dynamo-crds dynamo-crds-${DYNAMO_VERSION}.tgz -n default

helm fetch https://helm.ngc.nvidia.com/nvidia/ai-dynamo/charts/dynamo-platform-${DYNAMO_VERSION}.tgz
helm install dynamo-platform dynamo-platform-${DYNAMO_VERSION}.tgz \
  -n dynamo-system --create-namespace \
  --set grove.enabled=true \
  --set kai-scheduler.enabled=true

# 3. Prometheus Stack (Chart 80.14.4, metrics collection)
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false

# 4. Cluster Autoscaler (App 1.34.2, for cloud environments)
# helm install cluster-autoscaler autoscaler/cluster-autoscaler \
#   -n kube-system --set cloudProvider=aws --set autoDiscovery.clusterName=CLUSTER_NAME
```

### Implementation Priorities

1. **Start with GPU Operator**: Foundation for all GPU workloads
2. **Add Dynamo Platform**: Includes operator, scheduling, and frontend
3. **Add Prometheus**: For metrics collection and alerting
4. **Add Cluster Autoscaler**: If running in cloud environment

### Key Considerations

1. **Minimum stack is sufficient**: 4 Helm releases cover all 9 MUST requirements
2. **Dynamo bundles multiple capabilities**: Operator, Grove, KAI Scheduler, Gateway API integration
3. **GPU Operator bundles GPU infrastructure**: Drivers, device plugin, DCGM Exporter
4. **All components are 100% open source**: Apache 2.0 licensed
5. **Kubernetes-native**: No external dependencies (etcd/NATS optional)
6. **Cluster Autoscaler is separate**: NOT bundled with Dynamo - requires separate installation or cloud provider autoscaling
7. **Pod Autoscaling uses standard HPA**: Combined with DCGM Exporter metrics; KAI Scheduler enhances scheduling but doesn't replace HPA

### Advanced NVIDIA Open Source Options

For production deployments requiring higher performance and scale:

| Component | Purpose | When to Consider |
|-----------|---------|------------------|
| **NVIDIA Dynamo** | Distributed inference framework | High-throughput LLM serving, disaggregated prefill/decode |
| **KAI Scheduler** | Advanced GPU scheduling | Multi-tenant clusters, fractional GPU sharing, topology-aware |
| **Grove** | Orchestration API | Complex multi-component AI systems, datacenter scale |

These components are 100% open source (Apache 2.0) and integrate with the basic conformance stack.

### Future Considerations

- **KAR-0001 (Driver/Runtime Management)**: Expect automated driver version verification
- **KAR-0003 (GPU Sharing)**: MIG and time-slicing support will become standard; KAI Scheduler provides this today
- **KAR-0004 (Virtualized Accelerators)**: vGPU support through DRA is emerging
- **Automated testing (v1.37+)**: New requirements will require automated conformance tests
- **Disaggregated Inference**: NVIDIA Dynamo/Grove patterns likely to influence future conformance requirements

### Estimated Resource Requirements

| Component | CPU Request | Memory Request | Notes |
|-----------|-------------|----------------|-------|
| GPU Operator | 100m | 200Mi | Per cluster |
| Envoy Gateway | 100m | 256Mi | Per replica |
| Kueue | 100m | 256Mi | Single controller |
| kube-prometheus-stack | 500m | 1Gi | Includes Prometheus, Grafana |
| KEDA | 100m | 256Mi | Controller + metrics adapter |
| KubeRay | 100m | 256Mi | Operator only |
| **Total Platform Overhead** | ~1 CPU | ~2.5Gi | Excluding workloads |

---

## Appendix A: Helm Values Files

### GPU Operator values.yaml
```yaml
driver:
  enabled: true
  repository: nvcr.io/nvidia
  version: "570.76.02"
toolkit:
  enabled: true
devicePlugin:
  enabled: true
dcgmExporter:
  enabled: true
  serviceMonitor:
    enabled: true
mig:
  strategy: single
```

### Kueue configuration
```yaml
# cluster-queue.yaml
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: gpu-queue
spec:
  namespaceSelector: {}
  resourceGroups:
  - coveredResources: ["cpu", "memory", "nvidia.com/gpu"]
    flavors:
    - name: default
      resources:
      - name: "nvidia.com/gpu"
        nominalQuota: 8
      - name: "cpu"
        nominalQuota: 64
      - name: "memory"
        nominalQuota: 256Gi
```

---

## Appendix B: References

### Official Documentation
- [CNCF k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance)
- [WG AI Conformance](https://github.com/kubernetes-sigs/wg-ai-conformance)
- [Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/)
- [Kueue Documentation](https://kueue.sigs.k8s.io/)
- [NVIDIA GPU Operator](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/)
- [KubeRay Documentation](https://ray-project.github.io/kuberay/)

### NVIDIA Open Source AI Infrastructure
- [NVIDIA Dynamo](https://github.com/ai-dynamo/dynamo) - Distributed inference framework
- [NVIDIA KAI Scheduler](https://github.com/NVIDIA/KAI-Scheduler) - GPU-native Kubernetes scheduler
- [NVIDIA Grove](https://github.com/NVIDIA/grove) - Kubernetes AI orchestration API
- [NVIDIA DRA Driver for GPUs](https://github.com/NVIDIA/k8s-dra-driver-gpu) - Dynamic Resource Allocation

### Helm Repositories
- NVIDIA GPU Operator: https://helm.ngc.nvidia.com/nvidia
- KAI Scheduler: ghcr.io/nvidia/kai-scheduler/kai-scheduler (OCI)
- Prometheus Community: https://prometheus-community.github.io/helm-charts
- KubeRay: https://ray-project.github.io/kuberay-helm/
- KEDA: https://kedacore.github.io/charts
- Envoy Gateway: docker.io/envoyproxy/gateway-helm (OCI)
- Volcano: https://volcano-sh.github.io/helm-charts

### Additional Resources
- [NVIDIA Dynamo on EKS](https://awslabs.github.io/ai-on-eks/docs/blueprints/inference/GPUs/nvidia-dynamo)
- [NVIDIA Dynamo on AKS](https://blog.aks.azure.com/2025/10/24/dynamo-on-aks)
- [KAI Scheduler + KubeRay Integration](https://developer.nvidia.com/blog/enable-gang-scheduling-and-workload-prioritization-in-ray-with-nvidia-kai-scheduler/)
- [Grove Technical Blog](https://developer.nvidia.com/blog/streamline-complex-ai-inference-on-kubernetes-with-nvidia-grove/)

---

*Document prepared for cloud-native-stack project*
*Based on analysis of CNCF AI Conformance v1.34 and 12 provider implementations*
*Includes NVIDIA open source advanced components: Dynamo, KAI Scheduler, Grove*
