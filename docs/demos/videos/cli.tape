# VHS Demo: Cloud Native Stack CLI Demo
# Run with: vhs docs/demos/videos/e2e.tape -o docs/demos/videos/e2e.gif

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
Type "# Cloud Native Stack: CLI Demo"
Enter
Sleep 2s
Ctrl+l
Sleep 500ms

# Intro
Type "# CLI help"
Enter
Sleep 500ms
Type "cnsctl -h"
Enter
Sleep 3s
Ctrl+l
Sleep 500ms

# Recipe from CLI
Type "# Recipe from CLI"
Enter
Sleep 500ms
Type "cnsctl recipe \"
Enter 
Type "  --service eks \"
Enter 
Type "  --accelerator gb200 \"
Enter 
Type "  --os ubuntu \"
Enter 
Type "  --intent training \"
Enter 
Type "  --output recipe.yaml"
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
