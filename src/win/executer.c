#include "executer.h"
  
#define MAX_PARAM_LEN       4096  
  
int main( int argc, char ** argv )  
{  
    int i, j = 0, len;  
    char command_buf[MAX_PARAM_LEN];  
  
    STARTUPINFO si;  
    PROCESS_INFORMATION pi;  
    DEBUG_EVENT de;  
    BOOL stop = FALSE;  
  
    ZeroMemory( &si, sizeof(si) );  
    si.cb = sizeof(si);  
    ZeroMemory( &pi, sizeof(pi) );  
  
    if (argc<2) {  
        printf("Usage: %s <app_name> [arguments ...]\n", argv[0]);  
        return 0;  
    }  
  
    // Combine the module name and params into one string.  
    for (i = 1; i < argc; ++i) {  
        len = strlen(argv[i]);  
        if (len >= MAX_PARAM_LEN - j - 1) {  
            printf("buffer overflow\n");  
            exit(-1);  
        }  
        j += _snprintf(command_buf + j, MAX_PARAM_LEN - j, "%s ", argv[i]);  
        command_buf[j] = '\0';  // just for sure  
    }  
  
    if( !CreateProcess(NULL, command_buf, NULL, NULL, FALSE,            
        DEBUG_ONLY_THIS_PROCESS, NULL, NULL, &si, &pi ) ) {  
            printf( "CreateProcess failed (%d).\n", GetLastError() );  
            exit(-1);  
    }  
  
    while (TRUE) {  
        WaitForDebugEvent (&de, INFINITE);  
  
        switch (de.dwDebugEventCode) {  
        case EXCEPTION_DEBUG_EVENT:         /* exception */  
            switch (de.u.Exception.ExceptionRecord.ExceptionCode) {   
            case   EXCEPTION_INT_DIVIDE_BY_ZERO:    /* #DE */  
                // Do what the parent process want to do when the child process gets #DE interrupt.  
                TerminateProcess(pi.hProcess,1);   
                break;   
            case   EXCEPTION_BREAKPOINT:            /* #BP */  
                // Do what the parent process want to do when the child process gets #BP interrupt.  
                break;  
  
            default:   
                printf("Unknown Exception\n");   
                break;  
            }      
  
            ContinueDebugEvent(de.dwProcessId,de.dwThreadId,DBG_EXCEPTION_HANDLED);  
            continue;  
  
        case CREATE_PROCESS_DEBUG_EVENT:        /* child process created */  
  
            // Do what the parent process want to do when the child process was created.  
            break;  
  
        case EXIT_PROCESS_DEBUG_EVENT:          /* child process exits */  
            stop = TRUE;  
  
            // Do what the parent process want to do when the child process exits.  
            break;  
  
        default:  
            printf("Unknown Event!\n");  
            break;  
        }  
  
        if (TRUE == stop) {  
            //printf("Process exit\n");  
            break;  
        }  
  
        ContinueDebugEvent (de.dwProcessId, de.dwThreadId, DBG_CONTINUE);  
  
    } // end of loop  
  
    assert(stop);  
  
    CloseHandle( pi.hProcess );  
    CloseHandle( pi.hThread );  
  
    return 0;  
}  