Go言語でつくるインタプリタ/コンパイラ

- 変数拘束
- 四則演算
- 関数呼び出し
- 条件分岐
- 配列
- ハッシュ
- builtin(len, puts)

```
$ go run main.go

> let a = 1;
> a;
// 1
> a + 2;
// 3
> "hello" + "world!";
// helloworld!

> let multiple = fn(a, b){ return a * b; };
> multiple(10, 2)
// 20
> fn(a, b){ return a / b; }(9, 3);
// 3

> if (1 > 2) { return 1; } else { return 2; };
// 2

> let arr = [1, 2];
> arr[0]
// 1

>> let hash = {"key": "value"};
>> hash["key"];
// "value"

> len(arr)
// 2
> puts(arr)
// [1, 2]
// null

>> exit
// bye!
```
