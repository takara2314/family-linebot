---
type: システム構成台帳
title: family-linebotの現行システム
description: コードとgcloud CLIから確認した、刷新判断に必要な既存構成の基準情報。
resource: https://console.cloud.google.com/home/dashboard?project=hamaguchi-family-linebot
tags: [family-linebot, google-cloud, app-engine, 現状調査]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T00:00:00+09:00" }
verified: { by: "process:gcloud-read-only-inventory", at: "2026-08-03T00:00:00+09:00" }
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
| Runtime | `go116` |
| Instance class | `F1` |
| Scaling | Automatic |
| 実行サービスアカウント | `hamaguchi-family-linebot@appspot.gserviceaccount.com` |

上表は本番の移行前基準であり、作業ツリーではGo 1.26への更新を開始している。Runtimeなど刷新によって変更される項目は、デプロイ完了までは本番の現在値と混同しない。

# 未デプロイの刷新実装

2026-08-03時点の作業ツリーでは、module pathを `github.com/takara2314/family-linebot`、runtimeを `go126` に更新した。LINE SDK v8、Translation Advanced v3、Gemini Enterprise Agent Platform上のGemini 3.5 Flash-Lite、Speech-to-Text v2 `chirp_3`、Secret Manager参照を実装済みである。元のGin handler、グローバルclient、`convertJpTh.go` と `post*Message.go` の構成を維持し、環境変数・既定値・Secret Manager参照だけを `config.go` に集約している。ビルドと静的解析は確認済みだが、まだ本番状態を表さない。

作業ツリーのApp Engine automatic scalingは、家族内アプリの利用規模と費用上限を優先し、`min_instances: 0`、`max_instances: 1` としている。

# Live stateの確認

配信中version、traffic split、version数など、デプロイごとに変わるlive stateはこのOKFに固定値として保存しない。必要な時点で次のread-onlyコマンドから取得する。

```bash
gcloud app services list --project=hamaguchi-family-linebot
gcloud app versions list --project=hamaguchi-family-linebot
```

このconceptには、刷新の背景説明やimport判断に必要で、比較的長く有効な構成だけを置く。

# 秘密情報

配信中versionには `LINEBOT_CHANNEL_SECRET` と `LINEBOT_CHANNEL_TOKEN` という環境変数がある。値は取得していない。Secret Managerは無効である。

# 認証と権限

| Principal | 現在のProject role |
|---|---|
| `hamaguchi-family-linebot@appspot.gserviceaccount.com` | Editor |
| `takara2314@hamaguchi-family-linebot.iam.gserviceaccount.com` | Owner |
| `takara2314.tk@gmail.com` | Owner |

個人名のサービスアカウントには、2020年と2022年に作成された無期限のユーザー管理キーが2個ある。無効化前に利用状況を調査する必要がある。

# デリバリー基盤

現時点ではWorkload Identity Pool、WIF provider、Cloud Build trigger、GitHub Actions workflow、Secret Manager secret、明示的なArtifact Registry repositoryは存在しない。

GitHubリポジトリはPublicである。不変なRepository IDは `324532124`、Owner IDは `26915490`、現在のdefault branchは `master` である。
