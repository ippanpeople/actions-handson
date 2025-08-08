# 03 sacloud_apprun_actions を使った Go アプリケーションのデプロイ
このステップでは、GitHub Actions を利用して Go アプリケーションをさくらの AppRun サービスにデプロイする方法を解説します。

> [!IMPORTANT]
> **目的**: GitHub Actions および Composite Actions の基礎を理解した上で、Composite Actions sacloud_apprun_actions を使い、データ永続化を含む Go アプリケーションを Sakura の AppRun サービスにデプロイする方法を学ぶ

> [!IMPORTANT]
> **ゴール**: GitHub Actions で Go アプリケーションを Sakura の AppRun にデプロイし、同時に Object Storage を作成してデータ永続化を実現する。これにより、元々手動で行っていたデプロイ作業を自動化し、効率的な開発フローを体験する

## sacloud_apprun_actions とは？
`sacloud_apprun_actions` は、Go アプリケーションをさくらの AppRun サービスにデプロイするときのワークフローを簡潔にするための Composite Actionsです。**アプリケーションのビルド、コンテナレジストリへのプッシュ、AppRun へのデプロイを自動化します**。また、**データ永続化のためのオブジェクトストレージバケットの作成**機能も含まれており、アプリケーションの再起動や再デプロイ後もデータの永続化ができます。

>[!NOTE]
> 今回利用する actions: [sacloud-apprun-action](https://github.com/ippanpeople/sacloud-apprun-action/blob/master/action.yaml)

## sacloud_apprun_actions の使い方
1. **ワークフロー例**: `sacloud_apprun_actions` ワークフローが含まれるリポジトリをフォークし、アプリケーションに合わせてワークフローファイルをカスタマイズします。
````yaml
      - name: Goアプリをデプロイ
        id: deploy
        uses: ippanpeople/sacloud-apprun-action@v0.0.4
        with:
          use-repository-dockerfile: false
          app-dir: ./03_sacloud_apprun_actions
          sakura-api-key: ${{ secrets.SAKURA_API_KEY }}
          sakura-api-secret: ${{ secrets.SAKURA_API_SECRET }}
          container-registry: ${{ secrets.REGISTRY }}
          container-registry-user: ${{ secrets.REGISTRY_USER }}
          container-registry-password: ${{ secrets.REGISTRY_PASSWORD }}
          port: '8080'
          # SQLite + Litestream を使う場合は以下も指定
          object-storage-bucket: ${{ secrets.STORAGE_BUCKET_NAME }}
          object-storage-access-key: ${{ secrets.STORAGE_ACCESS_KEY }}
          object-storage-secret-key: ${{ secrets.STORAGE_SECRET_KEY }}
          sqlite-db-path: ./data/app.db
          litestream-replicate-interval: 10s
````

> [!NOTE]
> 前ステップに紹介した　Composite Actions を利用することで、　元々[sacloud-apprun-action](https://github.com/ippanpeople/sacloud-apprun-action/blob/master/action.yaml)のように複数のステップを持つアクションを簡潔にまとめることができます。

2. **Secrets と Variables の設定**: リポジトリの設定で必要な GitHub Actions シークレットと変数を作成します。
   - コンテナレジストリの URL:
````
REGISTRY
````
   - コンテナレジストリのユーザー名:
````
REGISTRY_USER
````
   - コンテナレジストリのパスワード:
````
REGISTRY_PASSWORD
````
   - さくらの API キー:
````
SAKURA_API_KEY
````
   - さくらの API シークレット:
````
SAKURA_API_SECRET
````
   - オブジェクトストレージのバケット名:
````
STORAGE_BUCKET_NAME
````
   - オブジェクトストレージのアクセスキー:
````
STORAGE_ACCESS_KEY
````
   - オブジェクトストレージのシークレットキー:
````
STORAGE_SECRET_KEY
````
3. **ワークフローの実行**: `03 Sacloud Apprun Actions`を手動で実行します。

## データ永続化の実践方法
AppRun はステートレスなため、デプロイのたびにアプリケーションが再起動されます。データ永続化のため、SQLite と Litestream を利用し、アプリ再起動後もデータが保持されるようにします。
1. SQLite とは
   - SQLite は軽量なデータベースで、小規模アプリや組み込み用途に最適です。データは単一ファイルに保存され、管理やバックアップが容易です。
2. Litestream とは
   - Litestream は SQLite 用のレプリケーションツールで、SQLite データベースの変更をリアルタイムでリモートストレージ（S3, GCS など）に複製できます。これにより、アプリ再起動や再デプロイ後もデータが保持されます。
3. Object Storage の利用
   - 本実習では、さくらの Object Storage に SQLite データベースファイルを保存します。これにより、アプリ再起動後もデータが利用でき、バックアップやリストアも容易です。
4. SQLite + Litestream + Object Storage の組み合わせ
   - `sacloud_apprun_actions` でデプロイされたアプリは、起動時に Litestream が SQLite データベースの変更を Object Storage に自動複製します。SQLite の変更や一定間隔ごとにデータが Object Storage に同期されるため、アプリ再起動時もデータが復元されます。
   - AppRun の新バージョンをデプロイする際も、同じ SQLite ファイルパスと Litestream 設定を使うことで、データは新バージョンのアプリでも問題なく復元されます。

> [!WARNING]
> 注意: SQLite と Litestream を利用する際は、システム設計の妥当性に注意してください。SQLite は小規模用途に適していますが、高負荷や大量データではパフォーマンス上の制約があります。特に、**TPS (Transactions Per Second)** や、**QPS (Queries Per Second)**、などの観点で制約が発生しやすいです。Litestream も設定や運用方法によってはデータの安全性・一貫性 (**Consistency**) に注意が必要です。システム要件に応じて、適切なデータベースやストレージ方式の選定を検討してください。

## To Do リスト
- [ ] AppRun・オブジェクトストレージのデプロイ準備として Actions Secrets と Variables を設定
    - [ ] Actions Secret `REGISTRY` にコンテナレジストリの URL を登録
    - [ ] Actions Secret `REGISTRY_USER` にコンテナレジストリのユーザー名を登録
    - [ ] Actions Secret `REGISTRY_PASSWORD` にコンテナレジストリのパスワードを登録
    - [ ] Actions Secret `SAKURA_API_KEY` にさくらの API キーを登録
    - [ ] Actions Secret `SAKURA_API_SECRET` にさくらの API シークレットを登録
    - [ ] Actions Secret `STORAGE_BUCKET_NAME` にオブジェクトストレージのバケット名を登録
    - [ ] Actions Secret `STORAGE_ACCESS_KEY` にオブジェクトストレージのアクセスキーを登録
    - [ ] Actions Secret `STORAGE_SECRET_KEY` にオブジェクトストレージのシークレットキーを登録
- [ ] 03 Sacloud Apprun Actionsを実行し、AppRun がデプロイされ Slack に成功メッセージが送信されることを確認
    - [ ] リポジトリの Actions タブに移動
    - [ ] 対象ワークフローを選択し「Run workflow」をクリック
    - [ ] Slack チャンネルにメッセージが表示されることを確認
- [ ] AppRun のデプロイを確認
    - [ ] Slack メッセージ内の URL を確認
    - [ ] ブラウザで URL にアクセスし、デプロイ済みアプリを表示
    - [ ] ボードにメッセージを追加し、機能をテスト
- [ ] データ永続化のテスト
    - [ ] `main.go` ファイルを修正
    - [ ] ワークフローで再デプロイ
    - [ ] 再度 URL にアクセスし、以前追加したメッセージが表示されることを確認
