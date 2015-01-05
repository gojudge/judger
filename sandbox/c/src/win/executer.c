/**
 * File Name: executer.c
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2015年1月2日 星期五 20时54分12秒
 */

#include "executer.h"

#define VERSION "1.0.0"

size_t PATH_LEN = 1024;
long max_time;           // max time
int max_mem;             // max memory
int judger_model = 2;    // assert as default
char* input = NULL;      // input file path
char* output = NULL;     // output file path
char* executable = NULL; // executable path
FILE* fd = NULL;         // debug file
DWORD pid;               // process id

/** kill process by pid */
BOOL KillProcess(DWORD ProcessId){
    HANDLE hProcess=OpenProcess(PROCESS_TERMINATE,FALSE,ProcessId);
    if(hProcess==NULL)
        return FALSE;
    if(!TerminateProcess(hProcess,0))
        return FALSE;
    return TRUE;
}

/** Process Exit */
void ProcessExit(const char* exit_mark){
    FILE* run_result = NULL;

    if (!KillProcess(pid)){
        return;
    }

    run_result = fopen("RUNRESULT", "w");
    fprintf(run_result, "%s", exit_mark);
    fclose(run_result);

    if (1 != judger_model){
        printf("[%s]", exit_mark);
    }
    
    dprintf(fd, "Process Exited! [%s]\n", exit_mark);
    exit(0);
}

/** Check Memory */
void CheckMemory(HANDLE hProcess){
    PROCESS_MEMORY_COUNTERS pmc;
    int mem = 0;

    GetProcessMemoryInfo(hProcess, &pmc, sizeof(pmc));  
    mem = pmc.PagefileUsage/1024;
     
    if (max_mem < mem){
        ProcessExit("POM");
    }
}

/** get current time */
int CurrentTime(){
    SYSTEMTIME t;
    int millisec = 0;

    GetLocalTime(&t);
    
    millisec = (t.wHour * 3600 + t.wMinute * 60 + t.wSecond) * 1000 + t.wMilliseconds;

    return millisec;
}

/** Parse Command Args */
void parse_args(int argc, char *argv[]){
    int i = 0;
    int len = 0;
    char *arg = NULL;
    char *buff = NULL;
    const int BUF_LEN = 128;
    char *tag_name = NULL;
    char *tag_value = NULL;

    buff = (char *)malloc(sizeof(char) * BUF_LEN);

    for (i = 1; i < argc; i++) {
        memset(buff, 0, sizeof(char) * BUF_LEN);
        strncpy(buff, argv[i], strlen(argv[i]) + 1);

        if (buff[0] == '-') {   // options
            tag_name = strtok(buff + 1, "=");   // string time
            tag_value = strtok(NULL, "=");  // decemal time value

            if (!strcmp(tag_name, "t")) {   // time
                int tmp_time = atoi(tag_value);
                if (tmp_time > 0) {
                    max_time = tmp_time;
                    dprintf(fd,"[max time] %d\n", max_time);
                } else {
                    dprintf(fd, "invalid time [%d], use default.\n", tmp_time);
                }
            } else if (!strcmp(tag_name, "m")) {    // memory
                int tmp_mem = atoi(tag_value);
                if (tmp_mem > 0) {
                    max_mem = tmp_mem;
                    dprintf(fd,"[max memory] %d\n", max_mem);
                } else {
                    dprintf(fd, "invalid memory [%d], use default.\n", tmp_mem);
                }
            } else if (!strcmp(tag_name, "j")) {
                dprintf(fd, "[judger model] %s\n", tag_value);

                if (!strcmp(tag_value, "io")) {
                    judger_model = 1;
                } else {
                    judger_model = 2;
                }
            } else if (!strcmp(tag_name, "-stdin")) {
                dprintf(fd, "[input] %s\n", tag_value);

                len = strlen(tag_value);
                memset(input, 0, PATH_LEN);
                strncpy(input, tag_value, len);
                input[len] = 0;
            } else if (!strcmp(tag_name, "-stdout")) {
                dprintf(fd, "[output] %s\n", tag_value);

                len = strlen(tag_value);
                memset(output, 0, PATH_LEN);
                strncpy(output, tag_value, len);
                input[len] = 0;
            }

        } else {                // executable path, just one
            int len = strlen(argv[i]);

            // empty string will cause err, filted
            if (len >= 2){
                memset(executable, 0, PATH_LEN);
                strncpy(executable, argv[i], len);
                executable[len] = 0;

                dprintf(fd, "[executable] %s\n", executable);
            }
            
        }

    }
}

