package hashids

import (
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"
	"github.com/samber/lo"

	"github.com/speps/go-hashids/v2"
)

type HashID struct {
	config   *Config
	instance *hashids.HashID
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	if config.MinLength == 0 {
		config.MinLength = 5
	}

	if config.Alphabet == "" {
		config.Alphabet = DefaultAlphabet
	}

	if config.Salt == "" {
		config.Salt = "default-salt-key"
	}

	return container.Container.Provide(func() (*HashID, error) {
		data := hashids.NewData()
		data.MinLength = int(config.MinLength)
		data.Salt = config.Salt
		data.Alphabet = config.Alphabet

		hashid, err := hashids.NewWithData(data)
		if err != nil {
			return nil, err
		}

		return &HashID{
			instance: hashid,
			config:   &config,
		}, nil
	}, o.DiOptions()...)
}

func (h *HashID) EncodeInt64(id int64) (string, error) {
	return h.instance.EncodeInt64([]int64{id})
}

func (h *HashID) MustEncodeInt64(id int64) string {
	return lo.Must1(h.instance.EncodeInt64([]int64{id}))
}

func (h *HashID) EncodeWithSalt(salt string, id int64) (string, error) {
	ins, err := hashids.NewWithData(&hashids.HashIDData{
		Salt:      salt,
		Alphabet:  h.config.Alphabet,
		MinLength: int(h.config.MinLength),
	})
	if err != nil {
		return "", err
	}
	return ins.EncodeInt64([]int64{id})
}

func (h *HashID) MustEncodeWithSalt(salt string, id int64) string {
	return lo.Must1(h.EncodeWithSalt(salt, id))
}

func (h *HashID) Decode(hash string) ([]int64, error) {
	return h.instance.DecodeInt64WithError(hash)
}

func (h *HashID) MustDecode(hash string) []int64 {
	return lo.Must1(h.Decode(hash))
}

func (h *HashID) DecodeWithSalt(salt, hash string) ([]int64, error) {
	ins, err := hashids.NewWithData(&hashids.HashIDData{
		Salt:      salt,
		Alphabet:  h.config.Alphabet,
		MinLength: int(h.config.MinLength),
	})
	if err != nil {
		return nil, err
	}
	return ins.DecodeInt64WithError(hash)
}

func (h *HashID) MustDecodeWithSalt(salt, hash string) []int64 {
	return lo.Must1(h.DecodeWithSalt(salt, hash))
}
