# TouNetCore Permission System - Quick Reference

## Permission Hierarchy
```
admin (4) > trusted (3) > user (2) > disableduser (1)
```

## Key Endpoints

| Action | Method | Endpoint | Auth Level |
|--------|--------|----------|------------|
| Create User | POST | `/api/v1/admin/users` | Admin |
| Update User Permissions | POST | `/api/v1/admin/users/{id}/update` | Admin |
| Get Current User | GET | `/api/v1/user/me` | User+ |
| List User Apps | GET | `/api/v1/user/apps` | User+ |
| Generate NKey | POST | `/api/v1/nkey/generate` | User+ |
| Validate NKey | POST | `/api/v1/nkey/validate` | Public |
| List All Apps | GET | `/api/v1/admin/apps` | Admin |
| Create App | POST | `/api/v1/admin/apps` | Admin |

## JavaScript Utils

```javascript
// Permission checking utility
const hasPermission = (userLevel, requiredLevel) => {
  const levels = { disableduser: 1, user: 2, trusted: 3, admin: 4 };
  return levels[userLevel] >= levels[requiredLevel];
};

// Promote user to trusted
const promoteToTrusted = async (userId, adminToken) => {
  await fetch(`/api/v1/admin/users/${userId}/update`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({ status: 'trusted' })
  });
};

// Generate NKey for apps
const generateNKey = async (userToken, appIds) => {
  const response = await fetch('/api/v1/nkey/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${userToken}`
    },
    body: JSON.stringify({ app_ids: appIds })
  });
  return response.json();
};
```

## React Components

```jsx
// Permission Guard Component
const PermissionGuard = ({ userLevel, requiredLevel, children }) => {
  const hasAccess = hasPermission(userLevel, requiredLevel);
  return hasAccess ? children : <div>Access Denied</div>;
};

// Usage
<PermissionGuard userLevel={user.status} requiredLevel="trusted">
  <TrustedFeature />
</PermissionGuard>
```

## Security Rules

✅ **DO:**
- New users are always created as "user" level
- Promote users through update endpoint only
- Validate permissions on both frontend and backend
- Show clear permission indicators to users

❌ **DON'T:**
- Don't create users directly with admin/trusted status
- Don't allow self-demotion for admins
- Don't trust frontend-only permission checks
- Don't show features users can't access

## Error Codes

| Code | Message | Meaning |
|------|---------|---------|
| 400 | cannot demote your own admin privileges | Self-demotion attempt |
| 403 | insufficient permissions | User lacks required permission |
| 403 | no permission for app | User cannot access specific app |
| 404 | user not found | Invalid user ID |
| 409 | username already exists | Duplicate username |

## Pre-configured Apps

| App ID | Name | Level Required |
|--------|------|----------------|
| `searchall` | Global Search | user |
| `segaasstes` | Sega Assets | user |
| `advanced_analytics` | Advanced Analytics | trusted |
| `livecontent_admin` | Live Content Admin | admin |

## Common Patterns

```javascript
// Check if user can access app
const canAccessApp = (userLevel, appRequiredLevel) => {
  return hasPermission(userLevel, appRequiredLevel);
};

// Filter apps by user permission
const getAccessibleApps = (allApps, userLevel) => {
  return allApps.filter(app => 
    hasPermission(userLevel, app.required_permission_level)
  );
};

// Handle permission upgrade
const handlePermissionChange = async (userId, newLevel) => {
  try {
    await updateUserPermission(userId, newLevel);
    alert(`User promoted to ${newLevel}. They need to log in again.`);
  } catch (error) {
    alert('Permission update failed: ' + error.message);
  }
};
```
