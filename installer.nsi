; 该脚本使用 易量安装(az.eliang.com) 向导生成
; 安装程序初始定义常量
!define PRODUCT_NAME "Goj-Judger"
!define PRODUCT_VERSION "1.0"
!define PRODUCT_PUBLISHER "duguying"
!define PRODUCT_WEB_SITE "http://oj.duguying.net"
!define PRODUCT_DIR_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\judger.exe"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"

SetCompressor lzma

; 提升安装程序权限(vista,win7,win8)
RequestExecutionLevel admin

; ------ MUI 现代界面定义 ------
!include "MUI2.nsh"

; MUI 预定义常量
!define MUI_ABORTWARNING
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP "etc\header\install.bmp"
!define MUI_HEADERIMAGE_UNBITMAP "etc\header\uninstall.bmp"
!define MUI_ICON "etc\icon\install.ico"
!define MUI_UNICON "etc\icon\uninstall.ico"
!define MUI_WELCOMEFINISHPAGE_BITMAP "etc\wizard\install.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "etc\wizard\uninstall.bmp"

; 欢迎页面
!insertmacro MUI_PAGE_WELCOME
; 许可协议页面
!insertmacro MUI_PAGE_LICENSE "LICENSE"
; 安装目录选择页面
!insertmacro MUI_PAGE_DIRECTORY
; 安装过程页面
!insertmacro MUI_PAGE_INSTFILES
; 安装完成页面
!define MUI_FINISHPAGE_RUN "$INSTDIR\judger.exe"
!insertmacro MUI_PAGE_FINISH

; 安装卸载过程页面
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

; 安装界面包含的语言设置
!insertmacro MUI_LANGUAGE "SimpChinese"

; ------ MUI 现代界面定义结束 ------

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "setup.exe"
;ELiangID 统计编号     /*  安装统计项名称：【Goj-Judger】  */
InstallDir "$PROGRAMFILES\Goj-Judger"
InstallDirRegKey HKLM "${PRODUCT_UNINST_KEY}" "UninstallString"
ShowInstDetails show
ShowUninstDetails show
BrandingText "Goj-Golang based Online Judge"

;安装包版本号格式必须为x.x.x.x的4组整数,每组整数范围0~65535,如:2.0.1.2
;若使用易量统计,版本号将用于区分不同版本的安装情况,此时建议用户务必填写正确的版本号
;!define INSTALL_VERSION "2.0.1.2"

;VIProductVersion "${INSTALL_VERSION}"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "ProductName"      "Goj-Judger"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "Comments"         "Goj-Judger(duguying)"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "CompanyName"      "duguying"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "LegalCopyright"   "duguying(http://duguying.net)"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "FileDescription"  "Goj-Judger"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "ProductVersion"   "${INSTALL_VERSION}"
;VIAddVersionKey /LANG=${LANG_SimpChinese} "FileVersion"      "${INSTALL_VERSION}"

Section "MainSection" SEC01
  SetOutPath "$INSTDIR"
  SetOverwrite ifnewer

  File "/oname=judger.exe" "judger.exe"

  SetOutPath "$INSTDIR\bin"
  File "/oname=executer.exe" "sandbox\c\build\executer.exe"
  File "/oname=executer.json" "sandbox\c\build\executer.json"

  SetOutPath "$INSTDIR\conf"
  File "/oname=config.ini" "conf\config.ini"

  CreateDirectory "$SMPROGRAMS\Goj-Judger"
  CreateShortCut "$SMPROGRAMS\Goj-Judger\Goj-Judger.lnk" "$INSTDIR\judger.exe"
  CreateShortCut "$DESKTOP\Goj-Judger.lnk" "$INSTDIR\judger.exe"
SectionEnd

Section -AdditionalIcons
  WriteINIStr "$INSTDIR\${PRODUCT_NAME}.url" "InternetShortcut" "URL" "${PRODUCT_WEB_SITE}"
  CreateShortCut "$SMPROGRAMS\Goj-Judger\访问Goj-Judger主页.lnk" "$INSTDIR\${PRODUCT_NAME}.url"
  CreateShortCut "$SMPROGRAMS\Goj-Judger\卸载Goj-Judger.lnk" "$INSTDIR\uninst.exe"
SectionEnd

Section -Post
  WriteUninstaller "$INSTDIR\uninst.exe"
  WriteRegStr HKLM "${PRODUCT_DIR_REGKEY}" "" "$INSTDIR\judger.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayName" "$(^Name)"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\uninst.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\judger.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
SectionEnd

/******************************
*  以下是安装程序的卸载部分  *
******************************/

Section Uninstall
  Delete "$INSTDIR\${PRODUCT_NAME}.url"
  Delete "$INSTDIR\uninst.exe"
  Delete "$INSTDIR\judger.exe"
  Delete "$INSTDIR\bin\executer.exe"
  Delete "$INSTDIR\bin\executer.json"

  RMDir "$INSTDIR\bin"

  Delete "$INSTDIR\conf\config.ini"

  RMDir "$INSTDIR\conf"

  Delete "$SMPROGRAMS\Goj-Judger\访问Goj-Judger主页.lnk"
  Delete "$SMPROGRAMS\Goj-Judger\卸载Goj-Judger.lnk"
  Delete "$SMPROGRAMS\Goj-Judger\Goj-Judger.lnk"
  Delete "$DESKTOP\Goj-Judger.lnk"

  RMDir "$SMPROGRAMS\Goj-Judger"

  RMDir "$INSTDIR"

  DeleteRegKey ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}"
  DeleteRegKey HKLM "${PRODUCT_DIR_REGKEY}"
SectionEnd

/* 根据 NSIS 脚本编辑规则,所有 Function 区段必须放置在 Section 区段之后编写,以避免安装程序出现未可预知的问题. */

Function un.onInit
FunctionEnd

Function un.onUninstSuccess
FunctionEnd
