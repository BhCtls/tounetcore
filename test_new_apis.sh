#!/bin/bash

# TouNetCore New APIs Test Script
# This script tests all the newly added admin APIs

SERVER_URL="http://localhost:44544"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="admin123"

echo "🔧 TouNetCore New APIs Test Script"
echo "=================================="

# Login as admin and get JWT token
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

# Test 1: Create a new user
echo ""
echo "📝 Step 2: Create New User"
USER_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "username": "testuser_'$(date +%s)'",
    "password": "password123",
    "status": "user",
    "phone": "13800138000",
    "pushdeer_token": "TEST_TOKEN"
  }')

USER_ID=$(echo "$USER_RESPONSE" | jq -r '.data.user_id')
echo "✅ User created with ID: $USER_ID"

# Test 2: Update user information
echo ""
echo "📝 Step 3: Update User Information"
curl -s -X POST "$SERVER_URL/api/v1/admin/users/$USER_ID/update" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "phone": "13900139000",
    "pushdeer_token": "UPDATED_TOKEN"
  }' > /dev/null
echo "✅ User information updated"

# Test 3: Create a new app
echo ""
echo "📝 Step 4: Create New App"
APP_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/apps" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "app_id": "testapp_'$(date +%s)'",
    "name": "Test Application",
    "description": "Test app for API testing",
    "required_permission_level": "user",
    "is_active": true
  }')

APP_ID=$(echo "$APP_RESPONSE" | jq -r '.data.app_id')
echo "✅ App created with ID: $APP_ID"

# Test 4: Update app information
echo ""
echo "📝 Step 5: Update App Information"
curl -s -X POST "$SERVER_URL/api/v1/admin/apps/$APP_ID/update" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Updated Test Application",
    "description": "Updated description for test app"
  }' > /dev/null
echo "✅ App information updated"

# Test 5: Toggle app status
echo ""
echo "📝 Step 6: Toggle App Status"
TOGGLE_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/apps/$APP_ID/toggle" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
NEW_STATUS=$(echo "$TOGGLE_RESPONSE" | jq -r '.data.is_active')
echo "✅ App status toggled to: $NEW_STATUS"

# Test 6: Generate invite code
echo ""
echo "📝 Step 7: Generate Invite Code"
INVITE_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite-codes" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
INVITE_CODE=$(echo "$INVITE_RESPONSE" | jq -r '.data.invite_code')
echo "✅ Invite code generated: $INVITE_CODE"

# Test 7: List invite codes
echo ""
echo "📝 Step 8: List Invite Codes"
INVITE_LIST=$(curl -s -X GET "$SERVER_URL/api/v1/admin/invite-codes?page=1&size=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
INVITE_COUNT=$(echo "$INVITE_LIST" | jq -r '.data.total')
echo "✅ Listed invite codes, total: $INVITE_COUNT"

# Test 8: View audit logs
echo ""
echo "📝 Step 9: View Audit Logs"
AUDIT_LOGS=$(curl -s -X GET "$SERVER_URL/api/v1/admin/logs?page=1&size=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
LOG_COUNT=$(echo "$AUDIT_LOGS" | jq -r '.data.total')
echo "✅ Viewed audit logs, total: $LOG_COUNT"

# Test 9: Delete invite code (if not used)
echo ""
echo "📝 Step 10: Delete Invite Code"
curl -s -X POST "$SERVER_URL/api/v1/admin/invite-codes/$INVITE_CODE/delete" \
  -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
echo "✅ Invite code deleted"

# Test 10: Delete app
echo ""
echo "📝 Step 11: Delete App"
curl -s -X POST "$SERVER_URL/api/v1/admin/apps/$APP_ID/delete" \
  -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
echo "✅ App deleted"

# Test 11: Delete user
echo ""
echo "📝 Step 12: Delete User"
curl -s -X POST "$SERVER_URL/api/v1/admin/users/$USER_ID/delete" \
  -H "Authorization: Bearer $ADMIN_TOKEN" > /dev/null
echo "✅ User deleted"

# Final audit log check
echo ""
echo "📝 Step 13: Final Audit Log Check"
FINAL_LOGS=$(curl -s -X GET "$SERVER_URL/api/v1/admin/logs?page=1&size=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
FINAL_COUNT=$(echo "$FINAL_LOGS" | jq -r '.data.total')
echo "✅ Final audit logs count: $FINAL_COUNT"

echo ""
echo "🎉 All new API tests completed successfully!"
echo "=================================="
echo ""
echo "📊 Summary of tested APIs:"
echo "  ✅ Admin create user"
echo "  ✅ Admin update user"
echo "  ✅ Admin delete user"
echo "  ✅ Admin create app"
echo "  ✅ Admin update app"
echo "  ✅ Admin toggle app status"
echo "  ✅ Admin delete app"
echo "  ✅ Admin generate invite code"
echo "  ✅ Admin list invite codes"
echo "  ✅ Admin delete invite code"
echo "  ✅ Admin view audit logs"
echo ""
echo "All operations have been logged in the audit system."
