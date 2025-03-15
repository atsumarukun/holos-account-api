# アカウントAPI

## 開発環境

以下のDesignDocに詳細を記載.<br />
[development-environment.md](https://github.com/atsumarukun/holos-account-api/blob/main/docs/design-docs/development-environment/development-environment.md)

## ルール

### ブランチ名

Issue番号を用いたブランチ名にする.

```
issue-${ISSUE_NUMBER}   // issue-1
```

### コミットメッセージ

以下の表に記載するタグを利用する.

| タグ | 説明 |
| --- | --- |
| create | 機能作成 |
| update | 機能更新 |
| remove | 機能削除 |
| refactor | リファクタリング |
| fix | 不具合修正 |

```
${TAG}: ${MESSAGE}   // create: ユーザー削除機能を作成.
```
