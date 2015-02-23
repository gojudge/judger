package net.duguying.sandbox;

import java.io.FilePermission;
import java.lang.reflect.ReflectPermission;
import java.security.Permission;
import java.security.SecurityPermission;
import java.util.PropertyPermission;

class SandboxSecurityManager extends SecurityManager {

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
