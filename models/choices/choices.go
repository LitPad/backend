package choices

type AccType string

const (
	ACCTYPE_READER AccType = "READER"
	ACCTYPE_WRITER AccType = "WRITER"
)

type PaymentType string

const (
	PTYPE_GPAY   PaymentType = "GOOGLE PAY"
	PTYPE_STRIPE PaymentType = "STRIPE"
	PTYPE_PAYPAL PaymentType = "PAYPAL"
)

type PaymentStatus string

const (
	PSPENDING   PaymentStatus = "PENDING"
	PSSUCCEEDED PaymentStatus = "SUCCEEDED"
	PSFAILED    PaymentStatus = "FAILED"
	PSCANCELED  PaymentStatus = "CANCELED"
)

type AgeType int

const (
	ATYPE_FOUR     AgeType = 4
	ATYPE_TWELVE   AgeType = 12
	ATYPE_SIXTEEN  AgeType = 16
	ATYPE_EIGHTEEN AgeType = 18
)

type ChapterStatus string

const (
	CS_DRAFT     ChapterStatus = "DRAFT"
	CS_PUBLISHED ChapterStatus = "PUBLISHED"
	CS_TRASH     ChapterStatus = "TRASH"
)

type RatingChoice int

const (
	RC_1 RatingChoice = 1
	RC_2 RatingChoice = 2
	RC_3 RatingChoice = 3
	RC_4 RatingChoice = 4
	RC_5 RatingChoice = 5
)

type NotificationTypeChoice string

const (
	NT_LIKE          NotificationTypeChoice = "LIKE"
	NT_REPLY         NotificationTypeChoice = "REPLY"
	NT_FOLLOWING     NotificationTypeChoice = "FOLLOWING"
	NT_BOOK_PURCHASE NotificationTypeChoice = "BOOK_PURCHASE"
	NT_GIFT          NotificationTypeChoice = "GIFT"
	NT_REVIEW        NotificationTypeChoice = "REVIEW"
	NT_VOTE          NotificationTypeChoice = "VOTE"
)

type NotificationStatus string

const (
	NS_CREATED NotificationStatus = "CREATED"
	NS_DELETED NotificationStatus = "DELETED"
)

type ContractTypeChoice string

const (
	CT_EXCLUSIVE      ContractTypeChoice = "EXCLUSIVE"
	CT_NON_EXCLUSIVE  ContractTypeChoice = "NON-EXCLUSIVE"
	CT_ONLY_EXCLUSIVE ContractTypeChoice = "ONLY-EXCLUSIVE"
)

type ContractIDTypeChoice string

const (
	CID_DRIVERS_LICENSE ContractIDTypeChoice = "DRIVERS-LICENSE"
	CID_GOVERNMENT_ID   ContractIDTypeChoice = "GOVERNMENT-ID"
	CID_PASSPORT        ContractIDTypeChoice = "PASSPORT"
)

type ContractStatusChoice string

const (
	CTS_PENDING  ContractStatusChoice = "PENDING"
	CTS_APPROVED ContractStatusChoice = "APPROVED"
	CTS_DECLINED ContractStatusChoice = "DECLINED"
	CTS_UPDATED  ContractStatusChoice = "UPDATED"
)

func (c ContractStatusChoice) IsValid() bool {
	if c == "" {
		return true
	}

	switch c {
	case CTS_PENDING, CTS_APPROVED, CTS_DECLINED, CTS_UPDATED:
		return true
	}

	return false
}