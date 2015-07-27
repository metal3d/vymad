# Vymad

Vymad is a markdown generator from "vym" file format. Vym (View Your Mind) is a very nice Mindmapping software for \*Nix environments. Writing a book, I was searching a way to easilly write content in Vym and to generate something that cas be used by Pandoc.

I decided to build my own in Go and to share this little tool to make your life easier :)

# Installation

## Built package

I provide binary files for Linux, OSX and FreeBSD. You can download the specific binary for you environment from the release page:

## With Go

If you've got "golang" (go) on you computer, you may use:

```
go get -u github.com/metal3d/vymad
```

## Test

Test if you correctly installed "vymad" by typing this in a terminal:

```
vymad -version
```

If no version appears or if you have an error, please check how you've installed vymad.

# Usage

Simply pass a vym file as argument:

```
vymad myfile.vym
```

If the given file is a correct vym file, so the markdown content will be displayed.

If you passed a zip file, vymad will search the first xml file at the root of the archive. If he find one, so it will try to parse.

Note that vymad will never write something on the given file.

To keep the output in a new file, simply use shell redirection:

```
vymad myfile.vym > myfile.md
```

You may use specific shell syntax to build pdf or other format on-the-fly:

```
# With bash, sh, or zsh
pandoc --toc --chapter <(vymad myfile.vym) -o book.pdf

# With fish
pandoc --toc --chapter (vymad myfile.vym | psub) -o book.pdf

# generic and universal AFAIK
vymad myfile.vym | pandoc --toc --chapter -o book.pdf

```


