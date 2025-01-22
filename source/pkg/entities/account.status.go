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
		Suspend: 4,
	}
}

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
		Analysis: 3,
		Suspend:  4,
		Inactive: 5,
	}
}
