# NVIDIA GPU Operator Configuration Guide
## Knowledge Summary and Best Practices

---

## Table of Contents
1. [GPU Operator Versioning](#gpu-operator-versioning)
2. [Component Versions](#component-versions)
3. [GPU Driver Management](#gpu-driver-management)
4. [Container Device Interface (CDI) Configuration](#container-device-interface-cdi-configuration)
5. [CDI Runtime Considerations](#cdi-runtime-considerations)
6. [Common Issues and Solutions](#common-issues-and-solutions)

---

## GPU Operator Versioning

### Version Pattern
The GPU Operator follows a **YY.MM.PP** versioning scheme:
- **YY.MM**: Major version indicating initial release date (e.g., 23.6, 24.9)
- **PP**: Patch version for critical bugs, CVEs, and minor features

### Supported Versions
| GPU Operator Version | Status |
|---------------------|--------|
| 24.9.x | GA (Generally Available) |
| 24.6.x | Maintenance |
| 24.3.x and lower | EOL (End of Life) |

### Upgrade Policy
- ✅ **Supported**: Upgrades within major releases (e.g., 24.9.0 → 24.9.2)
- ✅ **Supported**: Upgrades to next major release (e.g., 24.6.x → 24.9.x)
- ❌ **Unsupported**: Skipping major versions

---

## Component Versions

The following table shows validated component versions for **GPU Operator 24.9.2**:

| Component | Version |
|-----------|---------|
| **NVIDIA GPU Operator** | v24.9.2 |
| **NVIDIA GPU Driver** | 570.86.15 (recommended)<br>565.57.01<br>560.35.03<br>550.144.03 (default)<br>550.127.08<br>535.230.02<br>535.216.03 |
| **NVIDIA Driver Manager** | v0.7.0 |
| **NVIDIA Container Toolkit** | 1.17.4 |
| **NVIDIA Device Plugin** | 0.17.0 |
| **DCGM Exporter** | 3.3.9-3.6.1 |
| **Node Feature Discovery** | v0.16.6 |
| **NVIDIA GPU Feature Discovery** | 0.17.0 |
| **NVIDIA MIG Manager** | 0.10.0 |
| **Validator** | v24.9.2 |

> **Note**: GKE with Container-Optimized OS (COS) manages GPU drivers separately from GPU Operator.

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

### CDI Settings Overview
CDI provides a standardized way to expose specialized devices to containers.

#### Key Configuration Parameters
- `cdi.enabled`: Enables CDI support across the stack
- `cdi.default`: Makes CDI the default GPU injection mechanism

### Configuration Scenarios

#### Scenario 1: Legacy Mode (Default)
```yaml
cdi:
  enabled: false
```
- Uses traditional GPU allocation
- No runtime class specification needed
- Proven stability

#### Scenario 2: CDI Available, Legacy Default
```yaml
cdi:
  enabled: true
  default: false
```
- **Default behavior**: Legacy GPU allocation
- **Opt-in to CDI**: Set `runtimeClassName: nvidia-cdi` in pod spec
- **Known Issue**: ⚠️ **Avoid with `DEVICE_LIST_STRATEGY=volume-mounts`** (causes deployment issues)

#### Scenario 3: CDI Default Mode
```yaml
cdi:
  enabled: true
  default: true
```
- **Default behavior**: CDI-based GPU allocation
- **Opt-out to legacy**: Set `runtimeClassName: nvidia` in pod spec
- **Recommended for**: AWS deployments, newer clusters

### Platform-Specific Recommendations

| Platform | CDI Configuration | Notes |
|----------|------------------|--------|
| **AWS EKS** | `enabled: true, default: true` | Recommended configuration |
| **GKE** | `enabled: true` (explicit) | Must explicitly enable CDI |
| **OCI with CRI-O** | `enabled: false, default: false` | CRI-O has built-in CDI support |

### Device List Strategy
When CDI is enabled, configure:
```yaml
devicePlugin:
  env:
    - name: DEVICE_LIST_STRATEGY
      value: "cdi-cri,cdi-annotations,volume-mounts"
```

---

## CDI Runtime Considerations

### Version-Specific Changes

#### GPU Operator v25.10.0+
- `cdi.default` parameter becomes **no-op** (ignored)
- CDI behavior determined by other configuration parameters

### Runtime Differences

#### Container Runtime Support
- **CRI-O**: CDI enabled by default in runtime
- **containerd v1.7.x**: CDI disabled by default, enabled by toolkit when in CDI mode
- **Abstraction**: GPU Operator toolkit abstracts runtime differences

#### GKE-Specific Requirements
- CDI must be **explicitly enabled** (`cdi.enabled: true`)
- Cannot rely on default CDI settings

---

## Common Issues and Solutions

### Issue 1: CDI Configuration Bug
**Problem**: Combination of `cdi.enabled=true`, `cdi.default=false` with `DEVICE_LIST_STRATEGY=volume-mounts`

**Impact**: Deployment failures on GB200 AWS clusters

**Solution**:
```yaml
cdi:
  enabled: true
  default: true  # Set both to true
```

**Status**: Fixed in upcoming GPU Operator release (not backported to 25.3.2)

### Issue 2: EFA Compatibility on EKS H100
**Problem**: H100 GPUs require newer EFA version (1.43.3) with R580 driver

**Impact**: Network performance issues, job failures

**Solution**: Upgrade to AMI images that include EFA version 1.43.3+

### Issue 3: Missing GPU Feature Discovery Pods
**Problem**: Pods not created due to missing tolerations/node selectors

**Root Cause**: Nodes lack required labels for GPU Operator component scheduling

**Solution**: Ensure nodes have appropriate labels:

#### Required Node Labels (HIPPO GKE Example):
```yaml
# Node group classification
nodeGroup: customer-gpu | customer-cpu | system-cpu

# Workload dedication
node.dgxc.nvidia.com/dedicated: customer-workload | system-workload

# GKE-specific requirement for node-feature-discovery
gke-no-default-nvidia-gpu-device-plugin: "true"
```

#### Verification Steps:
1. Check node labels: `kubectl get nodes --show-labels`
2. Verify GPU Operator component scheduling requirements
3. Apply missing labels to nodes
4. Restart affected pods if necessary

---

## Key Takeaways

### 🎯 **Configuration Best Practices**
1. **Use CDI default mode** (`enabled: true, default: true`) for new deployments
2. **Explicitly configure CDI** on GKE clusters
3. **Validate node labels** before GPU Operator deployment
4. **Match EFA versions** with GPU driver requirements on EKS

### 🔧 **Troubleshooting Checklist**
- [ ] Verify GPU Operator version compatibility
- [ ] Check CDI configuration for known issues
- [ ] Validate node labels and tolerations
- [ ] Confirm driver version compatibility
- [ ] Review platform-specific requirements

### 📚 **Reference Resources**
- [GPU Operator Documentation](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/)
- [CDI Specification](https://github.com/container-orchestrated-devices/container-device-interface)
- [NVIDIA GPU Driver Versions](https://www.nvidia.com/drivers)

---

*This document consolidates practical knowledge and lessons learned from GPU Operator deployments across various Kubernetes platforms. Keep this guide updated as new versions and configurations are validated.*