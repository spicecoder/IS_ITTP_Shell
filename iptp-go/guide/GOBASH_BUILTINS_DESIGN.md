# gobash Built-in Commands Design

## Philosophy
Keep bash compatibility, add IPTP-aware built-ins that scripts can use.

## Built-in Commands

### 1. field-query
Query the current field state.

```bash
# Check if pulse exists with TV
gobash field-query "deployment_active=Y"
# Exit code 0 if match, 1 if no match

# Get pulse value
VALUE=$(gobash field-query "deployment_active" --get-response)
echo "Deployment by: $VALUE"

# Get all pulses for process
gobash field-query --process "auth" --json
# Output: [{"name": "...", "TV": "Y", "response": "..."}]
```

### 2. emit-pulse
Emit a pulse to the field.

```bash
# Simple emit
gobash emit-pulse "deployment_active" "Y"

# With response
gobash emit-pulse "deployment_active" "Y" --response "alice@prod"

# From JSON
gobash emit-pulse --json '{"name": "test_passed", "TV": "Y", "response": "100%"}'
```

### 3. emit-intention
Send an intention to another process/object.

```bash
# Emit intention with signal
gobash emit-intention \
    --name "notify_team" \
    --target "SlackNotifier" \
    --signal '[{"name": "deployment_complete", "TV": "Y"}]'

# From file
gobash emit-intention --from-file notification.json
```

### 4. field-match
Wait for field conditions to be satisfied.

```bash
# Wait for condition (blocking)
gobash field-match "preprocessing_done=Y" --wait

# Wait with timeout
gobash field-match "data_ready=Y" --wait --timeout 300

# Check multiple conditions
gobash field-match "auth_ready=Y,db_ready=Y" --all
```

### 5. process-list
List active processes and their states.

```bash
# List all processes
gobash process-list

# List with field state
gobash process-list --show-field

# Filter by pulse
gobash process-list --pulse "deployment_active=Y"

# JSON output
gobash process-list --json
```

### 6. field-absorb
Manually absorb signals into field (advanced).

```bash
# Absorb from JSON
gobash field-absorb --json '[
    {"name": "test_passed", "TV": "Y"},
    {"name": "coverage", "TV": "Y", "response": "95%"}
]'

# Absorb from file
gobash field-absorb --from-file results.json
```

### 7. trigger-check
Check if trigger conditions would activate.

```bash
# Check trigger
gobash trigger-check --intention "deploy_to_prod" --conditions '[
    {"name": "tests_passed", "TV": "Y"},
    {"name": "approved", "TV": "Y"}
]'

# Output: "READY" or "WAITING: tests_passed=N"
```

### 8. field-watch
Monitor field changes (for debugging).

```bash
# Watch all field changes
gobash field-watch

# Watch specific pulses
gobash field-watch --pulse "deployment_active,test_status"

# Output format:
# [2025-01-15 10:30:00] deployment_active: N -> Y (alice@prod)
# [2025-01-15 10:32:15] test_status: U -> Y (passed)
```

## Example: Complete Deployment Script

```bash
#!/bin/bash
# smart-deploy.sh - IPTP-aware deployment

set -e

DEPLOYMENT="deploy_$(date +%s)"

echo "🚀 Smart Deployment: $DEPLOYMENT"

# 1. Check if another deployment is active
if gobash field-query "deployment_active=Y"; then
    ACTIVE=$(gobash field-query "deployment_active" --get-response)
    echo "❌ Deployment already active: $ACTIVE"
    exit 1
fi

# 2. Wait for prerequisites
echo "⏳ Waiting for prerequisites..."
gobash field-match "tests_passed=Y,approval_given=Y" --wait --timeout 600 || {
    echo "❌ Prerequisites not met"
    exit 1
}

# 3. Claim deployment slot
echo "🔒 Claiming deployment slot..."
gobash emit-pulse "deployment_active" "Y" --response "$USER@$DEPLOYMENT"
trap "gobash emit-pulse 'deployment_active' 'N'" EXIT

# 4. Notify team
gobash emit-intention \
    --name "notify_team" \
    --target "SlackNotifier" \
    --signal '[{"name": "deployment_started", "TV": "Y", "response": "'$DEPLOYMENT'"}]'

# 5. Run deployment
echo "📦 Building..."
npm run build

echo "🧪 Testing..."
npm run test

echo "🚀 Deploying..."
npm run deploy

# 6. Update field with results
gobash emit-pulse "deployment_active" "N"
gobash emit-pulse "last_deployment" "Y" --response "$DEPLOYMENT@success@$(date)"

echo "✅ Deployment complete: $DEPLOYMENT"

# 7. Notify success
gobash emit-intention \
    --name "notify_team" \
    --target "SlackNotifier" \
    --signal '[{"name": "deployment_complete", "TV": "Y", "response": "'$DEPLOYMENT'"}]'
```

## Implementation in Go

```go
// In commands.go, add built-in command handling

func (sh *Shell) executeCommand(line string) {
    parts := strings.Fields(line)
    if len(parts) == 0 {
        return
    }

    cmd := parts[0]
    args := parts[1:]

    switch cmd {
    case "cd":
        sh.cmdGoto(args)
    
    // ... existing cases ...
    
    // New built-ins for IPTP
    case "field-query":
        sh.cmdFieldQuery(args)
    case "emit-pulse":
        sh.cmdEmitPulse(args)
    case "emit-intention":
        sh.cmdEmitIntention(args)
    case "field-match":
        sh.cmdFieldMatch(args)
    case "process-list":
        sh.cmdProcessList(args)
    case "field-absorb":
        sh.cmdFieldAbsorb(args)
    case "trigger-check":
        sh.cmdTriggerCheck(args)
    case "field-watch":
        sh.cmdFieldWatch(args)
    
    default:
        sh.cmdExec(parts)
    }
}
```

## Benefits of This Approach

1. ✅ **Bash Compatible**: All existing scripts work
2. ✅ **IPTP Enhanced**: Scripts can use field/intentions
3. ✅ **Gradual Adoption**: Start with bash, add gobash commands as needed
4. ✅ **No New Language**: Developers know bash already
5. ✅ **Composable**: Mix bash and gobash commands freely

## Future: Optional .gsh Language

Once the built-ins are solid, you could add .gsh as a **compiled layer**:

```go
// hello.gsh compiles to:
#!/bin/bash
for i in {0..9}; do
    echo "Iteration $i"
done

PIDS=$(ps aux | grep nginx | awk '{print $2}')
echo "PIDs: $PIDS"
```

But that's **Phase 2**. Phase 1 is the built-in commands.
