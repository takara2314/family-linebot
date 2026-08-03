---
type: 費用モデル
title: Google Cloud月額費用の見込み
description: 刷新後の単一本番家族アプリケーションに関する費用見込み。
tags: [cost, free-tier, google-cloud]
generated: { by: "openai-codex/gpt-5", at: "2026-08-03T13:06:45+09:00" }
verified: { by: "process:official-pricing-review", at: "2026-08-03T00:00:00+09:00" }
status: stable
stale_after: 2026-09-03
sources:
  - id: appengine-pricing
    resource: https://docs.cloud.google.com/appengine/docs/standard/quotas
    title: App Engine Standard quota
  - id: translation-pricing
    resource: https://cloud.google.com/products/translate/pricing
    title: Cloud Translation料金
  - id: vision-pricing
    resource: https://cloud.google.com/vision/pricing
    title: Cloud Vision料金
  - id: vertex-pricing
    resource: https://cloud.google.com/vertex-ai/generative-ai/pricing
    title: Gemini Enterprise Agent Platform生成AI料金
  - id: speech-pricing
    resource: https://cloud.google.com/speech-to-text/pricing
    title: Speech-to-Text料金
  - id: secrets-pricing
    resource: https://cloud.google.com/secret-manager/pricing
    title: Secret Manager料金
  - id: storage-pricing
    resource: https://cloud.google.com/storage/pricing
    title: Cloud Storage料金
---

# 基本的な見込み

家族内の利用量では、App Engine F1、Translation v3 NMT、Secret Manager、WIF、CI/CDは無料または実質無視できる費用に収まる見込みである。Gemini/Visionによるスタンプ処理とSpeech-to-Textの変動課金は許容する。

# 現行料金の判断材料

* App Engine Standardには1日28 F1 instance-hoursが含まれる。[^appengine-pricing]
* Translation NMTは月50万文字まで含まれ、その後は100万文字あたりUSD 20である。[^translation-pricing]
* Vision text detectionは月1,000 unitsまで含まれ、その後は1,000 unitsあたりUSD 1.50である。[^vision-pricing]
* Gemini Enterprise Agent PlatformのGeminiは画像入力と構造化text出力のtoken課金であり、正確な月額は選択model、画像resolution、thinking設定、応答長に依存する。[^vertex-pricing]
* Speech-to-Text v2 standard recognitionは最初の1分から処理1分あたりUSD 0.016である。[^speech-pricing]
* Secret Managerにはactive version 6個と月10,000 accessesが含まれる。[^secrets-pricing]
* Cloud Storage Always Freeの保存領域は指定された米国regionに限られるため大阪のtfstate bucketは課金対象だが、小さいstate bundleは月USD 0.01を大幅に下回る見込みである。[^storage-pricing]

# 試算例

月300分の音声と2,000枚のVision OCR画像では、為替換算前の変動API費用は約USD 6.30である。スタンプ評価でmodelが未選定のため、この例はVision基準でありGeminiをまだ含まない。

* Speech: `300 × 0.016 = USD 4.80`
* OCR: `(2,000 - 1,000) / 1,000 × 1.50 = USD 1.50`

この費用は認識品質との交換条件として許容する。可視化のためbilling budgetとalertを設定する。budget alertは利用量を自動停止しない。

# 費用管理

App EngineはF1 automatic scaling、`min_instances: 0`、`max_instances: 1` を維持する。不要なsecret versionをdestroyし、App Engine管理build imageの保存量を定期的に確認する。

[^appengine-pricing]: 2026-08-03に確認したApp Engine Standard無料quota。
[^translation-pricing]: 2026-08-03に確認したCloud Translation料金。
[^vision-pricing]: 2026-08-03に確認したCloud Vision料金。
[^vertex-pricing]: 2026-08-03に確認したGemini Enterprise Agent PlatformのGemini料金。既定modelは `gemini-3.5-flash-lite` とし、変更時に再確認する。
[^speech-pricing]: 2026-08-03に確認したSpeech-to-Text v2料金。
[^secrets-pricing]: 2026-08-03に確認したSecret Manager無料利用量。
[^storage-pricing]: 2026-08-03に確認したCloud Storage料金と無料region制限。
