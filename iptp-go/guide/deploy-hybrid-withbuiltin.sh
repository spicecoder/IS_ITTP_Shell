#!/bin/bash
# deploy.sh - Bash script with gobash built-ins

# Standard bash
echo "Starting deployment..."

# gobash built-in: check field
if gobash field-match "deployment_active=Y"; then
    echo "❌ Deployment already active"
    exit 1
fi

# gobash built-in: emit pulse
gobash emit-pulse \
    --name "deployment_active" \
    --tv "Y" \
    --response "$(whoami)@$(date +%s)"

# gobash built-in: emit intention
gobash emit-intention \
    --name "notify_team" \
    --target "SlackNotifier" \
    --signal '[
        {"name": "deployment_started", "TV": "Y", "response": "v2.3"}
    ]'

# Standard bash
npm run build
npm run test

if [ $? -eq 0 ]; then
    # Success
    npm run deploy
    
    # gobash built-in: update field
    gobash emit-pulse \
        --name "deployment_active" \
        --tv "N"
    
    gobash emit-pulse \
        --name "deployment_complete" \
        --tv "Y" \
        --response "success@$(date)"
else
    # Failure
    gobash emit-pulse \
        --name "deployment_failed" \
        --tv "Y" \
        --response "tests failed"
    
    exit 1
fi

echo "✅ Deployment complete!"
