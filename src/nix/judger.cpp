#include "judger.h"

int main(int argc, char *argv[])
{
  struct user_regs_struct reg;
  int syscall;

  pid_t child;
  long orig_eax;
  child = fork();

  ptrace(PTRACE_GETREGS,child,NULL,&reg);
  #ifdef __i386__
  syscall = reg.orig_eax;
  #else
  syscall = reg.orig_rax;
  #endif

  printf("syscall is %d\n", syscall);

  if(child == 0) {
     ptrace(PTRACE_TRACEME, 0, NULL, NULL);
     execl("/bin/ls", "ls", NULL);
  }else{
     wait(NULL);
     printf("the child pid is %d\n", child);
     
     orig_eax = ptrace(PTRACE_PEEKUSER, 
                       child, 4 * ORIG_EAX, 
                       NULL);

     printf("The child made a system call %ld\n", orig_eax);
     ptrace(PTRACE_CONT, child, NULL, NULL);
  }
  return 0;
}
