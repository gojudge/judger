#include "executer.h"

#define VERSION "1.0.1"

pid_t child;
long begin_time;
char* executable = NULL;
int fd = 0;

/* print error */
//void PRTERR(){
//  extern int errno;
//  char* message;
//  
//  printf("errno [%d]\n", errno);
//  message = strerror(errno);
//  printf("Mesg: %s\n", message);
//}

/* now, get now time of microsecond */
long t_now(){
  long tmp_now = 0;
  struct timeval tv;

  memset(&tv, 0, sizeof(struct timeval));
  gettimeofday(&tv, NULL);

  tmp_now = (long)tv.tv_sec * 1000 + (long)tv.tv_usec / 1000;

  return tmp_now;
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
      printf("over time [%lu], killed!\n", now_time);
      kill(child,SIGKILL);
      exit(-1);
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
      strncpy(executable, argv[i], strlen(argv[i]));
    }

  }

  return;
}

int main(int argc, char *argv[])
{
  long orig_eax;
  int EXE_LEN = 1024;

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
    execl(executable, executable, (char*)NULL);
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
          printf("Runtime Error\n");
          exit(-exitcode);
        }
        //normal exit
        dprintf(fd, "Exit Normally.\n");
        exit(0);
      }
      else if (WIFSIGNALED(runstat))
      {
        // call kill(pid, SIGKILL)
        // Ignore
        exit(-1);
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
            dprintf(fd, "[WIFSTOPPED>SIGTRAP] Syscall [%d] is Forbidden.\n", syscall);
            printf("Syscall [%d] is Forbidden.\n", syscall);

            kill(child,SIGKILL);
            return -1;
          }
          
        }else if(signal == SIGUSR1){
          // Ignore
        }else if(signal == SIGXFSZ){
          dprintf(fd, "[WIFSTOPPED>SIGXFSZ] Output Limit Exceed.\n");
          printf("Output Limit Exceed.\n");

          exit(-1);
        }else{
          dprintf(fd, "[WIFSTOPPED>SIGXFSZ] Runtime Error.\n");
          printf("Runtime Error.\n");

          exit(-1);
        }

      }

      ptrace(PTRACE_SYSCALL, child, NULL, NULL);

    }
    
    close(fd);

  }

  return 0;
}

