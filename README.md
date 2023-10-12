# Wakatime to Slack Profile
![Go Tests](https://github.com/walnuts1018/wakatime-to-slack-profile/actions/workflows/go-test.yaml/badge.svg)
[![Code Coverage](https://img.shields.io/codecov/c/github/walnuts1018/wakatime-to-slack-profile/master.svg)](https://codecov.io/github/walnuts1018/wakatime-to-slack-profile?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/walnuts1018/wakatime-to-slack-profile)](https://goreportcard.com/report/github.com/walnuts1018/wakatime-to-slack-profile)
![Last Image Build](https://github.com/walnuts1018/wakatime-to-slack-profile/actions/workflows/docker-image.yaml/badge.svg)
[![Latest Image](https://ghcr-badge.egpl.dev/walnuts1018/wakatime-to-slack-profile/latest_tag?trim=major&label=Latest%20Docker%20Image&color=ROYALBLUE&ignore=test-*)](https://ghcr-badge.egpl.dev/walnuts1018/wakatime-to-slack-profile/latest_tag?trim=major&label=Latest%20Docker%20Image&color=ROYALBLUE&ignore=test-*)

このプログラムは、Wakatime経由でユーザーが現在書いているコードを取得し、Slackのカスタムステータスに絵文字として設定するプログラムです。

![image](https://github.com/walnuts1018/wakatime-to-slack-profile/assets/60650857/e6044d30-5008-40b8-a0ba-8c0952fe2cee)

## Getting Started

### PostgreSQL
いい感じに用意してください。

### 環境変数
|env|sample|detail|
| --- | --- | --- |
|GIN_MODE|release|gin用release mode設定|
|WAKATIME_APP_ID|**********|Wakatime APIのApp ID|
|WAKATIME_CLIENT_SECRET|**********|Wakatime APIのClient Secret|
|COOKIE_SECRET|*************|Cookie用のSecret: 64Byteのランダム文字列を入れてください|
|PSQL_ENDPOINT|postgres-release-postgresql.databases.svc.cluster.local|Postgresqlのendpoint|
|PSQL_PORT|5432|PostgreSQLのポート|
|PSQL_DATABASE|wakatime_to_slack|PostgreSQLデータベース名|
|PSQL_USER|user|PostgreSQLユーザー名|
|PSQL_PASSWORD|**********|PostgreSQL ユーザーパスワード|
|SLACK_ACCESS_TOKEN|xoxp-********|SlackのACCESS Token|

### emoji.jsonを用意

`emoji.json`はSlackの絵文字名と言語名が対応していない場合に手動で指定するためのファイルです。

実行パスに配置
`Wakatimeにおける言語名`: `Slackの絵文字ID`
```json
{
  "Go": "gopher",
  "YAML": "k8s",
  "SQL": "postgresql"
}
```

絵文字は、

`emoji.jsonでの手動指定`→`言語名そのまま`→`言語名を全て小文字にしたもの`→`(絵文字が見つからなかった場合)❓`

の順番に探されます。

また、過去十分間にコードを書いた履歴がない場合は🦥になります。

### Start with Docker
対応arch: `amd64`, `arm64`
```bash
docker run -p 8080:8080 ghcr.io/walnuts1018/wakatime-to-slack-profile:latest
```

### ログイン
初回起動時にはブラウザでのログインが必要です。
ブラウザで [http://localhost:8080/signin](http://localhost:8080/signin) を開きます（他のURLで公開している場合は適宜指定してください。）



