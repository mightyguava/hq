# hq
HTML query CLI using CSS selectors, built on top of [Cascadia](https://github.com/andybalholm/cascadia).


## Installation

```
go get github.com/mightyguava/hq
```

## Usage

```
hq 'div.class_name' test.html
cat test.html | hq 'div.class_name'
```

## Special pseudo-selectors

Use the `::text` pseudo-selector to recursively print the text nodes of matching nodes instead of rendering their HTML. Its behavior is the same as setting `--text`.

Use the `::attr(<attr>)` pseuo-selector to print the value of a matching attribute in the matching nodes instead of than rendering the their html. For example, `::attr(href)` prints the `href` value of the matching nodes. Each match is printed on a separate line. Non-matching nodes are skipped.
