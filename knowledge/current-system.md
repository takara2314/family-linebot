---
type: システム構成台帳
title: family-linebotの現行システム
description: コードとgcloud CLIから確認した、刷新後の本番構成に関する基準情報。
resource: https://console.cloud.google.com/home/dashboard?project=hamaguchi-family-linebot
tags: [family-linebot, google-cloud, app-engine, 現状調査]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T13:06:45+09:00" }
verified:
  - { by: "process:gcloud-read-only-inventory", at: "2026-08-03T13:06:45+09:00" }
  - { by: "human:takara2314", at: "2026-08-03T13:06:45+09:00" }
status: stable
sources:
  - id: repository
    resource: https://github.com/takara2314/family-linebot
    title: family-linebot GitHubリポジトリ
  - id: gcp-project
    resource: https://console.cloud.google.com/home/dashboard?project=hamaguchi-family-linebot
    title: hamaguchi-family-linebot Google Cloudプロジェクト
---

# アプリケーション

| 項目 | 現在値 |
|---|---|
| Google Cloud Project ID | `hamaguchi-family-linebot` |
| Project number | `962903877963` |
| App Engineリージョン | `asia-northeast2` |
| App Engine環境 | Standard |
| Service | `default` |
| Runtime | `go126` |
| Instance class | `F1` |
| Scaling | Automatic |
| Scaling上限 | `min_instances: 0`、`max_instances: 1` |
| 実行サービスアカウント | `family-linebot-runtime@hamaguchi-family-linebot.iam.gserviceaccount.com` |

# アプリケーション機能

module pathは `github.com/takara2314/family-linebot` である。LINE SDK v8、Translation Advanced v3 NMT、Gemini Enterprise Agent Platform上のGemini 3.5 Flash-Lite、Speech-to-Text v2 `chirp_3` を利用する。Speechは `us` multi-regionで `ja-JP` だけを受け付ける。LINE認証情報はSecret Managerから起動時に取得する。

元のGin handler、グローバルclient、`convertJpTh.go` と `post*Message.go` の構成を維持し、環境変数・既定値・Secret Manager参照は `config.go` に集約している。

# Live stateの確認

配信中version、traffic split、version数など、デプロイごとに変わるlive stateはこのOKFに固定値として保存しない。必要な時点で次のread-onlyコマンドから取得する。

```bash
gcloud app services list --project=hamaguchi-family-linebot
gcloud app versions list --project=hamaguchi-family-linebot
```

このconceptには、刷新の背景説明やimport判断に必要で、比較的長く有効な構成だけを置く。

# 秘密情報

`linebot-channel-secret` と `linebot-channel-token` をSecret Managerで管理する。App Engineにはsecret payloadを置かず、secret名だけを環境変数で渡す。payloadはTerraform state、GitHub Actions、リポジトリに保存しない。

# 認証と権限

| Principal | 現在のProject role |
|---|---|
| `hamaguchi-family-linebot@appspot.gserviceaccount.com` | Editor |
| `takara2314@hamaguchi-family-linebot.iam.gserviceaccount.com` | Owner |
| `takara2314.tk@gmail.com` | Owner |

個人名のサービスアカウントにはユーザー管理キーが2個残っている。現行のApp Engine実行とGitHub Actionsデプロイでは使用していないが、無効化前に他用途の利用状況を確認する。

# デリバリー基盤

TerraformはApp Engine application、必要API、runtime/deployerサービスアカウント、IAM、GitHub Actions用WIF pool/provider、Secret Manager secretを管理する。state bucket `hamaguchi-family-linebot-tfstate` はbootstrap resourceとしてTerraform外で管理する。apply後のplanは差分0である。

`master` pushはGitHub Actionsでbuild、vet、脆弱性検査、WIF認証を行い、成功後にApp Engineへpromoteする。任意のArtifact Registry repositoryは作成していない。`gae-standard` repositoryはApp Engineがbuild image保存用に自動作成した内部実装である。

GitHubリポジトリはPublicである。不変なRepository IDは `324532124`、Owner IDは `26915490`、現在のdefault branchは `master` である。
