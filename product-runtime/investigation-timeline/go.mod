module github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline

go 1.26.9

require (
	github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection v0.0.0
	github.com/tobianahillel-afk/cmdr/product-runtime/event-search v0.0.0
)

replace github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection => ../event-inspection

replace github.com/tobianahillel-afk/cmdr/product-runtime/event-search => ../event-search