int main(int argc, char ** argv){
    STARTUPINFO si;
    PROCESS_INFORMATION pi;
    DEBUG_EVENT de;
    BOOL stop = FALSE;
    int StartTime = 0;

    ZeroMemory(&de, sizeof(de));
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    ZeroMemory(&pi, sizeof(pi));

    // alloc memory for input/output, executable path
    input = (char*)malloc(PATH_LEN);
    output = (char*)malloc(PATH_LEN);
    executable = (char*)malloc(PATH_LEN);

    ZeroMemory(input,PATH_LEN);
    ZeroMemory(output,PATH_LEN);
    ZeroMemory(executable,PATH_LEN);
  
    fd = dopen("executer.debug");

    // show help
    if (argc<2) {
        printf("Unsafely Executer for Windows\n");
        printf("Usage: %s [arguments ...] <app_name>\n", argv[0]);
        printf("Options:\n"
               "  -m=mem          max memory\n"
               "  -t=time         max time\n"
               "  -j=mode         mode[io/assert]\n"
               "  --stdin=path    stdin file path\n"
               "  --stdout=path   stdout file path\n"
               "Version " VERSION "\n"
            );
        return 0;
    }else{
        parse_args(argc, argv);
    }

    if(!CreateProcess(NULL, executable, NULL, NULL, FALSE,
        CREATE_DEFAULT_ERROR_MODE, NULL, NULL, &si, &pi)){
            dprintf(fd, "CreateProcess [%s] failed (%d).\n", executable, GetLastError());
            printf("CreateProcess failed (%d).\n", GetLastError());
            exit(-1);
    }else{
        // get process id
        pid = pi.dwProcessId;
        
        StartTime = CurrentTime();

        dprintf(fd, "Process [%s] Created.\n", executable);
        dprintf(fd, "Start Time [%d]\n", StartTime);

    }
  
    while (TRUE) {  
        int ct = 0;

        WaitForDebugEvent (&de, INFINITE); 

        if (max_mem>0){
            CheckMemory(pi.hProcess);
        }
        

        if (max_time > 0){
            // time check
            ct = CurrentTime();
            if (ct - StartTime > max_time){
                printf("Over Time [%d]\n", ct);
                ProcessExit("POT");
            }
        }
        
        if (de.dwDebugEventCode>0){
            dprintf(fd, "Trace DebugEventCode [%x]\n", (de.dwDebugEventCode));
        }
        
        switch (de.dwDebugEventCode) {  
            case EXCEPTION_DEBUG_EVENT:         /* exception */  
                switch (de.u.Exception.ExceptionRecord.ExceptionCode) {   
                    case   EXCEPTION_INT_DIVIDE_BY_ZERO:    /* #DE */  
                        // Do what the parent process want to do when the child process gets #DE interrupt.  
                        ProcessExit("PRE");
                        break;
                    case   EXCEPTION_BREAKPOINT:            /* #BP */  
                        // Do what the parent process want to do when the child process gets #BP interrupt.  
                        ProcessExit("PBP");
                        break;
          
                    default:   
                        // printf("Unknown Exception\n"); 
                        ProcessExit("PRE");
                        break;
                }      
      
                ContinueDebugEvent(de.dwProcessId,de.dwThreadId,DBG_EXCEPTION_HANDLED);
                continue;
      
            case CREATE_PROCESS_DEBUG_EVENT:
                break;

            case CREATE_THREAD_DEBUG_EVENT:
                dprintf(fd, "[CREATE_THREAD_DEBUG_EVENT]\n");
                continue;
            case EXIT_PROCESS_DEBUG_EVENT:
                dprintf(fd, "[EXIT_PROCESS_DEBUG_EVENT]\n");
                stop = TRUE;
                break;
            case EXIT_THREAD_DEBUG_EVENT:
                dprintf(fd, "[EXIT_THREAD_DEBUG_EVENT]\n");
                continue;
            case LOAD_DLL_DEBUG_EVENT:
                dprintf(fd, "[LOAD_DLL_DEBUG_EVENT]\n");
                continue;
            case OUTPUT_DEBUG_STRING_EVENT:
                dprintf(fd, "[OUTPUT_DEBUG_STRING_EVENT]\n");
                continue;
            case RIP_EVENT:
                dprintf(fd, "[RIP_EVENT]\n");
                continue;
            case UNLOAD_DLL_DEBUG_EVENT:
                dprintf(fd, "[UNLOAD_DLL_DEBUG_EVENT]\n");
                continue;
      
            default:  
                // printf("Unknown Event!\n");
                // ProcessExit("PEN");
                break;
        }  
  
        if (TRUE == stop) {  
            //printf("Process exit\n");
            ProcessExit("PEN");
            break;
        }  
  
        ContinueDebugEvent(de.dwProcessId, de.dwThreadId, DBG_CONTINUE);
  
    } // end of loop  
  
    assert(stop);
  
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
  
    dclose(fd);
    
    return 0;  
}  