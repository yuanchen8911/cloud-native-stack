# GPU Operator Configuration Guide

A practical guide for deploying and configuring the NVIDIA GPU Operator across Kubernetes platforms.

---

## Table of Contents

1. [GPU Operator Versioning](#gpu-operator-versioning)
2. [Component Versions](#component-versions)
3. [GPU Driver Management](#gpu-driver-management)
4. [Container Device Interface (CDI) Configuration](#container-device-interface-cdi-configuration)
5. [CDI and Container Runtime Compatibility](#cdi-and-container-runtime-compatibility)
6. [Common Issues and Solutions](#common-issues-and-solutions)

---

## GPU Operator Versioning

### Version Pattern

The GPU Operator follows a **YY.MM.PP** calendar versioning scheme:

| Segment | Description | Example |
|---------|-------------|---------|
| **YY.MM** | Major version (release date) | 25.10, 25.3 |
| **PP** | Patch version | 0, 1, 2 |

### Product Lifecycle

| Phase | Description |
|-------|-------------|
| **GA** | Active development with new features |
| **Maintenance** | Critical bugs and CVE fixes only |
| **EOL** | No longer supported |

### Current Supported Versions

| GPU Operator Version | Status |
|---------------------|--------|
| 25.10.x | GA (Generally Available) |
| 25.3.x | Maintenance |
| 24.9.x and earlier | EOL (End of Life) |

### Upgrade Policy

| Upgrade Path | Supported |
|--------------|-----------|
| Within major release (25.10.0 → 25.10.1) | ✅ Yes |
| To next major release (25.3.x → 25.10.x) | ✅ Yes |
| Skipping major versions (24.9.x → 25.10.x) | ❌ No |

---

## Component Versions

The following table shows validated component versions for **GPU Operator 25.10.1** (latest):

| Component | Version |
|-----------|---------|
| **NVIDIA GPU Operator** | v25.10.1 |
| **NVIDIA GPU Driver** | 580.105.08 (default, recommended)<br>580.95.05<br>580.82.07<br>575.57.08<br>570.195.03<br>550.163.01<br>535.274.02 |
| **NVIDIA Driver Manager** | v0.9.1 |
| **NVIDIA Container Toolkit** | 1.18.0 |
| **NVIDIA Device Plugin** | 0.18.1 |
| **DCGM Exporter** | v4.4.2-4.7.0 |
| **Node Feature Discovery** | v0.18.2 |
| **NVIDIA GPU Feature Discovery** | 0.18.1 |
| **NVIDIA MIG Manager** | 0.13.1 |
| **DCGM** | 4.4.2-1 |
| **Validator** | v25.10.0 |
| **NVIDIA KubeVirt GPU Device Plugin** | v1.4.0 |
| **NVIDIA vGPU Device Manager** | v0.4.1 |
| **NVIDIA GDS Driver** | 2.26.6 ⚠️ |
| **NVIDIA Kata Manager** | v0.2.3 |
| **NVIDIA Confidential Computing Manager** | v0.1.1 |
| **NVIDIA GDRCopy Driver** | v2.5.1 |

> **Notes**:
> - GKE with Container-Optimized OS (COS) manages GPU drivers separately from GPU Operator
> - ⚠️ **GDS Driver**: Requires NVIDIA Open GPU Kernel module driver
> - **Known Issue**: Drivers 570.124.06, 570.133.20, 570.148.08, and 570.158.01 have MIG scheduling issues. Upgrade to 580.65.06+ recommended

---

## GPU Driver Management

### Driver Versioning Strategy

#### Release Types
- **Major Feature Release**: New branch (X number), released quarterly
- **Production Branch**: Enterprise/datacenter qualified, supported for 1 year
- **Long Term Support Branch (LTSB)**: Extended 3-year support with critical updates

#### GPU Operator Driver Policy
- **Recommended Version**: Latest driver available after GPU Operator release
- **Default Version**: Latest driver available at GPU Operator release time
- **Validation**: All supported drivers undergo GPU Operator QA validation

### Platform-Specific Considerations

#### Standard Kubernetes
GPU Operator manages driver installation and lifecycle.

#### GKE with Container-Optimized OS (COS)
GPU drivers are managed by GKE, not GPU Operator.

**Option 1: GKE-Managed Installation**
Configure during cluster/node pool creation:
- `INSTALLATION_DISABLED`: Manual installation required
- `DEFAULT`: Standard GPU driver for COS/Ubuntu
- `LATEST`: Latest available GPU driver

**Option 2: Manual Installation**
1. Disable auto-installation during node pool creation
2. Install specific driver version post-creation using `gcp-nvidia-drivers/nvidia-driver-installer`

**Useful References:**
- [GKE to COS version mapping](https://www.gstatic.com/gke-image-maps/gke-to-cos.json)
- [COS release notes](https://cloud.google.com/container-optimized-os/docs/release-notes/m113)
- [Driver versions per COS](https://storage.googleapis.com/cos-tools/18244.236.88/lakitu/gpu_driver_versions.textproto)

---

## Container Device Interface (CDI) Configuration

### What is CDI?

CDI (Container Device Interface) is a specification that standardizes how container runtimes access specialized hardware devices. For GPUs, CDI provides:

- Declarative device configuration via JSON spec files
- Runtime-agnostic device injection
- Cleaner separation between device vendors and container runtimes

### Key Concepts

Understanding CDI requires distinguishing two separate concerns:

| Concern | Description | Controlled By |
|---------|-------------|---------------|
| **CDI Spec Generation** | Creating JSON files that describe GPU devices | GPU Operator (`cdi.enabled`) |
| **CDI Spec Consumption** | Runtime reading and applying CDI specs | Container runtime (containerd/CRI-O) |

Both must be functional for CDI to work.

### GPU Operator CDI Parameters

| Parameter | Description | Status |
|-----------|-------------|--------|
| `cdi.enabled` | Enables CDI spec generation by nvidia-container-toolkit | Active |
| `cdi.default` | Makes CDI the default injection mechanism | **Deprecated in v25.10.0+** (no-op) |

### Configuration Scenarios

#### Scenario 1: Legacy Mode (No CDI)

```yaml
cdi:
  enabled: false
```

- Uses traditional nvidia-container-runtime hook
- No CDI specs generated
- Works with all container runtimes
- Proven stability, but less flexible

#### Scenario 2: CDI Enabled

```yaml
cdi:
  enabled: true
```

- Generates CDI specs in `/var/run/cdi/`
- Runtime must support CDI (see compatibility section below)
- **Recommended for new deployments**

> **Note (v25.10.0+)**: The `cdi.default` parameter is deprecated and ignored. CDI behavior is determined by runtime configuration and RuntimeClass selection.

### Platform-Specific Recommendations

| Platform | Configuration | Rationale |
|----------|---------------|-----------|
| **AWS EKS** | `cdi.enabled: true` | containerd with CDI support |
| **GKE** | `cdi.enabled: true` (required) | Must be explicit; cannot rely on defaults |
| **Azure AKS** | `cdi.enabled: true` | containerd with CDI support |
| **On-prem (containerd 2.0+)** | `cdi.enabled: true` | CDI enabled by default in runtime |
| **On-prem (containerd 1.7)** | `cdi.enabled: true` | Toolkit configures runtime automatically |
| **On-prem (CRI-O)** | `cdi.enabled: true` | CRI-O has native CDI support |

### Device List Strategy

The device plugin uses `DEVICE_LIST_STRATEGY` to determine how GPUs are communicated to containers:

```yaml
devicePlugin:
  env:
    - name: DEVICE_LIST_STRATEGY
      value: "cdi-cri,cdi-annotations,volume-mounts"
```

| Strategy | Description |
|----------|-------------|
| `cdi-cri` | Passes CDI device IDs via CRI (preferred with CDI) |
| `cdi-annotations` | Uses pod annotations for CDI devices |
| `volume-mounts` | Legacy volume mount method |

The value is a comma-separated fallback chain—the device plugin tries each strategy in order.

> **Known Issue**: Combining `cdi.enabled=true` with `DEVICE_LIST_STRATEGY=volume-mounts` alone can cause issues. Use the full fallback chain when CDI is enabled.

---

## CDI and Container Runtime Compatibility

### How CDI Works with Container Runtimes

```
┌─────────────────────┐     ┌──────────────────────┐     ┌─────────────────┐
│  GPU Operator       │     │  CDI Spec Files      │     │  Container      │
│  (nvidia-toolkit)   │────▶│  /var/run/cdi/*.json │────▶│  Runtime        │
│                     │     │                      │     │  (reads specs)  │
│  cdi.enabled: true  │     │  Device definitions  │     │                 │
└─────────────────────┘     └──────────────────────┘     └─────────────────┘
        GENERATES                   STORED                    CONSUMES
```

### Container Runtime CDI Support Matrix

| Runtime | Version | CDI Support | Configuration Needed |
|---------|---------|-------------|---------------------|
| **CRI-O** | ≥1.23 | Native (enabled by default) | None—works out of the box |
| **CRI-O** | 1.20–1.22 | Available (flag required) | Enable via `--enable-cdi` |
| **containerd** | ≥2.0 | Native (enabled by default) | None—works out of the box |
| **containerd** | 1.7.x | Available (disabled by default) | Toolkit auto-configures when `cdi.enabled: true` |
| **containerd** | <1.7 | Not available | Use legacy mode (`cdi.enabled: false`) |

### Runtime-Specific Notes

#### CRI-O

CRI-O has **native CDI support** enabled by default since version 1.23. This means:

- The runtime can consume CDI specs without any configuration changes
- GPU Operator's `cdi.enabled: true` is still required to **generate** the CDI specs
- No runtime configuration or restart needed

```yaml
# Correct configuration for CRI-O
cdi:
  enabled: true  # Required to generate CDI specs
```

#### containerd 2.0+

containerd 2.0 enables CDI by default:

- CDI spec consumption works out of the box
- GPU Operator's `cdi.enabled: true` generates the specs
- Simplest CDI deployment path

#### containerd 1.7.x

containerd 1.7 has CDI support but disabled by default:

- Requires `enable_cdi = true` in containerd config
- **nvidia-container-toolkit automatically configures this** when GPU Operator sets `cdi.enabled: true`
- Runtime restart may be required after initial setup

### GKE-Specific Requirements

GKE requires explicit CDI configuration:

```yaml
# Required for GKE deployments
cdi:
  enabled: true  # Must be explicitly set
```

- Cannot rely on default values
- GKE manages containerd configuration differently
- Always verify CDI is enabled in your Helm values

---

## Common Issues and Solutions

### Issue 1: CDI + Volume Mounts Strategy Conflict

**Symptom**: GPU pods fail to start or devices not visible in containers.

**Cause**: Mismatch between CDI enablement and device list strategy.

**Affected Configurations**:
```yaml
# Problematic combination (GPU Operator <25.10.0)
cdi:
  enabled: true
  default: false
devicePlugin:
  env:
    - name: DEVICE_LIST_STRATEGY
      value: "volume-mounts"  # Missing CDI strategies
```

**Solution**:
```yaml
cdi:
  enabled: true
devicePlugin:
  env:
    - name: DEVICE_LIST_STRATEGY
      value: "cdi-cri,cdi-annotations,volume-mounts"
```

**Status**: Resolved in GPU Operator v25.10.0+ where `cdi.default` is deprecated.

---

### Issue 2: EFA Compatibility on EKS with H100

**Symptom**: Network performance issues or job failures on H100 instances.

**Cause**: H100 GPUs with R580 driver require EFA version 1.43.3+.

**Solution**: Upgrade to AMI images with EFA 1.43.3 or later.

**Verification**:
```bash
# Check EFA version
modinfo efa | grep version
```

---

### Issue 3: Missing GPU Feature Discovery Pods

**Symptom**: `nvidia-gpu-feature-discovery` pods not scheduled.

**Cause**: Nodes missing required labels or tolerations not configured.

**Diagnosis**:
```bash
# Check pending pods
kubectl get pods -n gpu-operator -l app=nvidia-gpu-feature-discovery

# Check node labels
kubectl get nodes --show-labels | grep -E 'nvidia|gpu'

# Check pod events
kubectl describe pod <pod-name> -n gpu-operator
```

**Solution**: Ensure nodes have appropriate labels:

```yaml
# Common required labels
nvidia.com/gpu.present: "true"

# GKE-specific
gke-no-default-nvidia-gpu-device-plugin: "true"
```

---

### Issue 4: CDI Specs Not Generated

**Symptom**: No files in `/var/run/cdi/` or `/etc/cdi/`.

**Diagnosis**:
```bash
# Check nvidia-container-toolkit pod logs
kubectl logs -n gpu-operator -l app=nvidia-container-toolkit-daemonset

# Verify CDI directory on node
ls -la /var/run/cdi/
```

**Common Causes**:
1. `cdi.enabled: false` in GPU Operator values
2. Toolkit pod not running or failing
3. Insufficient permissions to write CDI specs

---

### Issue 5: Runtime Not Consuming CDI Specs

**Symptom**: CDI specs exist but containers don't see GPUs.

**Diagnosis**:
```bash
# Check containerd config for CDI
cat /etc/containerd/config.toml | grep -A5 cdi

# For CRI-O, check if CDI is enabled
crictl info | jq '.config.cdi'
```

**Solution for containerd 1.7.x** (if toolkit didn't auto-configure):
```toml
# /etc/containerd/config.toml
[plugins."io.containerd.grpc.v1.cri"]
  enable_cdi = true
  cdi_spec_dirs = ["/etc/cdi", "/var/run/cdi"]
```

Then restart containerd: `systemctl restart containerd`

---

## Key Takeaways

### Configuration Best Practices

| Practice | Rationale |
|----------|-----------|
| Use `cdi.enabled: true` for new deployments | CDI is the modern, standardized approach |
| Always set CDI explicitly on GKE | Cannot rely on defaults |
| Use full device list strategy fallback | `cdi-cri,cdi-annotations,volume-mounts` |
| Validate node labels before deployment | Prevents scheduling issues |
| Match EFA versions with driver on EKS | Required for H100+ GPUs |
| Upgrade to v25.10.0+ | Simplifies CDI configuration |

### Quick Reference: CDI Configuration by Runtime

| Runtime | GPU Operator Config | Runtime Config |
|---------|--------------------|--------------------|
| CRI-O ≥1.23 | `cdi.enabled: true` | None needed |
| containerd 2.0+ | `cdi.enabled: true` | None needed |
| containerd 1.7.x | `cdi.enabled: true` | Auto-configured by toolkit |
| containerd <1.7 | `cdi.enabled: false` | N/A (use legacy) |

### Troubleshooting Checklist

```
[ ] GPU Operator version compatible with cluster?
[ ] CDI enabled if using modern runtime?
[ ] DEVICE_LIST_STRATEGY includes CDI options?
[ ] Node labels present for GPU scheduling?
[ ] CDI specs generated in /var/run/cdi/?
[ ] Container runtime configured to consume CDI?
[ ] Driver version compatible with GPU hardware?
```

### Reference Resources

| Resource | URL |
|----------|-----|
| GPU Operator Docs | https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/ |
| CDI Specification | https://github.com/cncf-tags/container-device-interface |
| NVIDIA Driver Downloads | https://www.nvidia.com/drivers |
| containerd CDI Docs | https://github.com/containerd/containerd/blob/main/docs/cdi.md |

---

*Last updated: January 2025*
