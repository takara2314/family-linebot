# Terraform

`hamaguchi-family-linebot` の本番環境だけを管理する単一root moduleです。環境別directoryやworkspaceは使用しません。

App Engine application本体も管理します。既存projectへ導入する場合は、作成済みapplicationを一度だけimportします。

```bash
terraform import google_app_engine_application.default PROJECT_ID
```

App Engine applicationは作成後にlocation変更や単体削除ができないため、resourceには `prevent_destroy` を設定しています。

## State backend

state bucketはTerraform自身のstateへ含めず、bootstrap resourceとして手動管理します。このprojectでは `hamaguchi-family-linebot-tfstate` を作成済みです。通常は次だけで初期化できます。

```bash
terraform init
```

別projectへ横展開する場合は、Terraformを初期化する前にbucketを作成します。

```bash
gcloud storage buckets create gs://BUCKET_NAME \
  --project=PROJECT_ID \
  --location=REGION \
  --uniform-bucket-level-access \
  --public-access-prevention

gcloud storage buckets update gs://BUCKET_NAME --versioning
```

作成後、`backend.tf` のbucket名を合わせて `terraform init` を実行します。bucket削除防止、soft delete、古いversionのlifecycle設定もGoogle Cloud側で維持します。

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
