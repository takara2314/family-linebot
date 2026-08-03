---
type: 実装計画
title: family-linebot刷新計画
description: 検証点とrollback手段を含む、本番環境だけの順序付き刷新作業。
tags: [plan, modernization, production]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T00:00:00+09:00" }
status: draft
stale_after: 2026-11-03
sources:
  - id: current
    resource: /current-system.md
    title: 現行システム台帳
  - id: decisions
    resource: /architecture-decisions.md
    title: アーキテクチャ意思決定
---

# ガードレール

本番環境は1つで、staging環境はない。変更を伴う各stepには、review済みplanまたは明示的な検証点を設ける。刷新versionのデプロイが成功するまで、既存本番trafficは現在のApp Engine versionに残す。

# Phase 1: ナレッジと検証

1. [完了] 意思決定の変更に合わせてこのOKF bundleを更新する。
2. この小規模な既存リポジトリには単体テストを新設せず、ビルド、静的解析、本番前の手動確認で検証する。

# Phase 2: Terraform bootstrap

1. [完了] 単一rootの `infra/` 構成を追加する。
2. [完了] local stateで `hamaguchi-family-linebot-tfstate` を作り、stateをGCSへ移行する。
3. [完了] 必要API、identity、IAM、secretを追加型resourceで管理する。既存App Engine applicationとlive stateは管理対象外とする。
4. [完了] apply後の `terraform plan` が0差分で、既存IAMを削除しないことを確認する。

# Phase 3: キーレスデプロイとsecret

1. [完了] runtimeとdeployerサービスアカウントを作る。
2. [完了] Repository ID `324532124`、Owner ID `26915490`、`master` に制限したGitHub WIF trustを作る。
3. [resource作成完了・payload未投入] Secret Manager resourceを作り、Terraform外からLINE secret payloadを投入する。
4. GitHub Actionsのビルド・静的解析とApp Engineデプロイを追加する。

# Phase 4: アプリケーション刷新

1. [実装・ビルド確認完了] Go 1.26、GitHub形式module path、LINE SDK v8へ更新する。
2. [一部完了] Google Cloud clientは起動時に生成して再利用する。既存コードとの差分を抑えるため、panicベースのerror処理の全面変更は行わない。
3. [実装・ビルド確認完了] TranslationをAdvanced v3 NMTへ移行する。
4. [実装・ビルド確認完了] Gemini Enterprise Agent PlatformのGemini 3.5 Flash-Liteによる構造化スタンプ解釈を第一経路として実装する。
5. 日本語・タイ語スタンプの評価fixtureを作り、Vision OCRを比較基準・fallbackとして追加するか判断する。
6. [実装・ビルド確認完了] Speechをv2 `chirp_3`、`ja-JP` のみへ移行する。
7. `app.yaml` に専用runtimeアカウントを指定する。

# Phase 5: 本番デプロイ

1. GitHub Actionsから新しいApp Engine versionをデプロイする。
2. root health、LINE署名検証、text翻訳、スタンプ理解、日本語音声を検証する。
3. log、latency、API error、Artifact Registryの挙動、billing signalを確認する。
4. rollback用に以前の配信versionを保持する。

# Phase 6: 後片付け

1. 古いユーザー管理サービスアカウントキーを調査して無効化する。
2. 新runtime identityの動作確認後にだけApp Engine defaultアカウントからEditorを外す。
3. 個人名Ownerサービスアカウントは削除検討前に無効化する。
4. 不要なApp Engine versionを削除し、保存build imageの状況から必要な場合だけArtifact Registry cleanup policyを追加する。

# 完了条件

* GitHubが長期Google credentialなしでデプロイできる。
* アプリがGo 1.26と最小権限runtime identityで稼働する。
* LINE secretがsource、deployment YAML、GitHub secrets、Terraform stateに存在しない。
* 日本語text、日本語speech、タイ語text、スタンプ理解が本番検証に合格する。
* importとデプロイ後のTerraform planが安定する。
* サービス停止なしで古いユーザー管理キーを廃止する。
