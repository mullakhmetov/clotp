package config

import (
	"errors"
	"reflect"
	"testing"
)

type stubMapper struct {
	ItemsToRead []*Item
	WriteBuffer []*Item
}

func (m stubMapper) Read() ([]*Item, error) {
	return m.ItemsToRead, nil
}
func (m *stubMapper) Write(items []*Item) error {
	m.WriteBuffer = items
	return nil
}

func TestConfigRead(t *testing.T) {
	wantItems := []*Item{
		{Name: "n1", Key: "k1"},
		{Name: "n2", Key: "k2"},
	}
	config := &Config{mapper: &stubMapper{ItemsToRead: wantItems}}

	err := config.Read()
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(config.Items, wantItems) {
		t.Errorf("wrong config items, want: %+v != got: %+v", wantItems, config.Items)
	}
}

func TestConfigWrite(t *testing.T) {
	items := []*Item{
		{Name: "n1", Key: "k1"},
		{Name: "n2", Key: "k2"},
	}
	mapper := &stubMapper{}
	config := &Config{mapper: mapper, Items: items}

	err := config.Write()
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(mapper.WriteBuffer, items) {
		t.Errorf("wrong config items, want: %+v != got: %+v", items, mapper.WriteBuffer)
	}
}

func TestItemValidate(t *testing.T) {
	for _, c := range []struct {
		name string
		item Item
		want bool
	}{
		{
			name: "no name",
			item: Item{Key: "1"},
			want: false,
		},
		{
			name: "no key",
			item: Item{Name: "n"},
			want: false,
		},
		{
			name: "negative step",
			item: Item{Name: "n", Step: -1},
			want: false,
		},
		{
			name: "valid",
			item: Item{Name: "n", Key: "k"},
			want: true,
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := c.item.Validate(); got != c.want {
				t.Errorf("wrong validation result, should be %t", c.want)
			}
		})
	}
}

func TestConfigAdd(t *testing.T) {
	config := Config{}

	for _, c := range []struct {
		name string
		item *Item
		err  error
	}{
		{
			name: "no name",
			item: &Item{Key: "1"},
			err:  ErrInvalidItem,
		},
		{
			name: "no key",
			item: &Item{Name: "n"},
			err:  ErrInvalidItem,
		},
		{
			name: "negative step",
			item: &Item{Name: "foo", Key: "bar", Step: -1},
			err:  ErrInvalidItem,
		},
		{
			name: "valid #1",
			item: &Item{Name: "n", Key: "k"},
		},
		{
			name: "duplicate",
			item: &Item{Name: "n", Key: "k2"},
			err:  ErrItemAlreadyExists,
		},
		{
			name: "valid #2",
			item: &Item{Name: "n2", Key: "k"},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			err := config.Add(c.item)
			if err != nil {
				if c.err == nil {
					t.Errorf("unwanted error: %v", err)
					return
				}

				if !errors.Is(err, c.err) {
					t.Errorf("error should match with %v", c.err)
					return
				}

				return
			}

			if c.err != nil {
				t.Error("error shouldn't be nil")
				return
			}
		})
	}

	want := []*Item{
		{Name: "n", Key: "k"},
		{Name: "n2", Key: "k"},
	}

	// function type is incomparable
	for _, i := range config.Items {
		i.digest = nil
	}

	if !reflect.DeepEqual(want, config.Items) {
		t.Errorf("wrong config items state, want: %+v != got: %+v", want, config.Items)
	}
}

func TestNewFromConfigItem(t *testing.T) {
	item := &Item{
		Name:      "n",
		Issuer:    "issuer",
		Key:       "12345678901234567890",
		Algorithm: "sha1",
		Digits:    6,
		Step:      30,
	}

	totp := NewFromConfigItem(item)

	if item.Digits != totp.Digits {
		t.Errorf("wrong digits value")
	}

	if item.Key != totp.Secret {
		t.Errorf("wrong secret value")
	}

	if item.Step != totp.TimeStep {
		t.Errorf("wrong timestep value")
	}
}
