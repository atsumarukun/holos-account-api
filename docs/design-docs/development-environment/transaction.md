# 概要

データベースのトランザクションを実装する.

# 対象範囲

## 達成基準

- トランザクション関数を呼び出せる状態

## 除外項目

- 実際にDBを操作するリポジトリの作成は行わない

# 利用方法

Usecase層から以下関数を呼び出すことでトランザクションを行う.

```golang
err := transactionObject.Transaction(ctx, func(ctx context.Context) error {
  // ここに実装
})
```

Infrastructure層での利用方法は以下の通り.

```golang
driver := getDriver(ctx, r.db)
// ここに実装
```

# 詳細設計

トランザクションオブジェクトをcontextに保持させることでInfrastructure層に渡す.<br />
コンテキストにトランザクションオブジェクトが含まれる場合はトランザクションオブジェクトを、含まれない場合はデータベースオブジェクトをドライバーとして返却する.

# その他の手法

# 参考文献

# 変更履歴

| 変更日 | 変更者 | 変更内容 |
| --- | --- | --- |
| 2025/03/16 | @atsumarukun | 初版 |
