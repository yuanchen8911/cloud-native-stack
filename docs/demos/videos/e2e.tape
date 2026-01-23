# VHS Demo: Cloud Native Stack End-to-End Workflow
# Based on docs/demos/tape-script.md
#
# Run with: vhs docs/demos/videos/e2e.tape -o docs/demos/videos/e2e.gif

Output docs/demos/videos/e2e.gif

# Terminal settings
Set Shell "bash"
Set FontSize 18
Set Width 1400
Set Height 850
Set Padding 10
Set WindowBar Colorful
Set Margin 10
Set MarginFill "#674EFF"
Set BorderRadius 5

# Animation settings
Set TypingSpeed 50ms
Set Framerate 30
Set LoopOffset 0%

# Welcome
Type "# Cloud Native Stack: End-to-End Demo"
Enter
Sleep 2s
Ctrl+l
Sleep 500ms

# Recipe from CLI
Type "# Recipe from CLI"
Enter
Sleep 500ms
Type "cnsctl recipe --service eks --intent training"
Enter
Sleep 3s
Ctrl+l
Sleep 500ms

# Recipe from API
Type "# Recipe from API"
Enter
Sleep 500ms
Type 'curl -s "https://cns.dgxc.io/v1/recipe?service=eks&accelerator=gb200&intent=training" | jq .'
Enter
Sleep 3s
Ctrl+l
Sleep 500ms

# Make Snapshot (deploy agent)
Type "# Make Snapshot from Cluster"
Enter
Sleep 500ms
Type "cnsctl snapshot \"
Enter
Type "    --deploy-agent \"
Enter
Type "    --namespace gpu-operator \"
Enter
Type "    --image ghcr.io/mchmarny/cns:latest \"
Enter
Type "    --node-selector nodeGroup=customer-gpu"
Enter
Sleep 8s
Ctrl+l
Sleep 500ms

# Check Snapshot in ConfigMap
Type "# Check Snapshot in ConfigMap"
Enter
Sleep 500ms
Type "kubectl -n gpu-operator get cm cns-snapshot -o jsonpath='{.data.snapshot\.yaml}' | yq ."
Enter
Sleep 5s
Ctrl+l
Sleep 500ms

# Recipe from Snapshot
Type "# Recipe from Snapshot"
Enter
Sleep 500ms
Type "cnsctl recipe \"
Enter
Type "  --snapshot cm://gpu-operator/cns-snapshot \"
Enter
Type "  --intent training \"
Enter
Type "  --output recipe.yaml"
Enter
Sleep 5s
Ctrl+l
Sleep 500ms

# Recipe Constraints
Type "# Recipe Constraints"
Enter
Sleep 500ms
Type "yq .constraints recipe.yaml"
Enter
Sleep 2s
Ctrl+l
Sleep 500ms

# Validate Recipe
Type "# Validate Constraints against Cluster"
Enter
Sleep 500ms
Type "cnsctl validate \"
Enter
Type "  --recipe recipe.yaml \"
Enter
Type "  --snapshot cm://gpu-operator/cns-snapshot | yq ."
Enter
Sleep 5s
Ctrl+l
Sleep 500ms

# Bundle from Recipe
Type "# Create Bundle from Recipe"
Enter
Sleep 500ms
Type "cnsctl bundle \"
Enter
Type "  --recipe recipe.yaml \"
Enter
Type "  --output ./bundles \"
Enter
Type "  --system-node-selector nodeGroup=system-pool \"
Enter
Type "  --accelerated-node-selector nodeGroup=customer-gpu \"
Enter
Type "  --accelerated-node-toleration nvidia.com/gpu=present:NoSchedule"
Enter
Sleep 7s
Ctrl+l
Sleep 500ms

# Bundle from API
Type "# Bundle from API"
Enter
Sleep 500ms
Type 'curl -s "https://cns.dgxc.io/v1/recipe?service=eks&accelerator=h100&intent=training" | \'
Enter
Type '  curl -X POST "https://cns.dgxc.io/v1/bundle" \'
Enter
Type '    -H "Content-Type: application/json" -d @- -o bundles.zip'
Enter
Sleep 2s
Ctrl+l
Sleep 500ms

# Chart dependencies
Type "# Review Bundle and validate chart dependencies"
Enter
Sleep 500ms
Type "cd ./bundles && tree ."
Enter
Sleep 1s
Type "helm dependency update"
Enter
Sleep 13s
Ctrl+l
Sleep 500ms

# Done
Type "# Demo Complete!"
Enter
Sleep 5s
Hide
