#### 1. Install Requirements

```bash
  $ pip install -r requirements.txt
```

#### 2. Start a go tool trace web server

``` bash
  $ go tool trace trace.out
  2018/06/13 21:07:40 Parsing trace...
  2018/06/13 21:07:42 Serializing trace...
  2018/06/13 21:07:43 Splitting trace...
  2018/06/13 21:07:46 Opening browser
```

#### 3. Run the script
For example, you url is http://127.0.0.1:36627/.
```bash
  $ python gotrace_parser.py -u http://127.0.0.1:36627/
  go goroutine execution time percentage report:
  This go trace log have totally 100 goroutines, max execution time is 47848111956 ns.
  0% --------- 10%, 2
  10% --------- 20%, 3
  20% --------- 30%, 5
  30% --------- 40%, 7
  40% --------- 50%, 9
  50% --------- 60%, 1
  60% --------- 70%, 3
  70% --------- 80%, 10
  80% --------- 90%, 10
  90% --------- 100%, 50
```
