package choices

type AccType string

const (
	ACCTYPE_READER AccType = "READER"
	ACCTYPE_WRITER AccType = "WRITER"
)
func (a AccType) IsValid() bool {
	switch a {
	case ACCTYPE_READER, ACCTYPE_WRITER:
		return true
	}
	return false
}

type PaymentType string

const (
	PTYPE_GPAY   PaymentType = "GOOGLE PAY"
	PTYPE_STRIPE PaymentType = "STRIPE"
	PTYPE_PAYPAL PaymentType = "PAYPAL"
)
func (p PaymentType) IsValid() bool {
	switch p {
	case PTYPE_GPAY, PTYPE_STRIPE, PTYPE_PAYPAL:
		return true
	}
	return false
}

type PaymentPurpose string

const (
	PP_COINS   PaymentPurpose = "COINS"
	PP_SUB PaymentPurpose = "SUBSCRIPTION"
)
func (p PaymentPurpose) IsValid() bool {
	switch p {
	case PP_COINS, PP_SUB:
		return true
	}
	return false
}

type PaymentStatus string

const (
	PSPENDING   PaymentStatus = "PENDING"
	PSSUCCEEDED PaymentStatus = "SUCCEEDED"
	PSFAILED    PaymentStatus = "FAILED"
	PSCANCELED  PaymentStatus = "CANCELED"
)
func (p PaymentStatus) IsValid() bool {
	switch p {
	case PSPENDING, PSSUCCEEDED, PSFAILED, PSCANCELED:
		return true
	}
	return false
}

type AgeType int

const (
	ATYPE_FOUR     AgeType = 4
	ATYPE_TWELVE   AgeType = 12
	ATYPE_SIXTEEN  AgeType = 16
	ATYPE_EIGHTEEN AgeType = 18
)
func (a AgeType) IsValid() bool {
	switch a {
	case ATYPE_FOUR, ATYPE_TWELVE, ATYPE_SIXTEEN, ATYPE_EIGHTEEN:
		return true
	}
	return false
}

type ChapterStatus string

const (
	CS_DRAFT     ChapterStatus = "DRAFT"
	CS_PUBLISHED ChapterStatus = "PUBLISHED"
	CS_TRASH     ChapterStatus = "TRASH"
)
func (c ChapterStatus) IsValid() bool {
	switch c {
	case CS_DRAFT, CS_PUBLISHED, CS_TRASH:
		return true
	}
	return false
}

type RatingChoice int

const (
	RC_1 RatingChoice = 1
	RC_2 RatingChoice = 2
	RC_3 RatingChoice = 3
	RC_4 RatingChoice = 4
	RC_5 RatingChoice = 5
)
func (r RatingChoice) IsValid() bool {
	switch r {
	case RC_1, RC_2, RC_3, RC_4, RC_5:
		return true
	}
	return false
}

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
func (n NotificationTypeChoice) IsValid() bool {
	switch n {
	case NT_LIKE, NT_REPLY, NT_FOLLOWING, NT_BOOK_PURCHASE, NT_GIFT, NT_REVIEW, NT_VOTE:
		return true
	}
	return false
}

type NotificationStatus string

const (
	NS_CREATED NotificationStatus = "CREATED"
	NS_DELETED NotificationStatus = "DELETED"
)
func (n NotificationStatus) IsValid() bool {
	switch n {
	case NS_CREATED, NS_DELETED:
		return true
	}
	return false
}

type ContractTypeChoice string

const (
	CT_EXCLUSIVE      ContractTypeChoice = "EXCLUSIVE"
	CT_NON_EXCLUSIVE  ContractTypeChoice = "NON-EXCLUSIVE"
	CT_ONLY_EXCLUSIVE ContractTypeChoice = "ONLY-EXCLUSIVE"
)
func (c ContractTypeChoice) IsValid() bool {
	switch c {
	case CT_EXCLUSIVE, CT_NON_EXCLUSIVE, CT_ONLY_EXCLUSIVE:
		return true
	}
	return false
}

type ContractIDTypeChoice string

const (
	CID_DRIVERS_LICENSE ContractIDTypeChoice = "DRIVERS-LICENSE"
	CID_GOVERNMENT_ID   ContractIDTypeChoice = "GOVERNMENT-ID"
	CID_PASSPORT        ContractIDTypeChoice = "PASSPORT"
)
func (c ContractIDTypeChoice) IsValid() bool {
	switch c {
	case CID_DRIVERS_LICENSE, CID_GOVERNMENT_ID, CID_PASSPORT:
		return true
	}
	return false
}

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

type SubscriptionTypeChoice string

const (
	ST_MONTHLY  SubscriptionTypeChoice = "MONTHLY"
	ST_ANNUAL SubscriptionTypeChoice = "ANNUAL"
)

func (s SubscriptionTypeChoice) IsValid() bool {
	switch s {
	case ST_MONTHLY, ST_ANNUAL:
		return true
	}
	return false
}