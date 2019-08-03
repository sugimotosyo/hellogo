# AWS Githubからソースを取得し、ESC(Fargate or EC2)にデプロイするパイプラインを作成する。（ALBの作成まで）


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


# route53でドメインを取得して、証明書を発行する。
かくこと


# VPC周り
[VPCとVPCサブネットの作成、EC2とRDSの構築、ELB設定まで \| Qrunch（クランチ）](https://qrunch.net/@hikey/entries/iZgLeVdEHqKiz2cu)

ユーザ：develop-deploy

## VPC作成
![vpc1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc1.png)



## VPCサブネットを作成

名前タグ：vpc-web-app-subnet-web-a
VPC:上記で作成したVPC
アベイラビリティーゾーン：ap-northeneast-1a
IPv4 CIDR ブロック：10.0.0.0/24

![vpc-subnet1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-subnet1.png)



名前タグ：vpc-web-app-subnet-web-c
VPC:上記で作成したVPC
アベイラビリティーゾーン：ap-northeneast-1c
IPv4 CIDR ブロック：10.0.1.0/24

![vpc-subnet2.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-subnet2.png)



名前タグ：vpc-web-app-subnet-db-a
VPC:上記で作成したVPC
アベイラビリティーゾーン：ap-northeneast-1a
IPv4 CIDR ブロック：10.0.10.0/24


![vpc-subnet3.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-subnet3.png)



名前タグ：vpc-web-app-subnet-db-c
VPC:上記で作成したVPC
アベイラビリティーゾーン：ap-northeneast-1c
IPv4 CIDR ブロック：10.0.11.0/24

![vpc-subnet4.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-subnet4.png)




## セキュリティグループを追加する

セキュリティグループ名：web-app-lb-sg
説明：ALB
VPC：上記で作成したVPC
※作成後一覧よりNameにセキュリティグループ名を設定

![vpc-sg1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-sg1.png)

同様に下記設定でも作成する

セキュリティグループ名：web-app-web-sg
説明：ECS-WEB
VPC：上記で作成したVPC
※作成後一覧よりNameにセキュリティグループ名を設定


セキュリティグループ名：web-app-db-sg
説明：DB
VPC：上記で作成したVPC
※作成後一覧よりNameにセキュリティグループ名を設定


セキュリティグループ名：web-app-ssh-sg
説明：ECS-SSH 
VPC：上記で作成したVPC
※作成後一覧よりNameにセキュリティグループ名を設定

![vpc-sg2.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-sg2.png)



## セキュリティグループにインバウンドルールを追加する。
一覧の下部インバンンドルールのタブより追加する。

web-app-lb-sg
![vpc-sg3.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-sg3.png)



web-app-web-sg
タイプ：カスタムTCPルール
プロトコル：TCP
ポート範囲：8080
ソース：カスタム　web-app-lb-sgのグループID

![vpc-sg4.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-sg4.png)

web-app-db-sg
タイプ：MYSQL/Aurora
プロトコル：TCP
ポート範囲：3306
ソース：カスタム　web-app-web-sgのグループID
![vpc-sg5.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-sg5.png)



web-app-ssh-sg
タイプ：SSH
プロトコル：TCP
ポート範囲：22
ソース：カスタム　0.0.0.0/0, ::/0
※ステップサーバー等を置いている場合はソースに適切なIPを設定してください。
![vpc-sg6.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/vpc-sg6.png)



# RDS周りの設定（今回のサンプルアプリでは使わない予定）

## RDSサブネットグループを作成
Amazon RDSのサイドバーからサブネットグループを選択

名前：web-app-rds-subnet-grp
VPC：上記で作成したVPC
サブネット：vpc-web-app-subnet-db-a、vpc-web-app-subnet-db-c ※サブネットIDしか出てこないのでVPCのサブネットからサブネットIDを確認して設定

![db-subnet1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/db-subnet1.png)



## RDSインスタンスを作成する
データベースの作成より作成する
エンジンのタイプ：Mysql ※お好きなものを選択してください。
テンプレート：無料利用枠　※私の場合はこれでOK
DBインスタンス識別子：web-app-db-1 ※適当に設定
ユーザー名とパスワードを入力
VPC：上記で作成したVPC
サブネットグループ：web-app-rds-subnet-grp
パブリックアクセス可能：あり　※適宜設定
VPCセキュリティグループ：既存の選択
VPCのセキュリティーグループ：web-app-db-sg ※defaultも設定しておくが、いらないと思う。
AZ：指定なし
データベースポート：3306


# インターネットゲートウェイとルートテーブルの設定
## インターネットゲートウェイの作成
VPCのサイドバーからインターネットゲイトウェイを選択し、作成する。
名前タグ：web-app-igw
![igw1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/igw1.png)

## インターネットゲートウェイをVPCにアタッチする
インターネットゲートウェイの一覧より上記で作成したインターネットゲートウェイを選択し、アクションよりVPCにアタッチする。
VPC：上記で作成したVPC

## ルートテーブルを作成する

名前タグ：web-app-rtb-global
VPC：上記で作成したVPC
![rtb1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/rtb1.png)


名前タグ：web-app-rtb-local
VPC：上記で作成したVPC

![rtb2.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/rtb2.png)


## ルートの設定
作成したルートテーブルを選択し、下部ルートタブより編集で設定

web-app-rtb-global
送信先：0.0.0.0/0
ターゲット：web-app-igw

web-app-rtb-local
送信先：10.0.0.0/16	※デフォルトで設定されていると思う。
ターゲット：local　※デフォルトで設定されていると思う。


## ルートテーブルをサブネットサブネットの関連を設定
作成したルートテーブルを選択し、下部サブネットの関連付けタブより編集で設定

web-app-rtb-global
VPCサブネット：vpc-web-app-subnet-web-a、vpc-web-app-subnet-web-c


web-app-rtb-local
VPCサブネット：vpc-web-app-subnet-db-a、vpc-web-app-subnet-db-c


# ターゲットグループの作成（アプリが増えるごとにやる必要がある？）
EC2のサイドバーからターゲットグループを選択し作成する。

ターゲットグループ名：web-app-tg-1
ターゲット種類：IP
プロトコル：HTTP
ポート：80
VPC：上記で作成したVPC

![tg1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/tg1.png)




# ALBを設定

EC2のサイドバーからロードバランサーを選択し作成する。

ロードバランサーの種類：Application Load Balancer

## 手順1ロードバランサーの設定
名前：web-app-alb
リスナー(HTTP)：プロトコル :HTTP ロードバランサーのポート：80
リスナー(HTTPS)：プロトコル :HTTPS ロードバランサーのポート：443
VPC：上記で作成したVPC
アベイラビティーゾーン（ap-northeast-1a）： サブネット：vpc-web-app-subnet-web-a
アベイラビティーゾーン（ap-northeast-1c）： サブネット：vpc-web-app-subnet-web-c
![alb1.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/alb1.png)

## 手順2セキュリティ設定の構成
証明書タイプ：ACMから証明書を選択する
証明書の名前：証明書を選択
※Route53よりドメインを発行し、作成した証明書を選択。証明書の作り方は別途記載しておく。


## 手順3セキュリティグループの設定
既存のセキュリティグループより「web-app-lb-sg」を選択する。

## 手順4ルーティングの設定
ターゲットグループ：既存のターゲットグループ
名前：web-app-tg-1
![alb2.png](https://raw.githubusercontent.com/sugimotosyo/hellogo/master/sample-image/alb2.png)

## 手順5ターゲットの登録
この時点ではECSインスタンスを生成していないので、登録済みのターゲットがないので、設定しない。
ECSで自動でターゲット登録される。




# ここまでの参考は
[VPCとVPCサブネットの作成、EC2とRDSの構築、ELB設定まで \| Qrunch（クランチ）](https://qrunch.net/@hikey/entries/iZgLeVdEHqKiz2cu)

[AWSのサービスでドメインを取得しALBでSSLで接続出来るようにする - Qiita](https://qiita.com/keitakn/items/4b2db95eae81044a779c)



















