package net.duguying.sandbox;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.io.PrintStream;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryMXBean;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.Callable;
import java.util.concurrent.FutureTask;

public final class Sandbox {
    private static String VERSION = "0.0.1";

    //用来收集线程的内存使用量
    private static MemoryMXBean _memoryBean = ManagementFactory.getMemoryMXBean();

    //定向输出
    private static ByteArrayOutputStream _baos = new ByteArrayOutputStream(1024);

    //执行提交程序线程
    private static Thread _thread = null;

    //系统提供的默认的classpath
    private static String _classPath = null;

    private static long 	_timeStart = 0;
    //程序运行时间
    private static long 	_timeUsed = 0;
    private static int		_baseMemory = 0;
    //程序运行空间
    private static int 		_memoryUsed = 0;

    //测评结果
    private static String _result = null;

    /**
     * 	核心执行函数<br>
     * 用于执行指定Main.class的main函数.<br>
     * 根据线程ID可获取运行时间.
     * */
    private static String process(int runId, final int timeLimit) throws Exception {
        FutureTask<String> task = null;
        final Timer timer = new Timer(true);

        System.out.printf("[qid] %d\n", runId);

        try {
            final Method mainMethod = getMainMethod(runId);

            System.out.printf("[main method] %s\n", mainMethod.toString());

            task = new FutureTask<String>(new Callable<String>() {
                public String call() throws Exception {
                    try {
                        _timeStart = System.currentTimeMillis();
                        //启动计时器
                        timer.schedule(getTimerTask(),timeLimit + 1);
                        // 执行main方法
                        mainMethod.invoke(null, new Object[] { new String[0] });
                    } catch (InvocationTargetException e) {
                        Throwable targetException = e.getTargetException();
                        // 超过内存限制
                        if (targetException instanceof OutOfMemoryError) {
//                            setResult(JudgeResult.MEMORY_LIMIT_EXCEED);
                            // 非法操作处理
                        }else if (targetException instanceof SecurityException || targetException instanceof ExceptionInInitializerError) {
//                            setResult(JudgeResult.RESTRICT_FUNCTION);
                            // 运行超时
                        } else if (targetException instanceof InterruptedException) {
//                            setResult(JudgeResult.TIME_LIMITE_EXCEED);
                        } else {
                            // 运行时错误处理
                            if (e.getCause().toString().contains("Output Limit Exceed")) {
//                                setResult(JudgeResult.OUTPUT_LIMIT_EXCEED);
                            }else{
//                                setResult(JudgeResult.RUNTIME_ERROR, e.getCause().toString());
                                ;
                            }
                        }

                        throw new RuntimeException("Runtime Exception");
                    } finally {
                        timer.cancel();
                    }
                    return _baos.toString();
                }
            });
        } catch (Exception e) {
            throw new RuntimeException("Initalization Error");
        }

        _baseMemory = (int) _memoryBean.getHeapMemoryUsage().getUsed();
        _thread = new Thread(task);
        _thread.start();
        return task.get();
    }

    /**
     * 	启动计时器,对执行程序进行时间限制.<br>
     * 若超时,则中断执行线程<br>
     * @return TimerTask
     * */
    private static TimerTask getTimerTask(){
        return new TimerTask(){
            public void run() {
                if (_thread  != null)
                    _thread.interrupt();
            }
        };
    }

    private static int _outputSize = 0;

    private static SandboxClassLoader _classLoader  = null;
    /**
     * 	获取指定路径下Main.class类的main入口函数.
     * @param runId 指定类路径
     * @return Method 返回的main方法
     * */
    private static Method getMainMethod(int runId) throws Exception{
        _classLoader = new SandboxClassLoader(_classPath);

        System.out.println("[base path] "+_classLoader.getClassPath() + "q" + runId);

        Class<?> targetClass = _classLoader.loadClass(_classLoader.getClassPath() + "q" + runId, "Main");

        System.out.println("[target class] " + targetClass.toString());

        Method mainMethod = null;
        mainMethod = targetClass.getMethod("main", String[].class);

        if(!Modifier.isStatic(mainMethod.getModifiers()))
            throw new Exception("Method Of Main Is Not Static");

        mainMethod.setAccessible(true);
        return mainMethod ;
    }

    /**
     * 	测评接口.
     * 运行接收到的Java程序.
     *
     * @param runId 执行id
     * @param timeLimit 时间限制
     * @param memoryLimit 空间限制
     * */
    public static void run(int runId, int timeLimit, int memoryLimit) {
        _timeUsed = 0;
        _memoryUsed = 0;
        _baos.reset();
        _outputSize = 0;
//        setResult(JudgeResult.WRONG_ANSWER);
        //定向输入流
//        System.setIn(new BufferedInputStream(new ByteArrayInputStream(standardInput.getBytes())));

        String output = null;

        try {
            // 必须在执行前对垃圾回收,否则不准确.
            System.gc();
            output = process(runId, timeLimit);
            // 将程序输出与标准输出作比较
//            setResult(matchOutput(standardOutput.getBytes(), output.getBytes()));
            // 获取程序运行时间和空间
            _timeUsed = System.currentTimeMillis() - _timeStart;
            _memoryUsed = (int) ((_memoryBean.getHeapMemoryUsage().getUsed() - _baseMemory) / 1000);
        }  catch (Exception e) {
            if (e.getMessage().equals("Initalization Error")){
//                setResult(JudgeResult.WRONG_ANSWER);
            }
        }

        if (_memoryUsed > memoryLimit) {
//            setResult(JudgeResult.MEMORY_LIMIT_EXCEED);
        }

    }

    /**
     * 	沙盒入口
     * 传入参数	: <br>
     * 			classPath -- args[0] ------ 保存class的classpath<br>
     * 			port 		  -- args[1] ------- 监听端口
     * */
    public static void main(String[] args) throws Exception{
        if(args.length < 1){
            System.out.printf("Sandbox for Java\n");
            System.out.printf("Usage:\n");
            System.out.printf("  sandbox <classpath>\n");
            System.out.printf("VERSION %s\n",VERSION);

            return;
        }

        System.out.print("[classpath] " + args[0] + "\n");

        _classPath = args[0];
        run(12,6000,10240);


        SecurityManager security = System.getSecurityManager();
        if (security == null) {
            System.setSecurityManager(new SandboxSecurityManager());
        }
    }
}
