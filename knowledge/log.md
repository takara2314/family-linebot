# ナレッジバンドル更新履歴

## 2026-08-03

* **Terraform基盤適用**: `infra/` の単一root moduleを追加し、state bucket、GCS backend、必要API、runtime/deployerサービスアカウント、追加型IAM、WIF、Secret Manager resourceを作成した。apply後planは0差分で、ユーザー管理サービスアカウントキーは0件である。
* **Secret移行**: 配信中App Engine versionの環境変数から値を出力せず、LINE channel secret/tokenをSecret Managerのversion 1へ移行した。
* **設定集約**: `config.go` を追加し、環境変数、既定値、Secret ManagerからのLINE認証情報取得を `loadConfig` にまとめた。`main.go` はclient初期化とroutingを担当する。
* **Gemini更新**: Vertex AIからGemini Enterprise Agent Platformへの公式名称変更に合わせ、Go SDK backendを `BackendEnterprise` へ変更した。スタンプmodelはGemini 2.5 Flashと同額で新しい `gemini-3.5-flash-lite` を既定値とした。
* **スケーリング上限**: 家族内利用と費用上限を優先し、App Engine automatic scalingを `min_instances: 0`、`max_instances: 1` とした。
* **翻訳レイアウト保持**: 構造化された文章を行単位で一括翻訳し、空行、インデント、末尾空白、最終改行を復元する処理を追加した。CRLFはLINE表示用にLFへ正規化する。
* **差分最小化**: 追加していたapplication層、設定層、service interfaceを廃止し、元のGin handler、グローバルclient、`convertJpTh.go` と `post*Message.go` の構成へ実装を戻した。
* **アプリ刷新の実装**: module pathを `github.com/takara2314/family-linebot` に変更し、Go 1.26、LINE SDK v8、Translation v3、Geminiによるスタンプ解釈、Speech-to-Text v2 `chirp_3`、Secret Manager参照へ更新した。これは未デプロイの作業ツリー状態である。
* **品質確認**: ビルドと `go vet` を通し、`govulncheck` が指摘したgRPCと `x/text` の間接依存を修正版へ引き上げた。既存方針に合わせ、単体テストは新設しない。
* **情報設計更新**: 配信中App Engine version、traffic、version数などのlive stateを手書きOKFから除外し、read-onlyコマンドによる取得方法へ置き換えた。
* **表記変更**: OKF frontmatterと本文を日本語へ統一した。
* **意思決定更新**: スタンプ処理を文書OCRではなくマルチモーダル解釈として扱い、Geminiを第一候補、Vision OCRを比較基準・障害時fallbackとした。
* **初期作成**: 現行システム、意思決定、Terraform、WIF、言語サービス、費用、実施計画をOKF v0.2で記録した。
