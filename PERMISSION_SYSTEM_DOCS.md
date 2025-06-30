# TouNetCore Permission System Documentation

## Overview

TouNetCore implements a hierarchical permission system with four distinct user levels and application-based access control. This document provides frontend developers with the necessary information to implement proper permission handling in their applications.

## User Permission Levels

The system uses a four-tier hierarchical permission structure:

```
admin > trusted > user > disableduser
```

### Permission Level Details

| Level | Value | Numeric Priority | Description |
|-------|-------|-----------------|-------------|
| **Admin** | `"admin"` | 4 | Full system access, can manage users, apps, and system settings |
| **Trusted** | `"trusted"` | 3 | Enhanced user with access to trusted applications and features |
| **User** | `"user"` | 2 | Standard user access to basic applications |
| **Disabled User** | `"disableduser"` | 1 | Minimal access, restricted to specific applications |

## API Endpoints for Permission Management

### 1. User Creation (Admin Only)

**Endpoint:** `POST /api/v1/admin/users`

**Security Note:** All new users are created with `"user"` status by default. Permission elevation must be done through the update endpoint.

```javascript
// Create a new user
const createUser = async (userData) => {
  const response = await fetch('/api/v1/admin/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({
      username: userData.username,
      password: userData.password,
      phone: userData.phone,
      pushdeer_token: userData.pushDeerToken
      // Note: status is automatically set to "user"
    })
  });
  
  return response.json();
};
```

### 2. Permission Elevation (Admin Only)

**Endpoint:** `POST /api/v1/admin/users/{user_id}/update`

```javascript
// Promote user to trusted status
const promoteUserToTrusted = async (userId) => {
  const response = await fetch(`/api/v1/admin/users/${userId}/update`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({
      status: 'trusted'
    })
  });
  
  return response.json();
};

// Promote user to admin status
const promoteUserToAdmin = async (userId) => {
  const response = await fetch(`/api/v1/admin/users/${userId}/update`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({
      status: 'admin'
    })
  });
  
  return response.json();
};

// Demote user (e.g., from trusted to user)
const demoteUser = async (userId, newStatus) => {
  const response = await fetch(`/api/v1/admin/users/${userId}/update`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({
      status: newStatus // 'user' or 'disableduser'
    })
  });
  
  return response.json();
};
```

### 3. Get User Information

**Endpoint:** `GET /api/v1/user/me`

```javascript
// Get current user's information including permission level
const getCurrentUser = async (userToken) => {
  const response = await fetch('/api/v1/user/me', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${userToken}`
    }
  });
  
  const data = await response.json();
  return data.data; // Contains user info including 'status' field
};
```

### 4. List Available Applications

**Endpoint:** `GET /api/v1/user/apps`

```javascript
// Get applications available to current user based on their permission level
const getUserApps = async (userToken) => {
  const response = await fetch('/api/v1/user/apps', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${userToken}`
    }
  });
  
  const data = await response.json();
  return data.data; // Array of available applications
};
```

## Application Permission Levels

Applications are assigned required permission levels that determine which users can access them:

### Pre-configured Applications

| Application ID | Name | Required Level | Description |
|---------------|------|---------------|-------------|
| `searchall` | Global Search | `user` | Basic search functionality |
| `segaasstes` | Sega Assets | `user` | Asset management |
| `dxprender` | DXP Render | `user` | Rendering service |
| `CardPreview` | Card Preview | `user` | Card preview functionality |
| `livecontent_basic` | Live Content Basic | `user` | Basic content access |
| `advanced_analytics` | Advanced Analytics | `trusted` | Analytics and reporting |
| `livecontent_admin` | Live Content Admin | `admin` | Administrative content access |

### Creating Applications with Permission Levels

**Endpoint:** `POST /api/v1/admin/apps`

```javascript
// Create an application with specific permission requirements
const createApp = async (appData) => {
  const response = await fetch('/api/v1/admin/apps', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({
      app_id: appData.appId,
      name: appData.name,
      description: appData.description,
      required_permission_level: appData.requiredLevel, // 'user', 'trusted', or 'admin'
      is_active: true
    })
  });
  
  return response.json();
};
```

## NKey Authorization System

The system uses temporary authorization keys (NKeys) that expire after 15 minutes for secure application access.

### Generate NKey

**Endpoint:** `POST /api/v1/nkey/generate`

```javascript
// Generate NKey for specific applications
const generateNKey = async (userToken, appIds) => {
  const response = await fetch('/api/v1/nkey/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${userToken}`
    },
    body: JSON.stringify({
      app_ids: appIds // Array of application IDs
    })
  });
  
  return response.json();
};
```

### Validate NKey

**Endpoint:** `POST /api/v1/nkey/validate`

```javascript
// Validate NKey for application access
const validateNKey = async (nkey, appId) => {
  const response = await fetch('/api/v1/nkey/validate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      nkey: nkey,
      app_id: appId
    })
  });
  
  return response.json();
};
```

## Frontend Permission Handling

### Permission Check Utility Functions

```javascript
// Utility functions for permission checking
const PermissionUtils = {
  // Permission level hierarchy
  PERMISSION_LEVELS: {
    'disableduser': 1,
    'user': 2,
    'trusted': 3,
    'admin': 4
  },
  
  // Check if user has required permission level
  hasPermission: (userLevel, requiredLevel) => {
    return PermissionUtils.PERMISSION_LEVELS[userLevel] >= 
           PermissionUtils.PERMISSION_LEVELS[requiredLevel];
  },
  
  // Check if user can access specific application
  canAccessApp: (userLevel, appRequiredLevel) => {
    return PermissionUtils.hasPermission(userLevel, appRequiredLevel);
  },
  
  // Get permission level name
  getPermissionLevelName: (level) => {
    const names = {
      'admin': 'Administrator',
      'trusted': 'Trusted User',
      'user': 'Standard User',
      'disableduser': 'Disabled User'
    };
    return names[level] || 'Unknown';
  }
};
```

### React Component Examples

```jsx
// Permission-based component rendering
import React, { useState, useEffect } from 'react';

