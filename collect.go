package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

const swedenTotalCode = "00"

type RegionMap map[string]string

func Regions(api API) (RegionMap, error) {
	md, err := api.VotingRatesMetadata()
	if err != nil {
		return nil, err
	}
	for _, vs := range md.Variables {
		if vs.Text == "region" {
			r := make(map[string]string)
			for i, v := range vs.Values {
				if v != swedenTotalCode {
					r[v] = vs.ValueTexts[i]
				}
			}
			return r, nil
		}
	}
	return nil, errors.New("No region variable found in metadata")
}

type Rate struct {
	Regs []string
	Year string
	Pct  float64
}

func Rates(api API, regs RegionMap) ([]Rate, error) {
	tr, err := api.VotingRatesQuery()
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Rate)
	for _, d := range tr.Data {
		reg := d.Key[0]
		if reg == swedenTotalCode || d.Values[0] == ".." {
			continue
		}
		reg, ok := regs[reg]
		if !ok {
			return nil, fmt.Errorf("Region %s does not exist in metadata", d.Key[0])
		}
		year := d.Key[1]
		pct, err := strconv.ParseFloat(d.Values[0], 64)
		if err != nil {
			return nil, err
		}
		if e, ok := m[year]; ok {
			if pct == e.Pct {
				e.Regs = append(e.Regs, reg)
				sort.Strings(e.Regs)
			} else if pct > e.Pct {
				e.Regs = []string{reg}
				e.Pct = pct
			}
		} else {
			m[year] = &Rate{
				Regs: []string{reg},
				Year: year,
				Pct:  pct,
			}
		}
	}
	rates, idx := make([]Rate, len(m)), 0
	for _, v := range m {
		rates[idx] = *v
		idx++
	}
	rs := RateSlice(rates)
	sort.Sort(rs)
	return rates, nil
}

type RateSlice []Rate

func (s RateSlice) Len() int {
	return len(s)
}

func (s RateSlice) Less(i, j int) bool {
	return s[i].Year < s[j].Year
}

func (s RateSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
