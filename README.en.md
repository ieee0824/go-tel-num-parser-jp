# go-tel-num-parser-jp

English | [日本語](README.md)

A Go package for identifying the format and type of Japanese telephone numbers and extracting them from text. Requires Go 1.27.1 or later.

```go
ok, kind := tnp.IsTelNumber("060-1234-5678") // true, tnp.MobilePhone
number, err := tnp.CropTelNumber("Contact: 0800-123-4567") // "0800-123-4567", nil
```

`IsTelNumber` checks the entire input. `CropTelNumber` returns the first number found in a string. Parenthesized numbers such as `03(5321)1111` are treated as `03-5321-1111`. Use `SetIgnoreTypes(tnp.MobilePhone)` to exclude a type. Each call replaces the exclusions; `SetIgnoreTypes()` clears them.

## Supported number types

The types below follow the numbering plan and allocation information available as of September 2026. The examples illustrate formats, and numbers without hyphens are also supported.

| `TelType` | Type | Example formats |
| --- | --- | --- |
| `FixedLinePhone` | Geographic fixed-line | `03-5321-1111`, `011-896-4081` |
| `M2M` | Data-only mobile | `020-123-45678` (11 digits), `0200-12345-67890` (14 digits) |
| `PocketBell` | Paging | `020-412-34567` |
| `IPPhone` | 050 IP phone | `050-1234-5678` |
| `MobilePhone` | Voice mobile | `060-1234-5678`, `070-1234-5678`, `080-1234-5678`, `090-1234-5678` |
| `IncomingCharge` | Toll-free / called-party pays | `0120-123-456` (10 digits), `0800-123-4567` (11 digits) |
| `UnifiedNumber` | Unified number | `0570-123-456` |
| `InformationCharge` | Information-fee collection | `0990-123-456` |
| `FMC` | Fixed-mobile convergence | `0600-123-4567` |

This package checks only the number's format and type. It does not verify whether a number has been allocated to a carrier or is in use, or whether a geographic area code exists. Unallocated `0170` and `0180` ranges, short codes such as emergency numbers, and international numbers are outside its scope. `060` is classified according to the numbering plan, regardless of whether carriers have started offering numbers from that range.

## License

This project is licensed under the [MIT License](LICENSE).

## References

- [Ministry of Internal Affairs and Communications (MIC): Telecommunications Number Allocations](https://www.soumu.go.jp/main_sosiki/joho_tsusin/top/tel_number/number_shitei.html) — number types and allocation status, published as of September 1, 2026 (Japanese).
- [e-Gov Data Portal: Telecommunications Number Allocations](https://data.e-gov.go.jp/data/ja/dataset/soumu_20160325_0034) — `020C` (11 digits), `0200` (14 digits), and service-number types (Japanese).
- [MIC: Addition of 060 Mobile Numbers](https://www.soumu.go.jp/main_sosiki/joho_tsusin/top/tel_number/060keitai.html) — 11-digit `060-1` through `060-9` numbers and their rollout status (Japanese).
- [NTT East: Free Access Hikari Wide](https://business.ntt-east.co.jp/service/hikari_of/free.html) — `0120` and `0800` number formats (Japanese).
- [MIC Numbering System Materials on e-Gov](https://public-comment.e-gov.go.jp/pcm/download?seqNo=0000302542) — 2025 numbering-system diagram (Japanese).
- [MIC: Partial Amendment of the Telecommunications Numbering Plan](https://www.soumu.go.jp/main_content/000976977.pdf) — categories including `0204`, `0600`, and service numbers (Japanese).
