#!/bin/bash 

#OAR -q production 
#OAR -l host=1/gpu=1
#OAR -l walltime=3:00:00
#OAR -p gpu-16GB AND gpu_compute_capability_major>=5
#OAR -O OAR_%jobid%.out
#OAR -E OAR_%jobid%.err 

# display some information about attributed resources
hostname 
nvidia-smi 
 
# make use of a python torch environment
module load conda
conda activate pytorch_env
python3 -c "import torch; print(torch.cuda.is_available()); print(torch.cuda.get_device_name(0))";