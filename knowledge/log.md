# ナレッジバンドル更新履歴

## 2026-08-03

* **本番検証完了**: オーナーがtext翻訳、改行保持、スタンプ理解、日本語音声認識を確認し、刷新を完了とした。
* **音声認識確定**: Speech-to-Text v2の `chirp_3`、`ja-JP` を採用した。アジアリージョンの実APIで拒否されたため、動作確認できた `us` multi-regionを利用する。
* **アプリ刷新**: Go 1.26、LINE SDK v8、Translation Advanced v3、Gemini 3.5 Flash-Lite、Speech-to-Text v2、Secret Managerへ更新した。元コードの構造を保ち、設定だけを `config.go` に集約した。
* **キーレスCD**: GitHub ActionsからWIFで認証し、build、vet、脆弱性検査後にApp Engineへpromoteするデプロイを本番運用した。
* **Terraform適用**: 単一root moduleでApp Engine application、必要API、runtime/deployer identity、IAM、WIF、Secret Managerを管理した。state bucketは手動bootstrapとし、apply後planは差分0である。
* **秘密情報移行**: 旧App Engine環境変数のLINE認証情報を値を表示せずSecret Managerへ移し、payloadをTerraform stateとGitHubへ保存しない構成にした。
* **OKF整備**: 現行構成、意思決定、Terraform、WIF、言語サービス、費用、刷新記録を日本語のOKF v0.2 bundleとして整理した。
