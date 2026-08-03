---
type: 実施記録
title: family-linebot刷新の実施記録
description: 単一本番環境で完了した刷新作業と、運用上の任意残件。
tags: [plan, modernization, production]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T13:06:45+09:00" }
verified: { by: "human:takara2314", at: "2026-08-03T13:06:45+09:00" }
status: stable
sources:
  - id: current
    resource: /current-system.md
    title: 現行システム台帳
  - id: decisions
    resource: /architecture-decisions.md
    title: アーキテクチャ意思決定
---

# 結果

刷新は2026-08-03に完了し、オーナーがLINEのtext翻訳、スタンプ理解、日本語音声認識を本番で確認した。環境は1つだけで、stagingは設けていない。

# Phase 1: ナレッジと検証

1. [完了] 意思決定の変更に合わせてこのOKF bundleを更新する。
2. [完了] 単体テストを新設せず、ビルド、静的解析、脆弱性検査、本番の手動確認で検証する。

# Phase 2: Terraform bootstrap

1. [完了] 単一rootの `infra/` 構成を追加する。
2. [完了] `hamaguchi-family-linebot-tfstate` を手動bootstrap resourceとして作り、Terraform管理外のGCS backendとして利用する。
3. [完了] 必要API、App Engine application、identity、IAM、secretを管理する。配信versionとtrafficはデプロイのlive stateとしてTerraform管理外にする。
4. [完了] apply後の `terraform plan` が0差分で、既存IAMを削除しないことを確認する。

# Phase 3: キーレスデプロイとsecret

1. [完了] runtimeとdeployerサービスアカウントを作る。
2. [完了] Repository ID `324532124`、Owner ID `26915490`、`master` に制限したGitHub WIF trustを作る。
3. [完了] Secret Manager resourceを作り、配信中versionから値を表示せずLINE secret payloadをversion 1へ移行する。
4. [完了] GitHub Actionsのビルド・静的解析・脆弱性検査とApp Engineデプロイを追加し、本番デプロイに成功する。

# Phase 4: アプリケーション刷新

1. [完了] Go 1.26、GitHub形式module path、LINE SDK v8へ更新する。
2. [完了] Google Cloud clientを起動時に生成して再利用する。既存コードとの差分を抑えるため、panicベースのerror処理は維持する。
3. [完了] TranslationをAdvanced v3 NMTへ移行し、翻訳時の改行とインデントを維持する。
4. [完了] Gemini 3.5 Flash-Liteによる構造化スタンプ解釈を第一経路として実装する。
5. [完了] Speechをv2 `chirp_3`、`us`、`ja-JP` のみへ移行する。
6. [完了] `app.yaml` に専用runtimeアカウントと `max_instances: 1` を指定する。

# Phase 5: 本番デプロイ

1. [完了] GitHub ActionsからApp Engineへデプロイし、新versionへtrafficを100%切り替える。
2. [完了] root health、LINE webhook、text翻訳、スタンプ理解、日本語音声を検証する。
3. [完了] Cloud Loggingで初回エラーを特定し、App Engine APIとSpeech locationを修正する。
4. [完了] rollback可能な以前のversionを保持する。

# 任意の運用改善

次は刷新の完了条件ではなく、影響範囲を確認してから別途判断する。

1. 古いユーザー管理サービスアカウントキーの他用途を調査し、不要なら無効化する。
2. App Engine defaultアカウントのEditorと個人名Ownerサービスアカウントの必要性を確認する。
3. rollback用の直前version以外を整理し、必要な場合だけArtifact Registryのcleanup policyを検討する。
4. GitHub Actionsにpath filterを追加し、ドキュメントやTerraformだけの変更による再デプロイを避けるか判断する。

# 完了条件

* GitHubが長期Google credentialなしでデプロイできる。
* アプリがGo 1.26と最小権限runtime identityで稼働する。
* LINE secretがsource、deployment YAML、GitHub secrets、Terraform stateに存在しない。
* 日本語text、日本語speech、タイ語text、スタンプ理解が本番検証に合格する。
* importとデプロイ後のTerraform planが安定する。
