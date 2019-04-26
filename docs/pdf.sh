#!/bin/bash
pandoc -t html5 -s README.md --pdf-engine=wkhtmltopdf -o report.pdf
