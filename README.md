# google-license-manager
Assign licenses for Google Workspace Subscriptions

Documents are available in [Traditional Chinese](README_zh_TW.md)

## Simple startup program 
```go
package main

import "github.com/linuzilla/google-license-manager"

func main() {
	google_license_manager.Main()
}
```

## Simple configuration (application.yml)
```yaml
name: Google License Manager

# log-level: debug, notice, warning, error, fatal
log-level: notice
database-file: licenses-users.db

#  Google Product and SKU IDs
#  https://developers.google.com/admin-sdk/licensing/v1/how-tos/products
google:
  domain: "example.com"
  admin: "admin@example.com"
  customer-id: "customer-id
  product-id: "101037"
  product-sku-id: "1010370001"
  credential-file: credential.json
  scopes:
    - "https://www.googleapis.com/auth/apps.licensing"
    - "https://www.googleapis.com/auth/admin.directory.user.readonly"
```
