# projet-SYSD

[https://systemes.pages.ensimag.fr/www-sysd-isi3a/](https://systemes.pages.ensimag.fr/www-sysd-isi3a/)

[https://grid5000.fr](https://grid5000.fr)

## To perf measure for pingpong/pingpong_IO

starting locally

`bash ssh <login@access.grid5000.fr>`

`bash grenoble`

### To measure pingpong

`bash cd pingpong`

### To measure pingpongIO

`bash cd pingpongIO`

### Command to create nodes

`bash ./create_nodes.sh`

### Command to install dependencies

`bash ./install_dep.sh`

### Command to measure performance

`bash ./measure_perf.sh`

### Results

To download results from server: ./retrieve_results.sh

To view raw data of latency and throughput `./pingpong/perf/logs' or `./pingpong_IO/perf/logs`

To view plots : `./pingpong/perf/perf_benchmarks/graphs` or `./pingpong_IO/perf/perf_benchmarks/graphs`

To view tables: './pingpong/perf/perf_benchmarks/tables`or`./pingpong_IO/perf/perf_benchmarks/tables`
