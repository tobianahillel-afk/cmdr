module github.com/tobianahillel-afk/cmdr/product-runtime/command-investigate-host

go 1.26.10

require (
	github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline v0.0.0
	github.com/tobianahillel-afk/cmdr/product-runtime/work-queue v0.0.0
)

replace github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline => ../investigation-timeline

replace github.com/tobianahillel-afk/cmdr/product-runtime/work-queue => ../work-queue
