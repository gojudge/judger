#include "executer.h"
#include <errno.h>

#define VERSION "1.0.1"

pid_t child;
long begin_time;
char* executable = NULL;
int EXE_LEN = 1024;
int fd = 0;

enum ecode{
  PEN,      // Exit Normally
  PRE,      // Runtime Error
  POM,      // Out of Memory
  POT,      // Out of Time
  POL,      // Output Limit Exceed
  PSF       // Syscall Forbidden
};

/* print error */
void PRTERR(){
  extern int errno;
  char* message;
  
  printf("errno [%d]\n", errno);
  message = strerror(errno);
  printf("Mesg: %s\n", message);
}

/* now, get now time of microsecond */
long t_now(){
  long tmp_now = 0;
  struct timeval tv;

  memset(&tv, 0, sizeof(struct timeval));
  gettimeofday(&tv, NULL);

  tmp_now = (long)tv.tv_sec * 1000 + (long)tv.tv_usec / 1000;

  return tmp_now;
}

/* process exit */
void pexit(enum ecode EC){
  if(EC == PEN){
    // Ignore
  }else if(EC == POT){
    printf("Out of Time.\n");
    kill(child,SIGKILL);
  }else if(EC == PSF){
    printf("Syscall Forbidden.\n");
    kill(child,SIGKILL);
  }else if(EC == POM){
    printf("Out of Memory.\n");
    kill(child,SIGKILL);
  }else if(EC == PRE){
    printf("Runtime Error.\n");
  }else if(EC == POL){
    printf("Output Limit Exceed.\n");
    kill(child,SIGKILL);
  }else{
    dprintf(fd, "EC [%d]\n", EC);
  }

  close(fd);
  exit(0);
}

/* if the syscall is forbidden, return 0 */
int check_syscall(int syscall){
  int i = 0;
  for(i = 0; i < array_len; i++){
    if(syscall==allow_syscall[i]){
      return 1; //true, the syscall matched one of the list, pass
    };
  }
  return 0; //false, not matched
}

/* timer, when over time, killed son and program exit */
void* time_watcher(void* unused){
  while (1){
    long now_time = t_now();
    if(now_time - begin_time - (long)max_time > 0){
      dprintf(fd, "over time [%lu], killed!\n", now_time);
      pexit(POT);
    }
  }
}

/* parse command args */
void parse_args(int argc, char *argv[]){
  int i = 0;
  char* arg = NULL;
  char* buff = NULL;
  const int BUF_LEN = 128;
  char* tag_name = NULL;
  char* tag_value = NULL;

  buff = (char*)malloc(sizeof(BUF_LEN));

  for(i = 1; i < argc; i++){
    memset(buff, 0, sizeof(BUF_LEN));
    strncpy(buff, argv[i], strlen(argv[i])+1);

    if(buff[0] == '-'){               // options
      tag_name = strtok(buff+1, "="); // string time
      tag_value = strtok(NULL, "=");  // decemal time value

      if(!strcmp(tag_name, "t")){     // time
        int tmp_time = atoi(tag_value);
        if(tmp_time > 0){
          max_time = tmp_time;
        }else{
          dprintf(fd, "invalid time [%d], use default.\n", tmp_time);
        }
      }else if(!strcmp(tag_name, "m")){  // memory
        int tmp_mem = atoi(tag_value);
        if(tmp_mem > 0){
          max_mem = tmp_mem;
        }else{
          dprintf(fd, "invalid memory [%d], use default.\n", tmp_mem);
        }
      }

    }else{ // executable path, just one
      int len = strlen(argv[i]);
      memset(executable, 0, sizeof(EXE_LEN));
      strncpy(executable, argv[i], len);
      executable[len]=0;
    }

  }

  return;
}

int main(int argc, char *argv[])
{
  long orig_eax;

  // alloc memory for path string
  executable = (char*)malloc(sizeof(EXE_LEN)); 
  memset(executable, 0, sizeof(EXE_LEN));

  if(argc<2){
    printf(
        "\033[0;39;1mSandbox for Linux Native\033[0m\n"
        "Usage: executer <option> <command>\n"
        "option:\n"
        "  \033[0;33m-t=time\033[0m     program max time\n"
        "  \033[0;33m-m=mem\033[0m      program max memory\n"
        "\033[0;32mversion "
        VERSION
        "\033[0m\n"
        );
    return 0;
  }else{
    parse_args(argc, argv);
  }

  child = fork();
  if(child == 0) {
    ptrace(PTRACE_TRACEME, 0, NULL, NULL);
    // must use execl for supporting segmentfault check
    //printf("exe [%s]\n", executable);
    execl(executable, "", (char*)NULL);
    exit(0);
  }else{
    struct rusage rinfo;
    int runstat, i=0;
    pthread_t thread_id;

    fd = 0;
    fd = open("executer.debug", O_WRONLY|O_CREAT);

    dprintf(fd, "the child pid is %d\n", child);

    begin_time = t_now();
    dprintf(fd, "begin time [%lu]\n", begin_time);
    dprintf(fd, "max_time [%lu]\nmax_mem [%d]\nexecutable path [%s]\n", 
        max_time, max_mem, executable
    );

    //read config
    char* config_string = read_config("executer.json");
    parse_config_json(config_string);
    free_config_buffer(config_string);

    // a new thread for timer, when over time, killed and exit
    pthread_create (&thread_id, NULL, &time_watcher, NULL);

    for(;;){
      //time_t now_time;
      wait4(child,&runstat,0,&rinfo);

      if (WIFEXITED(runstat))
      {
        int exitcode = WEXITSTATUS(runstat);
        dprintf(fd, "exitcode [%d]\n", exitcode);
        if (exitcode != 0)
        {
          //Runtime Error
          dprintf(fd, "Runtime Error\n");
          pexit(PRE);
        }
        //normal exit
        dprintf(fd, "Exit Normally.\n");
        pexit(PEN);
      }
      else if (WIFSIGNALED(runstat))
      {
        // call kill(pid, SIGKILL)
        // Ignore
        exit(0);
      }
      else if (WIFSTOPPED(runstat))
      {
        int signal = WSTOPSIG(runstat);

        if (signal == SIGTRAP){
          struct user_regs_struct reg;
          int syscall;
          static int executed = 0;
           
          ptrace(PTRACE_GETREGS,child,NULL,&reg);
          #ifdef __i386__
          syscall = reg.orig_eax;
          #else
          syscall = reg.orig_rax;
          #endif
          
          dprintf(fd, "syscall: %d\n", syscall);

          // syscall check 
          if(!check_syscall(syscall)){
            dprintf(fd, "Syscall [%d] is Forbidden.\n", syscall);

            pexit(PSF);
          }
          
        }else if(signal == SIGUSR1){
          // Ignore
        }else if(signal == SIGXFSZ){
          dprintf(fd, "Output Limit Exceed.\n");

          pexit(POL);
        }else{
          dprintf(fd, "Runtime Error.\n");

          pexit(PRE);
        }

      }

      ptrace(PTRACE_SYSCALL, child, NULL, NULL);

    }
    
    //close(fd);

  }

  return 0;
}

