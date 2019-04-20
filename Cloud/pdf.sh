#!/bin/bash
dot -Tpng docs/data_flow.dot -Gsize=3,5\! -Gdpi=100 -o docs/images/data_flow.png
pandoc README.md --pdf-engine=wkhtmltopdf -o docs/output.pdf
