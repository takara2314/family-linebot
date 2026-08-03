---
type: セキュリティ設計
title: GitHubデプロイの認証とIAM
description: Workload Identity Federationを使ったGitHub ActionsからApp Engineへのキーレスデプロイ。
tags: [github-actions, workload-identity-federation, iam, app-engine]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T00:00:00+09:00" }
verified: { by: "process:official-doc-review", at: "2026-08-03T00:00:00+09:00" }
status: stable
stale_after: 2027-02-03
sources:
  - id: wif-pipelines
    resource: https://docs.cloud.google.com/iam/docs/workload-identity-federation-with-deployment-pipelines
    title: デプロイパイプライン向けWIF設定
  - id: wif-practices
    resource: https://docs.cloud.google.com/iam/docs/best-practices-for-using-workload-identity-federation
    title: WIFのベストプラクティス
  - id: appengine-roles
    resource: https://docs.cloud.google.com/appengine/docs/standard/roles
    title: App Engineデプロイrole
  - id: service-accounts
    resource: https://docs.cloud.google.com/iam/docs/best-practices-service-accounts
    title: サービスアカウントのセキュリティベストプラクティス
---

# Identity設計

単一目的のサービスアカウントを2つ作る。

* `github-appengine-deployer@hamaguchi-family-linebot.iam.gserviceaccount.com`
* `family-linebot-runtime@hamaguchi-family-linebot.iam.gserviceaccount.com`

deployerには `gcloud app deploy` に必要なApp Engine Deployer、Cloud Build Editor、Storage Object Adminを付与する。Service Account UserはProject全体ではなくruntimeアカウントだけに付与する。[^appengine-roles]

runtimeアカウントには、アプリが必要とするAPI利用権限とsecret単位のアクセス権だけを付与する。Editorは付与しない。

# WIFの信頼条件

GitHub OIDC providerを作り、不変claimで次のように制限する。

```text
assertion.repository_id == "324532124" &&
assertion.repository_owner_id == "26915490" &&
assertion.ref == "refs/heads/master"
```

GitHubは共有マルチテナントissuerを使うため、spoofing防止には信頼するtenantまたはorganizationを制限するattribute conditionが必要である。[^wif-practices]

workflow権限は `contents: read` と `id-token: write` だけにし、`google-github-actions/auth` でdeployerをimpersonateする。長期Google credentialはGitHubに保存しない。

# 本番デプロイフロー

現在のデプロイbranchは `master` である。push時に次を実行する。

1. `go build`、`go vet`、`govulncheck` を実行する。既存方針に合わせて単体テストは新設しない。
2. GitHub OIDC tokenをWIFで交換する。
3. `gcloud app deploy app.yaml --quiet` を実行する。
4. 新しいApp Engine versionとLINE webhookの動作を監視する。

staging環境は作らない。ビルド・静的解析・デプロイ後の手動確認と、以前のApp Engine versionへtrafficを戻せることを本番の安全策とする。

# 古いキーの廃止

WIFデプロイ成功後、2個の古いユーザー管理キーの利用状況を調査する。サービスアカウントの推奨に従い、削除前に無効化して観察期間を設ける。[^service-accounts]

[^appengine-roles]: gcloud CLIデプロイに必要な現行App Engine role。
[^wif-practices]: Google Cloud公式WIF anti-spoofingガイダンス。
[^service-accounts]: Google Cloud公式サービスアカウントlifecycleガイダンス。
