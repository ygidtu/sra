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
        -h, --help      Show this help

Verbs:
    detail:
        -s, --study     Study ID to query
        -a, --accession Accession ID to query
        -r, --run       Run ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./output.json)
        -t, --threads   How many threads to use (default: 1)
    ebi:
        -s, --study     Study ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./ebi_output.csv)
    ena:
        -i, --input     KeyID or list of KeyID
        -u, --url       ENA的官方API链接 (default: https://www.ebi.ac.uk/ena/portal/api/filereport)
        -x, --proxy     代理链接地址，比如：http://127.0.0.1:7890
        -o, --output    输出文件夹 (default: ./output)
        -f, --fields    输出的文件信息 (default: study_accession,secondary_study_accession,sample_accession,secondary_sample_accession,experiment_accession,run_accession,submission_accession,tax_id,scientific_name,instrument_platforument_model,library_name,nominal_length,library_layout,library_strategy,library_source,library_selection,read_count,base_count,center_name,first_public,last_updated,experiment_title,study_title,study_alias,experiment_alias,run_alias,fastq_bytes,fastq_md5,fastq_ftp,fastq_aspera,fastq_galaxy,submitted_bytes,submitted_md5,submitted_ftp,submitted_aspera,submitted_galaxy,submitted_format,sra_bytes,sra_md5,sra_ftp,sra_aspera,sra_galaxy,cram_index_ftp,cram_index_aspera,cram_index_galaxy,sample_alias,broker_name,sample_title,nominal_sdev,first_created)
        -t, --threads   所使用的线程数 (default: 1)
        -c, --resume    Skip finished requests
    enap:
        -i, --input     KeyID or list of KeyID
        -u, --url       ENA的官方API链接 (default: https://www.ebi.ac.uk/ena/browser/api/summary/)
        -x, --proxy     代理链接地址，比如：http://127.0.0.1:7890
        -o, --output    输出文件夹 (default: ./output.csv)
        -t, --threads   所使用的线程数 (default: 1)
        -c, --resume    Skip finished requests
    search:
        -i, --input     RBP的list，csv格式，带列名
        -u, --url       SRA的官方链接 (default: https://www.ncbi.nlm.nih.gov/sra/)
        -x, --proxy     代理链接地址，比如：http://127.0.0.1:7890
        -o, --output    输出文件夹 (default: ./output)
        -p, --param     额外的查询参数 (default: "Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy]))
        -t, --timeout   Connection timeout in seconds (default: 10s)
        -e, --exec      path to chrome executable
            --open      是否打开chrome的图形化界面
    study:
        -s, --study     Study ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./output.csv)
            --open      是否打开chrome的图形化界面
        -t, --threads   How many threads to use (default: 1)
        -e, --exec      path to chrome executable
```


Build from source

```bash
flags="-s -w -X main.buildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.gitHash=`git rev-parse HEAD` -X 'main.goVersion=`go version`' -X main.version=v0.0.2"

env GOOS=windows GOARCH=amd64 go build -ldflags "$flags" -x -o sra_win . # && upx -9 sra_win
env GOOS=linux GOARCH=amd64 go build -ldflags "$flags" -x -o sra_linux . && upx -9 sra_linux
env GOOS=darwin GOARCH=amd64 go build -ldflags "$flags" -x -o sra_osx . # && upx -9 sra_osx
```