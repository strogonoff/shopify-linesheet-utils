NOTE: No support for arbitrary Shopify stores
and catalog designs. This probably won’t work for you.

Rationale
=========

If creating a fully custom catalog and stock Shopify
line sheet exporter doesn’t cut it, using this
in conjunction with Adobe InDesign
allows for quick iteration over catalog design
without manually entering every product’s data.

Functionality
=============

Takes a Shopify store CSV dump and produces
a CSV ready for data merge into InDesign-created catalog.

Workflow
========

1. Export Shopify data into CSV

2. Assuming Go runtime is installed on the machine::

      go run src/shopify2linesheet/*.go shopifydata.csv linesheetdata.csv  0.45 /Users/path/to/asset/directory

   where the arguments are:
   input file, output file, wholesale discount factor, and product image path.

   It won’t download product image when it sees that the file with the same
   name already exists.

3. In InDesign project, open Data Merge panel and select 
   ``linesheetdata.csv`` as data source.

4. Proceed as normal: create merged document, adjust layout, publish.

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
