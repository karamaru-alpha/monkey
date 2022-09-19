Go言語でつくるインタプリタ

- 四則演算
- 変数拘束
- 関数呼び出し

```
$ go run main.go

> let a = 1;
> a;
// 1

> let b = 2;
> a + b;
// 3

> let multiple = fn(a, b){ return a * b; };
> multiple(10, 2)
// 20

> fn(a, b){ return a / b; }(9, 3);
// 3
```
