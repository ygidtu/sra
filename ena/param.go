package ena

import (
	"encoding/json"
)

type Params struct {
	Key     string `goptions:"-i, --input, description='KeyID or list of KeyID'"` // obligatory,
	ENA     string `goptions:"-u, --url, description='ENA的官方API链接'"`
	Proxy   string `goptions:"-x, --proxy, description='代理链接地址，比如：http://127.0.0.1:7890'"`
	Output  string `goptions:"-o, --output, description='输出文件夹'"`
	Fields  string `goptions:"-f, --fields, description='输出的文件信息'"`
	Threads int    `goptions:"-t, --threads, description='所使用的线程数'"`
	Resume  bool   `goptions:"-c, --resume, description='Skip finished requests'"`
}

func (param *Params) String() string {
	str, _ := json.MarshalIndent(param, "", "    ")
	return string(str)
}

func DefaultParam() Params {
	return Params{
		ENA:     "https://www.ebi.ac.uk/ena/portal/api/filereport",
		Output:  "./output",
		Threads: 1,
		Fields:  "study_accession,secondary_study_accession,sample_accession,secondary_sample_accession,experiment_accession,run_accession,submission_accession,tax_id,scientific_name,instrument_platform,instrument_model,library_name,nominal_length,library_layout,library_strategy,library_source,library_selection,read_count,base_count,center_name,first_public,last_updated,experiment_title,study_title,study_alias,experiment_alias,run_alias,fastq_bytes,fastq_md5,fastq_ftp,fastq_aspera,fastq_galaxy,submitted_bytes,submitted_md5,submitted_ftp,submitted_aspera,submitted_galaxy,submitted_format,sra_bytes,sra_md5,sra_ftp,sra_aspera,sra_galaxy,sample_alias,broker_name,sample_title,nominal_sdev,first_created",
	}
}
