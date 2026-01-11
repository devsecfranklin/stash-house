AC_PREREQ([2.69])
AC_INIT([stash-house], [0.2], [smoooth.y62wj@passmail.net], [https://github.com/devsecfranklin/stash-house])
AC_CONFIG_MACRO_DIR([config/m4])
AC_CONFIG_AUX_DIR([config])
AM_INIT_AUTOMAKE

AC_CONFIG_COMMANDS([franklin-build],
 [echo "(⊃｡•‿•｡)⊃━⭑･ﾟﾟ･*:༅｡.｡༅:*ﾟ:*:✼✿ Configuring for $mypkgname ☽༓･*˚⁺‧"],
 [mypkgname=$PACKAGE_NAME])

AC_CONFIG_FILES([
  Makefile
])

dnl Checks for programs.
AC_PROG_CXX

AM_PATH_PYTHON(3.9) # minimum version of Python
PY39="python$PYTHON_VERSION" # define the python interpreter
dnl LDFLAGS="$LDFLAGS -l$PY39"
AC_SUBST(PY39, python$PYTHON_VERSION)

AM_CONDITIONAL(HAS_DOCKER, [test "docker" = "found"])
AM_CONDITIONAL(HAS_DOCKER, [test "gcloud" = "found"])

AC_OUTPUT
