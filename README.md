# Fumi-chan searcher
## 概要
二川二水ちゃん(@assault_lily)がツイッターに投稿しているアサルトリリィの設定を、簡単に検索するためのアプリです。

現状ツイッター検索しようとした時に起きる下記の課題を解決するために作りました。
1. ツイッターのツイート表示は3,200件の上限があり、上限以上のツイートを検索しようとすると高度な検索が必要
1. 二水ちゃんがツイートをコレクションにまとめてくれているが、閲覧できるツイートが各コレクションの間で抜け落ちている(ツイッターの仕様？)
1. コレクションから目的のツイートを探せない

## デモページ
https://fumichansearcher.com/

## 使い方
事前準備として、PostgreSQLをインストールし、DBを作成しておく
1. `.env_template`をリネームまたは、`.env`ファイルを新規作成してツイッターAPIのトークン等の必要情報を入れる。
2. 以下のコマンドを実行する。
``` golang
go run main.go
```
3. 登録したツイッターアカウントで取り込みたいツイートをいいねする
4. 一定周期でいいねしたツイートがアプリに取り込まれる