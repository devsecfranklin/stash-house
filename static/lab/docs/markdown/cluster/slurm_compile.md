# How to compile with C/C++

Follow these steps to compile a program with a C/C++ compiler.

1. Log on to the Grid.
2. Start an interactive job. Type: `srun -q debug -t 10:0 --pty bash`. This puts you on a compute
node for you to do compiling on. Never run jobs or do compiling on the head node Warrior.
3. Type: `nano hello.c`
4. Type the following lines into the editor:

    Image

5. Choose a C/C++ compiler. In this example we will use gcc. Load the module.

    Type: `module load gnu7/7.3.0`

6. To see the pathname of the gcc compiler you will be using type: `which gcc`
This shows you the first instance of it in your path.
7. Type: `gcc hello.c -o hello_c_example`
8. Now that you're done compiling, end the interactive job by pressing Ctrl + D.
9. Create the job script. Type: `nano c_job_script`
10. Enter the following into the editor:

    ```sh
    #!/bin/bash
    # Job name
    #SBATCH --job-name C Job
    # Submit to the primary QoS
    #SBATCH -q primary
    # Request one node
    #SBATCH -N 1

    # Total number of cores, in this example it will 1 node with 1 core each.
    #SBATCH -n 1

    # Request memory
    #SBATCH --mem=5G

    # Mail when the job begins, ends, fails, requeues
    #SBATCH --mail-type=ALL

    # Where to send email alerts
    #SBATCH --mail-user=xxyyyy@wayne.edu

    # Create an output file that will be output_<jobid>.out
    #SBATCH -o output_%j.out

    # Create an error file that will be error_<jobid>.out
    #SBATCH -e errors_%j.err

    # Set maximum time limit
    #SBATCH -t 1-:0:0

    cd $TMPDIR

    cp /wsu/home/zz/zz99/zz9991/hello_c_example $TMPDIR

    $TMPDIR/hello_c_example
    ```

    Save the file and exit the editor. Press Ctrl + X, then Y for Yes and keep the name c_job_script.

    NOTE: `zz9992` is the account used in this example, but you will have to modify the script using the AccessID and directories specific to your account.

11. Make sure that the file is executable. You can do this by changing the permissions.

    Type: `chmod +x c_job_script`

12. Submit the job. Type: `sbatch c_job_script`
13. You can check to see what jobs you're running by using the command qme. Once the job
is finished running you will have an error and output file.
Check your home directory by typing: `ls`
