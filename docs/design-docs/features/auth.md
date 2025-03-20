# 概要

認証認可機能を作成する.

# 対象範囲

## 達成基準

- ログイン用エンドポイントが作成されている
- ログアウト用エンドポイントが作成されている
- 認可用エンドポイントが作成されている

## 除外項目

- アカウント機能の対応は行わない

# 利用方法

| パス | メソッド | 備考 |
| --- | --- | --- |
| /signin | POST | ログイン |
| /signout | DELETE | ログイン |
| /authorization | GET | 認可 |

## シーケンス

```mermaid
sequenceDiagram
  participant client as クライアント
  participant server as サーバー

  client ->>+ server: ① ログイン
  server -->>- client: token
  client ->>+ server: ② ①のtokenを用いて認可
  Note over client, server: AuthorizationHeader: Session ${TOKEN}
  server -->>- client: account_id
```

# 詳細設計

## 要件

- パスワードを用いてログインを行う
- ログアウトを行える
- 認可を行える

## 仕様

- ログイン時にランダムな文字列のトークンを発行する
  - トークンは32文字
  - トークンの有効期限は1週間
- トークンを削除することでログアウトを行う
- トークンを用いて認可を行う
  - アカウントIDを返却する

## ドメインオブジェクト

| キー | 型 | 備考 |
| --- | --- | --- |
| id | uuid | |
| account_id | uuid | |
| token | string | 32文字 |
| expires_at | time | 1週間 |

## テーブル

| カラム名 | 型 | キー | null許容 | 備考 |
| --- | --- | --- | :---: | --- |
| id | char(36) | PK | | ID |
| account_id | char(36) | FK, UQ | | アカウントID |
| token | char(32) | UQ | | トークン |
| expires_at | datetime(6) | | | 有効期限 |
| created_at | datetime(6) | | | 作成日時 |
| updated_at | datetime(6) | | | 更新日時 |
| deleter_at | datetime(6) | | * | 削除日時 |

## テスト項目

### ログイン

| ステータス | テストケース | 入力値 |
| --- | --- | --- |
| 正常系 | ログイン成功 | account_name: name<br />password: password |
| 異常系 | アカウントが存在しない | account_name: NAME<br />password: password |
| 異常系 | パスワードが異なる | account_name: name<br />password: PASSWORD |

### ログアウト

| ステータス | テストケース | 入力値 |
| --- | --- | --- |
| 正常系 | ログアウト成功 | token |
| 異常系 | セッションが存在しない | token |

### 認可

| ステータス | テストケース | 入力値 |
| --- | --- | --- |
| 正常系 | 認可成功 | token |
| 異常系 | セッションが存在しない | token |

# その他の手法

# 参考文献

# 変更履歴

| 変更日 | 変更者 | 変更内容 |
| --- | --- | --- |
| 2025/03/16 | @atsumarukun | 初版 |
| 2025/03/20 | @atsumarukun | テーブル構造を変更 |
