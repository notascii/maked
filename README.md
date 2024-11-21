# projet-SYSD

[https://systemes.pages.ensimag.fr/www-sysd-isi3a/](https://systemes.pages.ensimag.fr/www-sysd-isi3a/)

[https://grid5000.fr](https://grid5000.fr)

## To perf measure for pingpong/pingpong_IO

starting locally

` ssh <login@access.grid5000.fr>`

`ssh grenoble`

### To measure pingpong

`cd pingpong`

### To measure pingpongIO

`cd pingpongIO`

### Command to create nodes

`./create_nodes.sh`

Make sure to change terminal and connect back your account by ssh on grid5000 Grenoble's access machine 

### Command to install dependencies

`./install_dep.sh`

### Command to measure performance

`./measure_perf.sh`

### Results

To download results from server: `./retrieve_results.sh`

To view raw data of latency and throughput `./pingpong/perf/logs` or `./pingpong_IO/perf/logs`

To view plots : `./pingpong/perf/perf_benchmarks/graphs` or `./pingpong_IO/perf/perf_benchmarks/graphs`

To view tables: `./pingpong/perf/perf_benchmarks/tables`or`./pingpong_IO/perf/perf_benchmarks/tables`
