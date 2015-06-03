package main

import (
	"fmt"
	"github.com/coreos/fleet/client"
	"net"
	"net/http"
	"net/url"
)

func getUnitMachineIP(c client.API, u string) (string, error) {
	id, err := getUnitMachineID(c, u)
	if err != nil {
		return "", err
	}
	unitMachineIP, err := getMachineIP(c, id)
	if err != nil {
		return "", err
	}
	return unitMachineIP, nil
}

func getUnitMachineID(c client.API, u string) (string, error) {
	unit, err := c.Unit(u)
	if err != nil {
		return "", err
	}
	if unit == nil || unit.CurrentState != "launched" {
		return "", fmt.Errorf("cannot get unit's machine id.")
	}
	return unit.MachineID, nil
}

// getMachineIP function takes a machine ID and returns its IP address
func getMachineIP(c client.API, id string) (string, error) {
	machines, err := c.Machines()
	if err != nil {
		return "", err
	}
	for _, m := range machines {
		if m.ID == id {
			return m.PublicIP, nil
		}
	}
	return "", fmt.Errorf("%s", "machine id not found.")
}

func getClient(endpoint string) (client.API, error) {
	ep, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	dialFunc := net.Dial
	if ep.Scheme == "unix" {
		ep.Scheme = "http"
		ep.Host = "domain-sock"
		sockPath := ep.Path
		ep.Path = ""
		dialFunc = func(network, addr string) (net.Conn, error) {
			// return net.Dial("unix", ep.Path)
			return net.Dial("unix", sockPath)
		}

	}
	c := &http.Client{
		Transport: &http.Transport{
			Dial:              dialFunc,
			DisableKeepAlives: true,
		},
	}
	return client.NewHTTPClient(c, *ep)
}
