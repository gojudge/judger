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

    //用来收集线程的内存使用量
    private static MemoryMXBean _memoryBean = ManagementFactory.getMemoryMXBean();

    //定向输出
    private static ByteArrayOutputStream _baos = new ByteArrayOutputStream(1024);

    //Socket通信
    private static Socket _socket = null;

    private static ServerSocket _serverSocket = null;

    private static ObjectInputStream _inputStream = null;

    private static ObjectOutputStream _outputStream = null;

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

        try {
            final Method mainMethod = getMainMethod(runId);

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
    /**
     * 	初始化
     * @param classPath	系统默认classPath路径
     * 	@param	port		socket服务器监听端口
     * */
    private static void inita(String classPath, int port) throws Exception{
        _classPath = classPath;

        _serverSocket = new ServerSocket(port);
        _socket = _serverSocket.accept();

        _outputStream = new ObjectOutputStream(_socket.getOutputStream());
        _inputStream = new ObjectInputStream(_socket.getInputStream());

        //重新定向输出流
        System.setOut(new PrintStream(new BufferedOutputStream(_baos) {
            public void write(byte[] b, int off, int len) throws IOException {
                _outputSize += len - off;
                try {
                    super.write(b, off, len);
//                    if (_outputSize > ConstantParam.OUTPUT_LIMIT) {
//                        throw new RuntimeException("Output Limit Exceed" + _outputSize);
//                    }
                } catch (IOException e) {
                    if(e.getMessage().equals("Output Limit Exceed")){
                        throw e;
                    }
                }
            }
        }));
    }

    private static SandboxClassLoader _classLoader  = null;
    /**
     * 	获取指定路径下Main.class类的main入口函数.
     * @param runId 指定类路径
     * @return Method 返回的main方法
     * */
    private static Method getMainMethod(int runId) throws Exception{
        _classLoader = new SandboxClassLoader(_classPath);
        Class<?> targetClass = _classLoader.loadClass(_classLoader.getClassPath() + runId, "Main");

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
     * @param standardInput 程序标准输入字符串
     * @param standardOutput 程序标准输出字符串
     * */
    public static void run(int runId, int timeLimit, int memoryLimit,
                           String standardInput, String standardOutput) {
        _timeUsed = 0;
        _memoryUsed = 0;
        _baos.reset();
        _outputSize = 0;
//        setResult(JudgeResult.WRONG_ANSWER);
        //定向输入流
        System.setIn(new BufferedInputStream(new ByteArrayInputStream(standardInput.getBytes())));

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

        try {
            //向主模块返回执行结果
            sendResult(runId, (int)_timeUsed, _memoryUsed, _result);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    /**
     * 	向主模块发送运行结果.
     *
     * @param runId 运行runId
     * 	@param timeUsed	代码运行时间(MS)
     * @param	 memoryUsed	代码运行空间(B)
     * 	@param result	 代码执行结果
     * */
    private static void sendResult(int runId,int timeUsed, int memoryUsed, String result) throws IOException{
        _outputStream.writeInt(runId);
        _outputStream.writeInt(timeUsed);
        _outputStream.writeInt(memoryUsed);
        _outputStream.writeUTF(result);
    }

    /**
     * 	接收运行参数
     *
     * */
    private static void receiveMsg() throws IOException{
        int runId = _inputStream.readInt();
        int timeLimit = _inputStream.readInt();
        int memoryLimit = _inputStream.readInt();
        String standardInput = _inputStream.readUTF();
        String standardOutput = _inputStream.readUTF();

        run(runId, timeLimit, memoryLimit, standardInput, standardOutput);
    }

    /**
     * 	关闭网络连接
     * */
    private static void close(){
        try {
            if (_inputStream != null)
                _inputStream.close();
            if (_outputStream != null)
                _outputStream.close();
            if (_socket != null)
                _socket.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    /**
     * 	沙盒入口
     * 传入参数	: <br>
     * 			classPath -- args[0] ------ 保存class的classpath<br>
     * 			port 		  -- args[1] ------- 监听端口
     * */
    public static void main(String[] args) throws Exception{
        inita(args[0], Integer.parseInt(args[1]));

        SecurityManager security = System.getSecurityManager();
        if (security == null)
            System.setSecurityManager(new SandboxSecurityManager());

        while (!_socket.isClosed()){
            receiveMsg();
        }

        close();
    }
}
