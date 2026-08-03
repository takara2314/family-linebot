---
type: アーキテクチャ意思決定記録
title: 刷新に関するアーキテクチャ意思決定
description: オーナーが承認した運用制約と主要な技術選択。
tags: [意思決定, 本番環境, app-engine, 品質]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T13:06:45+09:00" }
verified:
  - { by: "human:takara2314", at: "2026-08-03T00:00:00+09:00" }
  - { by: "human:takara2314", at: "2026-08-03T13:06:45+09:00" }
status: stable
sources:
  - id: owner-decisions
    resource: 2026-08-03 family-linebot刷新会話の範囲
    title: 刷新に関する会話でのオーナーの意思決定
---

# 背景

これは家族内専用のLINEアプリケーションである。エンタープライズ向けの環境分離や、すべてのAPI呼び出しを無料に保つことよりも、運用の単純さと出力品質を優先する。

# 意思決定

1. ユーザー管理コンテナを使わず、Google Cloud App Engine Standardを継続利用する。
2. 環境は1つだけとする。成功したデプロイは即本番とし、staging環境、環境別ディレクトリ、環境別Terraform workspaceは作らない。
3. インフラは単一のTerraform root moduleと単一のremote stateで管理する。
4. GitHub ActionsからGitHub OIDCとGoogle Cloud Workload Identity Federationを使ってデプロイする。CI用サービスアカウントJSONキーは作らない。
5. デプロイ専用サービスアカウントとApp Engine実行専用サービスアカウントを分離する。
6. LINEの認証情報はSecret Managerに保存するが、secret payloadをTerraform stateに保存しない。
7. 課金が発生してもスタンプ文字理解と音声認識の品質を優先する。スタンプ処理を文書OCRだけでなくマルチモーダル解釈として扱い、Geminiを第一候補、Vision OCRを比較基準または障害時fallbackとする。
8. 音声認識は日本語だけを受け付ける。`ja-JP` のみを指定し、日本語・タイ語の自動判定を有効にしない。
9. 翻訳はまずCloud Translation Advanced v3 NMTを使う。Translation LLMは既定値ではなく評価候補とする。
10. Artifact RegistryはApp Engineのビルド内部実装として扱い、任意のrepositoryは作らない。
11. App Engine automatic scalingは `min_instances: 0`、`max_instances: 1` とし、同時実行性能より家族内利用に十分な規模と費用上限を優先する。
12. 日本語Chirp 3は、アジアリージョンの実APIで拒否されたため `us` multi-regionを利用する。日本語限定と認識品質をデータ処理リージョンの近さより優先する。

# 結果とトレードオフ

* 唯一の本番環境に直接影響するため、ビルド・静的解析・デプロイ後の手動確認とrollback手段を必須とする。
* Speech-to-Text v2は最初の処理時間から課金されるが、認識品質のために許容する。
* タイ語音声は失敗または誤認識し得るが、音声の約99%が日本語なので許容する。
* Terraform構造はGoogle Cloudの複数環境向け参照構造より意図的に単純にする。

# 関連concept

* [Terraformとstate](/terraform-and-state.md)
* [認証とデプロイ](/identity-and-deployment.md)
* [言語・画像サービス](/language-services.md)
* [費用モデル](/cost-model.md)
