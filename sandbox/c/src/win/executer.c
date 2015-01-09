/**
 * File Name: executer.c
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2015年1月2日 星期五 20时54分12秒
 */

#include "executer.h"

#define VERSION "1.1.0"

size_t PATH_LEN = 1024;
long max_time;           // max time
int max_mem;             // max memory
int judger_model = 2;    // assert as default
char* input = NULL;      // input file path
char* output = NULL;     // output file path
char* executable = NULL; // executable path
FILE* fd = NULL;         // debug file
DWORD pid;               // process id
int StartTime = 0;       // process start time
BOOL pot = FALSE;        // process out of time
BOOL pom = FALSE;        // process out of memory

HANDLE hChildStdoutRd, hChildStdoutWr;
HANDLE hChildStdinRd, hChildStdinWr;
FILE* file_out = NULL;
FILE* file_in = NULL;

#define BUFSIZE 2048

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

    if (1 == judger_model){
        fclose(file_out);
    }

    if (!KillProcess(pid)){
        dprintf(fd, "Kill Process Failed!\n");
        return;
    }

    run_result = fopen("RUNRESULT", "w");
    fprintf(run_result, "%s", exit_mark);
    fclose(run_result);
    
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
        pom = TRUE;
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

void ThreadProc(void* arg){
    int ct = 0;

    while(TRUE)     
    {
        if (1 == judger_model){
            DWORD dwRead;

            CHAR chBufOut[BUFSIZE+1];

            ZeroMemory(chBufOut, BUFSIZE+1);


            ReadFile(hChildStdoutRd, chBufOut, BUFSIZE, &dwRead, NULL);
            chBufOut[BUFSIZE]=0;
            fprintf(file_out, "%s", chBufOut);

            
        }

        if (max_time > 0){
            // time check
            ct = CurrentTime();
            if (ct - StartTime > max_time){
                dprintf(fd, "Over Time [%d]\n", ct);
                pot = TRUE;
                ProcessExit("POT");
            }
        }
    }  
}

