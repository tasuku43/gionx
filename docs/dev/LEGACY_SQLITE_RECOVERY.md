---
title: "Legacy SQLite Recovery"
status: implemented
---

# Legacy SQLite Recovery

## 対象

開発中の旧バージョンで SQLite 由来ファイルを生成済みの環境を、Filesystem-first運用へ切り替える手順。

## 方針

- 真実は `KRA_ROOT` 配下のディレクトリと `.kra.meta.json`。
- SQLite state store は廃止済み。旧 `state.db` が残っていても runtime は参照しない。
- まずバックアップを作成し、段階的に切り替える。

## 手順

1. 対象 root を確定する
```sh
kra context current
```

2. 旧データをバックアップする
```sh
mkdir -p .kra-recovery-backup
cp -a workspaces archive .kra-recovery-backup/
```

3. 旧 SQLite ファイルを退避する（存在する場合）
```sh
LEGACY_SQLITE_FILES="$(find "${XDG_DATA_HOME:-$HOME/.local/share}/kra" -name '*.db' 2>/dev/null)"
if [ -n "$LEGACY_SQLITE_FILES" ]; then
  mkdir -p .kra-recovery-backup/legacy-index
  while IFS= read -r f; do
    [ -n "$f" ] && cp -a "$f" .kra-recovery-backup/legacy-index/
  done <<EOF
$LEGACY_SQLITE_FILES
EOF
fi
```

4. `registry.json` の旧 `state_db_path` は放置でよい
- 新しい実装は `root_path` / timestamps を正準に扱う。
- コマンド実行時に registry が更新される。

5. 動作確認
```sh
kra ws list
kra ws select --act go
```

6. 問題なければ旧 SQLite ファイルを削除（任意）
```sh
if [ -n "$LEGACY_SQLITE_FILES" ]; then
  while IFS= read -r f; do
    [ -n "$f" ] && rm -f "$f"
  done <<EOF
$LEGACY_SQLITE_FILES
EOF
fi
```

## トラブル時

- `registry.json` 破損エラーが出る場合:
  - `XDG_DATA_HOME/kra/registry.json` をバックアップして削除し、再実行。
- `ws --act reopen` が失敗する場合:
  - 対象 `archive/<id>/.kra.meta.json` の `repos_restore` を確認。
  - bare repo が不足している場合は `kra repo add` で再登録してから再実行。
