AC_PREREQ([2.69])
AC_INIT([lab-franklin], [0.2], [franklin@bitsmasher.net], [https://github.com/devsecfranklin/lab-franklin])
AC_CONFIG_MACRO_DIR([config/m4])
AC_CONFIG_AUX_DIR([config])
AM_INIT_AUTOMAKE
LT_INIT # Initialize Libtool for shared library support

AC_ARG_ENABLE([debug],
[  --enable-debug    Turn on debugging],
[case "${enableval}" in
  yes) debug=true ;;
  no)  debug=false ;;
  *) AC_MSG_ERROR([bad value ${enableval} for --enable-debug]) ;;
esac],[debug=false])
AM_CONDITIONAL([DEBUG], [test x$debug = xtrue])

LT_INIT

AC_CONFIG_COMMANDS([franklin-build],
 [echo "(⊃｡•‿•｡)⊃━⭑･ﾟﾟ･*:༅｡.｡༅:*ﾟ:*:✼✿ Configuring for $mypkgname ☽༓･*˚⁺‧"],
 [mypkgname=$PACKAGE_NAME])

AC_CONFIG_FILES([
  Makefile
  ansible/Makefile
  assembly/Makefile
  assembly/src/Makefile
  container/Makefile
  container/kerberos-containers/kdc-server/Makefile
  container/kerberos-containers/kdc-ssh-server/Makefile
  container/kerberos-containers/ssh-container/Makefile
  test/test_cluster/Makefile
  test/test_gui/Makefile
])

AC_OUTPUT

