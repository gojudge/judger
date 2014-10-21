#include "judger.h"

int main(int argc, char *argv[])
{
  pid_t child;
  long orig_eax;

  if(argc<2){
    printf("Usage: judger <command> \n");
    return 0;
  }

  child = fork();
  if(child == 0) {
    ptrace(PTRACE_TRACEME, 0, NULL, NULL);
    execvp(argv[1], argv+1);
    exit(0);
  }else{
    struct rusage rinfo;
    int runstat;

    printf("the child pid is %d\n", child);

    for(;;){
      wait4(child,&runstat,0,&rinfo);

      if (WIFEXITED(runstat))
      {
        exit(0);
      }
      else if (WIFSIGNALED(runstat))
      {
        //result(RS_ECR,SIGKILL);
      }
      else if (WIFSTOPPED(runstat))
      {
        if (WSTOPSIG(runstat) == SIGTRAP){
          struct user_regs_struct reg;
          int syscall;
          static int executed = 0;
          int forbidden = 0;
                
          ptrace(PTRACE_GETREGS,child,NULL,&reg);
          #ifdef __i386__
          syscall = reg.orig_eax;
          #else
          syscall = reg.orig_rax;
          #endif
                      
          printf("syscall: %d\n",syscall);
        }
      }

//      orig_eax = ptrace(PTRACE_PEEKUSER, child, 4 * ORIG_EAX, NULL);
//      printf("The child made a system call %ld\n", orig_eax);
      ptrace(PTRACE_SYSCALL, child, NULL, NULL);
    }
  }
  return 0;
}

