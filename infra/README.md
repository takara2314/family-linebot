# Terraform

`hamaguchi-family-linebot` の本番環境だけを管理する単一root moduleです。環境別directoryやworkspaceは使用しません。

## State backend

このprojectでは `hamaguchi-family-linebot-tfstate` をlocal stateでbootstrapした後、`backend.tf` のGCS backendへ移行済みです。通常は次だけで初期化できます。

```bash
terraform init
```

別projectへ横展開する場合は、最初だけ `backend.tf` を追加する前のlocal stateでstate bucketを作成し、その後backendを追加して移行します。

```bash
terraform apply -target=google_storage_bucket.terraform_state
terraform init -migrate-state
terraform plan
```

## Secret payload

TerraformはSecret ManagerのsecretとIAMだけを管理し、LINE認証情報の値はstateへ保存しません。値は初回apply後に、管理者が標準入力から追加します。

```bash
printf '%s' 'VALUE' | gcloud secrets versions add SECRET_ID \
  --project=hamaguchi-family-linebot \
  --data-file=-
```

## 通常操作

```bash
terraform fmt -check -recursive
terraform validate
terraform plan
```
