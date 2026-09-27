# go-tel-num-parser-jp

日本の電話番号を形式で判定し、文字列から抽出する Go パッケージです。

```go
ok, kind := tnp.IsTelNumber("060-1234-5678") // true, tnp.MobilePhone
number, err := tnp.CropTelNumber("連絡先: 0800-123-4567") // "0800-123-4567", nil
```

`IsTelNumber` は文字列全体を判定します。`CropTelNumber` は文字列内で最初に見つかった番号を返します。括弧を使った `03(5321)1111` は `03-5321-1111` に変換して扱います。`SetIgnoreTypes(tnp.MobilePhone)` で特定の区分を判定対象から外せます。設定は呼び出すたびに置き換わり、`SetIgnoreTypes()` で解除できます。

## 対応する番号区分

2026年9月時点の番号計画・指定状況を参照しています。表中の例は形式の例であり、ハイフンなしの形式にも対応します。

| `TelType` | 区分 | 対応する形式の例 |
| --- | --- | --- |
| `FixedLinePhone` | 固定電話 | `03-5321-1111`、`011-896-4081` |
| `M2M` | データ伝送携帯電話 | `020-123-45678`（11桁）、`0200-12345-67890`（14桁） |
| `PocketBell` | 無線呼出 | `020-412-34567` |
| `IPPhone` | 特定 IP 電話 | `050-1234-5678` |
| `MobilePhone` | 音声伝送携帯電話 | `060-1234-5678`、`070-1234-5678`、`080-1234-5678`、`090-1234-5678` |
| `IncomingCharge` | 着信課金 | `0120-123-456`（10桁）、`0800-123-4567`（11桁） |
| `UnifiedNumber` | 統一番号 | `0570-123-456` |
| `InformationCharge` | 情報料代理徴収 | `0990-123-456` |
| `FMC` | FMC 電話 | `0600-123-4567` |

このパッケージは番号の形式と区分を判定します。実際に事業者へ指定済みか、利用中か、固定電話の市外局番が実在するかは検証しません。指定のない `0170`・`0180`、緊急通報などの短縮番号、国際番号は対象外です。`060` は番号計画上の携帯電話番号として判定します。事業者による利用開始状況とは区別してください。

## 参考資料

- [総務省「電気通信番号指定状況」](https://www.soumu.go.jp/main_sosiki/joho_tsusin/top/tel_number/number_shitei.html) — 番号区分と指定状況（2026年9月1日時点の公表資料）。
- [e-Gov データポータル「電気通信番号指定状況」](https://data.e-gov.go.jp/data/ja/dataset/soumu_20160325_0034) — `020C`（11桁）、`0200`（14桁）および付加的役務番号の区分。
- [総務省「携帯電話番号への060番号の追加」](https://www.soumu.go.jp/main_sosiki/joho_tsusin/top/tel_number/060keitai.html) — `060-1～9` の11桁番号と利用開始時期。
- [NTT東日本「フリーアクセス・ひかりワイド」](https://business.ntt-east.co.jp/service/hikari_of/free.html) — `0120` と `0800` の番号形式。
- [総務省の番号制度資料（e-Gov パブリック・コメント掲載）](https://public-comment.e-gov.go.jp/pcm/download?seqNo=0000302542) — 2025年の番号体系図。
- [総務省「電気通信番号計画の一部変更等について」](https://www.soumu.go.jp/main_content/000976977.pdf) — `0204`、`0600`、付加的役務番号を含む区分表。
