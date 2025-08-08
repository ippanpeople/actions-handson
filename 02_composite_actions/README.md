# 02 Composite Actions の紹介
前のステップでは、GitHub Actions の基本設定を行い、`run` でランナー内にコマンドを実行しました。このステップでは、Composite Actions を活用し、再利用可能なアクションによってワークフローの可読性と保守性が向上することを学びます。

> [!IMPORTANT]
> **目的**: 実務では複数のステップから成る複雑な Workflow を扱うことが一般的です。Composite Actions を活用することで、可読性と保守性の高い Workflow を構築する方法を学びます。次のステップで利用する sacloud_apprun_actions も同様の手法で構築されているため、その違いを意識できるようにする

> [!IMPORTANT]
> **ゴール**: Composite Actions を使った場合と使わなかった場合の違いをイメージとして理解できるようになること

## Composite Actions とは？
Composite Actions は、複数のシェルステップを組み合わせて作成されるカスタムアクションです。YAML で記述され、GitHub Actions のワークフロー内で再利用可能なロジックを提供します。これにより、複雑な処理を簡潔にまとめ、再利用性を高めることができます。

>[!NOTE]
> 今回通知用に利用する actions: [slack-notify-action](https://github.com/ippanpeople/slack-notify-action/blob/master/action.yml)

## Composite Actions の基本構成要素
Composite Actions は、以下の要素で構成されます。
- **name**: アクションの名称
- **description**: アクションの説明
- **inputs**: アクションが受け取る入力パラメータ
- **runs**: アクションが実行するコマンドやスクリプト（シェルコマンドや他のアクションの呼び出しが可能）

## 前のステップとの比較
![前のステップと比較](compare.png)

## To Do リスト
- [ ] ワークフローを設定し、Slack へメッセージを送信できるようにする
    - [ ] ワークフロー内の message を自分のメッセージに更新
    - [ ] 変更をリポジトリにプッシュ
- [ ] ワークフローが Slack へメッセージを送信できるかテスト
    - [ ] リポジトリの Actions タブに移動
    - [ ] `02 Composite Actions`を選択し「Run workflow」をクリック
    - [ ] Slack チャンネルにメッセージが表示されることを確認


