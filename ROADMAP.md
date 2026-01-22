# Cloud Native Stack Roadmap

> **Motif:** Launch CNS v2

This roadmap defines the work required to launch Cloud Native Stack v2.

The roadmap is organized into two sections:
- **MVP** — Scope required to launch. All MVP initiatives are P0 (must-have).
- **Backlog** — Items not yet fully evaluated or deemed non-critical for MVP launch.

---

**Table of Contents**
- [MVP](#mvp)
  - [Initiative 1: GitHub Project Setup for Sustainable OSS Operations](#initiative-1-github-project-setup-for-sustainable-oss-operations)
  - [Initiative 2: CNS CLI/API/Agent Validated Features up to MVP Scope](#initiative-2-cns-cliapiagent-validated-features-up-to-mvp-scope)
  - [Initiative 3: Downstream Users Validated and Integrated](#initiative-3-downstream-users-validated-and-integrated)
- [Backlog](#backlog)
- [Revision History](#revision-history)

---

# MVP

All MVP initiatives are **P0 (must-have)** for launch.

---

## Initiative 1: GitHub Project Setup for Sustainable OSS Operations

**Expected State:** The CNS repository is configured for professional open-source operations with automated quality gates, security scanning, clear contribution guidelines, and a reliable release process.

---

### Epic 1.1: V2 Content Integration and Release

**User Story**

> As an CNS contributor,
> I want all CNS v2 content integrated into the repository and released,  
> So that attendees can immediately access and use the announced capabilities.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** CNS v2 code is merged to main, **When** a release tag is created, **Then** all binaries (cnsctl, cnsd) are published to GitHub Releases with SBOM and SLSA attestations | Verify release artifacts exist and attestations pass `gh attestation verify` |
| 2 | **Given** the release is published, **When** a user downloads cnsctl, **Then** `cnsctl version` displays the correct semantic version | Manual verification |
| 3 | **Given** documentation is complete, **When** a user visits the repository, **Then** README provides clear installation instructions and links to user guides | Documentation review |
| 4 | **Given** all MVP recipes exist, **When** the release is tagged, **Then** recipe YAML files are included in the release artifacts | Verify files in release tarball |

**Non-Functional Requirements**
- Release process completes in < 15 minutes
- All release artifacts are cryptographically signed
- Documentation renders correctly on GitHub

**Dependencies**
- Epic 1.2 (Repository operations) must be complete
- All MVP recipes finalized (Epic 2.1)

**Assumptions**
- GHCR and GitHub Releases remain available
- Ko, GoReleaser, and Syft tooling is stable

**Out of Scope**
- Marketing materials and blog posts
- Conference booth demos (separate track)
- Helm chart distribution (uses embedded templates)

**Definition of Done**
- [ ] Release tag created and CI/CD pipeline succeeds
- [ ] All binaries published with checksums and attestations
- [ ] Container images pushed to ghcr.io with SBOM
- [ ] README updated with v2 installation instructions
- [ ] Release notes document all MVP features
- [ ] At least one maintainer has verified end-to-end installation

---

### Epic 1.2: Repository Operations Formalized and Implemented

**User Story**

> As a CNS maintainer,  
> I want repository operations (CI/CD, code ownership, security scanning) formalized and automated,  
> So that contributions are consistently validated and releases are secure and reproducible.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** a PR is opened, **When** CI runs, **Then** all tests pass with race detector, linting passes, and security scan completes | PR cannot merge without green CI |
| 2 | **Given** DCO enforcement is enabled, **When** a commit lacks sign-off, **Then** the PR is blocked | Submit unsigned commit, verify block |
| 3 | **Given** CODEOWNERS is configured, **When** a PR modifies pkg/bundler/*, **Then** the appropriate team is auto-assigned for review | Verify GitHub assigns reviewers |
| 4 | **Given** a security vulnerability is detected, **When** Trivy scans run, **Then** results are uploaded to GitHub Security tab | Check Security tab after scan |
| 5 | **Given** a release tag is pushed, **When** on-tag workflow runs, **Then** SLSA Build Level 3 provenance is generated | Verify with `gh attestation verify` |

**Tasks**

| Task | Description |
|------|-------------|
| Repo Configuration | Branch protection rules, required status checks, auto-merge settings |
| DCO/Signing | Enable DCO bot, document sign-off requirements in CONTRIBUTING.md |
| CODEOWNERS | Define ownership for pkg/*, docs/*, .github/* |
| PR Qualification | on-push.yaml: test, lint, scan pipeline |
| Release Workflow | on-tag.yaml: build, attest, publish, deploy |
| Security Scanning | Trivy integration, SARIF upload, Dependabot alerts |
| S3C (Supply Chain Security) | SBOM generation, Sigstore signing, Rekor transparency |

**Non-Functional Requirements**
- CI pipeline completes in < 10 minutes for PRs
- Security scans detect MEDIUM+ vulnerabilities
- All workflows use pinned action versions (SHA, not tags)

**Dependencies**
- GitHub repository admin access
- GHCR write permissions for release workflow
- GCP Workload Identity Federation for Cloud Run deploy

**Assumptions**
- GitHub Actions remains the CI/CD platform
- Current composite action architecture is maintained

**Out of Scope**
- Self-hosted runners
- Alternative CI systems (Jenkins, GitLab CI)
- Code signing for Windows binaries

**Definition of Done**
- [ ] Branch protection enabled on main
- [ ] DCO enforcement active and documented
- [ ] CODEOWNERS file covers all critical paths
- [ ] on-push.yaml passes on sample PR
- [ ] on-tag.yaml successfully publishes test release
- [ ] Security scan results visible in GitHub Security tab
- [ ] CONTRIBUTING.md updated with workflow documentation

---

## Initiative 2: CNS CLI/API/Agent Validated Features up to MVP Scope

**Expected State:** CNS tooling (CLI, API server, Kubernetes agent) reliably generates recipes, produces valid deployment bundles, and validates cluster state for all MVP-scoped platforms and configurations.

---

### Epic 2.1: MVP Recipes and Components Finalized

**User Story**

> As a platform engineer evaluating CNS,  
> I want a clear, validated set of recipes covering my target platforms,  
> So that I can confidently deploy GPU workloads on supported configurations.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** the MVP recipe list is defined, **When** I query `cnsctl recipe --list`, **Then** all MVP recipes are enumerated with their target platforms | CLI output verification |
| 2 | **Given** each MVP recipe, **When** I inspect its components, **Then** all required component versions are specified and resolvable | Recipe YAML validation |
| 3 | **Given** each MVP recipe, **When** bundlers generate artifacts, **Then** generated Helm values match documented component versions | Diff generated values against expected |
| 4 | **Given** recipe documentation, **When** a user reads it, **Then** supported platforms, prerequisites, and limitations are clearly stated | Documentation review |

**Tasks**

| Task | Description |
|------|-------------|
| Define MVP recipe matrix | Platforms (EKS, GKE, AKS, OKE), Accelerators (H100, GB200, A100), Intents (training, inference) |
| Validate component versions | GPU Operator, Network Operator, Cert-Manager versions locked to known-good |
| Create recipe YAML files | One recipe per platform/accelerator/intent combination |
| Document each recipe | Prerequisites, deployment steps, known limitations |
| Establish recipe testing | Automated validation that recipes resolve correctly |

**Non-Functional Requirements**
- Recipe resolution completes in < 1 second
- Recipe YAML files pass schema validation
- All component versions have verified compatibility

**Dependencies**
- Bundler support for all referenced components (Epic 2.3)
- Access to target platforms for validation

**Assumptions**
- Component versions remain stable through launch
- Platform APIs do not change significantly

**Out of Scope**
- Self-managed Kubernetes (non-cloud)
- Deprecated GPU architectures (V100, P100)
- Windows node support

**Definition of Done**
- [ ] MVP recipe matrix documented and approved
- [ ] All MVP recipe YAML files created in examples/recipes/
- [ ] Schema validation passes for all recipes
- [ ] Each recipe has corresponding documentation
- [ ] At least one successful deployment per recipe on target platform

---

### Epic 2.2: Recipe v1alpha1 API Finalized and Validated

**User Story**

> As a CNS developer,  
> I want the Recipe v1alpha1 API schema finalized and validated against MVP requirements,  
> So that recipes are consistent, extensible, and support in-cluster validation.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** the v1alpha1 schema, **When** I validate MVP recipes, **Then** all recipes pass schema validation | `cnsctl validate --schema recipe.yaml` |
| 2 | **Given** the schema includes validation constraints, **When** a recipe defines constraints, **Then** validators can evaluate cluster state against them | Unit tests for constraint evaluation |
| 3 | **Given** synthetic workload references, **When** a recipe specifies a workload, **Then** the reference resolves to a valid workload definition | Workload resolution test |
| 4 | **Given** performance characteristics, **When** a recipe defines expected metrics, **Then** validators can compare actual vs expected | Integration test with mock metrics |

**Tasks**

| Task | Description |
|------|-------------|
| Finalize Recipe schema | apiVersion, kind, metadata, spec, constraints, validation |
| Design constraint syntax | Fully-qualified paths, comparison operators, version semantics |
| Add synthetic workload reference | Schema field for workload definition reference |
| Add performance characteristics | Schema fields for expected throughput, latency, utilization |
| In-cluster validation design | Document how recipes drive cluster-state validation |
| Schema versioning strategy | Document v1alpha1 → v1beta1 → v1 migration path |

**Non-Functional Requirements**
- Schema validation completes in < 100ms
- Schema is backward-compatible within v1alpha1
- OpenAPI 3.1 compliant for tooling integration

**Dependencies**
- Validator implementation (Epic 2.4)
- Understanding of CNCF AI conformance requirements

**Assumptions**
- v1alpha1 scope is sufficient for MVP
- Schema changes before launch are acceptable

**Out of Scope**
- Automatic schema migration tooling
- GraphQL schema generation
- CRD generation (recipes are not Kubernetes CRDs)

**Definition of Done**
- [ ] Recipe v1alpha1 schema documented in api/cns/v1/
- [ ] All MVP recipes validate against schema
- [ ] Constraint syntax documented with examples
- [ ] Synthetic workload reference field implemented
- [ ] Performance characteristics fields implemented
- [ ] Schema versioning strategy documented

---

### Epic 2.3: Bundlers Generate Valid Artifacts for All MVP Recipes

**User Story**

> As a platform engineer,  
> I want CNS bundlers to generate correct, deployable artifacts for all MVP recipes,  
> So that I can deploy validated GPU configurations using my preferred method (script or ArgoCD).

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** an MVP recipe, **When** I run `cnsctl bundle -r recipe.yaml -o ./out`, **Then** all component artifacts are generated without errors | CLI exit code 0, files exist |
| 2 | **Given** generated Helm values, **When** I diff against known-good values, **Then** differences are explainable and intentional | Manual review + automated diff |
| 3 | **Given** script deployer output, **When** I run install.sh on target cluster, **Then** components deploy successfully | E2E deployment test |
| 4 | **Given** ArgoCD deployer output, **When** I apply to ArgoCD, **Then** Applications sync successfully | ArgoCD sync verification |
| 5 | **Given** bundle checksums.txt, **When** I verify checksums, **Then** all files pass integrity check | `sha256sum -c checksums.txt` |

**Tasks**

| Task | Description |
|------|-------------|
| Validate GPU Operator bundler | Test against all MVP recipes, verify values.yaml correctness |
| Validate Network Operator bundler | Test RDMA/SR-IOV configurations for H100/GB200 |
| Validate Cert-Manager bundler | Test certificate configuration generation |
| Validate Skyhook bundler | Test node optimization configurations |
| Validate NVSentinel bundler | Test monitoring configuration |
| Script deployer E2E | Deploy via scripts on EKS, GKE, AKS |
| ArgoCD deployer E2E | Deploy via ArgoCD on at least one platform |
| Document deployment steps | Step-by-step guide for each deployer method |

**Non-Functional Requirements**
- Bundle generation completes in < 30 seconds
- Generated scripts are idempotent (safe to re-run)
- ArgoCD Applications support sync-wave ordering

**Dependencies**
- MVP recipes finalized (Epic 2.1)
- Access to target clusters for E2E testing

**Assumptions**
- Helm 3.x is the deployment mechanism
- ArgoCD v2.x is the GitOps target

**Out of Scope**
- Flux deployer validation (post-MVP)
- Automated rollback on failure
- Multi-cluster deployment orchestration

**Definition of Done**
- [ ] All bundlers generate without errors for all MVP recipes
- [ ] Script deployment validated on EKS, GKE, AKS
- [ ] ArgoCD deployment validated on at least one platform
- [ ] Deployment documentation complete with screenshots
- [ ] Checksums generated and verified for all bundles

---

### Epic 2.4: Validator Ascertains Cluster State Post-Deployment

**User Story**

> As a platform operator,  
> I want to validate that my cluster matches the expected state after deployment,  
> So that I can confirm the recipe was applied correctly and catch configuration drift.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** a deployed cluster and recipe, **When** I run `cnsctl validate -r recipe.yaml -s snapshot.yaml`, **Then** constraint pass/fail status is reported | CLI output verification |
| 2 | **Given** recipe smoke tests, **When** validator runs, **Then** critical components (GPU Operator, Network Operator) are verified running | Pod status checks |
| 3 | **Given** NCCL fabric validation, **When** I run `cnsctl validate --fabric`, **Then** GPU-to-GPU communication is verified (VACE or equivalent) | NCCL test job execution |
| 4 | **Given** CNCF AI conformance requirements, **When** I run `cnsctl validate --conformance ai`, **Then** conformance report is generated | Conformance output verification |
| 5 | **Given** validation failures, **When** validator reports issues, **Then** actionable remediation steps are provided | Output includes remediation guidance |

**Tasks**

| Task | Description |
|------|-------------|
| Recipe constraint validation | Evaluate snapshot against recipe constraints |
| Smoke test framework | Verify expected pods running, services available |
| NCCL fabric validation | Deploy test job, verify multi-GPU communication |
| CNCF AI conformance integration | Map conformance requirements to validation checks |
| Remediation guidance | Document common failures and fixes |
| (Nice-to-have) In-cluster validation | CronJob-based continuous validation |

**Non-Functional Requirements**
- Validation completes in < 5 minutes
- NCCL test uses minimal GPU time
- Conformance checks are idempotent

**Dependencies**
- Recipe schema with constraints (Epic 2.2)
- CNCF AI conformance specification (external)
- Access to deployed clusters

**Assumptions**
- Snapshots accurately capture cluster state
- CNCF AI conformance spec is stable

**Out of Scope**
- Performance benchmarking (separate epic)
- Automated remediation
- Historical validation tracking

**Definition of Done**
- [ ] `cnsctl validate` evaluates all recipe constraints
- [ ] Smoke tests verify GPU Operator and Network Operator health
- [ ] Fabric validation runs NCCL test successfully
- [ ] CNCF AI conformance report generated
- [ ] Validation documentation with example outputs
- [ ] Remediation guidance for top 10 failure modes

---

## Initiative 3: Downstream Users Validated and Integrated

**Expected State:** Key downstream consumers (HIPPO, ERA) have validated CNS integration and can demonstrate end-to-end workflows.

---

### Epic 3.1: HIPPO Runtime Integration Validated

**User Story**

> As a HIPPO runtime engineer,  
> I want to integrate CNS for recipe management and bundle generation,  
> So that HIPPO can deploy validated GPU configurations with augmented settings.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** HIPPO calls CNS API, **When** requesting a recipe, **Then** recipe is returned in expected format | API integration test |
| 2 | **Given** HIPPO provides value overrides, **When** generating bundles, **Then** overrides are applied correctly | Bundle diff verification |
| 3 | **Given** HIPPO-specific configuration, **When** bundle is deployed, **Then** augmented settings are active | Cluster inspection |
| 4 | **Given** integration documentation, **When** HIPPO team reads it, **Then** API usage is clear | Documentation review by HIPPO team |

**Tasks**

| Task | Description |
|------|-------------|
| API contract alignment | Confirm CNS API meets HIPPO requirements |
| Value override testing | Test --set flag with HIPPO configurations |
| Integration test suite | Automated tests for HIPPO→CNS flow |
| Documentation | API usage guide for HIPPO integration |

**Non-Functional Requirements**
- API response time < 500ms
- Bundle generation < 30 seconds
- No breaking API changes without notice

**Dependencies**
- CNS API server deployed and accessible
- HIPPO team availability for validation

**Assumptions**
- HIPPO uses HTTP/JSON for CNS integration
- HIPPO runs outside the target cluster

**Out of Scope**
- HIPPO-specific bundler development
- HIPPO UI integration
- HIPPO authentication/authorization

**Definition of Done**
- [ ] HIPPO can request recipes via CNS API
- [ ] HIPPO can generate bundles with custom overrides
- [ ] HIPPO can deploy bundles to target cluster
- [ ] HIPPO team confirms integration works
- [ ] Integration documented in CNS docs

---

### Epic 3.2: ERA Runtime Configuration Demo

**User Story**

> As an ERA engineer,  
> I want to demonstrate CNS-driven runtime configuration for training workloads,  
> So that ERA can showcase CNS integration.

**Acceptance Criteria**

| # | Criterion | Test Method |
|---|-----------|-------------|
| 1 | **Given** ERA uses CNS recipe, **When** deploying training workload, **Then** GPU cluster is configured correctly | Training job completes successfully |
| 2 | **Given** demo scenario, **When** presented at launch, **Then** CNS value proposition is clear | Dry-run with stakeholders |
| 3 | **Given** vanilla Kubernetes cluster, **When** ERA applies CNS recipe, **Then** cluster is transformed to GPU-ready state | Before/after comparison |
| 4 | (Nice-to-have) **Given** inference workload, **When** ERA deploys via CNS, **Then** inference serving is operational | Inference endpoint responds |

**Tasks**

| Task | Description |
|------|-------------|
| Demo environment setup | EKS cluster with GPU nodes |
| Training recipe validation | Verify H100 training recipe on ERA cluster |
| Demo script development | Step-by-step demo flow |
| Dry-run with ERA team | Validate demo before acceptance |
| (Nice-to-have) Inference demo | Add inference recipe demonstration |

**Non-Functional Requirements**
- Demo completes in < 15 minutes
- Demo is reproducible without presenter intervention
- Fallback plan if live demo fails

**Dependencies**
- ERA team availability
- Demo cluster provisioned and accessible
- MVP recipes validated (Epic 2.1)

**Assumptions**
- EKS is the demo platform
- H100 nodes are available

**Out of Scope**
- Full ERA product integration
- Production ERA deployment
- Multi-cloud demo

**Definition of Done**
- [ ] Demo environment provisioned
- [ ] Training recipe deploys successfully via CNS
- [ ] Demo script reviewed and approved
- [ ] Dry-run completed with ERA team
- [ ] Demo recorded as backup

---

# Backlog

Items below have not been fully evaluated or are not critical for MVP launch. Priorities indicate relative importance for post-launch work.

---

## P1 - Expand Breadth of Recipes Supported

**Description:** Extend recipe coverage to additional platforms, accelerators, and workload types.

**Epics:**
- Prioritized backlog of recipes with their components
- Self-managed Kubernetes support
- Additional cloud providers (Oracle OCI, Alibaba Cloud)
- Additional accelerator support (L40S, future architectures)

---

## P1 - CNS v1 Playbook Migration

**Description:** Migrate validated configurations and operational knowledge from CNS v1 playbooks.

**Epics:**
- NIM Operator Bundler — Support deployment of NVIDIA Inference Microservices
- KServe Bundler — Validated inference serving configurations
- Nsight Operator Bundler — Cluster-wide profiling and observability tooling
- Monitoring Bundler — Integrate metrics, logging, and alerting components
- Storage Bundler — Best-practice storage configurations for GPU workloads

---

## P1 - Make Recipe Creation Easier

**Description:** Simplify the process of creating and contributing new recipes.

**Epics:**
- Snapshot-to-recipe transformation (generate recipe from cluster measurements)
- Recipe builder CLI (interactive recipe creation)
- Recipe contribution workflow (PR template, validation gates)

---

## P2 - Cluster Configuration Drift Detection

**Description:** Detect when deployed clusters drift from their recipe-defined state.

**Epics:**
- `cnsctl diff` command for snapshot comparison
- Scheduled drift detection via CronJob
- Alerting integration for drift events

---

## P2 - Integrated Node Configuration via Skyhook

**Description:** Deeper integration with Skyhook for OS-level node optimization.

**Epics:**
- OS-specific overlays (Ubuntu, RHEL, Amazon Linux)
- Consistent Skyhook recipe per OS (sysctl, grub, kernel parameters)
- Automated node configuration validation

---

## P3 - Expand API Support

**Description:** Additional API interfaces for programmatic integration.

**Epics:**
- gRPC API for high-performance programmatic integration
- GraphQL API for flexible querying of configuration and validation data
- Multi-tenancy support for API server

---

# Revision History

| Date | Change |
|------|--------|
| 2026-01-17 | Restructured to JIRA format with Initiatives and Epics |
| 2026-01-06 | Updated structure |
| 2026-01-05 | Added Opens section based on architectural decisions |
| 2026-01-01 | Initial comprehensive roadmap |
