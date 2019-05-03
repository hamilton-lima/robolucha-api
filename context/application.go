package context

import (
	_ "gitlab.com/robolucha/robolucha-api/datasource"
	_ "gitlab.com/robolucha/robolucha-api/publisher"
)

type ApplicationContext struct {
	ds  *DataSource
	pub Publisher
}

var Context ApplicationContext
