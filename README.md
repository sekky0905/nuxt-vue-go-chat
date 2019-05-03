# Nuxt Vue Go Chat

[![Build Status](https://travis-ci.org/sekky0905/nuxt-vue-go-chat.svg?branch=master)](https://travis-ci.org/sekky0905/nuxt-vue-go-chat)

## 概要

Nuxt.js(Vue.js)とLayered Architectureのお勉強のために作成した簡単なチャットアプリ。

## 技術構成

SPA(Nuxt.js(Vue.js)) + API(Go) + RDB(MySQL)という形になっている。

### フロントエンド

#### 使用技術

- [Vue.js](https://jp.vuejs.org/index.html)
- [Nuxt.js](https://ja.nuxtjs.org/)
- [vuetify](https://github.com/vuetifyjs/vuetify)

主に使用しているもののみ記載している。

### サーバーサイド

- [Go](https://github.com/golang/go)
- [gin](https://github.com/gin-gonic/gin)

## その他

Update系の処理を画面側で実装していないが、サーバー側で実装しているのはあくまでもお勉強のため。

## アーキテクチャ
`Layered Architecture` をベースにする。
ただし、レイヤ間の結合度を下げるために各レイヤ間でDIPを行う。

```
├── interface
│   └── controller // サーバへの入力と出力を扱う責務。
├── application // 薄く保ち、やるべき作業の調整を行う責務。
├── domain
│   ├── model // ビジネスの概念とビジネスロジック。
│   ├── service // EntityでもValue Objectでもないドメイン層のロジック。
│   └── repository // infra/dbへのポート。
├── infra // 技術的なものの提供
│    ├── db // DBの技術に関すること。
│    └── router // Routingの技術に関すること。 
├── middleware // リクエスト毎に差し込む処理をまとめたミドルウェア
├── util 
└── testutil
```

## 使い方

### ローカルでの起動方法

```bash
cd server
make deps
make run
```

### テスト

```bash
cd server
make test
```

### 静的解析

```bash
cd server
make check
```

## 参考文献

### サーバーサイド

- InfoQ.com、徳武 聡(翻訳) (2009年6月7日) 『Domain Driven Design（ドメイン駆動設計） Quickly 日本語版』 InfoQ.com

- エリック・エヴァンス(著)、今関 剛 (監修)、和智 右桂 (翻訳) (2011/4/9)『エリック・エヴァンスのドメイン駆動設計 (IT Architects’Archive ソフトウェア開発の実践)』 翔泳社

- pospome『pospomeのサーバサイドアーキテクチャ』

### フロントエンド

- 花谷拓磨 (2018/10/17)『Nuxt.jsビギナーズガイド』シーアンドアール研究所 

- 川口 和也、喜多 啓介、野田 陽平、 手島 拓也、 片山 真也(2018/9/22)『Vue.js入門 基礎から実践アプリケーション開発まで』技術評論社




## 参考にさせていただいた記事

- [Goを運用アプリケーションに導入する際のレイヤ構造模索の旅路 | Go Conference 2018 Autumn 発表レポート - BASE開発チームブログ](https://devblog.thebase.in/entry/2018/11/26/102401)

- [ボトムアップドメイン駆動設計 │ nrslib](https://nrslib.com/bottomup-ddd/)

- [GoでのAPI開発現場のアーキテクチャ実装事例 / go-api-architecture-practical-example - Speaker Deck](https://speakerdeck.com/hgsgtk/go-api-architecture-practical-example)

- [GoのAPIのテストにおける共通処理 – timakin – Medium](https://medium.com/@timakin/go-api-testing-173b97fb23ec)

- [CSSで作る！吹き出しデザインのサンプル19選](https://saruwakakun.com/html-css/reference/speech-bubble)
    - commentの吹き出し部分をかなり参考にさせていただいている
