package net.duguying.sandbox;

import java.io.FilePermission;
import java.lang.reflect.ReflectPermission;
import java.security.Permission;
import java.security.SecurityPermission;
import java.util.PropertyPermission;

class SandboxSecurityManager extends SecurityManager {

//    public static final int EXIT = ConstantParam.RANDOM.nextInt();

    /**
     * 重写强行退出检测
     * 防止用户自行终止虚拟机的运行,但是调用程序端可以执行退出
     * */
//    public void checkExit(int status) {
//        if (status != EXIT) {
//            throw new SecurityException("Exit On Client Is Not Allowed!");
//        }
//    }

    /**
     * 	策略权限查看
     * 当执行操作时调用,如果操作允许则返回,操作不允许抛出SecurityException
     * */
    private void sandboxCheck(Permission perm) throws SecurityException {
        // 设置只读属性
        if (perm instanceof SecurityPermission) {
            if (perm.getName().startsWith("getProperty")) {
                return;
            }
        } else if (perm instanceof PropertyPermission) {
            if (perm.getActions().equals("read")) {
                return;
            }
        } else if (perm instanceof FilePermission) {
            if (perm.getActions().equals("read")) {
                return;
            }
        } else if (perm instanceof RuntimePermission || perm instanceof ReflectPermission){
            return;
        }

        throw new SecurityException(perm.toString());
    }

    @Override
    public void checkPermission(Permission perm) {
        this.sandboxCheck(perm);
    }

    @Override
    public void checkPermission(Permission perm, Object context) {
        this.sandboxCheck(perm);
    }

}