// Create STDIO pipe
BOOL redirect_stdio(HANDLE* pChildStdoutRd, HANDLE* pChildStdoutWr){
    SECURITY_ATTRIBUTES secattr; 
    ZeroMemory(&secattr, sizeof(secattr));
    secattr.nLength = sizeof(secattr);
    secattr.lpSecurityDescriptor = NULL;
    secattr.bInheritHandle = TRUE;

    if(!CreatePipe(pChildStdoutRd, pChildStdoutWr, &secattr, 0)){
        return FALSE;
    }

    return TRUE;
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
                output[len] = 0;
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
    
    ZeroMemory(&de, sizeof(de));
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

    if(!redirect_stdio(&hChildStdoutRd, &hChildStdoutWr)){
        dprintf(fd, "Set Output Redirect Pipe Failed");
    }

    if(!redirect_stdio(&hChildStdinWr, &hChildStdinRd)){
        dprintf(fd, "Set Input Redirect Pipe Failed");
    }

    // IO Redirection
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    
    if (1 == judger_model){
        DWORD dwWritten;
        CHAR chBufIn[BUFSIZE+1];
        BOOL rst = FALSE;
        int in_len = 0;

        // set default input file name
        if (input[0]){
            char* tmp_input = "input";
            int len_input = strlen(tmp_input);
            strncpy(input,tmp_input,len_input);
            input[len_input]=0;
        }

        // redirection input
        ZeroMemory(chBufIn, BUFSIZE+1);
        file_in = fopen(input,"r");
        free(input);
        fread(chBufIn,BUFSIZE,1,file_in);
        in_len = strlen(chBufIn);
        chBufIn[in_len]='\n';
        rst = WriteFile(hChildStdinRd, chBufIn, in_len+1, &dwWritten, NULL);
        
        if (rst){
            dprintf(fd, "[%d] char written into input file\n", dwWritten);
        }

        // set redirection for child process
        si.dwFlags = STARTF_USESTDHANDLES | STARTF_USESHOWWINDOW;
        si.hStdInput = hChildStdinWr;
        si.hStdOutput = hChildStdoutWr;
        si.hStdError = hChildStdoutWr;

        // set default output file name
        if (output[0]==0){
            char* tmp_output = "output";
            int len_output = strlen(tmp_output);
            strncpy(output,tmp_output,len_output);
            output[len_output]=0;
        }
        dprintf(fd, "output file [%s]\n", output);
        file_out = fopen(output,"w");
        free(output);

    }

    if(!CreateProcess(NULL, 
                    executable, NULL, NULL, TRUE, 
                    DEBUG_ONLY_THIS_PROCESS, NULL, NULL, &si, &pi)){
        CloseHandle(hChildStdoutRd);
        CloseHandle(hChildStdoutWr);
        dprintf(fd, "CreateProcess [%s] failed (%d).\n", executable, GetLastError());
        printf("CreateProcess failed (%d).\n", GetLastError());
        exit(-1);
    }else{
        LPDWORD tid;

        ZeroMemory(&tid, sizeof(tid));

        // get process id
        pid = pi.dwProcessId;
        
        StartTime = CurrentTime();

        dprintf(fd, "Process [%s] Created.\n", executable);
        dprintf(fd, "Start Time [%d]\n", StartTime);

        if (NULL==CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)&ThreadProc, (LPVOID)NULL, 0, tid)){
            dprintf(fd, "Created Timer Thread failed\n");
        }

        CloseHandle(hChildStdoutWr);
        CloseHandle(hChildStdinWr);
    }

    while (TRUE) {
        
        WaitForDebugEvent (&de, INFINITE); 

        if (max_mem>0){
            CheckMemory(pi.hProcess);
        }
        
        dprintf(fd, "Trace DebugEventCode [%x]\n", (de.dwDebugEventCode));
        
        switch (de.dwDebugEventCode){
            case EXCEPTION_DEBUG_EVENT:
                switch (de.u.Exception.ExceptionRecord.ExceptionCode) {
                    case   EXCEPTION_ACCESS_VIOLATION:
                        dprintf(fd,"EXCEPTION_ACCESS_VIOLATION\n");
                        ProcessExit("PRE");
                        break;
                    case   EXCEPTION_ARRAY_BOUNDS_EXCEEDED:
                        dprintf(fd,"EXCEPTION_ARRAY_BOUNDS_EXCEEDED\n");
                        break;
                    case   EXCEPTION_BREAKPOINT:
                        dprintf(fd,"EXCEPTION_BREAKPOINT\n");
                        // Ignore
                        break;
                    case   EXCEPTION_DATATYPE_MISALIGNMENT:
                        dprintf(fd,"EXCEPTION_DATATYPE_MISALIGNMENT\n");
                        break;
                    case   EXCEPTION_FLT_DENORMAL_OPERAND:
                        dprintf(fd,"EXCEPTION_FLT_DENORMAL_OPERAND\n");
                        break;
                    case   EXCEPTION_FLT_DIVIDE_BY_ZERO:
                        dprintf(fd,"EXCEPTION_FLT_DIVIDE_BY_ZERO\n");
                        ProcessExit("PRE");
                        break;
                    case   EXCEPTION_FLT_INEXACT_RESULT:
                        dprintf(fd,"EXCEPTION_FLT_INEXACT_RESULT\n");
                        break;
                    case   EXCEPTION_FLT_INVALID_OPERATION:
                        dprintf(fd,"EXCEPTION_FLT_INVALID_OPERATION\n");
                        break;
                    case   EXCEPTION_FLT_OVERFLOW:
                        dprintf(fd,"EXCEPTION_FLT_OVERFLOW\n");
                        break;
                    case   EXCEPTION_FLT_STACK_CHECK:
                        dprintf(fd,"EXCEPTION_FLT_STACK_CHECK\n");
                        break;
                    case   EXCEPTION_FLT_UNDERFLOW:
                        dprintf(fd,"EXCEPTION_FLT_UNDERFLOW\n");
                        break;
                    case   EXCEPTION_ILLEGAL_INSTRUCTION:
                        dprintf(fd,"EXCEPTION_ILLEGAL_INSTRUCTION\n");
                        break;
                    case   EXCEPTION_IN_PAGE_ERROR:
                        dprintf(fd,"EXCEPTION_IN_PAGE_ERROR\n");
                        break;
                    case   EXCEPTION_INT_DIVIDE_BY_ZERO:
                        dprintf(fd,"EXCEPTION_INT_DIVIDE_BY_ZERO\n");
                        ProcessExit("PRE");
                        break;
                    case   EXCEPTION_INT_OVERFLOW:
                        dprintf(fd,"EXCEPTION_INT_OVERFLOW\n");
                        break;
                    case   EXCEPTION_INVALID_DISPOSITION:
                        dprintf(fd,"EXCEPTION_INVALID_DISPOSITION\n");
                        break;
                    case   EXCEPTION_NONCONTINUABLE_EXCEPTION:
                        dprintf(fd,"EXCEPTION_NONCONTINUABLE_EXCEPTION\n");
                        break;
                    case   EXCEPTION_PRIV_INSTRUCTION:
                        dprintf(fd,"EXCEPTION_PRIV_INSTRUCTION\n");
                        break;
                    case   EXCEPTION_SINGLE_STEP:
                        dprintf(fd,"EXCEPTION_SINGLE_STEP\n");
                        break;
                    case   EXCEPTION_STACK_OVERFLOW:
                        dprintf(fd,"EXCEPTION_STACK_OVERFLOW\n");
                        break;
                    default:
                        printf("Unknown Exception [0x%x]\n", de.u.Exception.ExceptionRecord.ExceptionCode);
                        break;
                }

                if (de.u.Exception.ExceptionRecord.ExceptionCode==EXCEPTION_BREAKPOINT){
                    dprintf(fd, "EXCEPTION_BREAKPOINT\n");
                    ContinueDebugEvent(de.dwProcessId, de.dwThreadId, DBG_CONTINUE);
                    continue;
                }else{
                    ContinueDebugEvent(de.dwProcessId,de.dwThreadId,DBG_EXCEPTION_HANDLED);
                    continue;
                }
                
            case CREATE_PROCESS_DEBUG_EVENT:
                dprintf(fd, "[CREATE_PROCESS_DEBUG_EVENT]\n");
                break;
            case CREATE_THREAD_DEBUG_EVENT:
                dprintf(fd, "[CREATE_THREAD_DEBUG_EVENT]\n");
                break;
            case EXIT_PROCESS_DEBUG_EVENT:
                dprintf(fd, "[EXIT_PROCESS_DEBUG_EVENT]\n");
                stop = TRUE;
                break;
            case EXIT_THREAD_DEBUG_EVENT:
                dprintf(fd, "[EXIT_THREAD_DEBUG_EVENT]\n");
                break;
            case LOAD_DLL_DEBUG_EVENT:
                dprintf(fd, "[LOAD_DLL_DEBUG_EVENT]\n");
                break;
            case OUTPUT_DEBUG_STRING_EVENT:
                dprintf(fd, "[OUTPUT_DEBUG_STRING_EVENT]\n");
                break;
            case RIP_EVENT:
                dprintf(fd, "[RIP_EVENT]\n");
                break;
            case UNLOAD_DLL_DEBUG_EVENT:
                dprintf(fd, "[UNLOAD_DLL_DEBUG_EVENT]\n");
                break;
      
            default:  
                dprintf(fd, "Unknown Event!\n");
                break;
        }  
  
        if (TRUE == stop) {
            if (!pot&&!pom){
                ProcessExit("PEN");
            }else if(pot){
                ProcessExit("POT");
            }else if(pom){
                ProcessExit("POM");
            }
            
            break;
        }  
  
        ContinueDebugEvent(de.dwProcessId, de.dwThreadId, DBG_CONTINUE);
  
    }
  
    assert(stop);
  
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
  
    dclose(fd);
    
    return 0;  
}  