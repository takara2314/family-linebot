---
type: アプリケーション設計
title: 翻訳、スタンプ理解、音声認識
description: Googleの言語・画像サービスに関する品質、言語、API version、評価の意思決定。
tags: [translation, ocr, speech-to-text, japanese, thai]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T13:06:45+09:00" }
verified:
  - { by: "human:takara2314", at: "2026-08-03T00:00:00+09:00" }
  - { by: "process:official-doc-review", at: "2026-08-03T00:00:00+09:00" }
  - { by: "human:takara2314", at: "2026-08-03T13:06:45+09:00" }
status: stable
stale_after: 2026-11-03
sources:
  - id: translation-overview
    resource: https://docs.cloud.google.com/translate/docs/api-overview
    title: Cloud Translation API概要
  - id: translation-migration
    resource: https://docs.cloud.google.com/translate/docs/migrate-to-v3
    title: Cloud Translation BasicからAdvancedへの移行
  - id: vision-ocr
    resource: https://docs.cloud.google.com/vision/docs/how-to
    title: Cloud Vision OCRガイド
  - id: gemini-image
    resource: https://ai.google.dev/gemini-api/docs/image-understanding
    title: Gemini画像理解
  - id: gemini-structured-output
    resource: https://ai.google.dev/gemini-api/docs/structured-output
    title: Gemini構造化出力
  - id: vertex-data-governance
    resource: https://docs.cloud.google.com/vertex-ai/generative-ai/docs/vertex-ai-zero-data-retention
    title: Gemini Enterprise Agent Platformのデータガバナンスとzero data retention
  - id: speech-migration
    resource: https://docs.cloud.google.com/speech-to-text/docs/migration
    title: Speech-to-Text v1からv2への移行
  - id: chirp-release
    resource: https://docs.cloud.google.com/speech-to-text/docs/release-notes
    title: Speech-to-Text Chirp 3リリースノート
---

# 翻訳

Cloud Translation Basic v2からAdvanced v3へ移行済みである。カジュアルな会話に適し、IAMを利用でき、十分な月間無料枠があるため、NMTを本番の基準とする。将来切り替える前に、匿名化した家族の日本語・タイ語メッセージでTranslation LLMを評価する。

次の挙動を維持する。

* 日本語はタイ語へ翻訳する。
* タイ語は日本語へ翻訳する。
* その他の言語は日本語とタイ語の両方へ翻訳する。

元コードの振る舞いを維持するため、`gDetectLanguage.go` で言語判定し、`gtranslate.go` で翻訳する構成を継続する。構造化された文章は非空行を1回のTranslation APIリクエストにまとめ、空行、改行、行頭インデント、行末空白を元の位置へ戻す。

# スタンプ文字理解

画像処理を無料に保つことより品質を優先する。これは通常の文書OCRではない。意図されたメッセージは、キャラクターの表情、吹き出し、文字位置、縦書き、装飾文字、視覚的強調に依存し得る。そのためGemini画像理解を第一の実装候補とする。[^gemini-image]

`postStickerMessage.go` でLINEスタンプ画像をダウンロードし、そのbytesをGemini Enterprise Agent Platform上のGeminiへ送る。Google Cloud runtime identityとGoogle Cloudのデータガバナンス制御を使えるため、consumer向けGemini Developer API key連携よりAgent Platform APIを優先する。[^vertex-data-governance] HTTP timeoutと画像size制限は現在の最小差分実装には含めていない。

Geminiの応答は、現在コードで必要な判定値だけをJSON Schemaで構造化する。[^gemini-structured-output]

```json
{
  "has_visible_text": true,
  "normalized_message": "配置や絵の文脈を使って読み順と表記だけを整えた発話",
  "uncertain": false
}
```

promptでは画像に見えない単語の創作を禁止する。`normalized_message` はスタンプ全体から読み順や装飾表記を解決してよいが、見える文字にgroundingされていなければならない。`has_visible_text` がtrue、`uncertain` がfalse、かつ `normalized_message` が空でない場合だけ翻訳する。それ以外、または画像取得・Gemini処理に失敗した場合は、元コードに近い挙動として応答せずlogを残す。

Cloud Vision `TEXT_DETECTION` は現在コードには残していない。将来、次の必要性が確認された場合に追加を判断する。

* 評価データにおける決定論的な比較基準
* Geminiの障害、rate limit、safety filter時のfallback
* Geminiのhallucinationを調査するための第二の信号

既定modelは `gemini-3.5-flash-lite` とする。最新GAのFlash-Liteで画像入力と構造化出力に対応し、Gemini 2.5 Flashの置換候補として案内されている。料金はGlobal Standardで入力$0.30、text出力$2.50／100万tokenで、Gemini 2.5 Flashと同額である。最安の `gemini-3.1-flash-lite` は入力$0.25、text出力$1.50／100万tokenだが、スタンプ理解は品質優先のため既定値にはしない。

日本語・タイ語スタンプの代表的な評価データを作り、見える文字の完全一致、読み順、正規化メッセージ品質、存在しない文字の創作、latency、費用を測る。Vision `TEXT_DETECTION` との比較も必要に応じて行う。将来のmodel更新で再設計しなくてよいよう、model IDは環境変数で上書き可能にする。

# 音声認識

Speech-to-Text v2の `chirp_3` を本番利用する。品質のための課金は許容する。

recognizer設定では次だけを指定する。

```text
model: chirp_3
location: us
language_codes: [ja-JP]
```

言語自動判定を有効にせず、代替言語にタイ語を含めない。音声メッセージの約99%が日本語なので、タイ語音声は意図的に非対応とする。

非対応または認識不能の音声には、音声認識は現在日本語だけに対応することを明示したLINE応答を返す。LINEが返すMP3またはM4AをSpeech-to-Text v2のauto decodingへ渡す。v2はMP3、MP4/M4A AAC、WebM/Opusを自動判定できるため、sample rateやcodecを固定しない。`asia-northeast1` では実際のAPIが日本語Chirp 3を拒否したため、品質を優先して `us` multi-regionを利用する。

# クライアントのライフサイクル

環境変数、既定値、Secret ManagerからのLINE認証情報取得は `config.go` の `loadConfig` に集約する。Translation、Gemini、Speech clientは `main.go` でアプリ起動時に一度だけ作り、元コードと同様にグローバルclientとして各 `post*Message.go` から利用する。既存コードとの差分を抑えるため、handler構成やpanicベースのerror処理は全面変更しない。request timeout、payload size制限、構造化logは今回の最小変更には含めず、必要になった場合の改善項目とする。

[^gemini-image]: Gemini modelはマルチモーダル画像理解、visual question answering、画像文字処理をサポートする。
[^gemini-structured-output]: Geminiは対応するJSON Schemaで応答を制約できる。
[^vertex-data-governance]: Googleは、許可または指示なくAgent Platformの顧客データをmanaged modelの学習に使用しないとしている。参照URLに旧名称のVertex AIが残る場合がある。
