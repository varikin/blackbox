#!/bin/bash
dot -Tpng data_flow.dot -Gsize=3,5\! -Gdpi=100 -o images/data_flow.png
pandoc README.md --pdf-engine=wkhtmltopdf -o report.pdf
