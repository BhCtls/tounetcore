# Apps表URL字段添加 - API更新说明

## 更新概述
为apps表添加了URL字段，支持应用直接访问功能。更新了相关API以支持URL字段的创建、更新和查询。

## 数据库模型更新

### App 模型 (models.go)
- 添加了 `URL string json:"url"` 字段
- 支持应用URL的存储和访问

## API 更新

### 1. 创建应用 API
**Endpoint**: `POST /api/v1/admin/apps`

**新增字段**:
```json
{
  "app_id": "example_app",
  "name": "Example App",
  "description": "App description",
  "url": "https://example.com",  // 新增字段
  "required_permission_level": "user",
  "is_active": true
}
```

**响应包含URL字段**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "app_id": "example_app",
    "name": "Example App",
    "url": "https://example.com",  // 新增字段
    "secret_key": "generated_secret",
    "required_permission_level": "user",
    "is_active": true
  }
}
```

### 2. 更新应用 API
**Endpoint**: `PUT /api/v1/admin/apps/{app_id}` (修改为标准RESTful)

**新增字段**:
```json
{
  "name": "Updated App Name",
  "description": "Updated description",
  "url": "https://newurl.com",  // 新增字段
  "required_permission_level": "trusted",
  "is_active": false
}
```

### 3. 获取用户应用列表 API
**Endpoint**: `GET /api/v1/user/apps`

**响应包含URL字段**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "app_id": "example_app",
      "name": "Example App",
      "url": "https://example.com",  // 新增字段
      "enabled": true,
      "valid_until": null
    }
  ]
}
```

### 4. NKey 生成 API
**Endpoint**: `POST /api/v1/nkey/generate`

功能保持不变，无需修改。

## 数据库迁移

创建了 `migrations/add_url_to_apps.sql` 文件：
```sql
-- Add URL field to apps table
ALTER TABLE apps ADD COLUMN url TEXT DEFAULT '';

-- Add index on url field for better performance (optional)
CREATE INDEX IF NOT EXISTS idx_apps_url ON apps(url);

-- Update existing apps with empty URL (optional)
UPDATE apps SET url = '' WHERE url IS NULL;
```

## 路由更新

- 修改了应用更新API路由为标准RESTful风格：
  - 从 `POST /api/v1/admin/apps/:app_id/update` 
  - 改为 `PUT /api/v1/admin/apps/:app_id`
- 修改了应用删除API路由为标准RESTful风格：
  - 从 `POST /api/v1/admin/apps/:app_id/delete`
  - 改为 `DELETE /api/v1/admin/apps/:app_id`

## 前端集成流程

1. 用户在仪表板看到可用应用列表（包含URL）
2. 点击有URL的应用卡片
3. 前端调用 `POST /api/v1/nkey/generate` 为该应用生成临时密钥
4. 将NKey存储到cookie中：
   - 名称: `nkey`
   - 值: 生成的NKey字符串
   - 过期时间: 15分钟
   - 路径: `/`
   - SameSite: `Lax`
5. 在新标签页打开应用URL
6. 目标应用从cookie中读取NKey进行身份验证

## 安全考虑

- NKey在15分钟后自动过期
- Cookie使用SameSite=Lax提供CSRF保护
- 每次访问应用都会生成新的NKey
- NKey只对特定应用有效
- URL字段允许为空，向后兼容现有应用

## 测试建议

1. 测试创建应用时URL字段的保存
2. 测试更新应用时URL字段的修改
3. 测试用户应用列表API返回URL字段
4. 测试NKey生成和验证流程
5. 测试数据库迁移的执行
6. 测试新的RESTful API路由

## 更新完成时间
2025年8月1日