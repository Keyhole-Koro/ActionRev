# Python 正規化ツール サンドボックス
# セキュリティ設計:
#   - ネットワーク無効 (Cloud Run Jobs の vpc_access.egress=PRIVATE_RANGES_ONLY)
#   - 非 root ユーザーで実行
#   - インターネットアクセス不要な最小パッケージのみインストール

FROM python:3.12-slim

# 非 root ユーザー
RUN useradd -m -u 1000 sandbox
WORKDIR /sandbox

# 正規化ツールが使用できるライブラリ (ネットワーク不要なものに限定)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY runner.py .

USER sandbox

# entrypoint: Cloud Run Jobs から環境変数で制御される
# 環境変数:
#   GCS_SANDBOX_BUCKET  - 作業バケット
#   TOOL_SCRIPT_PATH    - 実行する Python スクリプトの GCS パス
#   INPUT_GCS_PATH      - 入力ファイルの GCS パス
#   OUTPUT_GCS_PATH     - 出力先 GCS パス
#   DRY_RUN             - "true" の場合は差分のみ出力し書き込まない

CMD ["python", "runner.py"]
