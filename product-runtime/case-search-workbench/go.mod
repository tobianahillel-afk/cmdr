module github.com/tobianahillel-afk/cmdr/product-runtime/case-search-workbench

go 1.26.10

require (
	github.com/tobianahillel-afk/cmdr/product-runtime/event-search v0.0.0
	github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline v0.0.0
)

replace github.com/tobianahillel-afk/cmdr/product-runtime/event-search => ../event-search

replace github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline => ../investigation-timeline
