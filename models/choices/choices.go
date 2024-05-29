package choices

type AccType string

const (
	ACCTYPE_READER AccType = "READER"
	ACCTYPE_WRITER AccType = "WRITER"
)

type PaymentType string

const (
	PTYPE_GPAY PaymentType = "GOOGLE PAY"
	PTYPE_STRIPE PaymentType = "STRIPE"
	PTYPE_PAYPAL PaymentType = "PAYPAL"
)

type PaymentStatus string

const (
	PSPENDING PaymentStatus = "PENDING"
	PSSUCCEEDED PaymentStatus = "SUCCEEDED"
	PSFAILED PaymentStatus = "FAILED"
	PSCANCELED PaymentStatus = "CANCELED"
)

type AgeType int

const (
	ATYPE_FOUR AgeType = 4
	ATYPE_TWELVE AgeType = 12
	ATYPE_SIXTEEN AgeType = 16
	ATYPE_EIGHTEEN AgeType = 18
)

type ChapterStatus string

const (
	CS_DRAFT ChapterStatus = "DRAFT"
	CS_PUBLISHED ChapterStatus = "PUBLISHED"
	CS_TRASH ChapterStatus = "TRASH"
)