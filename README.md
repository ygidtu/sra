# SRA

SRA data collector

Retrieve information from website like SRA and EBI, based on chromium browsers, only tested on macOS and Linux.


## Build from source

```bash
# flags is optional
flags="-s -w -X main.buildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.gitHash=`git rev-parse HEAD` -X 'main.goVersion=`go version`' -X main.version=v0.0.2"

env GOOS=windows GOARCH=amd64 go build -ldflags "$flags" -x -o sra_win .  # windows exe could be built, by never tested
env GOOS=linux GOARCH=amd64 go build -ldflags "$flags" -x -o sra_linux .
env GOOS=darwin GOARCH=amd64 go build -ldflags "$flags" -x -o sra_osx .
```

## Usage

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
        -u, --url       The official API of ENA (default: https://www.ebi.ac.uk/ena/portal/api/filereport)
        -x, --proxy     The proxy url, eg: http://127.0.0.1:7890
        -o, --output    The output directory (default: ./output)
        -f, --fields    The required fields (default: study_accession,secondary_study_accession,sample_accession,secondary_sample_accession,experiment_accession,run_accession,submission_accession,tax_id,scientific_name,instrument_platform,instrument_model,library_name,nominal_length,library_layout,library_strategy,library_source,library_selection,read_count,base_count,center_name,first_public,last_updated,experiment_title,study_title,study_alias,experiment_alias,run_alias,fastq_bytes,fastq_md5,fastq_ftp,fastq_aspera,fastq_galaxy,submitted_bytes,submitted_md5,submitted_ftp,submitted_aspera,submitted_galaxy,submitted_format,sra_bytes,sra_md5,sra_ftp,sra_aspera,sra_galaxy,sample_alias,broker_name,sample_title,nominal_sdev,first_created)
        -t, --threads   The number of threads to use, do not use too much threads (default: 1)
        -c, --resume    Skip finished requests
    enap:
        -i, --input     KeyID or list of KeyID
        -u, --url       The official api of ENA (default: https://www.ebi.ac.uk/ena/browser/api/summary/)
        -x, --proxy     The proxy url, eg: http://127.0.0.1:7890
        -o, --output    The output directory (default: ./output.csv)
        -t, --threads   The number of threads to use, do not use too much threads (default: 1)
        -c, --resume    Skip finished requests
    search:
        -i, --input     The list of RBP in csv format with header
        -u, --url       The official api of SRA (default: https://www.ncbi.nlm.nih.gov/sra/)
        -x, --proxy     The proxy url, eg: http://127.0.0.1:7890
        -o, --output    The output directory (default: ./output)
        -p, --param     The extra parameters (default: "Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy]))
        -t, --timeout   Connection timeout in seconds (default: 10s)
        -e, --exec      path to chrome executable
            --open      whether open the GUI of chrome
    study:
        -s, --study     Study ID to query
        -p, --proxy     Proxy
        -o, --output    Output json (default: ./output.csv)
            --open      whether open the GUI of chrome
        -t, --threads   How many threads to use (default: 1)
        -e, --exec      path to chrome executable
```


