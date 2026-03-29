"""
sandbox/runner.py - 正規化ツール サンドボックス実行エントリポイント

フロー:
  1. GCS から入力ファイルをダウンロード
  2. GCS からツールスクリプトをダウンロード
  3. 制限された exec 環境でスクリプトを実行
  4. dry_run=false の場合、結果を GCS にアップロード
  5. 差分サマリを標準出力に JSON で出力 (バックエンドがパース)
"""

import difflib
import json
import os
import sys
import tempfile
import traceback
from pathlib import Path

from google.cloud import storage

# 許可する組み込み関数・モジュール (ツールスクリプトに公開する名前空間)
_ALLOWED_BUILTINS = {
    "len", "str", "int", "float", "bool", "list", "dict", "set", "tuple",
    "range", "enumerate", "zip", "map", "filter", "sorted", "reversed",
    "print", "repr", "isinstance", "hasattr", "getattr", "type",
    "min", "max", "sum", "abs", "round",
}


def download_gcs(bucket_name: str, blob_path: str) -> bytes:
    client = storage.Client()
    bucket = client.bucket(bucket_name)
    return bucket.blob(blob_path).download_as_bytes()


def upload_gcs(bucket_name: str, blob_path: str, data: bytes) -> None:
    client = storage.Client()
    bucket = client.bucket(bucket_name)
    bucket.blob(blob_path).upload_from_string(data)


def run_tool(script: str, input_text: str) -> str:
    """
    ツールスクリプトを制限された名前空間で実行する。
    スクリプトは process(text: str) -> str 関数を定義する必要がある。
    """
    allowed_builtins = {k: __builtins__[k] for k in _ALLOWED_BUILTINS if k in __builtins__}  # type: ignore

    # 許可モジュールのみインポート可能
    import chardet
    import ftfy
    import unicodedata

    namespace = {
        "__builtins__": allowed_builtins,
        "chardet": chardet,
        "ftfy": ftfy,
        "unicodedata": unicodedata,
    }

    exec(script, namespace)  # noqa: S102

    if "process" not in namespace:
        raise ValueError("ツールスクリプトに process(text: str) -> str 関数が定義されていません")

    result = namespace["process"](input_text)
    if not isinstance(result, str):
        raise TypeError(f"process() は str を返す必要があります。実際: {type(result)}")
    return result


def compute_diff(original: str, modified: str) -> dict:
    diff = list(difflib.unified_diff(
        original.splitlines(keepends=True),
        modified.splitlines(keepends=True),
        fromfile="original",
        tofile="modified",
        n=3,
    ))
    return {
        "lines_changed": sum(1 for l in diff if l.startswith(("+", "-")) and not l.startswith(("+++", "---"))),
        "diff_preview": "".join(diff[:50]),  # 先頭50行のみ
    }


def main():
    bucket       = os.environ["GCS_SANDBOX_BUCKET"]
    tool_path    = os.environ["TOOL_SCRIPT_PATH"]
    input_path   = os.environ["INPUT_GCS_PATH"]
    output_path  = os.environ["OUTPUT_GCS_PATH"]
    dry_run      = os.environ.get("DRY_RUN", "false").lower() == "true"

    result = {
        "status": "completed",
        "dry_run": dry_run,
        "diff": None,
        "error": None,
    }

    try:
        script_bytes = download_gcs(bucket, tool_path)
        input_bytes  = download_gcs(bucket, input_path)

        script_text = script_bytes.decode("utf-8")
        input_text  = input_bytes.decode("utf-8", errors="replace")

        output_text = run_tool(script_text, input_text)
        result["diff"] = compute_diff(input_text, output_text)

        if not dry_run:
            upload_gcs(bucket, output_path, output_text.encode("utf-8"))

    except Exception as e:
        result["status"] = "failed"
        result["error"] = traceback.format_exc()
        print(json.dumps(result), flush=True)
        sys.exit(1)

    print(json.dumps(result), flush=True)


if __name__ == "__main__":
    main()
