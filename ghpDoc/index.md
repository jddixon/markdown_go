<h1 class="libTop">markdown_go</h1>

This is **Markdown** implemented in the Go language.

Markdown is a simple protocol for marking up text.  It is widely used for
formatting software documentation.  Basically Markdown supports the use of
less frequently used characters to lay out documents.  Normally Markdown is
used to generate HTML.  So for example, the first line of this document

	# Markdown_go

When this is rurn through Markdown it becomes

	\<H1\>Markdown_go\</H1\>

## Implementations

There have been many implementations of Markdown.  Possibly the best
description of their syntax has been John Gruber's
[Daring Fireball's Markdown Syntax](http://daringfireball.net/projects/markdown/syntax).

Recently a concerted effort by a diverse group has given us
[CommonMark](http://commonmark.org),
which provides a coherent syntax for Markdown and implementations of the
new standard in several languages.

There are alternatives to CommonMark, in particular
[Github Flavored Markdown](https://help.github.com/articles/github-flavored-markdown/),
**GFM**.
Because of Github's very strong position in the marke this is likely to remain
a strong alternative to CommonMark.  However, in the longer run it seems likely
that GFM and CommonMark will merge.

## Project Status

Markdown_go has had two objectives.  First and less importantly it is an
attempt to add a faster markdown to the Go programmer's toolkit.  A second
and more important objective is to provide a library which will convert
Markdown-formatted documents into an absstract syntax tree which in turn
can be converted into, among other things, HTML.

Markdown_go is reasonably far along as an implementation of the Daring
Fireball spec.  That is, it generates good HTML for all but a few problem
cases.  However, it does not as yet provide a decent programming interface
for those interested in generating a syntax tree for further manipulation.

Markdown_go does not as yet conform to the CommonMark spec.

