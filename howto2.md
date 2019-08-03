# AWS Githubからソースを取得し、ESC(Fargate or EC2)にデプロイするパイプラインを作成する。


概要をかく

経緯とかも書く

---

# やりたいこと


# こんなことをやればできるはず

1. GitHubのリポジトリを作成。
2. AWSユーザを作成する。
3. Route53でドメイン発行。SSLの対応も
4. VPCとサブネットを作成する
5. RDB周りの設定
6. インターネットゲートウェイとルートテーブルの設定
7. ターゲットグループの作成
8. ALBを設定
9. ソースの用意
10. Dockerイメージの作成
11. ECRに追加
12. ECSクラスタを作成
13. ECSのタスクを作成
14. ECSのサービスを作成
15. CodeBuildでBuildプロジェクトを作成
16. Route53でサブドメインを切って、リスナーに登録する。
17. パイプラインを作成（ソース取得）
18. パイプラインを作成（ビルド）
19. パイプラインを作成（ECSにデプロイ）


# ソースの用意
今回はgolangのhelloworldで
hellogoというプロジェクト名でmain.goを作成

main.go
```
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

```

docker-compose.yml
```
version: '3'
services:
  api:
    build: .
    command: go run main.go 
    volumes:
      - .:/hellogo 
    ports:
      - 8080:8080
```

Dockerfile
```
FROM golang:1.10.0


WORKDIR /go
ADD . /go

EXPOSE  8080

CMD ["go", "run", "main.go"]

```

buildspec.yml
```
version: 0.2

phases:
  pre_build:
    commands:
      - echo Logging in to Amazon ECR...
      - echo $AWS_DEFAULT_REGION
      - $(aws ecr get-login --no-include-email --region $AWS_DEFAULT_REGION)
  build:
    commands:
      - echo Build started on `date`
      - echo Building the Docker image...
      - docker build --no-cache -t $IMAGE_REPO_NAME:$IMAGE_TAG .
      - docker tag $IMAGE_REPO_NAME:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME:$IMAGE_TAG
  post_build:
    commands:
      - echo Build completed on `date`
      - echo Pushing the Docker image...$AWS_ACCOUNT_ID....$IMAGE_REPO_NAME:$IMAGE_TAG
      - docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME:$IMAGE_TAG
      - REPOSITORY_URI=$AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME
      - echo "[{\"name\":\"${CONTAINAR_NAME}\",\"imageUri\":\"${REPOSITORY_URI}:${IMAGE_TAG}\"}]" > imagedefinitions.json
artifacts:
  files: imagedefinitions.json
  
```


- portは8080:8080でListen
- githubにプッシュする



# ECRリポジトリを作成
リポジトリ名：hellogo-image-repo

![ecr1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/ecr1.png)


# CodeBuildでビルドプロジェクトを作成

## プロジェクトの設定
プロジェクト名：hellogo-build-project



## 送信元
ソースプロバイダ：GitHub
リポジトリ：GitHubアカウントリポジトリ
GitHubリポジトリ：hellogo

## 環境
環境イメージ：マネージド型イメージ
オペレーティングシステム：Ubuntu
ランタイム：Standard
イメージ：aws/codebuild/standard:1.0
イメージのバージョン：aws/codebuild/standard:1.0-1.8.0
特権付与：チェック
サービスロール：既存のサービスロール
ロール名：CodeBuildServiceRole ※ない場合は作る必要あり。
環境変数：以下を参照
AWS_DEFAULT_REGION：ap-northeast-1
AWS_ACCOUNT_ID：AWSユーザIDを設定
IMAGE_REPO_NAME：hellogo-image-repo
IMAGE_TAG：latest
CONTAINAR_NAME:hellogo-containar
※全てプレーンテキスト


## Buildspec
ビルド仕様：buildspecファイルを仕様する

## アーティファクト
タイプ：アーティファクトなし


# ECSタスクを作成
AmazonECSのサイドバータスク定義より新しいタスクの定義を作成する。


## 起動タイプの互換性の選択
FARGATE

## タスクとコンテナの定義の設定
タスク定義名：hellogo-task
タスクロール：なし
ネットワークモード：awsvpc
タスクの実行ロール：ecsTaskExectionRole ※ない場合は作成する。
タスクメモリ：0.5GB　※ビルドするだけのタスクなのでとりあえずミニマム
タスクCPU：0.25 vCPU　※ビルドするだけのタスクなのでとりあえずミニマム
コンテナの定義：
  コンテナ名：hellogo-containar
  イメージ：上記で作成したECRリポジトリのlatest
  ポートマッピング：8080

![tg1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/tg1.png)





# ECSクラスタ作成
クラスターテンプレートの選択：ネットワーキングのみ

## クラスタの設定
クラスタ名：hellogo-cluster

# サービスの作成
上記で作成したクラスタにサービスを作成する。

## サービスの設定
起動タイプ：FARGATE
タスク定義（ファミリー）：hellogo-task
リビジョン：latest
プラットフォームのバージョン：LATEST
クラスタ：hellogo-cluster
サービス名：hellogo-service
タスク数：1

## ネットワーク構成
クラスタVPC：前回作成したVPC
サブネット：全部設定しておく
セキュリティグループ：全部設定しておく
パブリックIPの自動割り当て：ENABLED
ロードバランサーの種類：Application Load Balancer
ロードバランサー名：web-app-alb
コンテナの選択：hellogo-containar:8080:8080
ロードバランサーに追加をクリック
ターゲットグループ名：web-app-tg-1
![service1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/service1.png)




# CodePipelineでパイプラインを作成する

## パイプラインの設定を選択する
パイプライン名：hellogo-pipeline
ロール名：既存のサービスロールからAWSCodePipelineServiceを選択


## ソースステージ
ソースプロバイダー：GitHub
リポジトリとブランチを選択しGitHubウェブフックにチェックをつける

## ビルドステージ
プロバイダーを構築する：AWSCodeBuild
リージョン：アジアンパシフィック（東京）
プロジェクト名：hellogo-build-project

## デプロイステージ
デプロイプロバイダー：AmazonECS
リージョン：アジアンパシフィック（東京）
クラスタ名：hellogo-cluster
サービス名：hellogo-service
![pipeline1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/pipeline1.png)



- 実行完了

![pipeline2.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/pipeline2.png)

- この時点でロードバランサのDNS 名でアクセスできます。
- 20分くらいかかってしまいました。



# サブドメインを設定し、リバースプロキシ対応する


## サブドメイン作成
Route53のサイドバーからホストゾーンを選択し、対象のドメインをクリック
レコードセットの作成をクリック
名前：サブドメインを入力　※今回はhellogoにしておく
エイリアス：はい
エイリアス先：web-app-alb　※前回作成したALBを設定


## ロードバランサの設定
EC2のサイドバーよりロードバランサを選択し、一覧上の対象のロードバランサの下部リスナータグより
ルールの表示編集にて設定をする。



![host-rule1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/host-rule1.png)


- これでHTTPでアクセスした場合、問題なくhelloworldされます。

- httpsも同様に設定する。
![host-rule2.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/host-rule2.png)


- httpできた場合はhttpsへリダイレクトさせたいのでｈttpの設定を変更
![host-rule3.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/host-rule3.png)




---
以上















