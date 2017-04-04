This program, codenamed shopify2linesheet, allows to create
a wholesale line sheet from a Shopify store in semi-automatic mode.

It processes Shopify product dump into a CSV suitable for a data merge
into an InDesign project with appropriate catalog layout.

Limitations
===========

As it is now, the program can’t work with arbitrary Shopify stores
and line sheet layouts.

It assumes that the incoming CSV has a specific set of columns
in particular order, and it’ll fail to produce expected output otherwise.

* Each Shopify product occupies its own page of line sheet
  (line sheet treats it more as a set of related products).
* Each page is divided into groups of product variants based on value of option 1
  (line sheet treats each group as a product in the set,
  and shows prices at this product level, not per each variant).
* Each individual product variant is defined by value of option 2
  (expected to contain color name, unique in given group of variants),
  SKU (unique across the whole catalog), and a photo.

Line sheet page layout is such that it expects not more than 4 products per set,
where the maximum number of variants for any product is 6;
or it expects not more than 8 products per set with a maximum of 1 variant
per product.

What this program does
======================

This program reads CSV from Shopify, where rows correspond to product variants,
and writes a CSV for InDesign data merge, where each row corresponds
to a page with one line sheet entry.

Workflow
========

1. Export your Shopify data into CSV

2. Assuming you have Go runtime installed::

      go run src/shopify2linesheet/*.go shopifydata.csv linesheetdata.csv  0.45 /Users/path/to/asset/directory

   where the arguments are:
   input file, output file, wholesale discount factor, and product image path.
   Product image path must be an existing directory.

3. In the InDesign project, open Data Merge panel and select 
   ``linesheetdata.csv`` as data source.

4. Proceed as normal: create merged document, adjust layout, publish.

To do
=====

* Separate configuration from hard-coded Shopify CSV parsing logic
  to allow supporting arbitrary Shopify stores.
* Separate configuration from hard-coded line sheet CSV export logic
  to allow supporting arbitrary line sheet layouts.