const PermissionGuard = ({ requiredLevel, userLevel, children, fallback }) => {
  const hasAccess = PermissionUtils.hasPermission(userLevel, requiredLevel);
  
  return hasAccess ? children : (fallback || <div>Access Denied</div>);
};

// Application list component
const ApplicationList = ({ userToken }) => {
  const [apps, setApps] = useState([]);
  const [userLevel, setUserLevel] = useState('');
  
  useEffect(() => {
    const fetchData = async () => {
      // Get user info
      const userInfo = await getCurrentUser(userToken);
      setUserLevel(userInfo.status);
      
      // Get available apps
      const userApps = await getUserApps(userToken);
      setApps(userApps);
    };
    
    fetchData();
  }, [userToken]);
  
  return (
    <div>
      <h2>Available Applications</h2>
      <p>Your Permission Level: {PermissionUtils.getPermissionLevelName(userLevel)}</p>
      
      {apps.map(app => (
        <div key={app.app_id} className="app-card">
          <h3>{app.name}</h3>
          <p>Status: {app.enabled ? 'Enabled' : 'Disabled'}</p>
          {app.valid_until && (
            <p>Valid Until: {new Date(app.valid_until).toLocaleDateString()}</p>
          )}
        </div>
      ))}
    </div>
  );
};

