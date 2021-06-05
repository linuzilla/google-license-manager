# Google License Manager

Google License 管理程式：部份 Google Workspace 的 Subscriptions 是以 licenses 計價，licenses 必需設定給使用者，則該使用者可以得到額外的功能。
例如: Google Workspace for Education: Teaching and Learning Upgrade。

Google 提供 [Admin Console](https://admin.google.com/) 讓管理者可以管理授權，而這程式能方便的對這些授權做管理，不需要有管理者權限的人透過 Admin Console 做設定。

## 開發背景

2021 年 5 月，因為疫情的關係，遠距教學變成日常所需，Google Meet 是一個選擇。但當老師們開始用 Google Meet 時，就開始希望有更多的功能。其中，
Google Workspace for Education: Teaching and Learning Upgrade 可以申請一次的 60 天，50 個 license 的免費試用。
所以很多學校在這時候使用這個試用，至少撐到暑假 ...

因此，我撰寫這個工具，方便自己管理。但學校有些人希望提供給他們使用，我這不是網頁版，考量風險沒有很高，所以有開給部份人員使用。

## 工作原理及風險

利用 [Google Cloud Platform](https://console.cloud.google.com) 建立 project，開啟 Enterprise License Manager API，然後產生一個 Service Account，
產生一組金鑰（以 json 格式下載），在 Admin Console 將此金鑰的 client id 註冊到網域的 domain-wide delegation，並附予
https://www.googleapis.com/auth/apps.licensing 的 scope

注意：本程式支援用戶名的取得，但這需額外的 API 及 scope，需要的 API 為 Admin SDK API，scope 為
https://www.googleapis.com/auth/admin.directory.user.readonly

當這些都準備好了，透過這程式，加上金鑰，就可以直接對授權的人員做管理

如果金鑰有遺失，到 Admin Console 把該筆記錄刪掉，並到 Gpogle Cloud Platform 刪掉這把金鑰即可。

## 程式運作模式

程式運作需有設定檔跟金鑰，另外會產生一個資料檔。資料檔是為了記錄異動用的，另外，資料檔也提供註解帳號的資訊，否則（如果不開 Admin SDK）時，帳號是無法顯示擁有者的姓名。

管理者可以把設定檔封裝在資料檔中，用 aes 加密保護，採用這個方式的話，程式開啟時，會詢問密碼，以解密設定跟鑰。用封裝後，就只要提供資料檔給其它管理者即可。

Go 是跨平台的，所以程式可以在 Windows, Linux, 跟 Mac 上執行，原本想用 sqlite 做資料庫，但因用到 cgc，在 Windows 上實現有點複雜，所以選用是純 GO 語言的程式庫。

程式使用介面是 command line 的，所以採下指令的方式運作。cli 支援 completion 跟 history 功能。另外，也可以結合 shell 來下指令，不需用程式提供的 cli。

## 設定檔

預設的設定檔名為 application.yml, 內容如下，其中 product-id 跟 product-sku-id 參見 
[Google 的文件](https://developers.google.com/admin-sdk/licensing/v1/how-tos/products)

其中，Google Workspace for Education: Teaching and Learning Upgrade 的 Product Id 為 101037
Product Sku ID 為 1010370001。

scopes 的部份，如果沒有授權 Admin SDK 的話，請把 "https://www.googleapis.com/auth/admin.directory.user.readonly" 拿掉。

如果程式執行時，無法正確讀取 application.yml 時，程式會假設跑在封裝模式，預設的資料檔為 licenses-users.db。

```yaml
name: Google License Manager

# log-level: debug, notice, warning, success, failed, error, fatal
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

## 指令

### list

列出目前授權的用戶 (資料檔如果有的話，會顯示用戶名稱），會顯示跟資料檔的差別，但不會回寫到資料檔

### sync

功能類似 list，唯一差別是會回寫，如果有新的帳戶，而且有開 Admin SDK 時，會把帳號的全名寫的描述資料中。

### add-user

授權使用者，後面需要帳號當參數，可多名使用者一起加

### revoke-user

對使用者移出授權，後面需要帳號當參數，可多名使用者一起移戈出授權

### describe

對使用者描述，需參數，第一個參數為帳號，接下來寫描述內容，例如：姓名

### dump-database

將資料庫中的使用者列出

### store-and-encode

這個功能是把設定檔跟金鑰檔包進資料庫中，只有在未封裝的模式啟動程式時才有這個指令，後面需要密碼做為參數

### change-password

修改封裝的密碼，只有在封裝模式啟動時，才有這個指令

## 命令列直接帶指令

如果想用 shell script 直接下指令做動作可以在程式後加兩個 dash, 再加指令及參數, 如 (注意, 需要雙引號)

```bash
./google-license-manager -- "add-user user01"
```
