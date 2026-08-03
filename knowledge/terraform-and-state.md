---
type: インフラ設計
title: Terraformとremote stateの設計
description: family-linebotを横展開できる単一環境のTerraform構成。
tags: [terraform, gcs, tfstate, infrastructure-as-code]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T00:00:00+09:00" }
verified: { by: "process:official-doc-review", at: "2026-08-03T00:00:00+09:00" }
status: stable
stale_after: 2027-02-03
sources:
  - id: tf-security
    resource: https://docs.cloud.google.com/docs/terraform/best-practices/security
    title: Google CloudにおけるTerraformセキュリティのベストプラクティス
  - id: state-storage
    resource: https://docs.cloud.google.com/docs/terraform/resource-management/store-state
    title: Cloud StorageへのTerraform state保存
  - id: iam-resources
    resource: https://cloud.google.com/docs/terraform/best-practices/working-with-resources
    title: Google CloudリソースをTerraformで扱うベストプラクティス
  - id: import
    resource: https://docs.cloud.google.com/docs/terraform/resource-management/import
    title: Google CloudリソースのTerraform stateへのimport
---

# ディレクトリ構成

`infra/` 配下の単一root moduleを使う。`environments/`、`production.tfvars`、`example.tfvars`、追加のTerraform workspaceは作らない。単一の `terraform.tfvars` に横展開時に変更する非secret値を置く。

想定ファイルは `versions.tf`、`providers.tf`、`backend.tf`、`variables.tf`、`terraform.tfvars`、目的別のresourceファイル、`outputs.tf`、`imports.tf`、`README.md` である。

# State bucket設計

bucket名は `hamaguchi-family-linebot-tfstate`、backend prefixは `default` とする。bootstrapとGCS backendへのstate移行は完了している。

`asia-northeast2` に次の設定で作る。

* Standard storage
* Uniform bucket-level access
* Public Access Preventionを強制
* Object Versioning
* 7日間のsoft delete
* 古い非現行versionのlifecycle rule
* Terraformの `prevent_destroy`

GoogleはCloud Storage remote backendを推奨し、build systemと高権限管理者だけにアクセスを制限するよう案内している。[^tf-security]

# 初回構築

最初だけlocal stateでstate bucketを作成し、GCS backendを有効化して `terraform init -migrate-state` で移行する。復旧方法とbootstrapコマンドを `infra/README.md` に記録する。

# Resourceの管理境界

既存App Engine application本体、配信version、trafficなどは今回のTerraform管理対象に含めない。App Engineのimmutableなapplication設定やlive stateをimportして管理範囲を広げず、デプロイ基盤に必要なAPI、identity、IAM、secret、state bucketだけを管理する。自動生成されたservice agentやGoogle内部APIも管理しない。

追加型の `google_*_iam_member` resourceを使う。Project全体に対するauthoritativeな `iam_policy` と `iam_binding` は、Google管理または残すべき手動memberを削除し得るため避ける。[^iam-resources]

Secret resourceとIAMはTerraformで管理するが、平文をstateへ入れないためsecret payload versionはTerraform外から追加する。[^tf-security]

[^tf-security]: Google Cloud公式Terraformセキュリティガイダンス。
[^import]: Google Cloud公式resource importガイダンス。
[^iam-resources]: Google Cloud公式IAM Terraform resourceガイダンス。
