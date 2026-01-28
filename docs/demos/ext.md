# CNS Extendability Demo

## Install

```shell
curl -sfL https://raw.githubusercontent.com/mchmarny/cloud-native-stack/main/install | bash -s --
```

Test CLI:

```shell
cnsctl -h
```

## Embedded Data

View embedded data files structure:

```shell
tree -L 2 pkg/recipe/data/
```

## Runtime Data Support

Generate recipe with external data:

```shell
cnsctl recipe \
  --service eks \
  --accelerator gb200 \
  --os ubuntu \
  --intent training \
  --data ./examples/data \
  --output recipe.yaml
```

Output shows:
* `7` embedded + `1` external = `8` merged components
* `dgxc-teleport` appears as Kustomize component
* Included in `deploymentOrder`

Now generate bundles:

```shell
cnsctl bundle \
  --recipe recipe.yaml \
  --output ./bundle \
  --data ./examples/data \
  --deployer argocd \
  --output oci://ghcr.io/mchmarny/cns-bundle \
  --system-node-selector nodeGroup=system-pool \
  --accelerated-node-selector nodeGroup=customer-gpu \
  --accelerated-node-toleration nvidia.com/gpu=present:NoSchedule
```

### Debug Mode

The `--debug` flag shows which files are loaded from external vs embedded sources:

```bash
cnsctl --debug recipe \
  --service eks \
  --accelerator gb200 \
  --data ./examples/data
```

## Links

* [Installation Guide](https://github.com/mchmarny/cloud-native-stack/blob/main/docs/user-guide/installation.md)
* [CLI Reference](https://github.com/mchmarny/cloud-native-stack/blob/main/docs/user-guide/cli-reference.md)
* [Data Reference](https://github.com/mchmarny/cloud-native-stack/blob/main/pkg/recipe/data/README.md)
