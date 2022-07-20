# mercari-gopherdojo-01

# ex00

## 概要

自作catコマンド

## build

```
go build
```

## usage

```
./ft_cat [file path...]
```

# ex01

## 概要

mercari-gopherdojo-00の画像形式変換コマンドと中身は同じ
ただし、テストをgoで実行

## build

```
go mod tidy
go build
```

## usage

```
./convert -i=[input format] -o=[output format] [image file path]
[input format]はpng, jpeg, gifに対応
[output format]はpng, jpeg, gif, pgm, ppmに対応
```
ppm, pgmはPNMの一種です https://ja.wikipedia.org/wiki/PNM_(%E7%94%BB%E5%83%8F%E3%83%95%E3%82%A9%E3%83%BC%E3%83%9E%E3%83%83%E3%83%88)

pgmに変換する際には画像が白黒になります。

## test

imageconvディレクトリの中で、

```
go test
```

-pで並列実行できます。
