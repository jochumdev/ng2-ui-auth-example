package registry

import (
	"errors"
	"log"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo"
)

var (
	ErrComponentNotFound = errors.New("Component not found.")
)

type Component interface {
	SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error
	SetupEcho(e *echo.Echo) error

	Shutdown() error

	RegistrySet(*Registry)

	NameGet() string
	WeightGet() int
}

var instance *Registry

func Instance() *Registry {
	if instance != nil {
		return instance
	}

	instance = &Registry{
		cnames: make(map[string](Component)),
	}

	return instance
}

type Registry struct {
	cnames map[string](Component)
	clist  []Component

	debug bool
}

func (r *Registry) Register(name string, c Component) {
	c.RegistrySet(r)
	r.cnames[name] = c

	idx := -1
	found := false
	for i, lcom := range r.clist {
		idx = i
		if lcom.WeightGet() < c.WeightGet() {
			continue
		} else {
			found = true
			break
		}
	}
	if !found {
		idx = idx + 1
	}

	if idx == 0 {
		tmp := r.clist
		r.clist = append([]Component{}, c)
		r.clist = append(r.clist, tmp...)
	} else if idx == len(r.clist) {
		r.clist = append(r.clist, c)
	} else {
		after := append([]Component{}, r.clist[idx:]...)
		r.clist = append(r.clist[:idx], c)
		r.clist = append(r.clist, after...)
	}
}

func (r *Registry) List() []Component {
	return r.clist
}

func (r *Registry) ComponentGet(name string) (Component, error) {
	var (
		result Component
		ok     bool
	)
	if result, ok = r.cnames[name]; !ok {
		return nil, ErrComponentNotFound
	}

	return result, nil
}

func (r *Registry) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	r.debug = debug

	for _, c := range r.clist {
		if r.debug {
			log.Printf("SetupFromIni (%d): %s\n", c.WeightGet(), c.NameGet())
		}
		if err := c.SetupFromIni(iniCfg, configFile, debug); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) SetupEcho(e *echo.Echo) error {
	for _, c := range r.clist {
		if err := c.SetupEcho(e); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) Shutdown() error {
	i := len(r.clist) - 1
	for true {
		c := r.clist[i]
		if r.debug {
			log.Printf("Shutdown (%d): %s\n", c.WeightGet(), c.NameGet())
		}
		if err := c.Shutdown(); err != nil {
			return err
		}

		i = i - 1

		if i < 0 {
			break
		}
	}

	return nil
}
