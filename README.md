# SRA

SRA data collector

```bash
https://www.ncbi.nlm.nih.gov/sra/?term=(%22knockdown%22%5BTitle%5DOR+%22knockdown%22%5BDescription%5D)+AND+(%22HNRNPA1%22%5BTitle%5D+OR+%22HNRNPA1%22%5BDescription%5D)+AND+%22Homo+sapiens%22%5Borgn%3A__txid9606%5D+AND(rna+seq%5BStrategy%5D)
("knockdown"[Title]OR "knockdown"[Description]) AND ("HNRNPA1"[Title] OR "HNRNPA1"[Description]) AND "Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy]) 
```

基于chrome浏览器来爬取相关数据，需保证安装有chrome，暂时无法调整默认下载目录，所以暂时写死以后只能用于macOS或者Linux运行

```bash
❯ ./sra -h
Usage: sra [global options] <verb> [verb options]

Global options:
        -v, --version   Show version
            --debug     Show debug info
            --help      Show this help

Verbs:
    detail:
        -s, --study     Study ID to query
        -a, --accession Accession ID to query
        -r, --run       Run ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./output.json)
        -t, --threads   How many threads to use (default: 1)
        -h, --help      Show this help
    ebi:
        -s, --study     Study ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./ebi_output.csv)
    search:
        -i, --input     RBP的list，csv格式，带列名 (*)
        -u, --url       SRA的官方链接 (default: https://www.ncbi.nlm.nih.gov/sra/)
        -x, --proxy     代理链接地址，比如：http://127.0.0.1:7890
        -o, --output    输出文件夹 (default: ./output)
        -p, --param     额外的查询参数 (default: "Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy]))
        -t, --timeout   Connection timeout in seconds (default: 10s)
            --open      是否打开chrome的图形化界面
            --debug     显示debug log
        -h, --help      Show this help
    study:
        -s, --study     Study ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./output.csv)
            --open      是否打开chrome的图形化界面
        -t, --threads   How many threads to use (default: 1)
```
