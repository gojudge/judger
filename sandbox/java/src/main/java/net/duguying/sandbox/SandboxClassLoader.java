package net.duguying.sandbox;

import java.io.File;
import java.io.FileInputStream;

class SandboxClassLoader extends ClassLoader{
    /**默认classPath*/
    private String _classPath;
    private String seperater = ".";

    /**
     * 	构造函数
     * @param classPath 类加载器默认classPath
     * */
    public SandboxClassLoader(String classPath) {
        this._classPath = classPath;
    }

    @Override
    protected Class<?> findClass(String className) throws ClassNotFoundException {
        return loadClass(_classPath, className);
    }

    /**
     * 	更改类加载器加载类的classpath,在制定路径下加载制定的类class文件
     * @param		classPath	要加载的类路径
     * @param		className 	要加载的类名
     * 				最为限定,只能加载不含包的类.
     * */
    public Class<?> loadClass(String classPath, String className) throws ClassNotFoundException{
        if(className.indexOf('.') >= 0) {
            throw new ClassNotFoundException(className);
        }

        File classFile = new File(classPath + seperater + className + ".class");
        byte[] mainClass = new byte[(int) classFile.length()];
        try {
            FileInputStream in = new FileInputStream(classFile);
            in.read(mainClass);
            in.close();
        } catch (Exception e) {
            //e.printStackTrace();
            throw new ClassNotFoundException(className);
        }

        return super.defineClass(className, mainClass, 0, mainClass.length);
    }

    /**
     * 	获取classPath
     * @return String		classPath
     * */
    public String getClassPath(){
        return _classPath + seperater;
    }
}