// Admin user management component
const UserManagement = ({ adminToken }) => {
  const [users, setUsers] = useState([]);
  
  const handlePromoteUser = async (userId, newStatus) => {
    try {
      if (newStatus === 'trusted') {
        await promoteUserToTrusted(userId);
      } else if (newStatus === 'admin') {
        await promoteUserToAdmin(userId);
      } else {
        await demoteUser(userId, newStatus);
      }
      
      // Refresh user list
      fetchUsers();
      alert('User permission updated successfully');
    } catch (error) {
      alert('Error updating user permission: ' + error.message);
    }
  };
  
  return (
    <div>
      <h2>User Management</h2>
      {users.map(user => (
        <div key={user.id} className="user-card">
          <h3>{user.username}</h3>
          <p>Current Level: {PermissionUtils.getPermissionLevelName(user.status)}</p>
          
          <div className="permission-controls">
            <button onClick={() => handlePromoteUser(user.id, 'trusted')}>
              Promote to Trusted
            </button>
            <button onClick={() => handlePromoteUser(user.id, 'admin')}>
              Promote to Admin
            </button>
            <button onClick={() => handlePromoteUser(user.id, 'user')}>
              Demote to User
            </button>
            <button onClick={() => handlePromoteUser(user.id, 'disableduser')}>
              Disable User
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};
```

## Security Considerations

### 1. Self-Demotion Protection

Administrators cannot demote their own privileges to prevent accidental lockouts:

```javascript
// This will return an error for self-demotion attempts
{
  "code": 400,
  "message": "cannot demote your own admin privileges"
}
```

### 2. Permission Validation

Always validate permissions on both frontend and backend:

```javascript
// Frontend permission check before showing UI elements
const shouldShowAdminPanel = PermissionUtils.hasPermission(userLevel, 'admin');

// Backend automatically validates permissions for all protected endpoints
```

### 3. Token Refresh

After permission changes, users need to log in again to get updated tokens:

```javascript
const handlePermissionChange = async () => {
  // After changing user permissions
  alert('Permission updated. Please log in again to access new features.');
  
  // Optionally force logout
  localStorage.removeItem('authToken');
  window.location.href = '/login';
};
```

## Error Handling

### Common Error Responses

```javascript
// Insufficient permissions
{
  "code": 403,
  "message": "insufficient permissions"
}

// User not found
{
  "code": 404,
  "message": "user not found"
}

// Invalid permission level
{
  "code": 400,
  "message": "invalid request data"
}

// Self-demotion attempt
{
  "code": 400,
  "message": "cannot demote your own admin privileges"
}
```

## Best Practices

### 1. Progressive Permission Disclosure

Show users only the features they have access to:

```javascript
const FeatureMenu = ({ userLevel }) => (
  <nav>
    <a href="/dashboard">Dashboard</a>
    
    <PermissionGuard requiredLevel="trusted" userLevel={userLevel}>
      <a href="/analytics">Advanced Analytics</a>
    </PermissionGuard>
    
    <PermissionGuard requiredLevel="admin" userLevel={userLevel}>
      <a href="/admin">Administration</a>
    </PermissionGuard>
  </nav>
);
```

### 2. Clear Permission Indicators

Always show users their current permission level:

```javascript
const UserStatus = ({ userLevel }) => (
  <div className={`user-status status-${userLevel}`}>
    <span className="status-badge">
      {PermissionUtils.getPermissionLevelName(userLevel)}
    </span>
  </div>
);
```

### 3. Graceful Permission Errors

Handle permission errors gracefully:

```javascript
const handleApiCall = async (apiFunction) => {
  try {
    return await apiFunction();
  } catch (error) {
    if (error.status === 403) {
      showNotification('You do not have permission to perform this action.', 'error');
    } else if (error.status === 401) {
      showNotification('Please log in again.', 'warning');
      redirectToLogin();
    } else {
      showNotification('An error occurred. Please try again.', 'error');
    }
  }
};
```

## Testing Scenarios

### 1. Permission Level Testing

```javascript
// Test different permission levels
const testPermissions = async () => {
  const testCases = [
    { userLevel: 'user', appLevel: 'user', shouldPass: true },
    { userLevel: 'user', appLevel: 'trusted', shouldPass: false },
    { userLevel: 'trusted', appLevel: 'user', shouldPass: true },
    { userLevel: 'trusted', appLevel: 'trusted', shouldPass: true },
    { userLevel: 'trusted', appLevel: 'admin', shouldPass: false },
    { userLevel: 'admin', appLevel: 'admin', shouldPass: true }
  ];
  
  testCases.forEach(test => {
    const result = PermissionUtils.hasPermission(test.userLevel, test.appLevel);
    console.assert(result === test.shouldPass, 
      `Permission test failed: ${test.userLevel} -> ${test.appLevel}`);
  });
};
```

### 2. NKey Flow Testing

```javascript
// Test complete NKey workflow
const testNKeyFlow = async (userToken) => {
  try {
    // 1. Generate NKey
    const nkeyResponse = await generateNKey(userToken, ['searchall']);
    const nkey = nkeyResponse.data.nkey;
    
    // 2. Validate NKey
    const validationResponse = await validateNKey(nkey, 'searchall');
    
    console.log('NKey validation result:', validationResponse.data.valid);
  } catch (error) {
    console.error('NKey flow test failed:', error);
  }
};
```

This documentation provides comprehensive guidance for frontend developers to implement proper permission handling in TouNetCore applications. Remember to always validate permissions on both client and server sides for maximum security.
