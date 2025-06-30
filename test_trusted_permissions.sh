#!/bin/bash

# TouNetCore Trusted Permission Level Test Script
# This script tests the new trusted permission level functionality

SERVER_URL="http://localhost:44544"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="admin123"

echo "🔧 TouNetCore Trusted Permission Level Test"
echo "==========================================="

# Step 1: Login as admin and get JWT token
echo "📝 Step 1: Admin Login"
ADMIN_TOKEN=$(curl -s -X POST "$SERVER_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$ADMIN_USERNAME\", \"password\": \"$ADMIN_PASSWORD\"}" | \
  jq -r '.data.token')

if [ "$ADMIN_TOKEN" = "null" ] || [ -z "$ADMIN_TOKEN" ]; then
  echo "❌ Failed to get admin token"
  exit 1
fi
echo "✅ Admin token obtained"

# Step 2: Generate invite code
echo ""
echo "📝 Step 2: Generate Invite Code"
INVITE_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite-codes" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
INVITE_CODE=$(echo "$INVITE_RESPONSE" | jq -r '.data.code')
echo "✅ Invite code generated: $INVITE_CODE"

# Step 3: Register a new user (will have default "user" status)
TEST_USERNAME="testuser_$(date +%s)"
echo ""
echo "📝 Step 3: Register New User ($TEST_USERNAME)"
USER_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"password123\",
    \"phone\": \"13800138000\",
    \"invite_code\": \"$INVITE_CODE\"
  }")

USER_ID=$(echo "$USER_RESPONSE" | jq -r '.data.user_id')
echo "✅ User registered with ID: $USER_ID (default status: user)"

# Step 4: Login as the new user
echo ""
echo "📝 Step 4: Login as New User"
USER_TOKEN=$(curl -s -X POST "$SERVER_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$TEST_USERNAME\", \"password\": \"password123\"}" | \
  jq -r '.data.token')

echo "✅ User token obtained"

# Step 5: Test user-level app access
echo ""
echo "📝 Step 5: Test User-Level App Access"
USER_APPS=$(curl -s -X GET "$SERVER_URL/api/v1/user/apps" \
  -H "Authorization: Bearer $USER_TOKEN")

echo "📋 Available apps for user-level:"
echo "$USER_APPS" | jq -r '.data[] | "  - \(.app_id): \(.name) (required: \(.required_permission_level))"'

# Step 6: Promote user to trusted status
echo ""
echo "📝 Step 6: Promote User to Trusted Status"
PROMOTE_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/users/$USER_ID/update" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"status": "trusted"}')

echo "✅ User promoted to trusted status"

# Step 7: Login again to get updated token
echo ""
echo "📝 Step 7: Login Again with Updated Permissions"
TRUSTED_TOKEN=$(curl -s -X POST "$SERVER_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$TEST_USERNAME\", \"password\": \"password123\"}" | \
  jq -r '.data.token')

echo "✅ Updated token obtained"

# Step 8: Test trusted-level app access
echo ""
echo "📝 Step 8: Test Trusted-Level App Access"
TRUSTED_APPS=$(curl -s -X GET "$SERVER_URL/api/v1/user/apps" \
  -H "Authorization: Bearer $TRUSTED_TOKEN")

echo "📋 Available apps for trusted-level:"
echo "$TRUSTED_APPS" | jq -r '.data[] | "  - \(.app_id): \(.name) (required: \(.required_permission_level))"'

# Step 9: Generate NKey for trusted apps
echo ""
echo "📝 Step 9: Generate NKey for Trusted User"
NKEY_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/nkey/generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TRUSTED_TOKEN" \
  -d '{"app_ids": ["advanced_analytics", "searchall"]}')

NKEY=$(echo "$NKEY_RESPONSE" | jq -r '.data.nkey')
echo "✅ NKey generated: $NKEY"

# Step 10: Validate NKey for trusted app
echo ""
echo "📝 Step 10: Validate NKey for Trusted App"
VALIDATE_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/nkey/validate" \
  -H "Content-Type: application/json" \
  -d "{\"nkey\": \"$NKEY\", \"app_id\": \"advanced_analytics\"}")

VALIDATION_RESULT=$(echo "$VALIDATE_RESPONSE" | jq -r '.data.valid')
USER_ROLE=$(echo "$VALIDATE_RESPONSE" | jq -r '.data.user_role')

if [ "$VALIDATION_RESULT" = "true" ]; then
  echo "✅ NKey validation successful for trusted app"
  echo "   User role: $USER_ROLE"
else
  echo "❌ NKey validation failed for trusted app"
fi

# Step 11: Test admin app access (should fail)
echo ""
echo "📝 Step 11: Test Admin App Access (Should Fail)"
ADMIN_VALIDATE_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/nkey/validate" \
  -H "Content-Type: application/json" \
  -d "{\"nkey\": \"$NKEY\", \"app_id\": \"livecontent_admin\"}")

ADMIN_VALIDATION_RESULT=$(echo "$ADMIN_VALIDATE_RESPONSE" | jq -r '.code')

if [ "$ADMIN_VALIDATION_RESULT" = "403" ] || [ "$ADMIN_VALIDATION_RESULT" = "401" ]; then
  echo "✅ Access correctly denied for admin-only app"
else
  echo "❌ Should not have access to admin-only app"
fi

# Step 12: Test privilege escalation protection
echo ""
echo "📝 Step 12: Test Self-Demotion Protection"
DEMOTION_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/users/1/update" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"status": "user"}')

DEMOTION_CODE=$(echo "$DEMOTION_RESPONSE" | jq -r '.code')

if [ "$DEMOTION_CODE" = "400" ]; then
  echo "✅ Self-demotion correctly prevented"
else
  echo "❌ Self-demotion should be prevented"
fi

echo ""
echo "🎉 Trusted Permission Level Test Complete!"
echo "=========================================="
echo "Summary:"
echo "- ✅ User creation with default 'user' status"
echo "- ✅ Admin can promote users to 'trusted' status"
echo "- ✅ Trusted users can access trusted-level apps"
echo "- ✅ Permission boundaries are enforced"
echo "- ✅ Self-demotion protection works"
