package utils

import "go.uber.org/fx"

func CreateJustProvider(concreteType interface{}) fx.Option {
	return fx.Provide(concreteType)
}

func CreateAnnotatedProvider(constructor any, returnType any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(returnType),
		),
	)
}

func CreateHandlerConsumer(constructor any, returnType any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(returnType),
			fx.ParamTags("", "", "", "", `group:"handler"`),
		),
	)
}

func CreateHandlerProvider(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.ResultTags(`group:"handler"`),
		),
	)
}
