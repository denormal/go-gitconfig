/*
Package gitconfig provides an interface to git configuration properties
as returned by "git config --list". gitconfig provides access to local,
global and system configuration, as well as the effective configuration
for the given git working copy. gitconfig attempts to use the locally installed
"git" executable using https://github.com/denormal/go-gittools.

See https://git-scm.com/docs/git-config for more information.
*/
package gitconfig
