package btn

import "github.com/a-h/templ"

const baseStyle = "cursor-pointer inline-flex justify-center items-center gap-2 whitespace-nowrap rounded-radius text-center font-medium tracking-wide transition hover:opacity-75 active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed focus-visible:outline-2 focus-visible:outline-offset-2"

const solidBaseStyle = "border transition"

const outlineBaseStyle = "bg-transparent border transition"

const ghostBaseStyle = "bg-transparent"

type sizeType string

const (
	SizeIcon sizeType = "p-0 p-y text-xs"
	SizeSm   sizeType = "px-4 py-2 text-xs"
	SizeMd   sizeType = "px-4 py-2 text-sm"
	SizeLg   sizeType = "px-4 py-2 text-base"
	SizeXl   sizeType = "px-4 py-2 text-lg"
)

type Appearance interface {
	GetType() string
	GetSize() string
	GetBase() string
}

type solid struct {
	vartype solidVariantType
	size    sizeType
	base    string
}

func (v *solid) GetType() string {
	return string(v.vartype)
}

func (v *solid) GetSize() string {
	return string(v.size)
}

func (v *solid) GetBase() string {
	return v.base
}

func NewSolidBtn(t solidVariantType, size sizeType) Appearance {
	return &solid{
		vartype: t,
		size:    size,
		base:    solidBaseStyle,
	}
}

var DefaultBtn = NewSolidBtn(SolidVariantPrimary, SizeMd)

type outline struct {
	vartype outlineVariantType
	size    sizeType
	base    string
}

func (v *outline) GetType() string {
	return string(v.vartype)
}

func (v *outline) GetSize() string {
	return string(v.size)
}

func (v *outline) GetBase() string {
	return v.base
}

func NewOutlineBtn(t outlineVariantType, size sizeType) Appearance {
	return &outline{
		vartype: t,
		size:    size,
		base:    outlineBaseStyle,
	}
}

type ghost struct {
	vartype ghostVariantType
	size    sizeType
	base    string
}

func (v *ghost) GetType() string {
	return string(v.vartype)
}

func (v *ghost) GetSize() string {
	return string(v.size)
}

func (v *ghost) GetBase() string {
	return v.base
}

func NewGhostBtn(t ghostVariantType, size sizeType) Appearance {
	return &ghost{
		vartype: t,
		size:    size,
		base:    ghostBaseStyle,
	}
}

type solidVariantType string

type outlineVariantType string
type ghostVariantType string

const (
	SolidVariantPrimary   solidVariantType = "bg-primary border border-primary text-on-primary focus-visible:outline-primary dark:bg-primary-dark dark:border-primary-dark dark:text-on-primary-dark dark:focus-visible:outline-primary-dark"
	SolidVariantSecondary solidVariantType = "bg-secondary border border-secondary text-on-secondary focus-visible:outline-secondary dark:bg-secondary-dark dark:border-secondary-dark dark:text-on-secondary-dark dark:focus-visible:outline-secondary-dark"
	SolidVariantAlternate solidVariantType = "bg-surface-alt border border-surface-alt text-on-surface-strong focus-visible:outline-surface-alt dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark-strong dark:focus-visible:outline-surface-dark-alt"
	SolidVariantInverse   solidVariantType = "bg-surface-dark border border-surface-dark text-on-surface-dark focus-visible:outline-surface-dark dark:bg-surface dark:border-surface dark:text-on-surface dark:focus-visible:outline-surface"
	SolidVariantInfo      solidVariantType = "bg-info border border-info text-on-info focus-visible:outline-info dark:bg-info dark:border-info dark:text-on-info dark:focus-visible:outline-info"
	SolidVariantDanger    solidVariantType = "bg-danger border border-danger text-on-danger focus-visible:outline-danger dark:bg-danger dark:border-danger dark:text-on-danger dark:focus-visible:outline-danger"
	SolidVariantWarning   solidVariantType = "bg-warning border border-warning text-on-warning focus-visible:outline-warning dark:bg-warning dark:border-warning dark:text-on-warning dark:focus-visible:outline-warning"
	SolidVariantSuccess   solidVariantType = "bg-success border border-success text-on-success focus-visible:outline-success dark:bg-success dark:border-success dark:text-on-success dark:focus-visible:outline-success"
)

const (
	OutlineVariantPrimary   = "border-primary text-primary focus-visible:outline-primary dark:border-primary-dark dark:text-primary-dark dark:focus-visible:outline-primary-dark"
	OutlineVariantSecondary = "border-secondary text-secondary focus-visible:outline-secondary dark:border-secondary-dark dark:text-secondary-dark dark:focus-visible:outline-secondary-dark"
	OutlineVariantAlternate = "border-outline text-outline focus-visible:outline-outline dark:border-outline-dark dark:text-outline-dark dark:focus-visible:outline-outline-dark"
	OutlineVariantInverse   = "border-surface-dark text-surface-dark focus-visible:outline-surface-dark dark:border-surface dark:text-surface dark:focus-visible:outline-surface"
	OutlineVariantInfo      = "border-info text-info focus-visible:outline-info dark:border-info dark:text-info dark:focus-visible:outline-info"
	OutlineVariantDanger    = "border-danger text-danger focus-visible:outline-danger dark:border-danger dark:text-danger dark:focus-visible:outline-danger"
	OutlineVariantWarning   = "border-warning text-warning focus-visible:outline-warning dark:border-warning dark:text-warning dark:focus-visible:outline-warning"
	OutlineVariantSuccess   = "border-success text-success focus-visible:outline-success dark:border-success dark:text-success dark:focus-visible:outline-success"
)

const (
	GhostVariantPrimary   = "text-primary focus-visible:outline-primary dark:text-primary-dark dark:focus-visible:outline-primary-dark"
	GhostVariantSecondary = "text-secondary focus-visible:outline-secondary dark:text-secondary-dark dark:focus-visible:outline-secondary-dark"
	GhostVariantAlternate = "text-outline focus-visible:outline-outline dark:text-outline-dark dark:focus-visible:outline-outline-dark"
	GhostVariantInverse   = "text-surface-dark focus-visible:outline-surface-dark dark:text-surface dark:focus-visible:outline-surface"
	GhostVariantInfo      = "text-info focus-visible:outline-info dark:text-info dark:focus-visible:outline-info"
	GhostVariantDanger    = "text-danger focus-visible:outline-danger dark:text-danger dark:focus-visible:outline-danger"
	GhostVariantWarning   = "text-warning focus-visible:outline-warning dark:text-warning dark:focus-visible:outline-warning"
	GhostVariantSuccess   = "text-success focus-visible:outline-success dark:text-success dark:focus-visible:outline-success"
)

type Props struct {
	Class      string
	Appearance Appearance
	Attrs      templ.Attributes
}
