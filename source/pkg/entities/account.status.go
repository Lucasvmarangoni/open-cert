package entities

type Status struct {
	Blocked  int
	Regular  int
	Suspend  int
	Analysis int
	Inactive int
}

type ReadUserStatus interface {
	GetBlocked() int
	GetRegular() int
	GetSuspend() int
}

func (s *Status) GetBlocked() int { return s.Blocked }
func (s *Status) GetRegular() int { return s.Regular }
func (s *Status) GetSuspend() int { return s.Suspend }

func User_GetStatus() ReadUserStatus {
	return &Status{
		Blocked: 1,
		Regular: 2,
		Suspend: 3,
	}
}

type UserStatus int
const (
	UserBlocked UserStatus = 1
	UserRegular UserStatus = 2
	UserSuspend UserStatus = 3
)

type ReadcertifierStatus interface {
	GetBlocked() int
	GetRegular() int
	GetSuspend() int
	GetAnalysis() int
	GetInactive() int
}

func (s *Status) GetAnalysis() int { return s.Analysis }
func (s *Status) GetInactive() int { return s.Inactive }

func Certifier_GetStatus() ReadcertifierStatus {
	return &Status{
		Blocked:  1,
		Regular:  2,
		Suspend:  3,
		Analysis: 4,
		Inactive: 5,
	}
}

type CertifierStatus int
const (
	CertifierBlocked  CertifierStatus = 1
	CertifierRegular  CertifierStatus = 2
	CertifierSuspend  CertifierStatus = 3
	CertifierAnalysis CertifierStatus = 4
	CertifierInactive CertifierStatus = 5
)
