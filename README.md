# FileBisect
Find the conflicting file (e.g. mod) in a folder. Doesn't work with nested folders (recursion), they are ignored.

## How to install it
1. Install Go
2. Run `go get -u github.com/comp500/filebisect`

## How to use it
1. Open the command prompt / terminal / whatever in the folder you want to bisect
2. Run `filebisect index` to generate an index of the files in the folder
3. Run `filebisect ignore [file]` to tell FileBisect which files you don't want to change (e.g. the mod you want to test compatibility with)
4. Run `filebisect bad` to state that the current set of files has a conflict. FileBisect will now remove some of the files (moved to a temporary directory) so you can test with only some of the files.
5. Now test the system (e.g. modpack) again. If there is a conflict, run `filebisect bad`. If not, run `filebisect good`.
6. When the bisection is done, FileBisect will say "Done!". The results are in `file-bisect-index.toml`, just Ctrl+F for `bad`.

Note: If some of the files need to be there at the same time (e.g. dependencies), just move them from the temporary directory (see `file-bisect-index.toml` for the location, it is in `%TEMP%` on Windows) after running `filebisect bad` or `filebisect good`. FileBisect will automatically recognise where the files are each time you run it.

## To Do
- Dependency checking
- CurseForge dependency import
- Recursion support (nested folders)
- `.disabled` support (instead of moving to a temporary folder)
- GUI
