#include "executer.h"

/* if the syscall is forbidden, return 0 */
int check_syscall(int syscall){
  int i = 0;
  for(i = 0; i < array_len; i++){
    if(syscall==forbidden_syscall[i]){
      return 0; //false, forbidden
    };
  }
  return 1; //true, pass
}

int main(int argc, char *argv[])
{
  pid_t child;
  long orig_eax;

  if(argc<2){
    printf("Usage: judger <command> <args>\n");
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

    int fd = 0 ,i = 0;
    fd = open("executer.debug", O_WRONLY|O_CREAT);

    dprintf(fd, "the child pid is %d\n", child);

    //read config
    char* config_string = read_config("executer.json");
    parse_config_json(config_string);
    free_config_buffer(config_string);

//    printf("time: %d\n", time);
//    printf("memory: %d\n", mem);
//    for(i = 0; i<array_len; i++){
//      printf("forbidden: %d\n", forbidden_syscall[i]);
//    }

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
          
          dprintf(fd, "syscall: %d\n", syscall);

          if(!check_syscall(syscall)){
            printf("syscall [%d] is forbidden\n", syscall);
            kill(child,SIGKILL);
            return -1;
          }
        }
      }

      ptrace(PTRACE_SYSCALL, child, NULL, NULL);

      close(fd);
    }
  }

  return 0;
}

