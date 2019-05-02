# nuxt-vue-go-chat

[![Build Status](https://travis-ci.org/sekky0905/nuxt-vue-go-chat.svg?branch=master)](https://travis-ci.org/sekky0905/nuxt-vue-go-chat)

Nuxt(Vue.js) + Goでの簡単なチャットアプリ。

# サーバサイド

## 技術

### 言語
Go

### フレームワーク
[Gin](https://github.com/gin-gonic/gin)


## アーキテクチャ
`Layered Architecture` をベースにする。
ただし、レイヤ間の結合度を下げるために各レイヤ間でDIPを行う。

```
├── interface
│   └── controller // サーバへの入力と出力を行う責務。
├── application // 薄く保ち、やるべき作業の調整を行う責務。
├── domain
│   ├── model // ビジネスの概念とビジネスロジック。
│   ├── service // EntityでもValue Objectでもないドメイン層のロジック。
│   └── repository // infra/dbへのポート。
└── infra // 技術的なものの提供
    ├── db // DBの技術に関すること。
    └── router // Routingの技術に関すること。
```

参考

InfoQ.com、徳武 聡(翻訳) (2009年6月7日) 『Domain Driven Design（ドメイン駆動設計） Quickly 日本語版』 InfoQ.com

エリック・エヴァンス(著)、 今関 剛 (監修)、 和智 右桂 (翻訳) (2011/4/9)『エリック・エヴァンスのドメイン駆動設計 (IT Architects’Archive ソフトウェア開発の実践)』 翔泳社

『pospomeのサーバサイドアーキテクチャ』pospome著

[Goを運用アプリケーションに導入する際のレイヤ構造模索の旅路 | Go Conference 2018 Autumn 発表レポート - BASE開発チームブログ](https://devblog.thebase.in/entry/2018/11/26/102401)

[ボトムアップドメイン駆動設計 │ nrslib](https://nrslib.com/bottomup-ddd/)

# フロントエンド

## 技術

### 言語
JavaScript

### ライブラリ/フレームワーク
主に使用しているもののみ記載。

- [Vue.js](https://jp.vuejs.org/index.html)
- [Nuxt.js](https://ja.nuxtjs.org/)
- [vuetify](https://github.com/vuetifyjs/vuetify)
