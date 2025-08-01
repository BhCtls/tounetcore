#!/bin/bash

# TouNetCore API Test Script
# Tests the new URL functionality

BASE_URL="http://localhost:44544/api/v1"
echo "Testing TouNetCore API with URL field support..."

# Test health check
echo "1. Testing health check..."
curl -s http://localhost:44544/health | jq .

echo -e "\n2. Note: To test the full API functionality, you need to:"
echo "   - Register a user with an invite code"
echo "   - Login to get JWT token"
echo "   - Use admin permissions to create/update apps"
echo "   - Test user app listing"

echo -e "\n3. Example API calls (replace {token} with actual JWT):"
echo ""
echo "Create app with URL:"
echo 'curl -X POST "$BASE_URL/admin/apps" \'
echo '  -H "Authorization: Bearer {token}" \'
echo '  -H "Content-Type: application/json" \'
echo '  -d "{"app_id":"test_app","name":"Test App","url":"https://example.com","description":"Test app with URL"}"'

echo ""
echo "Update app URL:"
echo 'curl -X PUT "$BASE_URL/admin/apps/test_app" \'
echo '  -H "Authorization: Bearer {token}" \'
echo '  -H "Content-Type: application/json" \'
echo '  -d "{"url":"https://updated-example.com"}"'

echo ""
echo "Get user apps (with URL field):"
echo 'curl -X GET "$BASE_URL/user/apps" \'
echo '  -H "Authorization: Bearer {token}"'

echo ""
echo "Generate NKey for app:"
echo 'curl -X POST "$BASE_URL/nkey/generate" \'
echo '  -H "Authorization: Bearer {token}" \'
echo '  -H "Content-Type: application/json" \'
echo '  -d "{"app_ids":["test_app"]}"'

echo -e "\nTest script completed. Server is running on http://localhost:44544"
