This program, codenamed shopify2linesheet,
processes Shopify product dump into a CSV suitable for a data merge
into InDesign project with appropriate catalog layout.

The original purpose is semi-automatic creation of linesheets.

Note: Right now the program is inflexible and assumes
an InDesign project with singular, very specific layout
and placeholder configuration.

What this program does
======================

This program takes Shopify’s CSV, which is a peculiarly formatted file,
creates products/sets/variants from it,
and then writes a CSV where each row corresponds
to a page with one line sheet entry
(this is where it only supports one specific layout).

Workflow
========

1. Export your Shopify data into CSV

2. Assuming you have Go runtime installed::

      go run src/shopify2linesheet/*.go shopifydata.csv linesheetdata.csv  0.45 /Users/path/to/asset/directory

   where the arguments are:
   input file, output file, wholesale discount factor, and product image path.

3. In the InDesign project, open Data Merge panel and select 
   ``linesheetdata.csv`` as data source.

4. Proceed as normal: create merged document, adjust layout, publish.
