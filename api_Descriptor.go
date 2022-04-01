package forGraphBLASGo

type (
	DescField int
	DescValue int
)

const (
	Outp DescField = iota
	Mask
	Inp0
	Inp1
	nDescFields
)

const (
	reserved DescValue = iota
	Replace
	Comp
	Tran
	Structure
)

type Descriptor []DescValue

var (
	DescT1      = Descriptor{Inp1: Tran}
	DescT0      = Descriptor{Inp0: Tran}
	DescT0T1    = Descriptor{Inp0: Tran, Inp1: Tran}
	DescC       = Descriptor{Mask: Comp}
	DescS       = Descriptor{Mask: Structure}
	DescCT1     = Descriptor{Mask: Comp, Inp1: Tran}
	DescST1     = Descriptor{Mask: Structure, Inp1: Tran}
	DescCT0     = Descriptor{Mask: Comp, Inp0: Tran}
	DescST0     = Descriptor{Mask: Structure, Inp0: Tran}
	DescCT0T1   = Descriptor{Mask: Comp, Inp0: Tran, Inp1: Tran}
	DescST0T1   = Descriptor{Mask: Structure, Inp0: Tran, Inp1: Tran}
	DescSC      = Descriptor{Mask: Structure | Comp}
	DescSCT1    = Descriptor{Mask: Structure | Comp, Inp1: Tran}
	DescSCT0    = Descriptor{Mask: Structure | Comp, Inp0: Tran}
	DescSCT0T1  = Descriptor{Mask: Structure | Comp, Inp0: Tran, Inp1: Tran}
	DescR       = Descriptor{Outp: Replace}
	DescRT1     = Descriptor{Outp: Replace, Inp1: Tran}
	DescRT0     = Descriptor{Outp: Replace, Inp0: Tran}
	DescRT0T1   = Descriptor{Outp: Replace, Inp0: Tran, Inp1: Tran}
	DescRC      = Descriptor{Outp: Replace, Mask: Comp}
	DescRS      = Descriptor{Outp: Replace, Mask: Structure}
	DescRCT1    = Descriptor{Outp: Replace, Mask: Comp, Inp1: Tran}
	DescRST1    = Descriptor{Outp: Replace, Mask: Structure, Inp1: Tran}
	DescRCT0    = Descriptor{Outp: Replace, Mask: Comp, Inp0: Tran}
	DescRST0    = Descriptor{Outp: Replace, Mask: Structure, Inp0: Tran}
	DescRCT0T1  = Descriptor{Outp: Replace, Mask: Comp, Inp0: Tran, Inp1: Tran}
	DescRST0T1  = Descriptor{Outp: Replace, Mask: Structure, Inp0: Tran, Inp1: Tran}
	DescRSC     = Descriptor{Outp: Replace, Mask: Structure | Comp}
	DescRSCT1   = Descriptor{Outp: Replace, Mask: Structure | Comp, Inp1: Tran}
	DescRSCT0   = Descriptor{Outp: Replace, Mask: Structure | Comp, Inp0: Tran}
	DescRSCT0T1 = Descriptor{Outp: Replace, Mask: Structure | Comp, Inp0: Tran, Inp1: Tran}
)

func DescriptorNew() Descriptor {
	return make(Descriptor, nDescFields)
}

func (d *Descriptor) Set(field DescField, value DescValue) error {
	switch field {
	case Outp:
		switch value {
		case Replace:
			if int(Outp) >= len(*d) {
				nd := make(Descriptor, nDescFields)
				copy(nd, *d)
				*d = nd
			}
			(*d)[Outp] = Replace
			return nil
		}
	case Mask:
		switch value {
		case Structure, Comp:
			if int(Mask) >= len(*d) {
				nd := make(Descriptor, nDescFields)
				copy(nd, *d)
				*d = nd
			}
			(*d)[Mask] |= value
			return nil
		}
	case Inp0, Inp1:
		switch value {
		case Tran:
			if int(field) >= len(*d) {
				nd := make(Descriptor, nDescFields)
				copy(nd, *d)
				*d = nd
			}
			(*d)[field] = Tran
			return nil
		}
	}
	return InvalidValue
}

func (d Descriptor) Is(field DescField, value DescValue) (bool, error) {
	switch field {
	case Outp:
		switch value {
		case Replace:
			if int(Outp) >= len(d) {
				return false, nil
			}
			return d[Outp] == Replace, nil
		}
	case Mask:
		switch value {
		case Structure, Comp:
			if int(Mask) >= len(d) {
				return false, nil
			}
			return d[Mask]&value != 0, nil
		}
	case Inp0, Inp1:
		switch value {
		case Tran:
			if int(field) >= len(d) {
				return false, nil
			}
			return d[field] == Tran, nil
		}
	}
	return false, InvalidValue
}
