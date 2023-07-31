package hashids

import (
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"

	"github.com/speps/go-hashids/v2"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func() (*hashids.HashID, error) {
		data := hashids.NewData()
		data.MinLength = int(config.MinLength)
		if data.MinLength == 0 {
			data.MinLength = 5
		}

		data.Salt = config.Salt
		if data.Salt == "" {
			data.Salt = "default-salt-key"
		}

		if config.Alphabet != "" {
			data.Alphabet = config.Alphabet
		}

		return hashids.NewWithData(data)
	}, o.DiOptions()...)
}
