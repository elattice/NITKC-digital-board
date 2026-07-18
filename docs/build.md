# macOS・Ubuntu向けビルド手順

Reactの本番ビルドを埋め込んだGoバイナリを、同じソースコードから次の2種類作成します。

```text
dist/digital-board-macos-arm64
dist/digital-board-linux-amd64
```

`modernc.org/sqlite` はPure Go実装のため、macOSとUbuntuのどちらをビルド環境にしても、Cコンパイラなしで両方のバイナリをクロスビルドできます。

## 対象環境

| 出力ファイル | 対象 |
| --- | --- |
| `digital-board-macos-arm64` | Apple Silicon搭載Mac |
| `digital-board-linux-amd64` | 一般的なIntel / AMD 64bit Ubuntu PC |

## 必要なもの

- Go 1.26以上
- Node.js 20以上
- npm

初回、または `frontend/package-lock.json` が更新されたときは、プロジェクト直下で依存パッケージをインストールします。

```bash
npm --prefix frontend ci
```

## 両方のバイナリを作成する

macOSまたはUbuntuのプロジェクト直下で、次を実行します。

```bash
./scripts/build-release.sh
```

スクリプトは次の順序で処理します。

1. `npm run build` でReactを `backend/internal/webui/dist/` へ出力する
2. `GOOS=darwin GOARCH=arm64` でApple Silicon Mac用バイナリを作る
3. `GOOS=linux GOARCH=amd64` でUbuntu用バイナリを作る

既存の同名ファイルは新しいビルド結果で置き換えられます。

## 出力の確認

```bash
ls -lh dist/digital-board-macos-arm64 dist/digital-board-linux-amd64
file dist/digital-board-macos-arm64 dist/digital-board-linux-amd64
```

`file` の期待結果:

- macOS版: `Mach-O 64-bit executable arm64`
- Ubuntu版: `ELF 64-bit LSB executable, x86-64`

## macOSで起動する

```bash
cd dist
./digital-board-macos-arm64
```

ブラウザで次を確認します。

- 掲示板: <http://localhost:8080/>
- 管理画面: <http://localhost:8080/admin>

## Ubuntuへ配置して起動する

`dist/digital-board-linux-amd64` をUbuntu PCの任意のディレクトリへコピーします。コピー方法によって実行権限が失われた場合は、権限を付け直します。

```bash
chmod +x digital-board-linux-amd64
./digital-board-linux-amd64
```

Reactのファイルはバイナリ内に埋め込まれているため、`frontend/` や `backend/internal/webui/dist/` をUbuntu PCへコピーする必要はありません。

## SQLite DBの保存場所

SQLite DBはバイナリには埋め込まれません。起動した作業ディレクトリを基準に、次の場所へ作成・保存されます。

```text
data/timetable.db
```

本番ではDBの保存場所が変わらないよう、毎回同じディレクトリから起動してください。systemdを使用する場合も `WorkingDirectory` を固定します。

## 配布前のチェックサム作成

macOS:

```bash
shasum -a 256 dist/digital-board-macos-arm64 dist/digital-board-linux-amd64
```

Ubuntu:

```bash
sha256sum dist/digital-board-macos-arm64 dist/digital-board-linux-amd64
```
